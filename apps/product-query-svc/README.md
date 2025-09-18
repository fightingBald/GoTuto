# product-query-svc 架构概览

本应用按“六边形/端口-适配器”组织：

- domain：领域模型与规则（不依赖外层）
- ports：核心对外边界
  - inbound（use case 接口）：`ProductService`
  - outbound（基础设施接口）：`ProductRepo`
- app：用例实现（编排业务，依赖 `ports`，不关心 HTTP/DB）
- adapters：适配器实现
  - inbound/http：实现 OpenAPI 生成的 `ServerInterface`，调用 `ports.ProductService`
  - outbound/inmem、outbound/postgres：实现 `ports.ProductRepo`
- backend/cmd/.../main.go：组装根，选择具体适配器（inmem 或 postgres），注入到 app 层，再挂到 HTTP。

依赖方向：`adapters -> app -> ports <- domain`（领域最内层，适配器最外层）。

## 依赖箭头图（以 GetProduct 为例）

以接口 GET `/products/{id}` 为例，调用路径如下：

```
Client
  ↓
HTTP Inbound Adapter
  apps/product-query-svc/adapters/inbound/http/handler_product_read.go (Server.GetProductByID)
  ↓ 依赖入站端口 ports.ProductService
Ports (Inbound)
  apps/product-query-svc/ports/inbound.go (interface ProductService)
  ↓ 由组装根注入 app 实现
Application (Use Case)
  apps/product-query-svc/app/product_service.go (ProductService.GetProduct)
  ↓ 依赖出站端口 ports.ProductRepo
Ports (Outbound)
  apps/product-query-svc/ports/outbound.go (interface ProductRepo)
  ↓ 由组装根选择并注入具体适配器
Outbound Adapters (Persistence)
  ├─ apps/product-query-svc/adapters/outbound/inmem/product_repository.go (InMemRepo.GetByID)
  └─ apps/product-query-svc/adapters/outbound/postgres/product_repository.go (PGProductRepo.GetByID)
      ↓
Domain
  apps/product-query-svc/domain/product.go (实体/校验)

Composition Root（组装根）
  backend/cmd/marketplace/product-query-svc/main.go
  - 读取配置，选择 inmem 或 postgres 作为 ProductRepo 的实现
  - 构造 app.ProductService，并作为 ports.ProductService 注入 HTTP 适配器
  - 启动 HTTP 服务器
```

示例（写入，用以体现 domain 的作用）：

```
Client
  ↓
HTTP Inbound Adapter
  apps/product-query-svc/adapters/inbound/http/handler_product_write.go (Server.CreateProduct)
  - 将 JSON DTO 映射为 domain.Product（美元转分，调用 NewProduct 校验不变式）
  ↓ 调用入站端口 ports.ProductService.CreateProduct
Application (Use Case)
  apps/product-query-svc/app/product_service.go (CreateProduct)
  - 调用 p.Validate / 富行为 → 通过出站端口持久化
  ↓
Ports (Outbound)
  apps/product-query-svc/ports/outbound.go (ProductRepo.Create)
  ↓
Outbound Adapter
  apps/product-query-svc/adapters/outbound/postgres/inmem (真正落库/内存存储)
  ↓
返回 HTTP（Created + JSON），领域错误映射为 400/404。
```
