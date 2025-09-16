# 项目搭建与开发指南

简短说明：本仓库包含一个示例后端服务 product-query-svc（支持 in-memory 与 Postgres），数据库迁移需通过 golang-migrate 执行（不再使用嵌入式迁移），以及用于本地开发的 Tilt + kind 配置与最小 Helm chart（已补全）。

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

说明：项目包含迁移文件（apps/product-query-svc/adapters/postgres/migrations），请统一使用 golang-migrate 工具管理数据库 schema。

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

本项目的迁移文件位于 `apps/product-query-svc/adapters/postgres/migrations`，支持三种方式执行迁移：

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

 - 相关文件：apps/product-query-svc/adapters/postgres/（migrations、migrations_embedded.go、product_repository.go）
- 建议 commit message："db: add Postgres migrations and adapters (migrations embedded via //go:embed)"

</details>

<details>
<summary>批次 3 — 核心应用层（domain、ports、service）</summary>

- 相关文件：apps/product-query-svc/domain/ apps/product-query-svc/ports/ apps/product-query-svc/adapters/service/
  - 建议 commit message："app: add domain models, service implementation and ports for product-query-svc"

</details>

<details>
<summary>批次 4 — 适配器：in-memory repo 与 HTTP handlers</summary>

- 相关文件：apps/product-query-svc/adapters/inmem/ apps/product-query-svc/adapters/http/
- 建议 commit message："feat: add in-memory repo and HTTP handlers for product endpoints"

</details>

<details>
<summary>批次 5 — 后端入口 / wiring / router</summary>

 - 相关文件：backend/cmd/marketplace/product-query-svc、apps/product-query-svc/adapters/http/
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
