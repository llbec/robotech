# V0.000 版本开发记录

## 1. 版本定义

| 项目 | 内容 |
| --- | --- |
| 版本号 | V0.000 |
| 版本状态 | IN DEVELOPMENT |
| 设计日期 | 2026-09-21 |
| 开发状态 | 已完成代码实现，等待实际地址端到端验收 |
| 版本性质 | 可独立运行、可实际使用的最小业务版本 |
| 功能模块 | 交易日志服务 |
| 数据范围 | Hyperliquid 合约成交与现货成交 |

V0.000 遵循敏捷迭代原则：版本完成后必须能够独立启动并交付实际功能，而不是只运行测试、验证接口或输出调研结论。

## 2. 本版本交付能力

用户可以通过 HTTP API 添加一个 Hyperliquid 地址。服务随后自动回补该地址近期的合约与现货成交，并持续接收实时成交。用户可以通过接口查看原始成交和标准化成交；新的主动成交可以及时转发给一个第三方 webhook。

完整业务闭环：

```text
确认独立 PostgreSQL 可用并由 Docker Compose 启动服务
→ HTTP API 添加监控地址
→ 回补该地址近期合约/现货成交
→ 建立 WebSocket 实时成交订阅
→ 保存原始成交
→ 标准化并保存成交事实
→ HTTP API 查询原始成交和标准成交
→ 将新的实时主动成交投递给第三方 webhook
```

完成 V0.000 后，即使后续版本尚未开发，用户也能使用它监控已知地址的实际成交活动。

## 3. 范围控制

### 3.1 本版本包含

- 一个持续后台运行的 Rust 服务 `trade-log-server`。
- 使用 PostgreSQL 持久化监控地址、原始成交、标准成交和 webhook 投递记录。
- 通过 HTTP API 添加、列出和查看监控地址。
- 相同地址重复添加时幂等返回原任务。
- 同时启用的监控地址最多为 10 个；达到官方用户级 WebSocket 地址限制后拒绝新增并返回明确错误。
- 通过 Hyperliquid 官方 `userFillsByTime` 回补近期历史成交。
- 通过 Hyperliquid 官方 WebSocket `userFills` 接收实时成交。
- 获取并保存 `meta`、`spotMeta`，正确区分合约和现货市场。
- 保存未经改写的来源响应和消息。
- 将成交转换为统一的 `TRADE` 标准事实。
- 对主动合约成交和现货买卖设置 `copy_eligible=true`。
- 只将进入实时阶段后首次观察到的可跟单成交发送至一个 webhook。
- WebSocket 断线自动重连，并通过 HTTP 重叠补取缺口。
- 提供原始成交、标准成交和 webhook 投递记录查询接口。
- 提供存活和就绪检查。
- 提供 Docker Compose 启动交易日志服务的单一部署路线，连接独立部署的 PostgreSQL。

### 3.2 本版本明确不包含

- 资金费、充值、提现、内部划转、奖励等账户账本事件。
- 合约状态、保证金、现货余额和账户价值快照。
- 历史价格保存和任意时点估值。
- 账户价值、可投资金额、净盈亏、交易次数和胜率计算。
- 清算事件的独立处理；若清算表现为成交，只保存原始成交并标记来源语义，不作为主动跟单成交。
- Nansen、节点数据流和其他第三方数据源。
- 数据完整审核和多来源对账。
- Kafka/Redpanda 和通用消息订阅系统。
- 多个 webhook 接收方及动态订阅配置。
- 跟单规则、风险控制、钱包签名和交易发送。
- 订单提交、挂单、撤单和拒单跟踪。
- Kubernetes 部署。
- Hyperliquid 之外的链和协议。

这些能力不是被系统总体目标删除，而是由后续版本逐项增加。每个后续版本仍应保持可部署、可运行，并兼容已有数据。

## 4. 数据来源

### 4.1 官方接口

| 用途 | 接口 |
| --- | --- |
| HTTP 历史成交 | `POST https://api.hyperliquid.xyz/info`，`type=userFillsByTime` |
| WebSocket 实时成交 | `wss://api.hyperliquid.xyz/ws`，`type=userFills` |
| 合约元数据 | `POST /info`，`type=meta` |
| 现货元数据 | `POST /info`，`type=spotMeta` |

官方参考：

- [Hyperliquid Info endpoint](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint)
- [Hyperliquid Spot info](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint/spot)
- [Hyperliquid WebSocket subscriptions](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/websocket/subscriptions)
- [Hyperliquid rate limits](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/rate-limits-and-user-limits)

### 4.2 历史范围

- 添加地址时可以提供 `start_time`。
- 未提供时默认回补最近 7 天。
- 回补结束时间固定为监控任务创建时间。
- 官方接口只能提供其允许访问的历史范围；超出范围时记录实际 `coverage_start` 和 `history_complete=false`。
- V0.000 不声称恢复地址成立以来的全部历史。

### 4.3 实时与断线恢复

- 历史回补完成后，任务状态进入 `LIVE`。
- WebSocket snapshot 只用于衔接，不触发历史 webhook。
- 实时连接中断时任务进入 `DEGRADED`，记录断线时间。
- 重连后使用 `userFillsByTime` 从最后成功水位前 2 分钟开始补取。
- HTTP 与 WebSocket 的重叠成交依靠稳定 `fact_id` 去重。
- 补偿完成后任务恢复 `LIVE`。

## 5. 标准成交事实

### 5.1 公共结构

```text
fact_id
revision = 1
fact_type = TRADE
schema_version = 1
chain_id = hyperliquid:mainnet
protocol = hyperliquid
account
account_key
occurred_at
observed_at
ordering_key
source = HYPERLIQUID_OFFICIAL
source_ref
raw_log_id
payload
```

`fact_id` 必须由稳定来源身份生成，使同一成交通过 HTTP、WebSocket snapshot 或实时消息到达时得到相同 ID。不能仅以完整 JSON 哈希作为业务身份。

### 5.2 成交字段

```text
instrument_type = PERPETUAL | SPOT
market
base_asset
quote_asset
side = BUY | SELL
position_effect = OPEN | INCREASE | DECREASE | CLOSE | REVERSE | NONE | UNKNOWN
trigger_type = USER | PROTOCOL | LIQUIDATION
copy_eligible
price
quantity
notional
fee
fee_asset
reported_realized_pnl
order_id
transaction_hash
```

规则：

- 使用成交发生时有效的 `meta/spotMeta` 版本解析市场，不按 `coin` 字符串外观猜测。
- 每个原始 fill 形成一个标准成交事实，不等待整张订单完成。
- 主动合约成交和主动现货买卖为 `copy_eligible=true`。
- 清算、强制减仓、交割、结算和修正数据为 `copy_eligible=false`。
- `copy_eligible` 只表示事件可以进入后续跟单规则，不表示必须执行。
- 金额、价格和数量使用十进制字符串，不经过二进制浮点数。

## 6. HTTP API

### 6.1 添加监控地址

`POST /api/v1/monitored-addresses`

```json
{
  "address": "0x...",
  "start_time": "2026-09-20T00:00:00Z"
}
```

- `address` 必填，规范化为小写。
- `start_time` 可选，必须为 UTC ISO 8601。
- 新建成功返回 HTTP 201。
- 地址已经存在时返回 HTTP 200 和原任务。
- 创建任务后立即返回，不等待历史回补完成。

### 6.2 查询接口

| 方法与路径 | 返回内容 |
| --- | --- |
| `GET /api/v1/monitored-addresses` | 监控地址列表 |
| `GET /api/v1/monitored-addresses/{address}` | 状态、覆盖范围、水位、最后接收时间和最近错误 |
| `GET /api/v1/monitored-addresses/{address}/raw-trades` | 该地址的原始成交 |
| `GET /api/v1/monitored-addresses/{address}/trades` | 该地址的标准成交事实 |
| `GET /api/v1/webhook-deliveries` | webhook 投递状态和尝试次数 |
| `GET /health/live` | 进程存活状态 |
| `GET /health/ready` | 数据库、配置和后台调度器就绪状态 |

成交查询使用游标分页，默认 100 条、最大 500 条，按 `occurred_at + ordering_key + sub_index` 稳定升序返回。

V0.000 不提供删除、暂停和恢复监控接口；服务重启后自动恢复所有监控任务。

## 7. Webhook 输出

V0.000 配置一个固定接收地址：

```text
POST {EVENT_WEBHOOK_URL}
Content-Type: application/json
Idempotency-Key: {event_id}
X-Robotech-Signature: sha256={HMAC_HEX}
```

发送条件：

```text
fact_type = TRADE
copy_eligible = true
任务状态 = LIVE
该 fact_id 第一次进入 outbox
```

投递语义：

- 标准成交与 outbox 在同一数据库事务写入。
- HMAC 使用 `EVENT_WEBHOOK_SECRET` 对原始 HTTP body 执行 SHA-256 签名。
- HTTP 2xx 表示成功。
- 网络错误、408、429 和 5xx 使用有界指数退避。
- 其他 4xx 标记永久失败并保留脱敏响应摘要。
- 服务提供至少一次投递，接收方必须使用 `Idempotency-Key` 幂等处理。
- 历史回补和 WebSocket snapshot 不发送，避免启动服务时产生历史补单。

## 8. 数据库

V0.000 交付数据库迁移，建立：

| 表 | 用途 |
| --- | --- |
| `monitored_addresses` | 地址、任务状态、覆盖范围、水位和错误 |
| `collection_runs` | 历史回补、断线补偿及执行统计 |
| `raw_trade_logs` | HTTP/WebSocket 原始成交和解析状态 |
| `market_metadata_versions` | 合约和现货市场映射版本 |
| `trade_fact_versions` | 标准成交事实版本 1 |
| `trade_facts_current` | 当前有效事实指针 |
| `outbox_events` | 待投递实时成交 |
| `webhook_deliveries` | webhook 投递过程和结果 |

核心约束：

- `(source, source_event_id)` 唯一。
- `(fact_id, revision)` 唯一。
- 一个 `raw_trade_log` 可以关联一个或多个标准事实。
- 原始 payload 和标准 payload 分开保存。
- 只有事务提交成功后才推进采集水位。

## 9. 程序模块

```text
services/trade-log/
├── src/
│   ├── main.rs              # 服务装配、后台任务和退出信号
│   ├── config.rs            # 环境配置及启动校验
│   ├── domain.rs            # 地址、来源 fill 和标准事实模型
│   ├── acquisition.rs       # Hyperliquid HTTP/WS 接入
│   ├── monitoring.rs        # 地址任务、回补、重连和补偿
│   ├── market_metadata.rs   # 合约与现货市场映射
│   ├── normalization.rs     # 成交标准化及稳定身份
│   ├── repository.rs        # PostgreSQL 持久化和查询
│   ├── publishing.rs        # outbox、HMAC 和 webhook
│   ├── api.rs               # HTTP 管理与查询接口
│   └── bin/
│       └── trade-log-migrate.rs
└── migrations/
    └── 0001_v0_000.sql
```

V0.000 采用单服务目录内聚结构；采集、标准化、存储、发布和 API 仍以独立模块隔离，后续复杂度上升时可以在不改变模块职责的前提下拆分子目录或 crate。

## 10. 交付件

V0.000 完成时必须同时交付以下内容，缺少任一必需项都不能宣布版本完成。

### 10.1 可执行程序

- Rust workspace 和可重复构建的依赖锁文件。
- `trade-log-server` 可执行程序。
- 启动时完成配置校验，支持 SIGTERM 优雅退出。

### 10.2 数据库

- PostgreSQL 建表和索引迁移文件。
- 独立迁移命令。
- 空数据库初始化和重复执行迁移测试。

### 10.3 接口

- 地址监控管理接口。
- 原始成交、标准成交和投递记录查询接口。
- OpenAPI 文档。
- 存活和就绪检查接口。

### 10.4 数据采集与输出

- Hyperliquid HTTP 客户端。
- Hyperliquid WebSocket 客户端及重连机制。
- 合约与现货市场元数据解析。
- 原始成交持久化和标准化。
- outbox 与签名 webhook 投递。

### 10.5 部署与配置

- `Dockerfile`。
- `compose.yaml`，只包含 `trade-log-server`；PostgreSQL 独立部署。
- `.env.example`，只包含字段说明和非敏感示例。
- 外部数据库连接、服务端口和健康检查配置。
- 启动、停止、查看日志和清理测试数据的操作说明。

### 10.6 测试与说明

- 领域规则单元测试。
- Hyperliquid HTTP/WebSocket fixture 测试。
- PostgreSQL repository 集成测试。
- 从添加地址到查询成交和 webhook 接收的端到端测试。
- README 运行说明、API 示例、限制和已知问题。
- 测试执行记录和构建版本信息。

## 11. 运行方式

V0.000 的唯一部署路线为 Docker Compose：

```text
复制 .env.example 为本地配置
→ 填写独立 PostgreSQL 连接地址、webhook 地址和密钥
→ 确认数据库可从容器网络访问
→ docker compose up -d
→ 等待 /health/ready 返回 200
→ 调用 API 添加监控地址
```

主要配置：

| 配置 | 默认值/要求 |
| --- | --- |
| `DATABASE_URL` | 必填，指向独立部署且容器可访问的 PostgreSQL |
| `HTTP_LISTEN_ADDR` | `0.0.0.0:8080` |
| `HYPERLIQUID_HTTP_URL` | `https://api.hyperliquid.xyz/info` |
| `HYPERLIQUID_WS_URL` | `wss://api.hyperliquid.xyz/ws` |
| `DEFAULT_HISTORY_LOOKBACK_SECONDS` | `604800`（7 天） |
| `HTTP_OVERLAP_SECONDS` | `120`（2 分钟） |
| `HTTP_RECONCILE_INTERVAL_SECONDS` | `30` |
| `MAX_MONITORED_ADDRESSES` | `10`，V0.000 不允许配置为更大值 |
| `EVENT_WEBHOOK_URL` | 必填 |
| `EVENT_WEBHOOK_SECRET` | 必填，不写入日志或数据库普通字段 |
| `RUST_LOG` | `info` |

## 12. 验收标准及解释

### 12.1 可以安装并持续运行

独立 PostgreSQL 可用时，执行 `docker compose up -d` 后 `trade-log-server` 成功启动并自动执行数据库迁移。等待初始化完成后，`/health/live` 和 `/health/ready` 返回 HTTP 200；服务连续运行，不出现反复退出和重启。数据库不可连接或权限不足时，服务启动失败并给出明确错误。

这证明交付件不是只能运行测试的代码，而是一个可部署服务。

### 12.2 可以添加监控地址

调用地址创建接口后立即获得任务 ID 和状态。相同地址重复提交不会创建重复任务。服务重启后任务仍存在并自动继续运行。

这证明监控任务已持久化，不依赖某次进程内存。

### 12.3 可以回补近期历史成交

添加存在历史交易的地址后，服务能够分页取得指定起点至任务创建时刻的成交，并保存实际覆盖范围。官方接口无法提供更早数据时，接口明确返回 `history_complete=false`，不能伪装完整。

这证明服务可以在开始监控后立即提供近期记录，并明确数据边界。

### 12.4 可以识别合约和现货

固定测试样本中的合约成交被标准化为 `PERPETUAL`，现货成交被标准化为 `SPOT`，市场和基础/计价资产均能通过对应元数据版本追溯。

这证明服务没有把所有 Hyperliquid fill 都错误当作合约。

### 12.5 可以持续接收实时成交

历史回补完成后任务进入 `LIVE`。测试地址发生新成交时，服务保存 WebSocket 原始消息和标准成交事实，并更新最后接收时间。

这证明服务不仅能导入历史文件，还能实际监控地址。

### 12.6 断线后不会明显漏数或重复

测试中主动中断 WebSocket，再恢复连接。服务进入 `DEGRADED`，使用 HTTP 重叠补取断线区间，然后恢复 `LIVE`。重叠成交只保存一个当前事实。

这证明网络短暂故障不会直接造成静默缺口，并验证水位与幂等机制。

### 12.7 可以查询原始成交和标准成交

使用同一地址分别调用 `raw-trades` 和 `trades` 接口，可以分页获得结果。每个标准成交可以通过 `raw_log_id` 找到来源记录，查询结果不会混入其他地址的数据。

这证明原始证据、标准结果和账户隔离均可用。

### 12.8 可以及时转发新的主动成交

任务进入 `LIVE` 后出现主动成交，数据库提交成功后生成 outbox，第三方 webhook 收到带稳定 `Idempotency-Key` 和有效 HMAC 签名的标准成交。历史回补、snapshot、清算或强制成交不触发普通跟单 webhook。

这证明 V0.000 已具备后续跟单服务所需的实时输入，但尚未执行任何跟单交易。

### 12.9 重复数据不会造成重复事实或重复业务事件

同一成交分别通过 HTTP 回补、WebSocket snapshot 和实时消息到达时，生成相同 `fact_id`。数据库只保留一个当前事实，outbox 不重复创建业务事件；网络重试导致的 webhook 重投使用同一个 `Idempotency-Key`。

这证明数据获取的至少一次语义不会转化为重复业务操作。

### 12.10 敏感信息不会泄露

Webhook 密钥不出现在 Git、API 响应、普通数据库字段和日志中。错误日志不得输出完整请求密钥或签名材料。

这证明版本满足最基本的凭证安全要求。

### 12.11 范围外功能确实未混入

本版本的 Compose 只启动交易日志服务，并连接独立 PostgreSQL；不启动账户分析、审核、钱包或跟单执行进程，也不宣称已经能够计算收益或任意时点账户状态。

这证明版本边界清楚，后续能力通过新版本迭代，而不是在 V0.000 内无限扩张。

## 13. 完成定义

只有同时满足以下条件，V0.000 才能从 `IN DEVELOPMENT` 改为 `RELEASED`：

1. 第 10 节全部交付件存在且文档可执行。
2. 第 12 节全部验收标准通过。
3. 自动化测试全部通过，没有被忽略的必需测试。
4. 使用一个实际 Hyperliquid 地址完成手工端到端验收。
5. 已知问题已记录，并且不存在导致数据静默丢失、错误触发 webhook 或密钥泄露的问题。
6. 版本镜像、Git revision、数据库 schema 版本和测试记录能够相互对应。

## 14. 后续迭代方向

后续版本按可运行的纵向能力逐步增加，具体版本号在各版本规划时确定：

1. 加入资金费、充值、提现、内部划转和其他账户账本事件。
2. 加入合约状态、现货余额和价格快照，提供当前账户资金状态。
3. 加入事件回放、历史状态和区间收益计算。
4. 加入 Nansen 或节点数据作为历史补充与完整性复核来源。
5. 增加独立跟单服务，消费可跟单成交并执行规则、签名和交易发送。

## 15. 变更记录

| 日期 | 状态 | 内容 |
| --- | --- | --- |
| 2026-09-07 | SUPERSEDED | 初版 Nansen 单次导入设计 |
| 2026-09-11 | SUPERSEDED | Nansen 数据覆盖验证设计 |
| 2026-09-20 | SUPERSEDED | 官方 API 全量账户事实采集设计，范围仍偏大 |
| 2026-09-21 | PLANNED | 按敏捷可运行版本重新收敛为合约/现货成交监控、查询和实时 webhook 闭环 |
| 2026-09-21 | IN DEVELOPMENT | 完成服务、数据库迁移、Docker Compose、OpenAPI 与运行文档，进入端到端验收阶段 |
