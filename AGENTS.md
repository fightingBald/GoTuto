# Repository Guidelines

## 0) GLOBAL GUARANTEES

* **MUST 输出**：`UNIFIED_DIFF` 或完整文件内容；能直接落地。
* **MUST 通过**：`go build ./...` 与 `go test ./...`。
* **MUST NOT**：无关格式化改动、私自新增第三方依赖、修改生成代码,using hand-written SQL strings
* **SHOULD**：保持行为兼容，除非任务明确允许破坏式变更。

---

## 1) CORE PHILOSOPHY & REALITY CHECK

* “Should work” ≠ “does work”。
* 我们不是堆代码，是解决问题。
* 未测试的代码只是猜测。

**30 秒自检（全部回答 YES）**

* 我是否构建/运行了代码？
* 我是否触发了**恰好**被改动的功能路径？
* 我是否亲眼看到期望结果（含 HTTP 状态/响应体）？
* 我是否检查了日志与错误分支？

**禁用措辞**

* “This should work now”“Try it now”“The logic is correct so...”“I’ve fixed it”（二次以后）。

**变更类型最小验收**

* UI 变更：实际点/点/点（如有 GUI）。
* API 变更：发真实 HTTP 请求验证。
* 数据变更：直连数据库验证行数/值。
* 逻辑变更：跑到具体业务场景断言结果。
* 配置变更：重启进程确认加载成功。

---

## 2) REPOSITORY LAYOUT & MODULES

| Path                                                           | Purpose                                                   |
| -------------------------------------------------------------- | --------------------------------------------------------- |
| `backend/cmd/marketplace/product-query-svc`                    | 进程入口与依赖注入（路由、仓库、配置）。                                      |
| `apps/product-query-svc/domain`                                | 领域聚合与不变式（`Product`, `Comment`, `User`）。无 http/sql/env 依赖。 |
| `apps/product-query-svc/application`                           | 用例编排（实现入站端口），只依赖 `ports` 与 `domain`。                      |
| `apps/product-query-svc/adapters`                              | 入站 HTTP handlers；出站持久化实现。禁止写业务规则。                         |
| `apps/product-query-svc/adapters/outbound/postgres/migrations` | SQL 迁移（使用 `migrate` 工具）。                                  |
| `apps/product-query-svc/api/openapi.yaml`                      | OpenAPI 单一事实源。                                            |
| `apps/product-query-svc/api/gen`                               | oapi-codegen 生成物（**禁止手改**）。                               |
| `test`                                                         | 端到端与集成测试（内存/PG 双路径）。                                      |
| `scripts`, `Makefile`                                          | 开发脚本、构建、DB 设置、集成流程。                                       |
| `charts`, `k8s`, `kind`                                        | 部署清单，配置变化时同步。                                             |

**分层约定**

* **Handler**：HTTP/JSON 校验与转换，调用 `application`，装配响应。**不得**嵌业务规则。
* **Application**：编排用例，调用仓库/缓存/消息等端口。业务规则落在 `domain`。
* **Domain**：纯对象与不变量，不依赖框架与存储。
* **Adapters**：实现出站端口；不跨适配器互相引用。

---

## 3) BUILD & DEV COMMANDS

| Command                           | Description                                    |
| --------------------------------- | ---------------------------------------------- |
| `make build`                      | 编译为 `bin/product-query-svc`。                   |
| `make run`                        | 本地启动服务做冒烟验证。                                   |
| `go test ./...`                   | 运行单测；沙箱可用 `GOCACHE=$(pwd)/.gocache`。           |
| `make test-integration-docker`    | 在 `test/http_pg` 下跑 Docker 集成套件。               |
| `make migrate-up MIGRATE_URL=...` | 执行迁移；回滚用 `make migrate-down`。                  |
| `make gen`                        | 在 `apps/product-query-svc/api` 内重生 OpenAPI 绑定。 |

---

## 4) DEPENDENCIES & VERSIONS

| Policy    | Details                                                |
| --------- | ------------------------------------------------------ |
| Go        | 目标 Go `1.24.x`，与 `go.mod/toolchain` 对齐。                |
| Pinning   | 锁到已发布 tag；用 `go get module@version` 升级并 `go mod tidy`。 |
| Review    | PR 审查传递依赖的差异；谨慎大版本升级。                                  |
| Security  | 优先打 Postgres/HTTP/OpenAPI 工具链相关 CVE。                   |
| Artifacts | `go.mod` 与 `go.sum` 一并提交；默认不 vendor。                   |

**Core Libraries**

* `github.com/go-chi/chi/v5`（路由）
* `github.com/jackc/pgx/v5`（Postgres 驱动/连接池）
* `github.com/Masterminds/squirrel`（SQL builder）
* `github.com/getkin/kin-openapi`（OpenAPI 校验）
* `github.com/oapi-codegen/runtime`（生成代码运行时）
* `github.com/testcontainers/testcontainers-go`（Docker 集成测试）

---

## 5) OPENAPI & CODEGEN

* **单一事实源**：`apps/product-query-svc/api/openapi.yaml`。
* **生成位置**：`apps/product-query-svc/api/gen`，**禁止手改**。
* **严格路由**：oapi-codegen 开启 `--strict-server`，不允许野路由。
`

---

## 6) CODING & NAMING

| Rule       | Details                                                |
| ---------- | ------------------------------------------------------ |
| Formatting | `make fmt`（`go fmt ./...`）提交前必跑。                       |
| Naming     | Go mixedCaps；只导出必要标识符；接口按能力命名（无 `I` 前缀）。               |
| JSON       | `snake_case` 与 OpenAPI 完全一致。                           |
| Context    | 所有外部 I/O 接收并透传 `context.Context`，设置合理超时。               |
| Errors     | 包级 sentinel，`fmt.Errorf("op: %w", err)` 包装；HTTP 层统一映射。 |
| Receivers  | 单字母稳定命名，如 `func (s *Server)`。                          |
| Forbidden  | `domain` 依赖 http/sql/env；`handler` 写业务；跨层循环依赖。         |

---

## 7) TESTING (MANDATORY STYLE)

* **表驱动**：`cases := []struct{ name string; in...; exp... }{...}`
* **子测试**：`t.Run(c.name, func(t *testing.T) { c := c; ... })`
* **断言**：`github.com/stretchr/testify/require`
* **覆盖面**：对外可观察行为为主；常见路径 + 边界 + 错误
* **集成**：优先 Testcontainers，本地与 CI 一致


---

## 8) CONCURRENCY, RESOURCES, SECURITY, OBSERVABILITY

* 所有外部调用**必须**使用 `context` 并设置超时。
* 禁止在 handler 内启动长生命周期 goroutine。
* 文件/连接/游标 `defer Close()`，避免泄漏。
* 输入校验在 handler 边界；不要信任客户端。
* 日志用 `log/slog` 键值对；禁止 `fmt.Println`。
* 传播 W3C `traceparent`（若存在）；曝光 P95/P99、error_rate、QPS。

---

## 9) ROLE BOUNDARIES（Agent 硬约束）

| 事项    | 必须                  | 禁止                  |
| ----- | ------------------- | ------------------- |
| 语言/版本 | Go 1.24 与标准库能力      | 过时 API、私有 fork      |
| 依赖    | 默认不新增依赖             | 未授权新增第三方库           |
| 输出    | 统一 diff 或完整文件，能编译测试 | 片段化、不可编译拼贴          |
| 变更范围  | 最小必要改动，保持 ABI/行为    | 大范围样式化改动            |
| 错误处理  | 有语义的 `error`，分级日志   | `panic`（测除外）、吞错、裸打印 |
| 并发    | 遵循取消/超时、无共享可变状态     | 忽视 `ctx`、数据竞争       |
| 安全    | 严格校验输入              | 信任外部输入              |
| 文档    | 关键导出符号有注释，公共约定更到本文  | 把约定只写在 PR 里         |

---

## 10) REFACTOR POLICY

**Allowed（必要时用于“高内聚低耦合”）**

* `EXTRACT_FUNC`：
* `MOVE_FILE`
* `RENAME_SYMBOL`：
* `SPLIT_FILE`：大文件按层或关注点拆分

**Prohibited**

* 修改 `/api/gen` 生成物
* 跨层互相引用导致循环
* 在 handler/adapter 内塞业务规则

**触发条件**

* handler 含业务决策或多步编排
* domain 与 infra 互相引用
* 单文件 > ~500 行且混杂多层


---

## 12) COMMIT, PR & CI GATES

| 检查项  | 要求                                              |
| ---- | ----------------------------------------------- |
| 构建   | `make build` 通过                                 |
| 格式   | `go fmt ./...`、`go vet ./...`                   |
| Lint | `staticcheck ./...`                             |
| 安全   | `govulncheck ./...`                             |
| 测试   | `make test` 与 `make test-integration-docker` 通过 |
| 生成   | 改了 `openapi.yaml` 则重生并提交 `api/gen`              |
| 文档   | 对外行为变化需更新本文件或 README                            |

**PR 提交流程（Checklist）**

* [ ] 本地 `make build && make test` 通过
* [ ] 如改 API：更新 `openapi.yaml` 并重生 `api/gen`
* [ ] 单测覆盖核心路径与边界用例
* [ ] 无无关格式化改动
* [ ] PR 写明：问题、方案、权衡、回滚策略、影响面
* [ ] 性能/缓存改动附基线与指标


