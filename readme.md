# 项目搭建与开发指南

简短说明：本仓库包含一个示例后端服务 product-query-svc（支持 in-memory 与 Postgres）、数据库迁移（嵌入支持 //go:embed）、以及用于本地开发的 Tilt + kind 配置与最小 Helm chart。

---

## 前置条件
- Go >= 1.22
- Docker
- 可选：kubectl、helm、tilt、psql 客户端、IDE（如 GoLand）

---

## 本地快速启动（单机，不用 k8s）
1. 启动 Postgres（示例）：

```sh
docker run --name marketplace-postgres \
  -e POSTGRES_USER=app \
  -e POSTGRES_PASSWORD=app_password \
  -e POSTGRES_DB=productdb \
  -p 5432:5432 -d postgres:16-alpine
```

2. 设置环境变量：

```sh
export DATABASE_URL="postgres://app:app_password@localhost:5432/productdb?sslmode=disable"
export HTTP_ADDRESS=":8080"
export LOG_LEVEL=debug
```

3. 运行服务（开发）：

```sh
cd backend/cmd/marketplace/product-query-svc
go run .
```

或构建后运行：

```sh
go build -o bin/product-query-svc ./backend/cmd/marketplace/product-query-svc
./bin/product-query-svc
```

说明：项目包含迁移文件（apps/product-query-svc/adapters/postgres/migrations）以及嵌入支持（migrations_embedded.go）。生产入口会在提供 DSN 时连接 Postgres 并可执行嵌入迁移（参见 flags: --db-dsn, --migrate）。

---

## 在 Kubernetes（kind + Tilt）上启动
- 启动本地 kind 集群（仓库脚本）：

```sh
bash scripts/kind-up.sh
```

- 启动 Tilt（构建镜像并部署）：

```sh
tilt up
```

- 访问服务：
  - Tilt 端口转发：http://localhost:8080
  - 可选 NodePort（若在 kind 配置中映射）：http://localhost:30080

- 连接数据库：

```sh
psql "postgres://app:app_password@localhost:5432/productdb"
```

---

## 使用 Helm
仓库包含最小 Helm chart：charts/product-query-svc，可用如下命令部署：

```sh
helm install product-query-svc ./charts/product-query-svc -f ./charts/product-query-svc/values.yaml
```

---

## 构建 Docker 镜像
```sh
docker build -t product-query-svc:dev .
```

注：Dockerfile 默认构建 backend/cmd/marketplace/product-query-svc 的二进制，用于镜像/部署。

---

## 代码生成（OpenAPI -> handlers）
<details>
<summary>展开查看生成说明</summary>


- 必要文件：api/openapi.yaml、api/paths/*、api/schemas/*、api/responses/*
- 配置文件：api/oapi-config.yaml
- 生成命令（示例）：

```sh
go generate ./api
# 或者根据 generate.go 的 //go:generate 指定路径
```

- 建议：将生成步骤写入 Makefile 或 CI，团队协同时要约定是否把生成产物纳入版本控制（两种策略均可）。

</details>

---

## 两个 main 的说明（prod vs dev）
- backend/cmd/marketplace/product-query-svc/main.go
  - 用于容器/生产路径。支持 --db-dsn 与 --migrate，连接 pgx 池并可运行嵌入迁移。
- cmd/server/main.go
  - 简化的本地开发入口，默认使用 in-memory repo，便于快速调试。

建议：保留两者以分离 dev/production 流程，或将 dev 入口重命名为 cmd/dev-server 以消除歧义。

---

## 分批提交建议（用于展示搭建进度）
将改动分为小而原子的一系列 commit，以下为推荐批次。点击展开查看每个批次应包含的文件及示例 commit message。

<details>
<summary>批次 1 — 基础 infra / tooling（Tilt / kind /Docker / scripts / Helm / .env）</summary>

- 相关文件：Tiltfile、kind/, k8s/, Dockerfile、.dockerignore、Makefile、.env.dev.example、charts/
- 建议命令：
  - git add Tiltfile kind/ k8s/ Dockerfile .dockerignore Makefile .env.dev.example charts/
  - git commit -m "infra: add Tilt, kind, k8s manifests, Dockerfile, helper scripts and Helm chart"

</details>

<details>
<summary>批次 2 — DB 配置与迁移</summary>

- 相关文件：apps/product-query-svc/adapters/postgres/（migrations、migrations_embedded.go、product_repository.go）、internal/adapters/db/pg_product_repo.go
- 建议 commit message："db: add Postgres migrations and adapters (migrations embedded via //go:embed)"

</details>

<details>
<summary>批次 3 — 核心应用层（domain、app、ports）</summary>

- 相关文件：apps/product-query-svc/domain/ apps/product-query-svc/app/ apps/product-query-svc/ports/
- 建议 commit message："app: add domain models, application service and ports for product-query-svc"

</details>

<details>
<summary>批次 4 — 适配器：in-memory repo 与 HTTP handlers</summary>

- 相关文件：apps/product-query-svc/adapters/inmem/ apps/product-query-svc/adapters/http/
- 建议 commit message："feat: add in-memory repo and HTTP handlers for product endpoints"

</details>

<details>
<summary>批次 5 — 后端入口 / wiring / router</summary>

- 相关文件：backend/cmd/marketplace/product-query-svc、internal/adapters/http/router.go、internal/adapters/http/handlers.go
- 建议 commit message："chore: add service main and HTTP wiring (router & handlers)"

</details>

<details>
<summary>批次 6 — 文档与测试</summary>

- 相关文件：readme.md、test/
- 建议 commit message："docs: add README startup steps and basic tests"

</details>

<details>
<summary>批次 7 — 依赖（go.mod / go.sum）</summary>

- 在本地运行：go mod tidy、go test ./...
- 提交：git add go.mod go.sum && git commit -m "chore: update go.mod and go.sum after tidy"

</details>

---

## 常见问题（FAQ，折叠）
<details>
<summary>常见错误与排查</summary>

- "FATAL: database \"app\" does not exist"：确认 Postgres 启动时环境变量 POSTGRES_DB 与服务的 DATABASE_URL 中数据库名一致（示例使用 productdb）；或手动创建数据库。
- Docker 构建报 "go.mod: unknown directive: tool"：确保 Docker 使用的 Go 版本 >= 1.22（Dockerfile 默认 golang:1.22-alpine）。
- Lens 中看不到资源：确认 Lens 使用的 kubeconfig 与 kubectl 当前上下文一致，并且查看正确命名空间（marketplace-dev）。

</details>

---

## 清理建议（可折叠）
<details>
<summary>哪些文件通常应从仓库中移除 / 忽略</summary>

- macOS / IDE 临时文件：.DS_Store、.idea/
- 生成产物（若团队决定不跟踪）：internal/adapters/http/marketplaceapi.gen.go、ports/... 的 generated 文件
- 重复配置文件（例如如果 kind.yaml 与 kind/kind-cluster.yaml 重复，保留一个）

建议操作：创建分支，更新 .gitignore，使用 git rm --cached 删除已跟踪的临时/生成文件，再提交。

</details>

---

## 下一步建议（简短）
- 选择是否保留两个 main（建议重命名 dev main 为 cmd/dev-server），并把该决定写进 README。
- 修复 Tiltfile 报错（如果需要，我可以替你逐项修复）。
- 根据上面的分批提交建议把改动按批次提交并推到远端。

---

最后：如需我直接在仓库里按上述建议执行分批提交或清理（创建分支、修改 .gitignore、提交），告诉我 "请开始执行" 即可。
