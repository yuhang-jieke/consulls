package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

// RegisterService 注册服务到 Consul（使用 TTL 健康检查）
// serviceID: 服务唯一标识
// serviceName: 服务名称
// address: 服务地址
// port: 服务端口
// tags: 服务标签
func (c *Client) RegisterService(serviceID, serviceName, address string, port int, tags []string) error {
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: address,
		Port:    port,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			TTL:                            "30s",
			DeregisterCriticalServiceAfter: "60s",
		},
	}

	if err := c.apiClient.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service %s: %w", serviceID, err)
	}

	return nil
}

// RegisterServiceWithTTL 使用 TTL 健康检查注册服务
// serviceID: 服务唯一标识
// serviceName: 服务名称
// address: 服务地址
// port: 服务端口
// tags: 服务标签
// ttl: TTL 时长，如 "30s"
func (c *Client) RegisterServiceWithTTL(serviceID, serviceName, address string, port int, tags []string, ttl string) error {
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: address,
		Port:    port,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			TTL:                            ttl,
			DeregisterCriticalServiceAfter: "60s",
		},
	}

	if err := c.apiClient.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service %s with TTL: %w", serviceID, err)
	}

	return nil
}

// DeregisterService 从 Consul 注销服务
// serviceID: 服务唯一标识
func (c *Client) DeregisterService(serviceID string) error {
	if err := c.apiClient.Agent().ServiceDeregister(serviceID); err != nil {
		return fmt.Errorf("failed to deregister service %s: %w", serviceID, err)
	}

	return nil
}

// UpdateTTL 更新服务的 TTL 健康检查
// serviceID: 服务唯一标识
// output: 健康检查输出信息
// status: 健康状态 (api.HealthPassing, api.HealthWarning, api.HealthCritical)
func (c *Client) UpdateTTL(serviceID, output, status string) error {
	if err := c.apiClient.Agent().UpdateTTL("service:"+serviceID, output, status); err != nil {
		return fmt.Errorf("failed to update TTL for service %s: %w", serviceID, err)
	}

	return nil
}

// PassTTL 标记服务 TTL 检查通过
func (c *Client) PassTTL(serviceID, output string) error {
	return c.UpdateTTL(serviceID, output, api.HealthPassing)
}

// WarnTTL 标记服务 TTL 检查警告
func (c *Client) WarnTTL(serviceID, output string) error {
	return c.UpdateTTL(serviceID, output, api.HealthWarning)
}

// FailTTL 标记服务 TTL 检查失败
func (c *Client) FailTTL(serviceID, output string) error {
	return c.UpdateTTL(serviceID, output, api.HealthCritical)
}
