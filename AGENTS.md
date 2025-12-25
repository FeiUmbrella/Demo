## 通用
### Feishu 文档输入处理
- 每当用户提供域名包含 `feishu.cn` 的链接，使用 `mcp__feishu2md__download` 工具加上路径参数`docs/prd`，资源下载到 `docs/prd/` 目录下，并在回答里说明该路径；
- 如果下载内容是需求文档，应先向用户确认是否要按“需求分析 → 方案设计（含接口/依赖/协议核查） → 代码实现（含测试） → Code Review → 分支提交 & MR”这套流程开展工作，再继续后续步骤。

### 交流语言
- 永远使用简体中文进行思考和对话。

### 构建与缓存规范
- 本地开发与测试阶段禁止覆盖 `GOMODCACHE`、`GOCACHE`，也不得创建 `.gomodcache`、`.gocache` 等临时缓存目录；一律使用 Go 默认缓存路径，避免额外清理及权限问题。

## 开发规范
整体流程：**需求分析 → 方案设计（含接口/依赖/协议核查） → 代码实现（含测试） → 集成测试 → Code Review**，默认不启动该流程，必须遵循“先确认、后进入下一阶段”的原则，默认不自动开启流程，需用户明确指示才进入下一阶段。仅当用户在对话中明确输入“开启流程”四个字时，方可视为授权进入该流程；未收到该指令前，禁止主动询问或进入任何阶段。**执行任一命令或行动前，必须先确认当前分支/目录正处于哪一个阶段（需求分析/方案设计/代码实现/集成测试/Code Review），若已有阶段进度需延续既有结论，而非重复从需求分析开始。**每个阶段需参照对应提示词：
- 需求分析：`.agents/01_requirement-analysis.md`
- 方案设计：`.agents/02_solution-design.md`
- 代码实现：`.agents/03_implementation.md`
- 集成测试：`.agents/04_integration_test.md`
- Code Review：`.agents/05_code-review.md`

协议/依赖更新需在方案设计阶段确认并纳入计划。每个需求需在 `docs/feature/<feature_name>/` 下维护四份文档（含版本号与变更记录）：
- `需求拆解.md`：记录输入需求、范围、依赖、确认结论；
- `技术方案.md`：说明实现思路、数据结构、测试策略/DDL；
- `接口文档.md`：基于需求拆解、技术方案生成，完成后生成飞书链接放入对应项目状态文档，格式如下【接口文档】：接口飞书文档链接；
- `项目状态.md`：此文档近记录需求文档链接、接口文档链接、taskId，不要记录其他额外信息。

### 分支与文档目录规范
- 每次处理新需求时，必须先基于 `master` 创建一个新的工作分支（建议遵循 `feature/<topic>`、`fix/<topic>` 等语义化命名），并在该分支上完成全部开发，禁止直接在 `master` 上提交改动。
- `docs/feature/<feature_name>/` 的 `<feature_name>` 必须与当前分支名一致：将分支名中的 `/`、`-` 等特殊符号转换为 `_` 后作为目录名。仅当当前分支为 `master` 时，才能按需求文档的语义命名目录。
- 如需在不同分支复用同一目录，应先确认目录归属并视情况重命名，避免多人并行造成覆盖。
- 若当前分支本身不是 `master`（即已经是一个工作分支），处理需求前需检查 `docs/feature/<feature_name>/` 是否存在：若不存在，应基于当前分支名创建对应目录及三份文档骨架，而不允许再从 `master` 重新拉取新的分支覆盖现有工作。

文档开头使用单行三列表格标注版本/更新时间/备注，并维护变更记录。
- `需求拆解.md`：拆解到页面/字段/权限级别，明确多语言要求与外部依赖范围。
- `技术方案.md`：包含接口请求/响应字段表（结合 `proto-go`、HTTP header）与示例、关键页面 ASCII 布局（标注依赖 `service/*`、`domain/*`）、关键流程/调用链时序图、多语言与权限方案、TODO/风险列表。
- `接口文档.md`：如果指定模版按模版生成，否则按一般接口文档规范生成。
- 接口说明需注明站点信息传参方式（路由 `{local}`、Header `AppLocal`、gRPC Metadata 等），并说明后续调用如何使用。
- `接口文档.md`：采用标准格式列出需求涉及的全部 Host（按环境拆分）、Method/Path、中间件/权限要求、Header/Query/Body 字段、完整请求示例与成功/失败响应示例、错误码说明与日志追踪字段，便于联调与测试。

### 阶段输出规范
- 每完成一个阶段（需求分析 / 方案设计 / 代码实现 / 集成测试 / Code Review），需向相关方同步：
    - `【阶段完成】<阶段名称>`
    - `- 核心结论：<关键成果>`
    - `- 待确认事项：<仍需确认或处理的要点，如无填“无”>`
    - `- 下一步建议：<建议的下一动作>`
- 阶段未完成时，仅可输出 `【阶段进度】<阶段名称>` 与阻塞信息，禁止提前进入下一阶段。

### 各阶段要点
1. **需求分析**（详见 `.agents/01_requirement-analysis.md`）
    - **实现目标**：明确需求范围、输入需求拆解（限定于本项目）、上下游依赖、系统交互流程及边界、数据结构变化与兼容方案；
    - **约束条件**：
        - 如涉及多系统交互，必须先确认完整交互流程和边界，再划定本项目责任范围；
        - 未确认影响范围前，不得修改协议或实现；
    - **验证方式**：将分析结论记录在任务或 `docs/` 的拆解文档中，确保可追踪并经相关方确认。

2. **方案设计**（详见 `.agents/02_solution-design.md`）
    - **实现目标**：形成清晰的实现方案（仅限定本项目改动）、数据结构设计与测试策略/DLL，并体现上下游交互流程；
    - **约束条件**：
        - 必须在 `docs/` 编写或更新设计文档，并记录核心用例；如涉及 DLL，需附详细 DLL 说明；
        - 方案设计阶段必须同步产出标准 `接口文档.md`，列明本需求涉及的所有 Host、Method/Path、请求/响应示例（成功+失败）及权限/中间件要求；
        - 方案需明确是否涉及协议/依赖调整（含 `@mi-grpc-proto` MR、`proto-go` 版本、Nacos/Redis/MQ 变更），在本阶段确认并纳入计划后方可进入实现；
    - **验证方式**：方案文档得到评审或相关人员确认，可作为后续实现和验收依据。

3. **代码实现（TDD）**（详见 `.agents/03_implementation.md`）
    - **实现目标**：按照方案逐步落地功能点；每个迭代需先写测试用例，再实现代码；
        - **约束条件**：
            - 只能在既有模块扩展逻辑，保持兼容性；
            - 禁止修改生成代码或无关模块；新增工具/结构必须符合项目命名与目录规范；
            - 代码风格遵循本仓规范（卫语句、错误包装、无冗余注释等）。
            - 若实现变更导致与已评审的 `技术方案.md` 或 `接口文档.md` 描述不一致，必须先回写并更新对应文档（含变更记录）后再提交代码。
    - **验证方式**：
        - 单功能开发流程：编写失败的测试 → 补充实现 → 确认测试通过 → 编写边界测试；
        - 执行 `go test` 覆盖改动涉及的包，确保新增/既有用例全部通过。

4. **集成测试**（详见 `.agents/04_integration_test.md`）
    - **实现目标**：通过测试平台（`http://oversea-ai.test.mioffice.cn`）在云端执行端到端用例，确保方案落地无回归。
    - **约束条件**：
        - 部署到 test 环境后，必须将当前 PRD/需求拆解与 `接口文档.md` 转换为飞书云文档，分别作为 `prdUrl`、`apiInfoUrl`，并把对应链接写入 `docs/feature/<feature_name>/测试用例.md` 以便后续复用；若仅有本地 Markdown，需使用飞书 MCP 工具上传并获取可分享链接。
        - 必须调用 `/ai/task/create` API 创建测试任务并记录 `taskId`（同样持久化在 `docs/feature/<feature_name>/测试用例.md`），随后轮询 `/ai/task/fetch`，直至状态为 `2`（成功）或 `3`（失败）。
        - 状态 `0/1` 需等待后重试；状态 `3` 必须根据返回的失败用例列表逐条分析，并结合 `logid` 调用 Hera MCP 拉取日志。
        - 平台输出、轮询日志、Hera 结果、taskId 均需沉淀在 `docs/feature/<feature_name>/测试用例.md`，若发现问题需回流“代码实现”阶段修复。
    - **验证方式**：以测试平台成功状态为准；若平台未覆盖全部用例，需在测试文档中补充人工验证记录。

5. **Code Review**（详见 `.agents/05_code-review.md`）
    - **实现目标**：通过结构化评审确认接口/依赖/测试/日志齐备，避免上线补丁。
    - **约束条件**：
        - Reviewer 核对接口文档、proto 变更、Nacos/Redis/RocketMQ 配置、依赖版本（`go.mod`）、测试覆盖、日志/监控，确认与 `.rules`、`docs/memory/architecture.md` 描述一致，并记录审查意见。
    - **验证方式**：审查清单覆盖 proto、配置、依赖文档，问题已关闭或备注；自检通过、关联需求/方案链接、测试命令及结果可追溯。

## 语言规范

### 代码风格（Go）
- 优先级：优先考虑代码的正确性和清晰性；速度和效率是次要优先级，除非另有说明或存在明确性能指标。
- 格式化与 Lint：`go fmt ./...`、`golangci-lint run`（由 pre-commit 提前拦截）。
- 错误处理：优先早返回，避免深嵌套；错误包装携带上下文（必要时使用 `%w`）。
- 控制流：使用卫语句；少用全局变量；避免隐式依赖。
- 注释：不要编写用于概括或总结代码的组织性注释或常规注释；仅当代码的编写方式存在棘手或不易理解的原因时，才应编写注释来解释“为什么”这样写；复杂逻辑可写简洁函数注释，避免赘述“做什么”。
- 泛型：仅在明确提升可读性与复用度时使用（例如通用分页响应）。

### 命名与文件组织
- 文件命名（snake_case）：使用 业务域_动作，如 `controller/csurvey/question_create.go`、`service/spageconfig/save_detail.go`、`dao/dpageconfig/page_config.go`。
- 包命名：与目录一致，语义清晰、全小写；避免与标准库冲突。
- 类型与函数：名词命名变量，动词/动宾短语命名函数，表达语义而非实现细节。
- 结构体 Tag：统一使用 `json:"field" form:"field"`，必要时加 `binding:"required"`。
- 文件组织：优先在现有文件中实现功能，除非它是一个新的逻辑组件；避免创建过多的小文件。

### 尺寸与复杂度控制
- 单函数行数：建议 ≤ 80 行（不含空行与注释），目标 50–80 行；硬上限 150 行，超限需拆分为私有 helper 或抽取到合适子模块。
- 圈复杂度：建议 ≤ 10，硬上限 15（参考 `gocyclo`/`cyclop`）。出现深层条件/分支优先用早返回与小函数重构。
- 单文件代码行数：建议 ≤ 400 行（不含 import/注释/空行）；硬上限 800 行，超限按领域或职责拆分。
    - 特例：路由聚合文件可适度更大（尽量 ≤ 600 行）；DAO 遵循“一表一文件”可适度超出建议值；生成代码/Swagger/Mock/Proto/GRPC 生成物不受此限制。

### Panic 与错误栈规范
- 禁止将 `panic` 用于业务流程控制；仅不可恢复的编程错误可 `panic`（需尽快修复）。
- 全局必须启用 `middleware.Recovery` 统一捕获 `panic`，返回 `500` 并记录 `runtime/debug.Stack`。
- 发生 `panic` 时日志字段至少包含：`trace_id`、`user`、`appLocal`、`method`、`path`、`err`（panic 值）、`stack`；避免记录敏感信息。
- 业务错误日志只在“决策边界”处记录一次（Controller/Service 顶层）；下层通过 `fmt.Errorf("...: %w", err)` 保留错误链路。
- goroutine 内需本地 `recover` 并记录错误栈，避免影响进程存活：
  ```go
  go func() {
      defer func() {
          if r := recover(); r != nil {
              stack := debug.Stack()
              logz.Error("goroutine panic", logz.F("err", r), logz.F("stack", string(stack)))
          }
      }()
      // ... do work
  }()
  ```

### 性能与并发（语言层建议）
- 避免在热路径频繁分配（可用缓冲/对象池，评估收益与复杂度）。
- 外部调用务必设置超时与重试；协程中传递 `context.Context` 并在退出时取消。

---

## 项目规范

### 依赖说明
- 项目依赖（数据库、缓存、RPC/HTTP、消息队列）和常用客户端请以 `docs/memory/architecture.md` 为准，负责配置对应的访问点与角色。
- 所有 HTTP/路由/Controller 责任范围，可同步参考 `docs/memory/dependency_api_full.md` 的接口列表以确认入口与调用关系。

### 适用范围
- 适用于主服务（`main.go`）、HTTP 接口（`controller/`、`router/`、`middleware/`）、业务逻辑（`service/`）、数据访问（`dao/`）、模型（`model/`）、作业任务（`job/`）、通用库（`lib/`、`utils/`）、配置（`conf/`）。

### 目录与分层
- `controller/`（前缀 c）：HTTP Handler，仅做参数解析与校验、调用 service、统一响应。
- `service/`（前缀 s）：业务编排与领域逻辑，聚合 DAO 与外部 SDK，控制事务边界与幂等。
- `model/`（前缀 m）：请求/响应/领域结构体，带 `json`/`form`/`binding`/Swagger 标签。
- `dao/`（前缀 d）：表到结构的 ORM 映射与数据访问（XORM，`store/db`）。
- `router/`：路由注册与分组，遵循域/模块/动作路由风格。
- `middleware/`：链路追踪、日志、鉴权、数据脱敏、CORS、Swagger 等。
- `lib/`：外部服务 SDK 封装、客户端与公共库（如 `lib/swag`）。
- `utils/`：通用工具（响应封装、日志上下文、时间/加解密/Excel/HTTP 等）。
- `job/`：批处理、定时、消费类任务，按业务域分子目录，独立 `main.go`。
- `conf/`：配置模板与静态资源（如 `conf/local.toml.offline`、`translate.json`）。

### 接口与路由规范
- 全局初始化：`router.InitRouter(r)`；公共路由 → 本地化路由 → 需要登录路由。
- 中间件顺序：
    - 全局：`CORS` → `Recovery` → `Logger` → `BaseParams`。
    - 登录：`CasLogin()` → `DataMasking()` → 业务路由。
- 本地化路由：以 `/{local}` 分组，站点集来自 `constant.LocalList()` 与渠道后缀组合。
- 路由风格：`/{domain}/{module}/{action}`，全小写、短横线连接（如 `/survey/survey-list`）。
- Swagger：非线上环境开启 `/swagger/*any`，注释遵循 `swag` 标注（见“文档与 Swagger”）。

### Controller 规范
- 仅负责：参数解析与校验、读取公共参数、调用 service、统一响应。
- 禁止：直连数据库/缓存；实现复杂业务或控制事务。
- 统一响应：
    - 成功：`utils.ResponseSuccess(c, data)`。
    - 失败：`utils.ResponseError(c, code.ErrorXXX, err)` 或 `ResponseErrorData`。
- 分页请求：
  ```go
  type PaginatedRequest struct {
      Page     int32 `form:"page" json:"page"`
      PageSize int32 `form:"page_size" json:"page_size"`
  }
  ```
- 分页响应：
  ```go
  type PageInfo struct { PageSize int32 `json:"page_size"`; PageNum int32 `json:"page_num"`; TotalSize int64 `json:"total_size"` }
  type ListResponse[T any] struct { List []T `json:"list"`; PageInfo PageInfo `json:"page_info"` }
  ```
- 默认值：`Page=1`、`PageSize=10` 在 Controller 兜底。

### Service 规范
- 职责：聚合领域逻辑、校验业务规则、调用 DAO/外部服务、控制事务与幂等。
- 事务：使用 XORM `Session`（`dao.GetSession()`）；`Begin`/`Commit`/`Rollback` 明确，边界在 Service 层。
- 幂等与重试：对外部调用和消息投递设计去重键与重试策略。
- 上下文：透传 `context.Context`，RPC/HTTP 调用添加必要 Header（见 `utils/grpc_header.go`）。

### DAO 规范（XORM）
- 仅做数据访问：表结构映射、查询、增删改；不写业务判断。
- 读写分离：读 `db.R()`/`dao.GetReadCon()`；写 `db.W()`/`dao.GetWriteCon()`；`ShowSQL` 由配置控制（线上默认关闭）。
- 软删除/状态：遵循表定义的 `is_delete`/`status` 等字段，统一常量与过滤条件。
- 更新/删除：检查影响行数，出现“partial”须报错（参考 `dao/dpageconfig/base.go`）。
- 批量操作：优先批量插入/更新，严格校验影响条数；大表查询走索引，必要时用断点续查，避免 N+1。

### Model 规范
- 结构放在 `model/{domain}`；`req.go`/`resp.go` 拆分。
- Tag：`json`、`form` 必备；分页、时间、枚举等字段标注清晰并配 Swagger 标签。
- 复用：分页响应复用 `ListResponse` 与 `PageInfo`，避免重复定义。

### 中间件与跨切关注
- TraceID：`middleware.SetTraceId` + `utils.SetTraceId` 生成与透传。
- 恢复：`middleware.Recovery` 捕获 panic，返回 500，打印 stack。
- 访问日志：`middleware.Logger` 记录参数、耗时、用户、IP，并可写入 DB（线上开启）。
- 登录与权限：`middleware.CasLogin()`、`middleware.RightsAuth`；本地开发可按 `conf/local.toml` 跳过。
- 数据脱敏：`middleware.DataMasking()` 按路由进行脱敏处理。
- CORS：统一通过 `middleware.CORS`。

### 错误处理与多语言
- 错误码：统一使用 `common/code` 枚举（`code.ErrorXXX`）。
- 返回：使用 `utils.Response*`，带上 `TraceID`；不要直接 `c.JSON` 拼结构。
- 多语言：`lib/translate.Translator().GetText` 根据 `appLocal` 派生 `lang`（从 `constant.AppLocalPure` 获取），查不到走 `code.ErrorMsg` 默认文案。

### 时间与时区
- 前端展示时间使用 `config.FormatTime(appLocal, timestamp)`（参考 `CLAUDE.md`），不要直接 `time.Unix().Format()`；保证各站点时区一致性。

### 外部服务与上下文
- 站点上下文：从 `gin.Context` 读取 `constant.AppLocal`/`AppLocalPure`；必要时用 `constant.GetAreaIdByAppLocal(appLocal)` 获取区域。
- 依赖调用：统一使用 `lib/` 客户端封装，严禁在业务代码中硬编码地址或协议细节。

#### Context 透传规范
- Controller → Service：所有对外调用（DAO、RPC/HTTP、MQ）均透传 `context.Context`，禁止在下层新建 `context.Background()`。
- Service → DAO：DAO 方法签名包含 `ctx context.Context`（历史代码逐步演进），用于超时、链路追踪与日志字段透传。
- 超时：外部 IO 调用使用 `context.WithTimeout`；默认 2–5s（按依赖 SLO），MQ/批量可更长；避免层层叠加过短超时。
- Trace 与用户：从 `gin.Context` 获取 TraceID/用户/站点信息写入日志；gRPC 使用 `utils/grpc_header.go` 的 `GetV1Header/WithLang`；HTTP 通过 Header 传递 `Trace-Id`、`User-Name`、`App-Local`。
- goroutine：使用上层 `ctx`（`c.Request.Context()` 或传入的 `ctx`），支持 `ctx.Done()` 取消，禁止 `context.TODO()`。

### 消息与任务（Job）
- MQ 生产者：通过 `queue/producer` 初始化与关闭，生命周期与 HTTP 服务一致（见 `main.go`）。
- 作业结构：`job/{domain}/.../main.go`；尽量沉淀公共逻辑到 `service/`、`lib/`。
- 配置与日志：沿用 `viper`；日志使用 `logz`/`xlog`。
- 退出：监听信号，优雅停机，释放资源（MQ、DB、HTTP 客户端等）。

### 配置与环境
- 配置入口：`-c conf/local.toml`，由 `bootstrap.Init` 加载并 watch；热更新通过 `viper.WatchConfig()`。
- 环境：`app.env ∈ {dev,test,online}`；非线上暴露 Swagger；线上 `gin.ReleaseMode` 并注册到服务发现（etcd）。
- 端口：默认 `:8080`，可通过配置覆盖。

### 文档与 Swagger
- 注释示例：
  ```go
  // @Summary 获取模板列表
  // @Tags template
  // @Accept json
  // @Produce json
  // @Param  page      query int false "页码"
  // @Param  page_size query int false "每页条数"
  // @Success 200 {object} common.Response{data=ListResponse[Template]}
  // @Router /template/list [get]
  ```
- 生成：本地/测试构建 `cd lib/swag && go install && cd ../.. && swag init` 或执行 `build.sh` 的 `swag init`；仅非线上开放 `/swagger/*any`。

### 日志规范（项目层）
- 统一使用 `xlog`/`logz` 结构化日志，带上 TraceID、user、appLocal；避免敏感数据输出；访问日志与操作日志分离（可通过 `IgnoreDBLog` 关闭 DB记录）。

### 安全规范
- 鉴权：需要登录的接口必须加 `CasLogin()` 与 `RightsAuth`（按功能设置权限标识）。
- 参数校验：对上传、导入等高风险入口做白名单/大小/类型校验。
- 输出安全：配合数据脱敏中间件与 `utils.Response*` 控制返回；不向客户端暴露错误栈。
- CORS：统一在中间件配置，禁止在 Handler 内部覆盖。

### 提交与分支（建议）
- 提交信息：简洁说明动机与影响（Add/Update/Fix/Refactor/Test/Docs 前缀）。
- 避免提交敏感信息与环境配置；配置信息通过 `conf/*.toml` 模板与环境变量管理。

### 常见示例
- 事务骨架：
  ```go
  s, err := dao.GetSession()
  if err != nil { return err }
  defer s.Close()
  if err := s.Begin(); err != nil { return err }
  // ... DAO 调用
  if err != nil { _ = s.Rollback(); return err }
  return s.Commit()
  ```
- Controller 骨架：
  ```go
  func List(c *gin.Context) {
      var req model.PaginatedRequest
      if err := c.ShouldBindQuery(&req); err != nil { utils.ResponseParamsError(c, err); return }
      if req.Page == 0 { req.Page = 1 }
      if req.PageSize == 0 { req.PageSize = 10 }
      resp, err := service.SomethingList(c, req)
      if err != nil { utils.ResponseError(c, code.ErrorDB, err); return }
      utils.ResponseSuccess(c, resp)
  }
  ```

### 落地检查清单（PR 自检）
- [ ] Controller 无业务逻辑/无直连 DB/统一响应/分页规范
- [ ] Service 控制事务边界/考虑幂等/透传 Context
- [ ] DAO 影响行数校验/无业务判断/读写分离
- [ ] Model 标签与 Swagger 完整/复用分页结构
- [ ] 中间件顺序正确/权限与脱敏覆盖
- [ ] 错误码与多语言文案齐全
- [ ] 日志包含 TraceID 与关键字段/无敏感信息
- [ ] 单测覆盖关键分支/Mock 合理
- [ ] Swagger 注释完整且可生成
- [ ] 代码通过 fmt 与 lint

---

## 单测规范

### 常用命令
- 运行全部：`go test ./...`
- 覆盖率：`go test -v -cover ./...`
- 运行指定：`go test -v -run TestFunctionName ./path/to/package`
- 执行 `go test` 时如需写缓存，统一使用 `GOCACHE=$(pwd)/.gocache <命令>`，命令结束后必须 `rm -rf .gocache`，避免缓存目录留在仓库。

### 工具与框架
- Mock：
    - 接口依赖：`gomock` + `mockgen` 生成接口的 Mock 实现。
    - 全局函数/第三方：`gomonkey` 打桩（慎用，只在无法接口化时使用）。
- 断言：`testify/assert` 或 `require`。

### 编写建议
- Controller：使用 `httptest` 构造请求，断言状态码与响应体；对分页与参数错误分支做用例。
- Service：以接口抽象 DAO/外部客户端，注入 Mock；覆盖事务成功/回滚、重试与幂等分支。
- DAO：可用事务回滚或测试库隔离；聚焦 SQL 条件、影响行数与索引命中。
- 时间与多语言：对时间/时区与多语言返回做最小必要断言，避免和具体文案强耦合。
- 日志与 Panic：不断言日志内容；对可能 `panic` 的边界用例使用 `require.NotPanics` 或本地 `recover`。

### 覆盖与质量
- 优先覆盖复杂条件分支与关键路径；对外部交互加错误与超时用例。
- 持续关注 `-cover` 指标；新增复杂功能应附带相应单测。
- 依赖管理
    - 默认使用 go.mod 中列出的 git.n.xiaomi.com/miopen/mit/mi-go/proto-go/sales/i18n/shop/... 子模块版本；仅需在 go.mod 调整版本号即可，不直接 go get 全量 proto，也不把远端 proto 落地到 lib。
    - 保留现有的 replace git.n.xiaomi.com/miopen/mit/mi-go/proto-go => ./lib/proto-go（以及必要子模块 replace），但这只对仓内已有的 proto 生成物生效；其它模块继续通过 go.mod 引用远端版本，避免手动复制。
    - 若需调试特定 proto，可在本地临时取消相应 replace 或更新 go.mod 版本号；完成后务必恢复到官方版本，保证团队在同一版本下开发。
- go get 与 go mod tidy
    - 禁止随意执行 go get git.n.xiaomi.com/miopen/mit/mi-go/proto-go 或 go mod repair 等命令，以免引入未知版本。
    - 运行 go mod tidy 时需确保 go.sum 中只包含本仓允许的依赖（即 go.mod 已列出的子模块）；若出现 “ambiguous import” 提示，优先检查 replace 与版本配置，而不是复制文件。
