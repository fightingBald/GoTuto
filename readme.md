# 项目搭建过程 

基础 infra / tooling（Tilt / kind / scripts / Docker / .env / Makefile / .dockerignore / Helm chart） 目标：把与集群/本地开发运行相关的文件先提交，展示搭建环境所需的基础资源

批次 2 — DB 配置与迁移（migrations、embedded 支持、postgres adapter） 目标：把数据库迁移文件和 postgres adapter 提交，说明已准备好 DB 层。 示例命令：
git add apps/product-query-svc/adapters/postgres/ apps/product-query-svc/adapters/postgres/migrations_embedded.go internal/adapters/db/pg_product_repo.go

批次 3 — 核心应用层（domain、app、ports） 目标：提交领域模型、用例/服务逻辑、ports 接口，展示业务代码骨架就绪。 示例命令：
git add apps/product-query-svc/domain/ apps/product-query-svc/app/ apps/product-query-svc/ports/

批次 4 — 适配器：in-memory repo 与 HTTP handlers 目标：先提交 inmem 实现与 HTTP 适配器（使服务能在不依赖真实 DB 下跑起来）。 示例命令：
git add apps/product-query-svc/adapters/inmem/ apps/product-query-svc/adapters/http/

批次 4 — 适配器：in-memory repo 与 HTTP handlers 目标：先提交 inmem 实现与 HTTP 适配器（使服务能在不依赖真实 DB 下跑起来）。 示例命令：
git add apps/product-query-svc/adapters/inmem/ apps/product-query-svc/adapters/http/


# 代码生成（Code Generation）

## 概述

要的文件：api/openapi.yaml（入口） + paths/ responses/ schemas（内容） + ports/marketplaceapi/oapi-config.yaml（配置） + ports/marketplaceapi/generate.go（命令钩��）。
要的命令：go mod init（一次）→ go get（一次）→ go generate ./ports/marketplaceapi（每次改规范后）。
流程：写 spec → 配置生成 → 跑生成 → 实现接口 → 起服务测通。

坑点 oapi-codegen v2， oapi-config等要遵守V2规范
要的命令（简要步骤）：
## 生成所需文件、命令与流程（简要）
3. 每次改 spec 后运行：
要的文件：
- api/openapi.yaml（入口）
- api/paths/、api/responses/、api/schemas/（规范拆分的内容）
- ports/marketplaceapi/oapi-config.yaml（oapi-codegen 配置）
- ports/marketplaceapi/generate.go（命令钩子，包含 //go:generate 注释）
典型流程（工程化建议）：
要的命令（简要步骤）：
1. go mod init （仅首次）
2. go get / go install 指定生成工具（仅首次或升级时）
3. 每次改 spec 后运行：
    ```shell
    go generate ./ports/marketplaceapi
    ```
典型流程（工程化建议）：
1. 写 OpenAPI spec（api/openapi.yaml + paths/ responses/ schemas/）
2. 配置生成器（ports/marketplaceapi/oapi-config.yaml）(注意用新版本)
3. 在 ports/marketplaceapi/generate.go 或源文件中添加 //go:generate，并运行 go generate
4. 实现生成的接口（实现 ServerInterface 等）
5. 启动服务并通过集成/端到端测试验证

说明：把生成步骤封装到 Makefile 或 generate.sh，有利于团队复现；关键生成产物（如 pb.go、marketplace.gen.go 等）可以考虑纳入版本控制或在 CI 中强制生成并校验差异。

## 项目启动（开发）

### 前置条件
- Go >= 1.22
- Docker
- 可选：kubectl、helm、tilt、psql 客户端

### 本地快速启动（不使用 k8s）
1. 启动 PostgreSQL（示例）：

```sh
docker run --name marketplace-postgres \
  -e POSTGRES_USER=app \
  -e POSTGRES_PASSWORD=app_password \
  -e POSTGRES_DB=productdb \
  -p 5432:5432 -d postgres:16-alpine
```

2. 设置环境变量（示例 `.env` 或导出到 shell）：

```sh
export DATABASE_URL="postgres://app:app_password@localhost:5432/productdb?sslmode=disable"
export HTTP_ADDRESS=":8080"
export LOG_LEVEL=debug
```

3. 运行服务（开发模式）：

```sh
cd backend/cmd/marketplace/product-query-svc
go run .
```

或构建后运行：

```sh
go build -o bin/product-query-svc ./backend/cmd/marketplace/product-query-svc
./bin/product-query-svc
```

> 说明：项目内有迁移文件位于 apps/product-query-svc/adapters/postgres/migrations，代码中有 migrations_embedded.go（//go:embed 支持）；服务启动时应会处理嵌入的迁移，若未自动执行，可在容器启动或 initContainer 中运行迁移脚本。

### 在 Kubernetes（kind + Tilt）上启动
1. 启动 kind 集群（本仓库有脚本）：

```sh
bash scripts/kind-up.sh
```

2. 启动 Tilt（会构建镜像并部署 charts/ 或 k8s/ 下的资源）：

```sh
tilt up
```

3. 访问服务：
- 通过 Tilt 的端口转发： http://localhost:8080
- 若使用 kind 的 NodePort 映射（见 kind/kind-cluster.yaml）： http://localhost:30080

4. 若需要直接连接数据库：

```sh
psql "postgres://app:app_password@localhost:5432/productdb"
```

### 使用 Helm（项目包含最小 Chart）
- charts/product-query-svc 下包含 chart，可用 helm install 或 helm upgrade --install 部署到集群。

### 构建 Docker 镜像（用于 Tilt 或手动部署）
```sh
docker build -t product-query-svc:dev .
```

### 常见问题
- 报错 "FATAL: database \"app\" does not exist"：确认 Postgres 启动时环境变量 POSTGRES_DB 与服务的 DATABASE_URL 中的数据库名一致（本示例使用 productdb）；或手动创建所需数据库。
- 构建时报 go.mod: unknown directive: tool：确保构建使用的 Go 版本 >= 1.22（Dockerfile 中使用 golang:1.22-alpine）。
- 在 Lens 看不到集群资源：确认 Lens 指向的 kubeconfig 与当前 kubectl 上下文一致（kubectl config current-context），以及资源所在命名空间（marketplace-dev）。
