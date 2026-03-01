# Consulls 微服务项目

基于 Consul + gRPC 的微服务架构 demo，实现服务注册、发现与调用。

## 项目目录结构


```
├─.idea 
└─srv
    ├─getaway
    │  ├─basic
    │  │  ├─cmd
    │  │  ├─config
    │  │  └─proto
    │  ├─client
    │  ├─handler
    │  │  └─request
    │  └─router
    ├─pkg
    │  └─consul
    └─user-server
        ├─basic
        │  ├─cmd
        │  ├─config
        │  └─inits
        ├─handler
        │  ├─proto
        │  └─server
        └─model
```

## 架构设计

```
┌─────────────┐     ┌─────────────┐     ┌─────────────────┐
│   客户端    │────▶│   网关     │────▶│   用户服务      │
│  (HTTP)     │     │  (getaway) │     │ (user-server)  │
└─────────────┘     └──────┬──────┘     └─────────────────┘
                          │
                          ▼
                   ┌─────────────┐
                   │   Consul   │
                   │  (服务注册) │
                   │  (服务发现) │
                   └─────────────┘
```

### 核心组件

1. **Consul**: 服务注册与发现中心
2. **user-server**: gRPC 服务端，提供商品管理接口
3. **getaway**: API 网关，接收 HTTP 请求并转发到 gRPC 服务

## 启动前准备

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置 Consul

确保 Consul 服务已启动，默认地址：`115.190.57.118:8500`

### 3. 配置 MySQL

在 `srv/dev.yaml` 中配置 MySQL 连接信息

## 启动步骤


### 1. 启动 user-server（服务提供者）

```bash
cd srv/user-server/basic/cmd
go run main.go
```

服务默认监听端口：8081

### 2. 启动 getaway（网关）

```bash
cd srv/getaway/basic/cmd
go run main.go
```

网关默认监听端口：8080

### 3. 验证服务注册

访问 Consul UI：`http://115.190.57.118:8500/ui/`

确认 `user-service` 已注册并健康

## 接口调用

### 添加商品

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "name": "商品名称",
    "price": 99.99,
    "stock": 100
  }'
```

### 响应示例

```json
{
  "code": 200,
  "msg": "商品添加成功",
  "data": {
    "message": "商品添加成功"
  }
}
```

## 接口列表

| 方法   | 路径         | 说明     |
|--------|--------------|----------|
| POST   | /orders      | 添加商品 |
| GET    | /orders/:id | 获取商品 |
| PUT    | /orders/:id | 修改商品 |
| DELETE | /orders/:id | 删除商品 |
| GET    | /orders     | 商品列表 |

## 配置说明

### dev.yaml

```yaml
# MySQL 配置
mysql:
  host: localhost
  port: 3306
  user: root
  password: password
  database: test

# Redis 配置
redis:
  host: localhost
  port: 6379
  password: ""
  database: 0
```
