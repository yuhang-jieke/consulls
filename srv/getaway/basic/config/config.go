package config

// Config 应用配置
type Config struct {
	Consul ConsulConfig
	Server ServerConfig
	Order  OrderConfig
}

// ConsulConfig Consul 配置
type ConsulConfig struct {
	Address string
}

// ServerConfig HTTP 服务器配置
type ServerConfig struct {
	Port int
}

// OrderConfig Order 服务配置
type OrderConfig struct {
	ServiceName string
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Consul: ConsulConfig{
			Address: "115.190.57.118:8500",
		},
		Server: ServerConfig{
			Port: 8080,
		},
		Order: OrderConfig{
			ServiceName: "user-service",
		},
	}
}
