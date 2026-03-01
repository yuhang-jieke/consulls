package consul

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

// ServiceInfo 服务信息结构体
type ServiceInfo struct {
	ID      string
	Name    string
	Address string
	Port    int
	Tags    []string
	Healthy bool
}

// DiscoverService 发现指定名称的健康服务
// serviceName: 服务名称
// 返回健康的服务实例列表
func (c *Client) DiscoverService(serviceName string) ([]ServiceInfo, error) {
	// 使用 Health().Service() 获取健康的服务实例
	services, _, err := c.apiClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover service %s: %w", serviceName, err)
	}

	var result []ServiceInfo
	for _, entry := range services {
		// 检查服务是否通过健康检查
		isHealthy := true
		for _, check := range entry.Checks {
			if check.Status != api.HealthPassing {
				isHealthy = false
				break
			}
		}

		service := ServiceInfo{
			ID:      entry.Service.ID,
			Name:    entry.Service.Service,
			Address: "127.0.0.1",
			Port:    entry.Service.Port,
			Tags:    entry.Service.Tags,
			Healthy: isHealthy,
		}
		result = append(result, service)
	}

	return result, nil
}

// DiscoverServiceWithTag 根据标签发现服务
// serviceName: 服务名称
// tag: 服务标签
func (c *Client) DiscoverServiceWithTag(serviceName, tag string) ([]ServiceInfo, error) {
	services, _, err := c.apiClient.Health().Service(serviceName, tag, true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover service %s with tag %s: %w", serviceName, tag, err)
	}

	var result []ServiceInfo
	for _, entry := range services {
		isHealthy := true
		for _, check := range entry.Checks {
			if check.Status != api.HealthPassing {
				isHealthy = false
				break
			}
		}

		service := ServiceInfo{
			ID:      entry.Service.ID,
			Name:    entry.Service.Service,
			Address: "127.0.0.1",
			Port:    entry.Service.Port,
			Tags:    entry.Service.Tags,
			Healthy: isHealthy,
		}
		result = append(result, service)
	}

	return result, nil
}

// GetAllServices 获取所有注册的服务
func (c *Client) GetAllServices() (map[string][]string, error) {
	services, err := c.apiClient.Agent().Services()
	if err != nil {
		return nil, fmt.Errorf("failed to get all services: %w", err)
	}

	result := make(map[string][]string)
	for id, service := range services {
		result[id] = service.Tags
	}

	return result, nil
}

// GetService 获取单个服务的详细信息
// serviceID: 服务唯一标识
func (c *Client) GetService(serviceID string) (*api.AgentService, error) {
	service, _, err := c.apiClient.Agent().Service(serviceID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get service %s: %w", serviceID, err)
	}

	return service, nil
}

// ServiceChangeHandler 服务变化回调函数类型
type ServiceChangeHandler func(services []ServiceInfo)

// WatchService 监听服务变化（使用 Watch 机制）
// serviceName: 服务名称
// handler: 服务变化时的回调函数
// 返回停止函数，调用可停止监听
func (c *Client) WatchService(ctx context.Context, serviceName string, handler ServiceChangeHandler) (stop func(), err error) {
	// 创建 Watch plan
	plan, err := watch.Parse(map[string]interface{}{
		"type":    "service",
		"service": serviceName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse watch plan: %w", err)
	}

	// 设置回调函数
	plan.Handler = func(idx uint64, raw interface{}) {
		if raw == nil {
			log.Printf("[Consul Watch] No data received for service %s", serviceName)
			handler(nil)
			return
		}

		// 解析服务列表
		services, ok := raw.([]*api.ServiceEntry)
		if !ok {
			log.Printf("[Consul Watch] Unexpected data type for service %s", serviceName)
			handler(nil)
			return
		}

		var result []ServiceInfo
		for _, entry := range services {
			// 只处理通过健康检查的服务
			isHealthy := true
			for _, check := range entry.Checks {
				if check.Status != api.HealthPassing {
					isHealthy = false
					break
				}
			}

			if isHealthy {
				service := ServiceInfo{
					ID:      entry.Service.ID,
					Name:    entry.Service.Service,
					Address: entry.Service.Address,
					Port:    entry.Service.Port,
					Tags:    entry.Service.Tags,
					Healthy: isHealthy,
				}
				result = append(result, service)
			}
		}

		log.Printf("[Consul Watch] Service %s updated, healthy instances: %d", serviceName, len(result))
		handler(result)
	}

	// 在后台运行 watch
	go func() {
		if err := plan.Run(c.config.Address); err != nil {
			log.Printf("[Consul Watch] Error watching service %s: %v", serviceName, err)
		}
	}()

	// 返回停止函数
	stop = func() {
		plan.Stop()
		log.Printf("[Consul Watch] Stopped watching service %s", serviceName)
	}

	return stop, nil
}
