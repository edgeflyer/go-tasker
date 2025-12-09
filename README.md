# Go Tasker

一个用 Go + Gin 构建的学习型任务管理服务，当前可运行基础健康检查，后续将按专业工程流程演进为“模块化单体”，并保留微服务友好的边界。

## 快速开始

- 前置条件：Go 1.21+（本地无需额外依赖）。
- 启动服务：`go run main.go`，默认监听 `:8080`。
- 健康检查：`GET /ping` 返回 `{"message":"pong"}`；`GET /` 返回欢迎文案。

## 技术与架构原则

- **架构策略：模块化单体** —— 先在单进程内划分清晰的领域边界，后续可平滑拆为独立服务。
- **分层职责**：`api`（HTTP 适配）/ `core`（领域）/ `infra`（数据与基础设施）/ `pkg`（通用工具）。
- **接口隔离**：业务以 interface 暴露能力，存储/传输实现可替换（内存 → DB → 远程服务）。

## 目录规划（渐进填充）

```text
go-tasker/
├── cmd/server/           # 入口（预留）
├── api/                  # HTTP handler, middleware
├── core/                 # 领域层（task, user 等）
├── infra/                # DB、配置、日志等实现
├── pkg/                  # 通用工具、响应与错误封装
└── main.go               # 当前入口（简版）
```

## 里程碑路线图

1. **Phase 0：骨架与规范**
   - 建立目录、基础路由。
   - 统一响应/错误结构，约定 HTTP 状态码语义。
2. **Phase 1：Task 模块（内存版）**
   - Task 领域模型与 Service interface。
   - Gin handler 实现 CRUD，业务与 HTTP 解耦。
3. **Phase 2：持久化与 Repository**
   - 引入 GORM/SQL，`core` 通过 Repository 接口访问数据。
   - `infra/db` 提供连接管理与具体实现。
4. **Phase 3：User/Auth**
   - 注册/登录、JWT 中间件，Task 绑定用户。
5. **Phase 4：工程化提升**
   - 配置管理（viper）、结构化日志（zap/logrus）、Swagger、Dockerfile/docker-compose。
6. **Phase 5：可选微服务演进**
   - 按 `core/task`、`core/user` 边界拆分，现有接口与基础设施可复用。

## 开发协作规范

- 分支示例：`feat/task-domain`、`feat/auth-api`、`chore/config`。
- 提交信息：`feat: add task service interface`，`fix: handle missing task error`。
- 质量保障：本地执行 `go fmt ./...`、`go vet ./...`，可按阶段补充单测/集成测试。
