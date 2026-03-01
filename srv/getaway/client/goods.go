package client

import (
	"context"
	"fmt"
	"github.com/yuhang-jieke/consulls/srv/pkg/consul"
	__ "github.com/yuhang-jieke/consulls/srv/user-server/handler/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
)

type OrderClient struct {
	clients     map[string]*grpc.ClientConn
	serviceName string
	consulAddr  string
	consul      *consul.Client
	mu          sync.RWMutex
	index       int
	stopWatch   func()
}

func NewOrderClient(consulAddr, serviceName string) (*OrderClient, error) {
	c := &OrderClient{
		clients:     make(map[string]*grpc.ClientConn),
		serviceName: serviceName,
		consulAddr:  consulAddr,
	}
	consulClient, err := consul.NewClient(consulAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}
	c.consul = consulClient

	services, err := consulClient.DiscoverService(serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to discover service %s: %w", serviceName, err)
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("no healthy instances found for service %s", serviceName)
	}

	c.updateClients(services)

	stopWatch, err := consulClient.WatchService(context.Background(), serviceName, func(services []consul.ServiceInfo) {
		log.Printf("[OrderClient] Received service change notification, updating clients...")
		c.updateClients(services)
	})
	if err != nil {
		log.Printf("[OrderClient] Failed to start watch: %v, continuing without watch", err)
	} else {
		c.stopWatch = stopWatch
		log.Printf("[OrderClient] Started watching service %s for changes", serviceName)
	}

	log.Printf("[OrderClient] Successfully initialized with %d service instances", len(c.clients))
	return c, nil
}

func (c *OrderClient) updateClients(services []consul.ServiceInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()

	currentClients := make(map[string]*grpc.ClientConn)

	for _, s := range services {
		addr := fmt.Sprintf("%s:%d", "127.0.0.1", s.Port)

		if conn, ok := c.clients[addr]; ok {
			currentClients[addr] = conn
			log.Printf("[OrderClient] Reusing existing connection: %s", addr)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		conn, err := grpc.DialContext(ctx, addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		cancel()

		if err != nil {
			log.Printf("1")
			log.Printf("[OrderClient] Failed to connect to %s: %v", addr, err)
			continue
		}

		currentClients[addr] = conn
		log.Printf("[OrderClient] Created new connection: %s", addr)
	}

	for addr, conn := range c.clients {
		if _, ok := currentClients[addr]; !ok {
			log.Printf("[OrderClient] Closing removed connection: %s", addr)
			conn.Close()
		}
	}

	c.clients = currentClients
	log.Printf("[OrderClient] Client pool updated, total connections: %d", len(c.clients))
}

func (c *OrderClient) getClient() (__.EcommerceServiceClient, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.clients) == 0 {
		return nil, fmt.Errorf("no available service instances")
	}

	addrs := make([]string, 0, len(c.clients))
	for addr := range c.clients {
		addrs = append(addrs, addr)
	}

	addr := addrs[c.index%len(addrs)]
	c.index++

	return __.NewEcommerceServiceClient(c.clients[addr]), nil
}

func (c *OrderClient) GetClient() (__.EcommerceServiceClient, error) {
	return c.getClient()
}

func (c *OrderClient) Close() {
	if c.stopWatch != nil {
		c.stopWatch()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for addr, conn := range c.clients {
		conn.Close()
		log.Printf("[OrderClient] Closed connection: %s", addr)
	}
	c.clients = make(map[string]*grpc.ClientConn)
}
