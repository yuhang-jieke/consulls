package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

// Client 封装 Consul API 客户端
type Client struct {
	apiClient *api.Client
	config    *api.Config
}

// NewClient 创建新的 Consul 客户端
// addr: Consul 服务器地址，如 "localhost:8500"
func NewClient(addr string) (*Client, error) {
	config := api.DefaultConfig()
	if addr != "" {
		config.Address = addr
	}

	apiClient, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	return &Client{
		apiClient: apiClient,
		config:    config,
	}, nil
}

// RawClient 返回底层的 consul/api 客户端
func (c *Client) RawClient() *api.Client {
	return c.apiClient
}

// Health 返回健康检查客户端
func (c *Client) Health() *api.Health {
	return c.apiClient.Health()
}

// Agent 返回 Agent 客户端
func (c *Client) Agent() *api.Agent {
	return c.apiClient.Agent()
}
