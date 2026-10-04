# v0.0 开发文档：程序基础

- 文档版本：0.0
- 状态：开发设计，尚未实现或完成验收
- 更新日期：2026-10-04
- 设计依据：[概要设计](overview-design.md)、[详细设计](detailed-design.md)
- 版本范围依据：[产品版本路线图](roadmap.md)

## 1. 版本目标

通过 Docker Compose 启动一个运行 `query-api` 的容器，通过配置文件控制运行参数，能够响应 HTTP 请求，并具备明确的启动、日志和退出行为。

本版建立后续版本共同使用的工程与运行基础。后续成交查询、账户分析和管理接口在同一架构中增加路由和业务模块，沿用配置机制、响应格式、错误处理和程序生命周期。

完成后的直接使用方式是：准备配置文件、执行 `docker compose up --build -d`、调用健康和版本接口、查看容器日志、停止容器。启动失败必须能够定位原因；正常停止必须能够释放监听端口。

代码结构遵循现有详细设计：HTTP 入口位于 `gateway/query-api`，配置文件位于 `config/`；后续业务逻辑归属相应服务，协议实现归属 `protocols/implementations/`。这一边界在首版设计中落实。

## 2. 功能范围

### 2.1 本版实现

| 功能 | 交付内容 |
| --- | --- |
| 容器启动 | Dockerfile、Docker Compose 配置、只读配置挂载及端口映射 |
| 工程基础 | Rust workspace、固定工具链、统一依赖、锁文件和 query-api package |
| 配置 | TOML 文件加载、指定文件路径、环境变量覆盖及启动校验 |
| HTTP 服务 | 按配置绑定地址与端口，提供健康和版本接口 |
| 响应规范 | JSON 成功信封、错误格式、schema 版本及请求 trace ID |
| 基础日志 | 启动、请求完成、错误和退出日志，支持级别与格式配置 |
| 生命周期 | 启动失败明确退出，SIGINT/SIGTERM 触发有界优雅退出 |
| 验证 | 可重复的构建检查、接口检查、配置与生命周期验收 |

### 2.2 实现边界

本版只启动 `query-api` 一个程序，不接入 Hyperliquid，不实现成交、账户或跟单业务。数据库、消息系统、缓存、对象存储及观测平台在实际需要的版本接入。

健康和版本接口是无账户数据的运行接口，本版允许匿名读取。不开放业务及管理接口；后续引入业务接口时按详细设计接入认证和授权。

配置在启动时加载一次，修改后重启生效。本版不提供热更新、配置写入 API 或配置内容查询接口。

## 3. 使用方式

### 3.1 Docker 启动入口

从仓库根目录运行：

```sh
docker compose up --build -d
docker compose ps
docker compose logs -f query-api
```

本版 Compose 只包含 `query-api`，后续按版本加入其他进程对应的服务。宿主机无需安装 Rust，构建在镜像构建阶段完成；开发调试仍可使用 Cargo 直接运行。

容器内的程序入口：

```text
query-api [--config <PATH>]
query-api --help
query-api --version
```

`--config` 默认为 `config/query-api.toml`，相对路径按进程当前工作目录解析。显式指定文件时以该文件替代默认文件，不把两份文件隐式合并。`--help` 和 `--version` 不依赖配置文件，不启动 HTTP 服务。

从仓库根目录运行：

```sh
cargo run --locked -p query-api -- --config config/query-api.toml
```

在其他目录运行构建产物时，应使用明确的配置路径：

```sh
/absolute/path/to/query-api --config /absolute/path/to/query-api.toml
```

### 3.2 请求入口

Compose 将容器的 8080 端口映射到宿主机 `127.0.0.1:8080`，运行后可调用：

```sh
curl -i http://127.0.0.1:8080/api/v1/health
curl -i http://127.0.0.1:8080/api/v1/version
```

成功返回 HTTP 200 和 JSON。不存在的路由返回 404；对已有路由使用不支持的方法返回 405，响应细节见第 5 节。

### 3.3 停止与重启

使用以下命令停止、清理或重新加载挂载配置：

```sh
docker compose stop
docker compose down
docker compose restart query-api
```

容器停止时程序接收 SIGTERM，停止接受新连接，并在配置时限内完成在途请求。修改挂载的 TOML 后重启生效；修改 Compose 的环境或端口映射后执行 `docker compose up -d --force-recreate query-api`。

文中的命令是实现后必须验证的使用契约，当前不代表这些代码或运行结果已经存在。

## 4. 配置来源与运行边界

### 4.1 配置文件

交付可直接运行的 `config/query-api.toml`：

```toml
config_version = 1

[server]
host = "127.0.0.1"
port = 8080
shutdown_timeout_seconds = 10

[logging]
level = "info"
format = "json"
```

| 配置项 | 含义与校验 |
| --- | --- |
| `config_version` | 必填，本版仅接受整数 1 |
| `server.host` | 必填，IPv4 或 IPv6 地址；本版不做域名解析 |
| `server.port` | 必填，整数 1–65535；端口 0 仅供测试内部使用，不接受为运行配置 |
| `server.shutdown_timeout_seconds` | 必填，整数 1–60，退出等待上限 |
| `logging.level` | 必填，`trace`、`debug`、`info`、`warn` 或 `error` |
| `logging.format` | 必填，`json` 或 `text` |

默认文件必须存在；缺失必填项、未知配置项、格式错误或不支持的版本均导致启动失败。配置只描述本版已有能力，不提前填入未使用的数据库和钱包配置。

### 4.2 加载与覆盖规则

顺序与详细设计一致：选定配置文件 → 环境变量覆盖 → 密钥引用解析（本版无此配置）→ 完整校验。先解析文件结构，再应用环境覆盖，最终值必须符合类型和范围。

| 环境变量 | 对应配置 |
| --- | --- |
| `ROBOTECH_SERVER_HOST` | `server.host` |
| `ROBOTECH_SERVER_PORT` | `server.port` |
| `ROBOTECH_SHUTDOWN_TIMEOUT_SECONDS` | `server.shutdown_timeout_seconds` |
| `ROBOTECH_LOG_LEVEL` | `logging.level` |
| `ROBOTECH_LOG_FORMAT` | `logging.format` |

环境变量未设置时使用文件值；设置为空或非法值时失败，不回退。`config_version` 由文件声明，不接受环境覆盖。不读取未在上表声明的环境变量作为配置，也不从 HTTP 请求获取运行配置。

示例：

```sh
ROBOTECH_SERVER_PORT=8081 cargo run --locked -p query-api -- --config config/query-api.toml
```

本版不增加 `--port` 等重复覆盖入口。启动日志只输出配置文件路径及允许公开的有效监听地址、日志级别、格式和退出时限；不得直接打印整份配置对象，避免后续新增敏感项时被自动泄露。

### 4.3 兼容与部署边界

后续版本新增配置须明确默认值或迁移说明；必填项或语义变化须更新 `config_version` 并提供可操作的迁移说明。不静默忽略拼错的字段。

本版通过 Docker 与 Compose 运行，镜像目标为 Linux x86_64，macOS Docker Desktop 可用于开发验证；Apple Silicon 上需验证原生 Linux arm64 构建或明确使用 linux/amd64 模拟运行。Kubernetes 部署在需要时加入，运行配置和程序入口保持一致。程序直接提供 HTTP；本版不配置 TLS 终止或外部代理。

## 5. 接口与数据定义

### 5.1 通用规则

路径沿用 `/api/v1`，响应 `Content-Type` 为 `application/json`。成功响应使用详细设计的 `data` 和 `meta` 信封：

```json
{
  "data": {},
  "meta": {
    "trace_id": "trace_...",
    "schema_version": 1
  }
}
```

每个请求在服务端生成唯一 `trace_id`，写入响应体、`X-Trace-Id` 响应头和请求日志。本版不将客户端提供的同名头直接作为服务端 trace ID。后续分布式链路传播通过相应中间件接入，不改变已有响应字段。

时间采用 UTC RFC3339，至少毫秒精度。API 路径版本、响应 schema 版本和程序版本分别表达接口兼容性、数据结构和发布版本，不随一次请求变化。

### 5.2 健康接口

`GET /api/v1/health`

```json
{
  "data": {
    "status": "ok",
    "service": "query-api",
    "started_at": "2026-10-04T02:00:00.000Z"
  },
  "meta": {
    "trace_id": "trace_example",
    "schema_version": 1
  }
}
```

HTTP 200 表示当前程序已完成初始化，HTTP 请求能够被处理。`started_at` 在程序初始化完成、开始服务时固定；每次请求返回相同值。

本版健康检查不检查不存在的外部依赖，也不表示未来数据源、数据库或业务结果完整。后续依赖健康与 readiness 按实际需要增加独立检查，不改变这个接口的基础存活语义。

### 5.3 版本接口

`GET /api/v1/version`

```json
{
  "data": {
    "service": "query-api",
    "version": "0.0.0"
  },
  "meta": {
    "trace_id": "trace_example",
    "schema_version": 1
  }
}
```

`version` 来自 package 编译时版本，首版为 `0.0.0`，不得由配置文件覆盖。CLI `--version`、启动日志和 HTTP 版本接口使用同一版本来源。本版不强制添加依赖构建机器状态的 Git revision 字段。

### 5.4 错误响应

沿用详细设计的 `code`、`message`、`trace_id` 和可选 `details`：

```json
{
  "code": "RESOURCE_NOT_FOUND",
  "message": "Route not found",
  "trace_id": "trace_example"
}
```

| 场景 | HTTP | code | 说明 |
| --- | --- | --- | --- |
| 未知路径 | 404 | `RESOURCE_NOT_FOUND` | 不返回框架默认纯文本 |
| 已有路径、不支持的方法 | 405 | `METHOD_NOT_ALLOWED` | 本版补充的运行接口错误码，返回正确 `Allow` 头 |
| 未预期的内部处理错误 | 500 | `INTERNAL_INVARIANT_VIOLATION` | 客户端获取概括错误，详细原因写入日志 |

本版声明的业务方法为 GET。框架按 HTTP 语义支持 HEAD，使用与 GET 一致的状态及响应头、不返回响应体；POST 等方法返回 405，`Allow` 包含 GET 和 HEAD。未来业务 API 的参数验证与授权复用已有错误处理机制。

示例 trace ID 和时间仅用于说明格式，不是实际运行结果。

## 6. 程序处理流程

### 6.1 启动

1. 解析启动参数；帮助和版本命令直接输出并正常退出。
2. 确定配置路径，读取并解析 TOML。
3. 应用环境覆盖，校验配置版本、必填项、未知字段和取值。
4. 按有效配置初始化日志，建立不可变运行配置。
5. 注册终止信号处理，创建应用状态、通用中间件和路由。
6. 绑定配置地址与端口；绑定失败记录明确原因并退出。
7. 记录固定启动时间和 `server_started` 日志，开始接受请求。

日志初始化前的参数与配置错误写入标准错误。只有成功绑定并完成初始化后才能记录启动成功。

### 6.2 请求

接收请求 → 生成 trace ID → 路由匹配 → handler 返回 DTO → 统一响应转换 → 写入 trace 响应头 → 记录请求完成日志。

404 和 405 同样经过通用错误转换和请求日志。handler 只负责运行接口的状态读取及响应；不能在 handler 中加入后续版本的协议请求、成交解析或业务计算。

请求完成日志包括 UTC 时间、级别、service、trace ID、HTTP 方法、路由模板、状态码和处理时长。未知路径使用固定 `unmatched` 标签；不记录完整 URL、查询参数、请求体或授权头。请求 2xx 使用 info，4xx 使用 warn，5xx 使用 error；日志级别配置控制实际输出。

### 6.3 退出

收到 SIGINT 或 SIGTERM → 记录 `shutdown_started` → 停止接受新连接 → 等待在途请求结束 → 释放监听与资源 → 记录 `shutdown_completed` → 返回退出码 0。

等待上限覆盖 HTTP draining 和必要资源收尾。到期后取消未完成任务，记录 `shutdown_timeout` 并以非零退出码结束。不要让无限等待的连接阻止程序退出。实现须通过 Tokio 任务与 Axum 的退出机制控制等待，而不是在信号处理函数中同步阻塞。

## 7. 程序结构设计

### 7.1 本版工程结构

```text
robotech/
├── Cargo.toml
├── Cargo.lock
├── rust-toolchain.toml
├── compose.yaml                       # 多进程容器启动入口，首版一个服务
├── .dockerignore
├── gateway/
│   └── query-api/
│       ├── Dockerfile                  # query-api 多阶段构建
│       ├── Cargo.toml
│       └── src/
│           ├── main.rs
│           ├── lib.rs
│           ├── bootstrap.rs
│           ├── config.rs
│           ├── state.rs
│           ├── lifecycle.rs
│           ├── logging.rs
│           └── http/
│               ├── mod.rs
│               ├── router.rs
│               ├── response.rs
│               ├── error.rs
│               ├── middleware.rs
│               └── handlers/
│                   ├── mod.rs
│                   ├── health.rs
│                   └── version.rs
├── config/
│   └── query-api.toml
├── tests/
│   ├── integration-tests/
│   │   └── query_api.rs
│   └── fixtures/
│       └── config/
└── docs/
    └── version0.0.md
```

根 workspace 首版只注册 `gateway/query-api`。其 package 同时包含 library 和 `query-api` binary，测试可以调用 library 的应用构建入口，避免启动进程才能验证所有接口。

根 `Cargo.toml` 统一管理依赖，提交 `Cargo.lock`；`rust-toolchain.toml` 固定实际验证过的 Rust stable 工具链版本。沿用 Tokio、Axum、Serde；TOML 解析、tracing 日志及 CLI 参数库按职责引入，具体兼容版本在实施时确定并锁定。

根 `tests/integration-tests/query_api.rs` 通过 query-api package 的显式 `[[test]]` 目标注册，确保 `cargo test --workspace` 会运行；不能仅创建文件而不注册测试。固定配置样本位于根 `tests/fixtures/config/`。

### 7.2 模块职责

| 模块 | 职责 |
| --- | --- |
| `main` | 调用 bootstrap，输出最终启动错误并返回退出码 |
| `bootstrap` | 参数解析、配置加载、日志初始化、应用构建及运行装配 |
| `config` | 强类型配置、环境覆盖和完整校验，不依赖 HTTP handler |
| `state` | 程序版本、固定启动时间等不可变应用状态 |
| `lifecycle` | 信号监听、服务运行、有界优雅退出 |
| `logging` | tracing subscriber 配置，统一字段与脱敏规则 |
| `http/router` | 路由及 fallback 注册 |
| `http/response` | 成功信封及公共 schema 元信息 |
| `http/error` | 错误 DTO、HTTP 状态与错误转换 |
| `http/middleware` | 请求 trace ID、响应头与请求完成日志 |
| `http/handlers` | 健康与版本接口的状态读取和 DTO 返回 |

配置加载入口接受明确的环境覆盖输入，便于在测试中验证覆盖规则而不修改进程全局环境。退出逻辑接受可注入的停止信号，生产连接 OS 信号，测试使用可控信号。

### 7.3 依赖与后续扩展

`main/bootstrap` 装配 `config`、`logging`、`lifecycle` 和 `http`；HTTP handler 读取 `state` 并返回响应。配置与生命周期不依赖具体 handler，HTTP 公共响应不依赖未来的协议 DTO 或数据库 row。

v0.1 的成交查询仍从 `gateway/query-api` 进入，数据获取和解析按详细设计放在 `services/trade-log` 与 `protocols/implementations/hyperliquid`。网关通过业务查询接口调用这些能力；业务 crate 不反向依赖 Axum 或网关。

v0.2 增加交易日志服务的存储与 migrations，v0.3 增加采集任务和相应运行入口。后续拆分独立服务进程沿用相同配置与生命周期约定，现有网关继续提供 HTTP 入口。

公共模块在出现跨进程实际复用需求时抽取为公共组件，保持配置字段和响应行为兼容。首版不创建未使用的业务 crate、空服务或未来全部七个进程，也不将运行配置塞入 `shared-types` 业务值对象包。

## 8. 异常处理

### 8.1 进程错误

| 退出码 | 场景 | 处理 |
| --- | --- | --- |
| 0 | 帮助、版本输出或正常停止 | 正常退出 |
| 2 | 参数或配置错误 | 标明参数或配置键及原因，不开始监听 |
| 3 | 初始化或监听失败 | 标明阶段、监听地址及系统原因，不输出启动成功 |
| 4 | 退出超时或运行期服务失败 | 记录原因，结束进程，避免假装仍正常运行 |

配置文件缺失、不可读、TOML 格式错误、未知字段和非法环境值均属于配置错误。端口占用与权限不足属于监听失败。日志初始化失败不得忽略。

错误信息应足够定位问题，但不直接附带整份配置、原始请求或未来凭证值。

### 8.2 HTTP 错误

404、405 和运行期可恢复处理错误按第 5 节返回 JSON，并具有 trace ID。可恢复请求错误不导致进程退出。影响整个 HTTP 服务运行的错误记录后退出，不能只留下日志却继续宣称服务正常。

本版没有外部数据请求，不实现来源重试机制。后续版本在业务服务或相应适配器内增加重试，统一转换为网关错误格式。

## 9. 验证与验收

### 9.1 必要自动验证

| 验证内容 | 关键断言 |
| --- | --- |
| 配置校验 | 正常文件通过；缺字段、未知字段、非法端口、非法日志级别和错误版本失败 |
| 环境覆盖 | 覆盖优先级确定，未设置沿用文件，空值或非法值报错 |
| 接口契约 | health/version 的状态码、JSON 信封、schema、字段类型和 Content-Type 正确 |
| Trace 关联 | 请求具有不同 trace ID，响应头、响应体和日志中的 ID 一致 |
| 错误与方法 | 未知路由 404；POST 已有路由 405 且 Allow 正确；HEAD 无响应体 |
| 程序状态 | 多次 health 返回同一启动时间；HTTP 与 CLI 版本一致 |
| 启动失败 | 缺失配置或端口占用返回预期退出码，不记录成功启动 |
| 优雅退出 | 可控在途请求在时限内完成，超时任务被终止，停止后可以重新绑定端口 |

接口与生命周期测试使用本地 socket，测试内部允许动态端口；固定端口只用于手工验收，避免测试彼此冲突。无需访问 Hyperliquid 或启动数据库。

### 9.2 实际运行验收

1. 从仓库根目录执行 `docker compose up --build -d`，确认容器正常运行、宿主机能访问映射端口，日志包含实际监听地址和版本。
2. 请求健康与版本接口，核对第 5 节契约及请求日志。
3. 请求未知路径，并对已有路径发 POST、HEAD，核对状态和错误格式。
4. 修改宿主机端口映射并重建容器，确认新端口工作、原端口释放；修改挂载配置并重启，确认配置生效。
5. 使用环境变量覆盖文件端口，确认最终监听值。
6. 使用无效配置、缺失文件和已占用端口启动，确认失败原因与退出码。
7. 使用 `docker compose stop` 验证 SIGTERM 正常退出并能重启；另验证程序 SIGINT 处理。
8. 用可控请求验证 draining 与退出超时，确认进程不会无限等待。

验收报告记录工具链、构建命令、程序 revision、配置路径、接口结果、退出结果及结论。Linux 基准环境须完成启动与信号验收；本地 macOS 结果注明环境，不能直接替代 Linux 发布验证。

### 9.3 完成标准

Docker 构建及 Compose 启停可实际使用，挂载配置和端口映射有效；配置和 HTTP 入口可以实际使用；正确与错误路径的行为符合文档；基本日志可用于定位请求；启动、退出和重启验收通过；工程依赖方向符合设计。未接入业务接口是本版范围，不用虚假的成交响应填充演示。

本文件为待实施设计，验收报告在实现后填写实际结果。

## 10. 交付与运行说明

### 10.1 交付物

- 根 workspace、锁文件与固定工具链。
- `gateway/query-api` 的 library、binary 及上述运行模块。
- 可直接使用的 `config/query-api.toml`。
- `gateway/query-api/Dockerfile`、根 `compose.yaml` 和 `.dockerignore`。
- 已注册的接口、配置和生命周期测试及固定样本。
- 本文档与实际运行验收报告。

构建产物加入忽略规则。运行日志输出到标准错误，便于终端查看和后续平台采集；不在本版自动创建日志文件或部署 Loki。

### 10.2 构建与检查

实现后执行：

```sh
cargo build --locked -p query-api
cargo test --locked --workspace
cargo fmt --all -- --check
cargo clippy --locked --workspace --all-targets -- -D warnings
cargo run --locked -p query-api -- --help
cargo run --locked -p query-api -- --version
```

发布构建：

```sh
cargo build --locked --release -p query-api
./target/release/query-api --config config/query-api.toml
```

运行要求是匹配目标平台的构建产物、可读取配置文件及可绑定的监听端口。构建可能需要获取依赖，运行基础接口无需外部网络连接。

### 10.3 Docker 构建与运行契约

Dockerfile 使用多阶段构建：构建阶段使用锁定工具链及 `cargo build --locked --release -p query-api`，运行阶段仅保留可执行文件及所需运行库。基础镜像版本在实施时固定，构建上下文为仓库根目录。运行镜像使用非 root 用户及只读根文件系统。

`ENTRYPOINT` 使用 exec 形式直接启动 `query-api`，避免 shell 截留终止信号。默认参数为 `--config /etc/robotech/query-api.toml`。运行镜像不内置部署配置；Compose 只读挂载仓库配置文件。镜像构建排除 `.git`、`target`、本地日志及敏感文件。

Compose 契约：

```yaml
services:
  query-api:
    image: robotech/query-api:0.0.0
    build:
      context: .
      dockerfile: gateway/query-api/Dockerfile
    environment:
      ROBOTECH_SERVER_HOST: "0.0.0.0"
      ROBOTECH_SERVER_PORT: "8080"
    ports:
      - "127.0.0.1:8080:8080"
    volumes:
      - type: bind
        source: ./config/query-api.toml
        target: /etc/robotech/query-api.toml
        read_only: true
        bind:
          create_host_path: false
    read_only: true
    stop_grace_period: 15s
    restart: "no"
```

本地直接运行的配置监听 `127.0.0.1`；Compose 通过环境覆盖使容器监听 `0.0.0.0`，宿主机访问通过端口映射进入。宿主机环境变量不会自动传入容器，修改应用覆盖值须在 Compose 的 `environment` 中声明。容器内部端口由 Compose 固定为 8080，更改时同步调整映射，不能只修改 TOML。

容器停止宽限时间须大于应用退出时限；默认分别为 15 秒和 10 秒。调整应用时限时同步调整 Compose。首版不自动重启，便于观察启动失败与退出码。容器运行状态不能代替接口健康判断，验收必须实际请求 health。

通过 `docker compose logs query-api` 读取日志；宿主机端口占用和配置挂载失败可能由 Docker 在程序启动前拒绝，此类错误使用 Compose 输出定位，不归入程序退出码。程序自身的错误可通过日志及容器退出状态检查。

构建验收包含 `docker compose config`、镜像构建、接口访问、配置重载、容器停止及重建；容器停止后退出码应为 0，不能依赖 SIGKILL 完成正常停止。镜像无需依赖宿主机 Rust 或构建产物。

Compose 字段依据：[Docker Compose 服务配置](https://docs.docker.com/reference/compose-file/services/)。以上为待实现契约，本次文档更新未创建镜像或 Compose 文件。

### 10.4 后续迭代约定

后续版本继续使用 `query-api` 入口、配置加载规则、路由版本、JSON 信封、trace ID 和退出机制。新增业务放入设计指定服务，通过查询接口接入网关；业务与协议代码保持独立于 HTTP 展示层。

配置或响应发生不兼容变化时明确提供迁移及版本说明。各版本增加已有程序的能力，不重新建立一套命令入口、配置或错误处理；必要的公共组件抽取与进程拆分保持已有使用方式可延续。

后续独立进程分别作为 Compose service 运行，每个容器运行一个职责明确的应用进程；只向宿主机暴露需要访问的 HTTP 入口，其余进程通过容器网络通信。按实际功能加入服务、依赖与数据卷，沿用同一个 Compose 启动入口。
