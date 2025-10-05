# 项目搭建与开发指南

[![Go CI](https://github.com/fightingBald/GoTuto/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/fightingBald/GoTuto/actions/workflows/go.yml)

简短说明：本仓库包含一个示例后端服务 product-query-svc（支持 in-memory 与 Postgres），数据库迁移需通过 golang-migrate 执行（不再使用嵌入式迁移），以及用于本地开发的 Tilt + kind 配置与最小 Helm chart（已补全）。
服务目前提供商品 CRUD 及评论功能，评论支持多用户查看、作者更新/删除，并通过 OpenAPI 严格校验暴露接口。

---

## 目录结构

项目采用按“应用 + 适配器”的分层组织，便于替换实现与独立演进。

```
.
├── api/                               # OpenAPI 定义与代码生成配置
│   ├── openapi.yaml                   # 主 OpenAPI 入口
│   ├── oapi-config.yaml               # oapi-codegen 配置
│   ├── generate.go                    # go generate 指令
│   ├── paths/                         # 路径/接口片段
│   ├── schemas/                       # 数据结构（Schema）
│   └── responses/                     # 响应体定义
├── apps/
│   └── product-query-svc/             # 应用层与适配器
│       ├── domain/                    # 领域模型与领域错误
│       ├── ports/                     # 端口（接口），抽象仓储与服务
│       ├── application/               # 应用服务实现（业务编排）
│       │   ├── product/               # 商品相关用例
│       │   └── comment/               # 商品评论用例
│       └── adapters/
│           ├── inbound/http/          # OpenAPI 严格服务 + 路由装配 + 轻量 handler
│           └── outbound/
│               ├── inmem/             # 内存仓储实现（开发/测试）
│               └── postgres/          # Postgres 仓储与迁移文件
├── backend/
│   └── cmd/marketplace/product-query-svc/  # 可执行入口（main.go），装配路由/依赖
├── charts/product-query-svc/          # 最小 Helm Chart（含迁移 Job 与 ConfigMap）
├── k8s/                               # 直接应用的 Kubernetes 清单（Service/Deployment/Postgres）
├── kind/                              # kind 本地集群配置
├── scripts/                           # 脚本（初始化迁移、kind 集群启动等）
├── test/                              # 单元/集成测试
├── bin/                               # 本地构建产物输出（make build）
├── Dockerfile                         # 服务镜像构建
├── Tiltfile                           # 本地开发编排（构建、端口转发、迁移）
├── Makefile                           # 常用命令封装（构建、代码生成、迁移）
└── .env.dev.example                   # 本地开发环境变量示例
```

说明：
- 目录遵循“端口与适配器”（Hexagonal/Clean Architecture）思路，`ports` 定义接口（inbound/outbound），`app` 为用例实现，`adapters/*` 提供适配器实现；`domain` 保持纯净，可复用。
- 数据库迁移统一放在 `apps/product-query-svc/adapters/outbound/postgres/migrations`，通过 golang-migrate 执行（脚本/Make/Tilt/Helm 均已支持）。
- 生产部署建议使用 Helm Chart；本仓库同时保留了 `k8s/` 便于直接 kubectl 应用与调试。

---

## HTTP 适配器设计（Strict Server）

- **代码生成统一使用 `oapi-codegen strict-server`**：`api/oapi-config.yaml` 只保留严格服务输出，避免手写 handler 接口。每次变更 OpenAPI 需执行 `go generate ./api` 重新生成 `marketplaceapi.gen.go`。
- **请求校验前移到 OpenAPI**：所有参数/请求体验证（`minimum`/`maxLength`/`enum` 等）写在 `api` 目录的 schema/parameter 中，由 `github.com/oapi-codegen/nethttp-middleware` 提供的 `OapiRequestValidator` 中间件统一拦截。
- **Handler 职责“三件套”**（`apps/product-query-svc/adapters/inbound/http/handler_*.go`）：
  1. 从生成的强类型 `RequestObject` 中取出入参（无需重复校验）；
  2. 调用对应的应用服务（`application/*`）；
  3. 利用 `response_helpers.go` 中的 `ok*/xxxError` 辅助函数返回严格的响应类型（仅 2xx/4xx）。
- **跨操作共享错误映射**：`response_helpers.go` 负责把领域错误映射成具体的 OpenAPI 响应类型，并封装标准错误载荷；新增业务错误时只需在此扩展。
- **统一路由出口**：`NewAPIHandler` 会加载内嵌的 Swagger、挂载必需的中间件（含请求校验）并包装 strict server；在 `main.go`、集成测试与 `internal/testutil` 中均通过该函数装配，保持行为一致。

## 验证服务是否可用（Tilt 本地）

- 端口转发就绪
  - 打开 Tilt UI，确认 `product-query-svc` 资源为绿色/Ready，且显示端口转发到 `http://localhost:8080`。

- 健康检查

```sh
curl -i http://localhost:8080/healthz
# 期望: HTTP/1.1 200 OK，响应体: ok
```

<details>
<summary>⚡ 快速 API 测试（可复制/点击运行）</summary>

前提：服务已监听 http://localhost:8080，已安装 curl 与 jq。

1) GET /healthz（健康检查）

```sh
curl -i http://localhost:8080/healthz
```

2) POST /products（创建商品）

```sh
curl -s -X POST http://localhost:8080/products \
  -H 'Content-Type: application/json' \
  -d '{"name":"Sample Plan","price":123.45}' | jq
```

3) PUT /products/{id}（整资源更新，示例使用已存在的 id）

```sh
curl -s -X PUT http://localhost:8080/products/1 \
  -H 'Content-Type: application/json' \
  -d '{"name":"Updated Plan","price":199.99}' | jq
```

4) GET /products/{id}（按 ID 查询，示例使用已种子或上一步创建/更新的 id）

```sh
# 如果使用迁移种子数据（Postgres），通常 1 为 Basic Plan
curl -s http://localhost:8080/products/1 | jq
```

5) GET /products/search（分页搜索；注意 q 至少 3 个字符）

```sh
curl -s 'http://localhost:8080/products/search?q=pro&page=1&pageSize=10' | jq
```

6) DELETE /products/{id}（删除；示例：先创建临时商品再删除）

```sh
ID=$(curl -s -X POST http://localhost:8080/products \
  -H 'Content-Type: application/json' \
  -d '{"name":"Temp Item","price":1.99}' | jq -r '.id'); \
echo "created id=$ID"; \
curl -i -X DELETE http://localhost:8080/products/$ID; \
echo; \
curl -i http://localhost:8080/products/$ID  # 期望 404
```

7) POST /products/{id}/comments（新增评论，需提供已有用户 ID）

```sh
COMMENT_ID=$(curl -s -X POST http://localhost:8080/products/1/comments \
  -H 'Content-Type: application/json' \
  -d '{"userId":1,"content":"Great product!"}' | jq -r '.id'); \
echo "comment id=$COMMENT_ID"
```

8) GET /products/{id}/comments（查看评论列表，默认按创建时间倒序）

```sh
curl -s http://localhost:8080/products/1/comments | jq
```

9) PUT /products/{id}/comments/{commentId}（更新评论内容，`userId` 需放在查询参数且与原作者一致）

```sh
curl -s -X PUT "http://localhost:8080/products/1/comments/${COMMENT_ID}?userId=1" \
  -H 'Content-Type: application/json' \
  -d '{"userId":1,"content":"Updated feedback"}' | jq
```

10) DELETE /products/{id}/comments/{commentId}（删除评论，同样需要 `userId` 查询参数）

```sh
curl -i -X DELETE "http://localhost:8080/products/1/comments/${COMMENT_ID}?userId=1"
```

</details>

<details>
<summary>🧪 使用临时 Docker Postgres 跑集成测试（避免 5432 端口冲突）</summary>

前置：本机已安装 Docker。

一键运行（自动起容器 → 迁移 → 运行带 Postgres 的集成测试 → 清理容器）：

```sh
make test-integration-docker
```

或直接运行脚本，并自定义 go test 目标/参数：

```sh
bash scripts/test-integration-docker.sh ./test -run Postgres
```

脚本要点：
- 使用 `docker run -P` 启动 postgres:16-alpine，随机映射宿主端口，避免与 Tilt 的 5432 冲突。
- 通过 `migrate/migrate` 容器在同一网络命名空间内执行迁移。
- 自动导出 `DATABASE_URL` 为宿主上的随机端口，并运行 go test。
- 需要单独验证仓储层（含评论 CRUD）的 Docker 集成测试时，可运行 `go test -tags docker ./apps/product-query-svc/adapters/outbound/postgres -run TestCommentRepository_WithDocker -count=1`，确保本机 Docker 可用；若暂不具备条件，可设置 `SKIP_DOCKER_TESTS=1` 跳过。

</details>

- 插入演示数据（Postgres）

```sh
# 连接数据库（Tilt 将 Postgres 转发到本机 5432）
psql "postgres://app:app_password@localhost:5432/productdb?sslmode=disable"

# 在 psql 中执行:
insert into products(name, price, tags)
values ('Basic Plan',9900,ARRAY['starter','subscription']),
       ('Pro Plan',19900,ARRAY['professional','subscription']),
       ('Enterprise Plan',49900,ARRAY['enterprise','subscription']);

# 再次验证
\q
curl -s 'http://localhost:8080/products/search?q=pro&page=1&pageSize=10' | jq
curl -s http://localhost:8080/products/1 | jq
```

- 迁移是否成功
  - 在 Tilt UI 查看 `db-migrate` 资源日志，确认 `up` 成功。
  - 或进入 psql 执行 `\dt` 检查是否存在 `products` 表。

- Pod/日志排查

```sh
kubectl -n marketplace-dev get pods
kubectl -n marketplace-dev logs deploy/product-query-svc
kubectl -n marketplace-dev logs statefulset/postgres
```

---

## 前置条件
- Go == 1.24
- Docker
- kubectl、helm、tilt、psql 客户端、IDE（

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

1. 设置环境变量：

```sh
export DATABASE_URL="postgres://app:app_password@localhost:5432/productdb?sslmode=disable"
export HTTP_ADDRESS=":8080"
export LOG_LEVEL=debug
```

1. 运行服务（开发）：

```sh
cd backend/cmd/marketplace/product-query-svc
go run .
```

或构建后运行：

```sh
go build -o bin/product-query-svc ./backend/cmd/marketplace/product-query-svc
./bin/product-query-svc
```

说明：项目包含迁移文件（apps/product-query-svc/adapters/outbound/postgres/migrations），请统一使用 golang-migrate 工具管理数据库 schema。

快捷初始化（迁移包含测试数据）：

```sh
# 使用脚本自动运行迁移（优先本机 migrate CLI，缺失则用 docker 镜像）
DATABASE_URL="postgres://app:app_password@localhost:5432/productdb?sslmode=disable" \
bash scripts/db-init.sh

# 或使用 Makefile 包装目标
make db-init
```

说明：迁移序列包含 `000002_seed_test_data.up.sql`，会插入示例数据（幂等）。回滚同名 `.down.sql` 可清理。

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

## 数据库迁移（Migration）说明与避坑

本项目的迁移文件位于 `apps/product-query-svc/adapters/outbound/postgres/migrations`，支持三种方式执行迁移：

- 手动本地执行（开发态）
  - 命令：`make migrate-up MIGRATE_URL="postgres://app:app_password@localhost:5432/productdb?sslmode=disable"`
  - 依赖：本机安装 `migrate` CLI（可选 Homebrew：`brew install golang-migrate`）。

- 通过 Tilt 自动迁移（当前默认，推荐开发态）
  - Tiltfile 配置了 `local_resource('db-migrate', ...)`，会在 `postgres` 端口转发就绪后，使用 Docker 运行 `migrate/migrate` 容器来执行迁移。
  - 无需在本机安装 `migrate` CLI，但需要本机 Docker 正常运行。
  - 连接串默认使用 `host.docker.internal:5432` 访问 Tilt 的本地端口转发（macOS/Windows Docker Desktop 默认可用）。
  - 如需修改连接串，可在 Tiltfile 顶部调整 `MIGRATE_URL`。

- 在 K8s/Helm 中执行（集群内）
  - 可选：用 Helm hook 或 Job 在集群内运行 `migrate/migrate`，`DATABASE_URL` 使用集群内 Service（例如 `postgres.marketplace-dev.svc.cluster.local`）。需要的话可以补充该 Job。

常见避坑：

- “migrate: command not found”
  - 原因：旧实现会在本机直接执行 `migrate` CLI；如果没安装就会报错。
  - 现状：Tiltfile 已改为通过 `docker run migrate/migrate ...` 执行，无需本机安装 CLI，但需要 Docker。

- Linux 下 `host.docker.internal` 不可用
  - 方案1：在 Docker 命令中添加 `--add-host=host.docker.internal:host-gateway`；或
  - 方案2：使用 `--network host` 并将 `MIGRATE_URL` 的主机名改为 `localhost`（注意该方案在 macOS/Windows 不通用）。

- 5432 端口占用冲突
  - Tilt 会把集群内 Postgres 端口转发到你本机 5432；如果你本地已有 Postgres 占用该端口，则端口转发失败。
  - 解决：停掉本地 Postgres，或修改端口转发/连接串（例如改为 15432，并同步调整 Tiltfile 的 `MIGRATE_URL`）。

- 凭据与安全
  - 目前开发态在 ConfigMap/值文件中包含了示例凭据，便于演示。
  - 生产/共享环境应将 `DATABASE_URL` 放入 Secret，并在迁移/应用中引用 Secret。

- Distroless 镜像与热更新
  - 当前服务镜像为 distroless，不支持容器内直接热重载，Tilt 将走重新构建/重启流程，这是预期行为。

---

## 使用 Helm
仓库包含最小 Helm chart：charts/product-query-svc，可用如下命令部署：

```sh
helm install product-query-svc ./charts/product-query-svc -f ./charts/product-query-svc/values.yaml
```

Helm 迁移 Job：
- Chart 内置了一个 `post-install, post-upgrade` 的迁移 Job（使用 `migrate/migrate` 镜像）。
- 迁移文件位于 Chart 下的 `migrations/` 目录（通过 ConfigMap 挂载到容器 `/migrations`）。
- 默认启用（`values.yaml: migrations.enabled=true`）。如需关闭，设置 `--set migrations.enabled=false`。
- 数据库连接串（DATABASE_URL）
  - 生产默认通过 Secret 注入：`values.database.secret.enabled=true`，并在 Deployment/Job 中 `secretKeyRef` 读取。
  - 如果已有现成 Secret，设置 `values.database.secret.name` 与 `values.database.secret.key` 即可。
  - 若需要 Chart 自动创建 Secret，设置 `values.database.secret.create=true` 并提供 `values.database.url`（作为 stringData）。
  - 开发环境下仍可回退读取 `values.env.DATABASE_URL`（不推荐用于生产）。

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

- 生成后的 `adapters/inbound/http/marketplaceapi.gen.go` **禁止手动修改**；需要调整校验或字段时改 OpenAPI 资源并重新生成。
- HTTP handler 只能依赖生成的 `StrictServerInterface`，其实现位于 `handler_*.go`，必须配合 `response_helpers.go` 和 `request_mappers.go` 使用。
- `NewAPIHandler` 会自动加载最新的 Swagger 并注册 `OapiRequestValidator` 中间件，生产/测试入口都应通过该函数获取路由。

- 建议：将生成步骤写入 Makefile 或 CI，团队协同时要约定是否把生成产物纳入版本控制（两种策略均可）。

</details>

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

 - 相关文件：apps/product-query-svc/adapters/outbound/postgres/（migrations、product_repository.go）
- 建议 commit message："db: add Postgres migrations and adapters (migrations embedded via //go:embed)"

</details>

<details>
<summary>批次 3 — 核心应用层（domain、ports、app）</summary>

- 相关文件：apps/product-query-svc/domain/ apps/product-query-svc/ports/ apps/product-query-svc/application/
  - 建议 commit message："app: add domain models, service implementation and ports for product-query-svc"

</details>

<details>
<summary>批次 4 — 适配器：in-memory repo 与 HTTP handlers</summary>

- 相关文件：apps/product-query-svc/adapters/outbound/inmem/ apps/product-query-svc/adapters/inbound/http/
- 建议 commit message："feat: add in-memory repo and HTTP handlers for product endpoints"

</details>

<details>
<summary>批次 5 — 后端入口 / wiring / router</summary>

 - 相关文件：backend/cmd/marketplace/product-query-svc、apps/product-query-svc/adapters/inbound/http/
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
- Docker 构建报 "go.mod: unknown directive: tool"：请使用与 go.mod 中 toolchain 对齐的 Go 版本（本项目使用 1.24）。
- Lens 中看不到资源：确认 Lens 使用的 kubeconfig 与 kubectl 当前上下文一致，并且查看正确命名空间（marketplace-dev）。

</details>
