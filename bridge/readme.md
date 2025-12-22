bridge/
 ├── cmd/
 │    └── bridge/
 │         └── main.go
 ├── internal/
 │    ├── registry/     // 服务注册中心
 │    ├── router/       // 路由选择
 │    ├── balancer/     // 权重轮询
 │    ├── limiter/      // per-service 限流
 │    ├── health/       // downstream health
 │    ├── httpx/        // HTTP server + proxy
 │    ├── grpcx/        // gRPC proxy
 │    └── metrics/
 ├── api/
 │    └── admin.go      // 注册 / 删除 / 查询
 └── go.mod

bridge 是一个基于 gRPC 的服务网格代理，它实现了服务发现、负载均衡、限流、健康检查等功能。

## 使用方法
1. 准备TLS证书： `server.crt`, `server.key`, `ca.crt`
2. 启动 bridge： 
   ```
   go run cmd/bridge/main.go
   ```
3. 注册服务： 
   ```
   curl "http://localhost:8080/admin/register?name=txstore&addr=localhost:8001&proto=grpc&weight=80&gray=0&qps=200"
   curl "http://localhost:8080/admin/register?name=txstore&addr=localhost:8002&proto=grpc&weight=20&gray=50&qps=200"
   ```

4. Prometheus metrics：访问 `http://localhost:9100/metrics`
5. HTTP → gRPC 或直接 gRPC 调用桥接 downstream