# 链上账户投资分析服务详细设计

## 1. 文档说明

### 1.1 目的

本文档将概要设计细化为可实施的服务、接口、数据模型和处理流程，作为编码、测试、部署和验收依据。

### 1.2 设计范围

本文覆盖：

- 交易日志获取服务。
- 账户数据分析服务。
- 数据完整审核服务。
- 跟单服务。
- 服务间事件、查询接口、存储、一致性、安全和可观测性。

系统面向多链、多协议和多交易类型设计。链和协议差异通过适配器隔离，不将特定链或协议的数据结构、签名方式和交易语义泄漏到通用核心逻辑。

### 1.3 关键原则

- 标准账户事实是账户分析的统一事实模型；标准交易记录是其中的 `TRADE` 类型，在跟单场景中同一条记录同时承担跟单信号作用，不另建内容重复的信号模型。
- 事件强调及时性，完成同步轻量校验后立即发布，不等待异步完整审核。
- 交易日志获取服务是原始日志和标准账户事实的唯一写入方。
- 数据完整审核服务只发现问题和提交修复请求，不直接修改交易日志数据。
- 账户分析结果是可重建的派生数据。
- 账户分析框架与具体计算口径分离；系统提供可用的默认计算方案，但允许按链、协议、账户类型和业务目标持续替换和演进。
- 跟单规则和风险限制属于跟单服务，不进入标准交易事件。
- 金额、价格和数量使用定点十进制数，不使用二进制浮点数。
- 所有跨服务消息按“至少一次”投递设计，生产方和消费方都必须幂等。

## 2. 基准技术架构

系统统一使用 **Rust** 实现。所有核心服务、链与协议适配器、数据处理 Worker、查询 API 和钱包适配器共享统一的 Rust 工程规范和基础库。基础设施组件可以替换，但替换后必须保持本文定义的数据所有权和消息语义。

| 能力 | 基准实现 | 说明 |
| --- | --- | --- |
| 编程语言 | Rust stable | 工作区统一锁定 toolchain |
| 异步运行时 | Tokio | 网络 I/O、定时任务和 Worker 调度 |
| 服务接口 | HTTP/JSON | 内部高吞吐场景可替换为 gRPC |
| 事件总线 | Kafka 或 Redpanda | 分区内有序、可重放、至少一次投递 |
| 业务数据库 | PostgreSQL | 各服务独立 schema 或独立实例 |
| 原始日志归档 | 对象存储 | 数据量较小时可先使用 PostgreSQL JSONB |
| 缓存 | Redis | 非强依赖，用于热点查询、锁和限流 |
| 任务调度 | 数据库任务表 + Worker | 后续可替换为工作流系统 |
| 监控 | OpenTelemetry + 指标/日志/链路平台 | 统一 trace_id、event_id 和 job_id |

### 2.1 Rust 工程约定

- 采用 Cargo workspace 管理服务和共享 crate。
- 公共事件类型、错误码、金额类型、链标识和可观测性封装为独立 crate。
- 网络、数据库和消息处理使用异步接口；CPU 密集型重算使用受控的独立任务池，不阻塞异步运行时。
- 核心域模型不依赖具体 Web、数据库、消息或链 SDK，外部能力通过 trait 和适配器接入。
- 禁止在业务路径中使用未审查的 `unsafe`；必须使用时集中封装并单独评审。
- 依赖使用 workspace 统一版本并提交锁文件，事件 schema 和数据库迁移随代码版本管理。

## 3. 服务与数据所有权

| 数据 | 唯一写入方 | 读取方 |
| --- | --- | --- |
| 采集任务与检查点 | 交易日志获取服务 | 数据完整审核服务 |
| 原始日志 | 交易日志获取服务 | 数据完整审核服务 |
| 标准账户事实及版本 | 交易日志获取服务 | 账户分析、审核服务；其中交易事实另供跟单服务消费 |
| 修复请求与审核结果 | 数据完整审核服务 | 交易日志获取服务、运维端 |
| 账户账本、持仓和分析结果 | 账户数据分析服务 | 查询端、数据完整审核服务 |
| 跟单配置和执行记录 | 跟单服务 | 用户端、运维端 |

禁止跨服务直接修改数据表。审核服务发现错误后创建修复请求，由交易日志获取服务回补、重新解析并发布新版本事件。

## 4. 标准账户事实与交易事件设计

### 4.1 事件信封

所有业务事件使用统一信封：

```json
{
  "event_id": "evt_...",
  "event_type": "account.fact.v1",
  "schema_version": 1,
  "occurred_at": "2026-09-04T01:02:03.456Z",
  "observed_at": "2026-09-04T01:02:04.120Z",
  "published_at": "2026-09-04T01:02:04.180Z",
  "producer": "trade-log-service",
  "trace_id": "...",
  "payload": {}
}
```

| 字段 | 约束 |
| --- | --- |
| `event_id` | 当前事实版本对应的投递事件 ID，修正后产生新 ID |
| `event_type` | 事件名称与主版本 |
| `schema_version` | 当前负载结构版本 |
| `occurred_at` | 链上交易发生时间 |
| `observed_at` | 系统首次收到数据的时间 |
| `published_at` | 事件写入事件总线的时间 |
| `producer` | 事件生产服务 |
| `trace_id` | 跨服务链路标识 |

### 4.2 标准账户事实公共负载

所有参与账户分析的事实使用统一公共负载，并通过 `fact_type` 区分具体内容：

```json
{
  "fact_id": "chain:network:source-position:0",
  "fact_type": "TRADE",
  "revision": 1,
  "change_type": "UPSERT",
  "confirmation_status": "OBSERVED",
  "chain_id": "namespace:network",
  "protocol": "protocol-name",
  "account": "chain-account",
  "account_key": "namespace:network:protocol-name:chain-account",
  "ordering_key": "source-defined-comparable-position",
  "sub_index": 0,
  "source": "THIRD_PARTY",
  "source_ref": "provider-event-id",
  "raw_log_id": "raw_...",
  "payload": {}
}
```

公共字段约束：

- `fact_id` 是业务事实的稳定 ID；同一事实的修正沿用该 ID 并递增 `revision`。
- `fact_type` 至少支持 `TRADE`、`TRANSFER`、`FEE`、`FUNDING`、`REWARD`、`LIQUIDATION_FEE` 和 `ACCOUNT_SNAPSHOT`。
- `account_key = chain_id + protocol + normalized_account`。
- `ordering_key + sub_index` 在同一 `account_key` 内形成确定顺序；同一源位置拆出的多个事实使用 `sub_index` 排序。
- `change_type=UPSERT` 表示新增或替换当前版本，`RETRACT` 表示撤销。
- 消费方按 `(fact_id, revision)` 幂等，忽略低于当前版本的消息。

### 4.3 标准交易事实 `TRADE`

`TRADE` 的 `payload` 为：

```json
{
  "market": "protocol:ETH-USD",
  "instrument_type": "PERPETUAL",
  "base_asset": "ETH",
  "quote_asset": "USDC",
  "action": "TRADE",
  "trigger_type": "USER",
  "side": "BUY",
  "position_effect": "INCREASE",
  "order_id": "order-or-null",
  "operation_id": "protocol-operation-id",
  "price": "3245.12000000",
  "quantity": "1.25000000",
  "notional": "4056.40000000",
  "fee": "2.03000000",
  "fee_asset": "USDC",
  "reported_realized_pnl": null,
  "reported_pnl_asset": null,
  "reported_pnl_includes_fee": null,
  "reported_pnl_includes_funding": null,
  "transaction_hash": "source-transaction-id",
  "block_number": null,
  "block_hash": null,
  "transaction_index": null,
  "log_index": null,
  "extension": {}
}
```

#### 唯一性与版本

- 交易事实的 `fact_id` 由链适配器使用该链上可重现的事件定位信息生成；在交易查询和跟单上下文中也可称为 `trade_event_id`，两者是同一个值，不是两套 ID。
- EVM 类链可使用 `chain_id + transaction_hash + log_index + sub_index`；其他链可使用区块内序号、成交 ID 或协议提供的稳定位置。
- 一条原始日志解析为多个交易动作时使用 `sub_index` 区分。
- `ordering_key` 由链适配器生成，必须在同一链和账户内可稳定比较。
- `transaction_hash`、`block_number`、`block_hash`、`transaction_index` 和 `log_index` 为可选的链上定位字段；不具备这些概念的链使用 `source_ref` 和 `ordering_key`。
- `revision` 从 1 开始递增，同一 `fact_id` 只允许一个当前版本。

#### 交易语义

通用交易语义：

| 字段 | 取值 | 含义 |
| --- | --- | --- |
| `action` | `TRADE` | 普通成交 |
| `action` | `LIQUIDATION` | 强制减仓或清算成交 |
| `trigger_type` | `USER/PROTOCOL/LIQUIDATION` | 交易由用户、协议机制或清算机制触发 |
| `side` | `BUY/SELL` | 成交的买卖方向 |
| `position_effect` | `OPEN` | 从无持仓到有持仓 |
| `position_effect` | `INCREASE` | 同方向增加仓位 |
| `position_effect` | `DECREASE` | 减少现有仓位 |
| `position_effect` | `CLOSE` | 仓位归零 |
| `position_effect` | `REVERSE` | 平掉原方向并建立相反方向仓位 |
| `position_effect` | `NONE/UNKNOWN` | 不适用持仓语义，或数据源无法直接确定 |

如果数据源不能可靠提供持仓影响，解析器应输出 `position_effect=UNKNOWN`，由账户账本结合前态推导，不得猜测。现货交易可使用 `NONE`。

#### 确认状态

| 状态 | 含义 |
| --- | --- |
| `OBSERVED` | 已观察到，尚未达到配置确认数 |
| `CONFIRMED` | 已达到配置确认数 |
| `ORPHANED` | 所在区块因重组失效 |

为满足跟单时效，`OBSERVED` 可以立即发布。跟单服务自行配置是否接受未确认事件；该配置属于跟单规则，不属于事件本身的风险提示。

#### 数值约束

- 数值在 JSON 中使用十进制字符串。
- 数据库存储使用 `NUMERIC(78, 30)` 或按资产精度收窄。
- `quantity > 0`，买卖方向由 `side` 表达，持仓影响由 `position_effect` 表达，不使用正负数量混合表达。
- 原始精度和标准化精度必须可追溯；不可无提示截断。

### 4.4 资金与费用事实

资金与费用事实使用相同公共负载，其 `payload` 为：

```json
{
  "asset": "USDC",
  "amount": "12.34000000",
  "direction": "CREDIT",
  "funding_type": "FUNDING",
  "related_market": "protocol:ETH-USD",
  "related_position_ref": "position-or-null",
  "related_trade_fact_id": "fact-or-null",
  "transaction_hash": "source-transaction-id",
  "extension": {}
}
```

- `direction=CREDIT/DEBIT` 从账户角度表示增加或减少。
- `funding_type` 根据 `fact_type` 取 `DEPOSIT`、`WITHDRAWAL`、`TRANSFER_IN`、`TRANSFER_OUT`、`TRADE_FEE`、`FUNDING`、`INTEREST`、`REBATE`、`REWARD` 或 `LIQUIDATION_FEE`。
- 能关联交易或市场时填写关联字段；无法可靠关联时保持为空，不得猜测。
- 这些事实供账户账本和审核使用，不作为跟单信号。

### 4.5 账户快照事实

`ACCOUNT_SNAPSHOT` 的 `payload` 为：

```json
{
  "snapshot_at": "2026-09-04T01:02:03.456Z",
  "valuation_currency": "USD",
  "account_value": "10000.00000000",
  "available_balance": "7000.00000000",
  "margin_used": "1000.00000000",
  "unrealized_pnl": "200.00000000",
  "balances": [],
  "positions": [],
  "extension": {}
}
```

快照中的协议报告值是账户状态输入和校验依据，不自动构成第二套分析结果。各链和协议能够提供的事实类型由能力描述声明；缺失时相关计算器必须降级并说明原因。

### 4.6 其他内部事件

#### 修复请求 `trade.repair.requested.v1`

```json
{
  "repair_request_id": "repair_...",
  "reason": "MISSING_BLOCK_RANGE",
  "chain_id": "eip155:42161",
  "source": "provider-a",
  "from_block": 123000,
  "to_block": 123100,
  "accounts": [],
  "requested_by_audit_id": "audit_..."
}
```

#### 重算请求 `account.recalculation.requested.v1`

```json
{
  "recalculation_id": "recalc_...",
  "account_key": "eip155:42161:protocol-name:0x...",
  "affected_from": "2026-09-04T01:02:03.456Z",
  "reason": "TRADE_EVENT_REVISED",
  "source_event_ids": ["evt_..."]
}
```

标准账户事实本身通过 revision 和 change_type 表达修正，不另建一套内容重复的“数据修正事件”。

## 5. 交易日志获取服务详细设计

### 5.1 组件

```text
Source Adapter → Collector → Raw Log Repository → Parser
                                             ↓
                                  Standard Event Repository
                                             ↓
                                           Outbox → Event Bus
```

#### Source Adapter

统一接口：

```text
capabilities(chain) -> SourceCapabilities
fetchRange(chain, positionRange, filter) -> RawLog[]
subscribe(chain, cursor, filter) -> Stream<RawLog>
getFinality(chain, position) -> FinalityStatus
getTransaction(chain, transactionRef) -> RawTransaction
health() -> SourceHealth
```

数据源适配器使用统一的 `positionRange`。其内部可以表示区块范围、时间范围、交易序号或供应商游标。适配器必须保存能够恢复和对账的源侧定位信息。

#### Collector

- 每条链维护一个实时水位和若干历史回补任务。
- 历史任务按链适配器定义的可恢复范围分片；范围完成后原子推进检查点。
- 实时订阅断开时，从持久化游标恢复；无法恢复时，从最后稳定检查点回拉。
- 同一采集范围允许重复执行，依赖原始日志唯一键实现幂等。

#### Parser

- 按 `(chain_id, protocol, parser_version)` 选择解析器。
- 解析器必须是确定性的：相同原始日志和版本产生相同结果。
- 解析失败进入死信表，不阻塞同批其他日志。
- 新解析器先对历史样本回放并对比，再切换为生产版本。

#### Repository 与 Outbox

标准账户事实和 outbox 消息必须在同一数据库事务中写入。独立发布进程将 outbox 发送到消息系统，成功后标记已发布，避免数据库成功但消息丢失。

### 5.2 主流程

1. Collector 从数据源收到原始日志。
2. 校验来源、链、区块定位和基本格式。
3. 使用原始日志唯一键幂等写入 `raw_logs`。
4. Parser 根据协议和版本解析。
5. 轻量校验标准账户事实的公共字段及对应 payload 的必填字段和数值。
6. 在一个事务中写入事件版本、更新当前版本并写 outbox。
7. Outbox Publisher 发布事件。
8. 更新发布状态和延迟指标。

### 5.3 数据表

#### `collection_jobs`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | UUID | 任务 ID |
| `chain_id` | TEXT | 链标识 |
| `source_id` | TEXT | 数据源 |
| `mode` | TEXT | `HISTORICAL/REALTIME/REPAIR` |
| `range_start/end` | JSONB | 链适配器定义的位置范围，可空 |
| `cursor` | TEXT | 第三方游标，可空 |
| `status` | TEXT | 任务状态 |
| `attempts` | INT | 尝试次数 |
| `heartbeat_at` | TIMESTAMPTZ | Worker 心跳 |
| `last_error` | TEXT | 最近错误摘要 |

#### `collection_checkpoints`

主键为 `(chain_id, source_id, partition_key)`，保存适配器定义的位置、稳定性证明、供应商游标和更新时间。

#### `raw_logs`

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | UUID | 内部 ID |
| `source_id` | TEXT | 数据来源 |
| `source_event_id` | TEXT | 来源事件 ID |
| `chain_id` | TEXT | 链标识 |
| `chain_position` | JSONB | 链特定定位信息 |
| `transaction_ref` | TEXT | 交易或成交引用，可空 |
| `ordering_key` | TEXT | 链内稳定排序键 |
| `payload` | JSONB | 原始内容 |
| `observed_at` | TIMESTAMPTZ | 接收时间 |
| `parse_status` | TEXT | 解析状态 |

唯一约束使用 `(chain_id, source_id, source_event_id)` 或链适配器生成的稳定原始日志 ID。不将 EVM 交易哈希和日志索引作为所有链的强制唯一键。

#### `account_fact_versions`

主键 `(fact_id, revision)`，保存事实类型、账户键、排序位置、完整负载、解析器版本、原始日志 ID 和创建时间。

#### `account_facts_current`

主键 `fact_id`，保存当前 revision、change_type、fact_type、account_key、发生时间、排序位置和确认状态。索引至少包括：

- `(account_key, ordering_key, sub_index)`。
- `(account_key, occurred_at, fact_id)`。
- 对 `fact_type=TRADE` 的市场和操作 ID 建立部分索引。
- 对具备交易哈希的链建立部分索引 `(transaction_hash)`。

#### `outbox_events`

保存 event_id、topic、partition_key、payload、状态、尝试次数和下次重试时间。所有账户事实发布到 `account.fact.v1`，并使用 `account_key` 作为分区键，保证同一协议账户的事实在同一分区内有序。跟单服务只处理其中 `fact_type=TRADE` 的记录。

### 5.4 最终性与链状态变更

- 实时观察到日志后发布 `OBSERVED` revision。
- 最终性 Worker 根据链适配器的稳定性规则发布新的 `CONFIRMED` revision。
- 链上事实因重组、共识回滚或源侧修正而失效时，发布 `RETRACT/ORPHANED` revision。
- 出现替代交易事实时，以新的稳定 `fact_id` 发布。
- 账户分析必须能够回滚旧 revision；跟单服务对确认状态的接受策略可配置。

### 5.5 错误处理

| 错误 | 处理 |
| --- | --- |
| 数据源超时 | 指数退避重试，超过阈值切换备用源并告警 |
| 游标失效 | 从最近持久化区块或时间窗口回拉 |
| 原始日志重复 | 唯一约束忽略，不重复解析发布 |
| 解析失败 | 进入死信表，保留原始数据和解析器版本 |
| 发布失败 | outbox 重试，不回滚已保存业务数据 |
| 未知协议版本 | 隔离并告警，不生成猜测事件 |

## 6. 账户数据分析服务详细设计

### 6.1 账户键与顺序

- `account_key = chain_id + protocol + normalized_account`。同一地址在不同协议中分别建账，避免保证金、余额和持仓语义互相污染。
- 地址按链规则规范化，但必须保留原始地址。
- 排序键为 `(ordering_key, sub_index, revision)`，其中 `ordering_key` 由链适配器保证在账户范围内可稳定比较。
- 同一账户的事件由消息分区保证正常情况下按序消费。
- 迟到事件或旧位置修正不直接套用当前状态，而是创建重算任务。

### 6.2 处理组件

```text
Event Consumer → Account Ledger → Calculation Scheduler
                                           ↓
                              Calculator Registry + Dependency Graph
                                           ↓
                         Versioned Calculation Results → Query API
```

账户分析服务只定义计算框架、输入输出契约和执行生命周期，不把持仓、盈亏、指标或策略规则固化在调度与存储逻辑中。

### 6.3 功能分层

账户分析按稳定性由低到高分为五层：

| 层级 | 职责 | 变更频率 |
| --- | --- | --- |
| 事实层 | 保存标准账户事实、市场数据及版本 | 低 |
| 账本与状态层 | 归并交易过程，还原余额、持仓、成本和账户净值 | 中 |
| 基础指标层 | 根据账户状态和完整持仓周期产生账户价值、净盈亏、成交量、交易次数、胜率及数据质量特征 | 中高 |
| 评分层 | 将多维指标转换为分项评分、综合评分和解释 | 高 |
| 选择层 | 根据评分、可跟随性和组合约束选择跟单地址 | 高 |

上层只依赖下层的版本化输出，不回写或修改下层事实。评分和选择规则可以高频调整，而无需改动采集、账本和存储逻辑。

### 6.4 账户账本

消费流程：

1. 根据 `(event_id)` 检查消费幂等表。
2. 比较 `fact_id` 的当前 revision。
3. 新事件按顺序写入账本。
4. 如果事件位置晚于账户计算水位，执行增量计算。
5. 如果事件是旧位置的新增、修正或撤销，创建从该位置开始的重算任务。
6. 同事务写入消费记录和任务/增量状态，提交后确认消息。

### 6.5 可扩展计算框架

#### 计算器契约

所有账户计算规则以 Rust trait 实现：

```rust
pub trait AccountCalculator: Send + Sync {
    fn descriptor(&self) -> CalculatorDescriptor;
    fn validate_input(&self, ctx: &CalculationContext) -> ValidationResult;
    fn calculate(&self, ctx: &CalculationContext) -> Result<CalculationOutput>;
}
```

`CalculatorDescriptor` 至少包含：

- `calculator_id`：全局稳定的计算器标识。
- `version`：规则版本，计算语义变更时必须升级。
- `input_types`：依赖的事件、市场数据和其他计算结果。
- `output_type`：输出类型及 schema 版本。
- `supported_scopes`：支持的链、协议、资产或交易类型范围。
- `execution_mode`：支持增量、全量或两者。
- `dependencies`：前置计算器及版本约束。

#### 计算上下文

`CalculationContext` 是只读快照，包含：

- 账户键、计算范围和数据水位。默认计算范围为截至该水位的完整历史，不要求时间窗口。
- 指定版本的当前有效账户事实。
- 需要的市场价格、资产元数据和汇率数据。
- 已完成的前置计算结果。
- 时区、计价币、精度和可配置参数。
- 数据质量状态和缺失数据说明。

计算器不能直接读取可变的外部状态，也不能直接写数据库，保证相同上下文和版本产生相同结果。

#### 计算器类型

框架预留以下类型，但不在本文档固定具体公式：

- **账本计算器**：将原子交易事件关联为完整交易过程。
- **状态计算器**：计算余额、持仓或其他账户状态。
- **盈亏计算器**：使用可配置成本法、结算规则和市场数据生成盈亏结果。
- **指标计算器**：从事件、状态和盈亏结果生成统计指标。
- **策略分析器**：根据其他计算结果生成标签、评分或解释。

#### 注册与选择

- 计算器由注册表管理，不通过硬编码分支选择。
- 调度器根据链、协议、交易类型、计算配置和版本选择计算器。
- 依赖关系组成有向无环图；注册时检测循环依赖和缺失依赖。
- 同一类计算可以并存多个版本或多种口径，查询时明确指定或使用已发布的默认组合。
- 链或协议特有的结算规则通过计算器实现，不在数据源适配器中隐式计算。

#### 扩展方式

- **参数扩展**：只变更阈值、计价币、统计范围或指标组合时，通过版本化配置发布，不修改代码。
- **规则扩展**：新的计算逻辑实现为独立 Rust crate，通过工厂注册到计算器注册表。
- **组合扩展**：已有计算器可以通过依赖图组成新的分析流程，不复制底层计算逻辑。
- Rust 稳定 ABI 不作为动态库插件契约；计算器默认随服务编译和发布。如需运行时加载不可信计算规则，应另行定义 WASM 沙箱契约，不直接加载本地动态库。

#### 输出契约

`CalculationOutput` 至少包含：

- 计算器 ID 和版本。
- 账户、计算范围和数据水位。
- 输出 schema 版本和结果负载。
- 依赖的输入版本及数据哈希。
- 数据质量、警告、样本数和不完整原因。
- 计算时间和可选的可解释证据。

### 6.6 默认计算方案

系统提供一套常用的默认计算器，使账户分析在没有自定义规则时仍可运行。这些计算器与自定义计算器遵守同一契约，不享有特殊的系统级硬编码。

#### 默认计算链

```text
交易与资金事实
        ↓
交易类型路由
    ┌───┴────────┐
合约持仓周期   现货资产批次
    └───┬────────┘
        ↓
账户价值与分类盈亏
        ↓
账户价值、净盈亏、成交量、分类交易次数和分类胜率
```

默认计算链只产出基础账户分析数据。收益率、时间窗口、风险特征、评分、可跟随性和自动选择地址均通过后续计算器或计算方案扩展，不是默认链的前置依赖。

#### 默认基础计算器

| 计算器 | 职责 | 主要输入 | 主要输出 |
| --- | --- | --- | --- |
| `derivative_episode` | 将合约成交归并为完整持仓周期 | 合约交易事件 | 已结束和进行中的持仓周期、成交集合 |
| `spot_inventory` | 维护现货资产数量和移动加权平均成本 | 现货交易与资金事实 | 资产批次、平均成本和剩余库存 |
| `spot_disposal` | 将卖出成交归并为逻辑卖出操作 | 现货成交、订单或协议操作标识 | 卖出操作及已实现盈亏 |
| `account_state` | 按协议语义还原账户状态 | 交易过程、资金事件 | 余额、持仓、成本、数据水位 |
| `capital_flow` | 区分外部资金流和投资收益 | 充值、提现、转账等资金事件 | 标准化净资金流 |
| `account_value` | 计算数据水位时点的账户价值 | 余额、当前持仓、市场价格 | 账户价值、未实现盈亏和估值依据 |
| `pnl` | 计算分类盈亏并汇总 | 合约持仓周期、现货卖出操作及相关费用 | 分类净盈亏和账户累计净盈亏 |
| `performance_features` | 汇总基础交易表现 | 合约持仓周期、现货卖出操作、成交和盈亏 | 分类成交量、交易次数、胜负次数和胜率 |

#### 预留扩展计算器

| 计算器 | 职责 | 主要输入 | 主要输出 |
| --- | --- | --- | --- |
| `return_series` | 排除外部资金流影响并生成收益序列 | 净值序列、资金流 | 分期收益和累计收益 |
| `risk_features` | 生成风险暴露和稳定性特征 | 净值、持仓、收益序列 | 回撤、波动、集中度等 |
| `behavior_features` | 分析账户交易行为 | 交易过程和账户状态 | 频率、持仓周期、方向和市场偏好 |
| `credibility_features` | 衡量分析结果的可信程度 | 样本、数据覆盖、完整度 | 样本充足度和数据可信度 |
| `followability` | 估算地址的交易是否容易复制 | 交易频率、持仓时间、市场流动性和时延 | 可跟随性特征和限制原因 |
| `dimension_score` | 按配置将基础特征转换为分项评分 | 多维特征、阈值和映射函数 | 分项分数和解释 |
| `composite_score` | 按版本化方案组合分项评分 | 分项分数、权重、资格条件 | 综合分数、资格结果和原因 |
| `address_selector` | 从候选地址中选择跟单对象 | 评分、可跟随性、相关性和组合约束 | 候选地址集合及选择依据 |

预留扩展计算器不属于默认启用范围，具体公式由后续计算器版本和计算方案决定。

#### 默认基础指标口径

- **分类计算原则**：标准事实和基础输出契约保持统一，但根据市场类型路由到不同计算器。合约使用持仓周期，现货使用资产库存和卖出操作；协议特殊语义由对应计算器实现。
- **合约完整持仓周期**：以 `(account_key, market)` 为归并范围。净持仓从零变为非零时开始，重新回到零时结束；加仓和减仓属于同一周期；持仓穿越零并反向时，先结束原周期，再以剩余反向数量开始新周期。
- **现货资产库存**：按 `(account_key, asset)` 使用移动加权平均成本维护数量和成本。买入增加数量和成本；卖出按卖出数量分摊成本并产生已实现盈亏；剩余资产按当前市场价格产生未实现盈亏。买卖手续费计入对应成本或处置结果。
- **现货交易对语义**：交易对中的基础资产是盈亏分析对象，计价资产只是本次买卖的支付或收入资产。例如买入 `ETH/USDC` 只更新 ETH 库存，不再把支付的 USDC 统计为一次卖出；稳定币兑换、多跳交换及缺少明确交易对的协议操作由协议适配器拆解并声明基础资产和计价资产，避免成交量、交易次数及盈亏重复计算。
- **现货资金变动**：充值和转入只改变资产库存，不增加交易次数；提现和转出只减少库存，不产生交易盈亏。转入资产没有可靠历史成本时，以转入时可靠市场价格建立估算成本并标记 `ESTIMATED`；后续补齐来源成本时触发重算。
- **现货卖出操作**：每个产生已实现盈亏的逻辑卖出操作计为一次现货交易。同一订单的多个 fill 必须合并；缺少订单标识时使用链交易、协议操作标识及子序号归并，不得直接把每个 fill 当成一次交易。
- **观察起点已有持仓或资产**：合约创建 `ESTIMATED_OPEN` 持仓周期，现货创建估算库存批次；均以首次可靠账户快照对应的市场价格作为估算成本，并保存基准价格、价格时间、来源和估算原因。它们可以计算观察期盈亏并参与汇总，但不是历史真实成本；结果必须同时返回估算交易或批次数量，且明细和汇总均标记 `ESTIMATED`。补齐更早历史后应触发重算并替换估算成本。
- **账户价值**：数据水位对应时点的账户净值，是时点状态，不按持仓周期统计。统一输出 `account_value`、`available_balance`、`margin_used`、`unrealized_pnl`、`valuation_currency` 和估值依据。协议直接提供且通过校验的账户价值时优先采用，否则由协议计算器依据余额、持仓和市场价格计算；输出必须携带价格来源、价格时间和数据水位。
- **净盈亏**：合约按已结束持仓周期、现货按卖出操作汇总已实现交易盈亏。合约结果扣除归属于周期的交易手续费和清算费用，并加上周期内资金费净收入；现货结果计入对应买卖手续费。进行中持仓或剩余库存的未实现盈亏单独输出并参与当前账户价值，但不计入已实现净盈亏。无持仓期间发生的资金费或其他资金变化作为独立账户资金事件。
- **唯一盈亏结果**：每个合约持仓周期或现货卖出操作只产生一个正式 `realized_net_pnl`。协议计算器依次选择可靠的协议结算值、基于完整成交的重建值或估算值，并输出 `calculation_method=PROTOCOL_REPORTED/RECONSTRUCTED/ESTIMATED`。协议报告字段只是输入，不形成第二套对外盈亏；手续费和资金费是否已包含由协议计算策略明确声明，通用计算器不得猜测。
- **成交量**：全部有效成交的绝对名义价值之和，开仓、加仓、减仓和平仓均计入。保留原始数值和计价币；账户汇总默认换算为 USD，并保存换算价格、时间和来源。无法可靠换算的成交不进入汇总值，同时将结果标为不完整。
- **交易次数**：合约 `trade_count` 为已结束的完整持仓周期数；现货 `trade_count` 为已归并且产生已实现盈亏的逻辑卖出操作数；两者均另保留 `fill_count`。
- **胜率**：各分类分别使用 `winning_trade_count / (winning_trade_count + losing_trade_count)`。净盈亏大于零为盈利，小于零为亏损，等于零为持平；持平交易计入 `trade_count` 和 `breakeven_trade_count`，但不进入胜率分母。没有盈利或亏损交易时胜率输出 `null`。
- **聚合边界**：账户价值、净盈亏、成交量和成交笔数可在完成统一计价后形成账户总计；合约与现货的交易次数及胜率语义不同，必须分别输出为 `derivative_metrics` 和 `spot_metrics`，不得计算混合总胜率或混合总交易次数。
- **盈亏输出**：统一返回 `realized_net_pnl`、`unrealized_pnl`、`total_net_pnl` 和 `estimated_pnl`，其中 `total_net_pnl = realized_net_pnl + unrealized_pnl`；`estimated_pnl` 是总盈亏中依赖估算成本的部分，只作质量披露，不得再次相加。
- **历史覆盖**：结果返回 `tracking_started_at`、`coverage_start`、`calculated_to` 和 `history_complete`。文档中的“完整历史”指系统在当前数据水位能够验证的全部历史，不暗示已经获得地址诞生以来的所有交易。
- **跨协议聚合**：分析与账本以 `account_key` 分开计算；需要查看同一自然地址的全局状态时，通过地址级聚合投影组合各协议结果，并保留每个来源结果和质量状态，不跨协议归并持仓周期或交易次数。
- **精度与 dust**：金额、价格、数量和盈亏全程使用十进制定点数，精度由协议及资产元数据确定。dust 阈值只影响展示和状态标签，不删除账本余额、不自动视为成交，也不用于伪造持仓归零。
- **清算归并**：清算造成的减仓或平仓作为 `trigger_type=LIQUIDATION` 的交易事实参与持仓周期；清算费用作为独立费用事实归属于该周期；清算后的余额、仓位和账户价值由状态快照校验。不得仅用一个资金事件代替清算成交。
- 默认指标覆盖账户纳入系统后可验证的完整历史，并返回事实覆盖起点、数据水位和完整度；不提供默认时间窗口。时间窗口属于未来可增加的计算范围扩展，不改变完整持仓周期定义。

#### 计算方案（Calculation Profile）

计算方案是一个不可变、版本化的计算器组合，定义：

- 启用的计算器及版本。
- 计算器依赖和执行顺序。
- 计价币、资产精度、成本法和数据源。
- 可选的统计范围；默认使用可验证的完整历史，未来可增加时间窗口。
- 指标参数、资格条件、分数映射、维度权重和缺失值政策。
- 是否启用地址选择以及选择约束。

参考 [8 维权重风险评分算法](8维权重.md) 中的维度、阈值和权重可以实现为一个独立计算方案，不作为系统内建且不可变更的唯一规则。

### 6.7 分析数据表

#### `account_ledger_events`

保存进入分析服务的账户事实版本：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `fact_id/revision` | TEXT/INT | 联合主键 |
| `fact_type/account_key` | TEXT | 事实类型和协议账户 |
| `ordering_key/sub_index` | TEXT/INT | 确定性顺序 |
| `occurred_at` | TIMESTAMPTZ | 事实发生时间 |
| `change_type/is_current` | TEXT/BOOL | 修订状态 |
| `payload` | JSONB | 版本化事实负载 |
| `consumed_event_id` | TEXT | 消费幂等与追踪 |

唯一约束 `(fact_id, revision)`；当前事实查询索引为 `(account_key, is_current, ordering_key, sub_index)`。

#### `positions_current`

可选的持仓查询投影。主键 `(account_key, market, calculator_id, calculator_version)`，字段由已启用的状态计算器输出 schema 定义。

#### `derivative_episodes`

保存 `episode_id`、账户、市场、方向、状态、开平时间、成交数量、成本、已实现盈亏、手续费、资金费、净盈亏、计算方法、估算标记、计算版本和数据水位。索引 `(account_key, market, status, opened_at)`。

#### `spot_inventories`

主键 `(account_key, asset, calculator_version)`，保存数量、总成本、移动平均成本、估算数量和数据水位。资产转出时按转出数量同比例移出库存成本，不产生交易盈亏。

#### `spot_disposals`

保存 `disposal_id`、账户、市场、order/operation ID、卖出数量、收入、分摊成本、手续费、净盈亏、计算方法、估算标记、计算版本和数据水位。逻辑操作 ID 在同一账户和版本内唯一。

#### `account_snapshots`

按账户、计算器、版本、快照时间和数据水位保存状态负载。快照频率由计算器或运行配置决定。

#### `account_metrics`

保存账户当前基础分析投影：账户价值、已实现/未实现/总净盈亏、估算盈亏、USD 成交量、成交笔数、`derivative_metrics`、`spot_metrics`、估算交易数、事实覆盖范围、数据水位、质量状态和计算版本。主键 `(account_key, calculation_version, data_watermark)`，另维护当前结果指针。

#### `strategy_profiles`

保存策略分析器输出的标签、评分或解释，并保留分析器版本及所有依赖结果版本。

#### `calculator_registry`

保存计算器 ID、版本、支持范围、输入输出 schema、依赖、状态和发布时间。

#### `calculation_results`

保存通用的不可变计算结果。特定查询表只是该结果的物化投影，可以从通用结果重建。

#### `calculation_profiles`

保存计算方案 ID、版本、计算器组合、统计范围、参数、权重、缺失值政策、支持范围和状态。已发布方案不可原地修改。

#### `recalculation_jobs`

保存账户、最早受影响位置、目标计算版本、状态、进度和错误。对同一账户的多个待处理请求合并到最早影响位置。

### 6.8 计算规则发布与兼容

- 新计算器版本先处于 `DRAFT`，通过固定数据集回放后发布为 `ACTIVE`。
- 已产生结果的版本不可原地修改；语义改变必须发布新版本。
- 可在同一数据水位上并行运行新旧版本并生成差异报告。
- 默认版本切换只改变新查询所指向的结果组合，不删除旧结果。
- 指标比较和排名必须使用相同计算器及版本、计价币、事实覆盖范围和数据质量标准。
- 计算器因输入不足无法计算时，应输出带原因的 `UNAVAILABLE` 结果，不使用零值或猜测值代替。

### 6.9 重算算法

1. 锁定账户的重算租约，确保同一账户同时只有一个重算 Worker。
2. 找到受影响位置之前最近的有效快照。
3. 从快照恢复状态；无快照时从账户首个事件开始。
4. 按确定顺序读取当前有效 revision 并重放。
5. 写入新的临时计算版本。
6. 校验水位、事件数量和关键余额不变量。
7. 原子切换账户当前计算版本。
8. 按计算依赖图重建所有受影响的下游结果。

查询 API 在重算期间继续返回旧版本，并标记 `recalculation_status=RUNNING`，避免暴露半成品。

## 7. 数据完整审核服务详细设计

### 7.1 审核层次

| 层次 | 检查内容 | 建议频率 |
| --- | --- | --- |
| 采集健康 | 心跳、游标推进、区块延迟、错误率 | 分钟级 |
| 范围完整 | 区块连续性、分页完整、事件数量异常 | 小时级 |
| 来源对账 | 第三方与链节点/备用源差异 | 抽样持续、全量按日 |
| 链最终性 | 链适配器定义的稳定性证明、回滚或重组 | 随链状态更新持续 |
| 业务一致 | 无持仓减仓、数量异常、无法闭合序列 | 账户计算时及每日 |

频率为默认建议，必须按链出块速度、第三方 SLA 和业务时效配置。

### 7.2 审核流程

```text
创建审核任务 → 读取检查点和数据范围 → 对账 → 生成差异
            → 分类异常 → 自动/人工决策 → 提交修复请求
            → 交易日志服务修复并发布新 revision → 复验关闭
```

### 7.3 异常分类

| 类型 | 示例 | 默认处理 |
| --- | --- | --- |
| `SOURCE_STALLED` | 游标长时间未推进 | 告警并切换/重启采集 |
| `MISSING_RANGE` | 区块区间缺失 | 自动创建范围回补 |
| `MISSING_EVENT` | 对账发现单个事件缺失 | 自动拉取交易回执并重解析 |
| `DUPLICATE_EVENT` | 同一事实多个记录 | 自动修复当前版本，保留审计记录 |
| `CONTENT_MISMATCH` | 第三方解析与链上不一致 | 以配置的权威源重解析 |
| `CHAIN_REORG` | 原区块不在主链 | 发布撤销 revision 并回拉新区块 |
| `PARSER_ERROR` | 无法识别协议版本 | 隔离并人工处理 |
| `BUSINESS_INCONSISTENCY` | 无前态却出现减仓 | 扩大范围回查，无法修复则标记账户数据降级 |

### 7.4 数据表

#### `audit_jobs`

保存审核类型、数据源、链、范围、调度原因、状态、进度和统计结果。

#### `audit_findings`

保存异常类型、证据、影响范围、严重级别、首次/最近发现时间和状态。

#### `repair_requests`

保存请求范围、预期操作、幂等键、状态、交易日志服务返回结果和复验状态。

### 7.5 数据质量状态

按 `chain + source + block range` 和 `account + time range` 维护：

| 状态 | 含义 |
| --- | --- |
| `UNKNOWN` | 尚未审核 |
| `CHECKING` | 正在审核 |
| `HEALTHY` | 当前审核规则下无异常 |
| `DEGRADED` | 存在问题但仍可部分使用 |
| `INVALID` | 数据不可用于可靠分析 |
| `REPAIRING` | 正在修复 |

账户分析查询返回对应范围的数据质量状态，但跟单主流程不等待该状态变为 `HEALTHY`。

## 8. 跟单服务详细设计

### 8.1 流程与状态机

```text
RECEIVED → RULE_REJECTED
         → APPROVED → SIGNING → SIGNED → SUBMITTING → SUBMITTED → CONFIRMED
                              ↘ FAILED              ↘ FAILED
```

终态包括 `RULE_REJECTED`、`CONFIRMED`、`FAILED` 和 `EXPIRED`。超时不等于失败：在重试签名或发送前必须先通过幂等键查询钱包或链上状态。

### 8.2 跟单配置

#### `copy_configs`

| 字段 | 说明 |
| --- | --- |
| `id/user_id` | 配置及用户 |
| `source_account` | 被跟随账户 |
| `wallet_ref` | 执行钱包引用 |
| `enabled` | 是否启用 |
| `markets/sides/actions` | 允许范围 |
| `sizing_mode` | 固定数量、固定金额或资金比例 |
| `sizing_value` | 仓位参数 |
| `max_event_age_ms` | 允许的最大事件延迟 |
| `accept_confirmation_status` | 可接受的链确认状态 |
| `max_single_notional` | 单笔上限 |
| `max_total_exposure` | 总敞口上限 |
| `max_slippage_bps` | 最大滑点 |
| `status/version` | 配置状态和版本 |

配置变更采用版本号。每次决策保存实际使用的配置版本，保证可追溯。

### 8.3 规则处理顺序

1. 事件幂等检查。
2. 事件年龄和确认状态检查。
3. 匹配启用的跟随配置。
4. 检查市场、动作和方向过滤条件。
5. 计算目标交易数量。
6. 校验单笔和总敞口限制。
7. 获取必要报价并校验滑点。
8. 生成交易意图 `trade_intent`。

同一标准交易事实可以命中多个用户配置，每个配置产生独立执行项。执行幂等键为 `fact_id + revision + copy_config_id`。

### 8.4 钱包接口

统一接口：

```text
capabilities(wallet_ref) -> { sign, submit, signAndSubmit, query }
sign(wallet_ref, unsigned_tx, idempotency_key) -> signed_tx
submit(signed_tx, idempotency_key) -> submission_result
signAndSubmit(wallet_ref, unsigned_tx, idempotency_key) -> submission_result
query(idempotency_key | transaction_ref) -> transaction_status
```

适配方式：

- MCP 钱包：通过 MCP 工具请求签名或签名并发送。
- 本地钱包：本地密钥服务签名，跟单服务通过 RPC 发送。
- 第三方托管钱包：调用供应商 API，优先使用其幂等请求能力。

钱包适配器不得向上层返回私钥。日志仅记录钱包引用、待签名内容哈希、请求 ID 和交易哈希，不记录私钥、助记词或完整敏感凭证。

### 8.5 跟单数据表

#### `wallet_refs`

保存钱包引用、类型、链、地址、外部资源 ID、状态和非敏感配置。敏感凭证存放在密钥管理系统，只保存引用。

#### `copy_executions`

保存执行幂等键、原事件 ID/revision、配置及版本、决策、拒绝原因、目标参数、钱包引用和当前状态。

#### `wallet_requests`

保存签名/发送请求 ID、内容哈希、幂等键、尝试次数、返回交易哈希和错误摘要。

#### `transaction_receipts`

保存链上交易哈希、提交时间、确认状态、区块位置、实际成交数量、价格、费用和失败原因。

### 8.6 重试策略

- 规则拒绝不重试。
- 钱包明确拒绝签名不重试。
- 网络超时先查询幂等键或交易哈希，确认未提交后才能重试。
- 链上交易序号冲突、报价过期或手续费变化需要重新构建交易，再次执行规则有效期检查。
- 超过 `max_event_age_ms` 后不再新建交易。
- 数据审核产生的历史修正默认不触发补单；如需补偿，必须由独立配置显式允许。

## 9. API 设计

### 9.1 通用约定

- 路径版本：`/api/v1`。
- 时间使用 UTC ISO 8601。
- 游标分页，不使用大偏移量分页。
- 写接口接受 `Idempotency-Key`。
- 错误响应包含 `code`、`message`、`trace_id` 和可选 `details`。

### 9.2 交易日志接口

| 方法与路径 | 用途 |
| --- | --- |
| `POST /api/v1/watch-accounts` | 添加关注账户和起始范围 |
| `DELETE /api/v1/watch-accounts/{id}` | 停止关注账户 |
| `GET /api/v1/trade-events` | 按账户、链、市场、时间或区块查询事件 |
| `GET /api/v1/trade-events/{fact_id}` | 查询交易事实当前版本和版本历史 |
| `GET /api/v1/collection-status` | 查询采集水位、延迟和数据源状态 |
| `POST /internal/v1/repair-requests` | 接收审核服务修复请求 |

### 9.3 账户分析接口

| 方法与路径 | 用途 |
| --- | --- |
| `GET /api/v1/accounts/{account_key}/overview` | 账户概览 |
| `GET /api/v1/accounts/{account_key}/positions` | 当前或历史持仓 |
| `GET /api/v1/accounts/{account_key}/trades` | 已结束/进行中合约持仓周期和现货卖出操作 |
| `GET /api/v1/accounts/{account_key}/pnl` | 盈亏和净值曲线 |
| `GET /api/v1/accounts/{account_key}/metrics` | 完整历史基础指标及其事实覆盖范围 |
| `GET /api/v1/accounts/{account_key}/strategy` | 策略标签和评价依据 |
| `GET /api/v1/accounts` | 多维筛选和排序 |
| `GET /api/v1/accounts/{account_key}/analysis-results` | 按计算器、方案、版本和统计范围查询通用结果 |
| `GET /api/v1/calculation-profiles` | 查询可用计算方案及版本 |
| `POST /internal/v1/recalculations` | 创建重算任务 |
| `GET /internal/v1/recalculations/{id}` | 查询重算状态 |

账户响应统一包含：

```json
{
  "data": {},
  "calculation_version": "pnl-v1/metrics-v1",
  "data_watermark": "eip155:42161:123456",
  "tracking_started_at": "2026-09-01T00:00:00Z",
  "coverage_start": "2026-09-01T00:00:00Z",
  "calculated_to": "2026-09-04T01:02:03Z",
  "history_complete": false,
  "data_quality": "HEALTHY",
  "updated_at": "2026-09-04T01:02:03Z",
  "recalculation_status": "IDLE"
}
```

`overview` 的 `data` 至少包含 `account_value`、`realized_net_pnl`、`unrealized_pnl`、`total_net_pnl`、`estimated_pnl`、`volume_usd`、`fill_count`、`estimated_trade_count`，以及分别统计交易次数、胜负次数、持平次数和胜率的 `derivative_metrics`、`spot_metrics`。

`positions` 分别返回当前合约持仓和现货库存，并携带成本、估算标记、估值依据和数据水位。`trades` 使用游标分页，返回合约持仓周期或现货卖出操作、正式净盈亏、`calculation_method`、计算版本和估算标记。

### 9.4 审核接口

| 方法与路径 | 用途 |
| --- | --- |
| `POST /api/v1/audits` | 创建专项审核 |
| `GET /api/v1/audits/{id}` | 查询审核进度和结果 |
| `GET /api/v1/findings` | 查询未关闭异常 |
| `POST /api/v1/findings/{id}/repair` | 人工确认修复 |
| `GET /api/v1/data-quality` | 查询范围或账户的数据质量状态 |

### 9.5 跟单接口

| 方法与路径 | 用途 |
| --- | --- |
| `POST /api/v1/copy-configs` | 创建跟单配置 |
| `PUT /api/v1/copy-configs/{id}` | 更新配置并生成新版本 |
| `POST /api/v1/copy-configs/{id}/enable` | 启用配置 |
| `POST /api/v1/copy-configs/{id}/disable` | 停用配置 |
| `GET /api/v1/copy-executions` | 查询规则判断和执行结果 |
| `GET /api/v1/copy-executions/{id}` | 查询单次执行详情 |

## 10. 消息主题与消费规则

| Topic | 分区键 | 保留建议 | 消费方 |
| --- | --- | --- | --- |
| `account.fact.v1` | `account_key` | 不少于最长重放周期 | 分析、审核；跟单只处理 `TRADE` |
| `trade.repair.requested.v1` | `chain_id` | 30 天以上 | 交易日志获取服务 |
| `account.recalculation.requested.v1` | `account_key` | 30 天以上 | 账户分析服务 |

消费要求：

- 每个消费组维护独立 offset。
- 业务写入与消费幂等记录在同一事务提交后再确认消息。
- 连续失败消息进入死信主题，保留原消息、错误码和尝试次数。
- 不允许通过跳过消息的方式永久绕过失败；必须重放或人工关闭。

## 11. 并发、一致性与幂等

### 11.1 交易日志

- 原始日志唯一键防止重复采集。
- 标准事件版本使用乐观锁：新 revision 必须等于当前 revision + 1。
- outbox 保证数据库写入和事件发布最终一致。

### 11.2 账户分析

- 单账户通过消息分区串行处理。
- 重算使用账户级租约锁。
- 增量任务发现账户处于重算状态时，仅写入账本并等待重算追上水位。

### 11.3 跟单

- 一个交易事实版本命中一个配置最多产生一个执行项。
- 签名和发送均携带稳定幂等键。
- 未知发送结果必须先查询，禁止直接重复广播或重新签名。

## 12. 安全设计

### 12.1 身份与权限

- 外部 API 使用用户身份认证；内部服务使用双向 TLS 或短期服务凭证。
- 查询、配置、修复、重算和钱包操作使用不同权限。
- 审核服务有权发起修复请求，但无权直接写交易日志表。
- 只有跟单服务的钱包适配进程具有签名权限。

### 12.2 密钥与敏感数据

- 本地私钥使用 KMS/HSM 或专用密钥存储加密，应用数据库只保存密钥引用。
- 私钥和助记词不得出现在 API 响应、消息、日志、错误堆栈或备份明文中。
- 第三方钱包凭证保存在密钥管理系统并定期轮换。
- 签名请求限定链、钱包、目标合约和金额范围。

### 12.3 审计日志

以下操作写入不可篡改审计日志：跟单配置变更、钱包绑定、签名请求、交易发送、人工修复、规则版本发布和重算版本切换。

## 13. 可观测性

### 13.1 核心指标

#### 交易日志获取

- 数据源最新区块与采集水位差。
- 原始日志接收率、解析成功率和死信数量。
- `observed_at - occurred_at` 获取延迟。
- `published_at - observed_at` 处理延迟。
- outbox 未发布数量及最老消息年龄。

#### 账户分析

- 消费积压、单事件计算耗时。
- 账户计算水位落后量。
- 重算任务数量、耗时和失败率。
- 缺失估值价格账户数量。

#### 数据审核

- 未审核区块范围、最近成功审核时间。
- 各类异常数量、修复平均时长和复验失败率。
- 第三方与权威源差异率。

#### 跟单

- 从事件发布时间到规则决策、签名和提交的分段延迟。
- 规则通过率、过期率、钱包失败率和链上失败率。
- 状态未知的提交数量。

### 13.2 告警

- 采集水位停止推进或超过链级延迟阈值。
- outbox 堆积持续增长。
- 解析失败率或来源对账差异率超过阈值。
- 账户分析消费积压超过 SLA。
- 数据质量进入 `INVALID`。
- 钱包签名失败率、交易失败率或未知状态数量异常。

## 14. 部署设计

### 14.1 进程划分

- `trade-collector`：采集和检查点管理。
- `trade-parser-publisher`：解析、持久化和 outbox 发布。
- `account-analyzer`：增量分析。
- `account-recalculator`：账户重算。
- `data-auditor`：审核和修复编排。
- `copy-trader`：规则判断和钱包调用。
- `query-api`：可按服务分别部署，也可使用统一网关聚合。

进程可以来自同一代码仓库，但使用独立运行身份和数据库写权限。

### 14.2 扩展方式

- 采集按链和数据源分片。
- 解析按原始日志 ID 分片，但同一原始日志只能由一个 Worker 成功提交。
- 账户分析按 `account_key` 分区。
- 审核按链和区块范围分片。
- 跟单按 `source_account` 或用户配置分片，同时保持执行幂等。

### 14.3 发布与迁移

- 数据库迁移向后兼容，先部署可读取新旧结构的消费者，再切生产者。
- 事件 schema 只允许新增可选字段进行小版本兼容；删除、改名或语义变化必须升级主事件版本和 topic。
- 解析器、盈亏公式、指标和规则均带版本，支持灰度和回滚。

## 15. 测试设计

### 15.1 单元测试

- 各协议日志到标准账户事实的固定样本测试。
- 每个计算器的固定样本、边界输入、缺失数据和精度测试。
- 完整持仓周期的加仓、分批减仓、归零、反向开仓、持平交易和清算归并测试。
- 观察起点已有持仓的估算成本、估算标记、历史补齐后替换及重算测试。
- 费用与资金费归属、成交量计价换算和不可换算数据降级测试。
- 计算依赖图、完整历史统计范围、版本选择和缺失值政策测试；扩展时间窗口启用后再增加对应测试集。
- 跟单规则和仓位计算测试。

### 15.2 集成测试

- 采集断开后从检查点恢复且无数据遗漏。
- 数据库提交后发布进程失败，outbox 恢复后事件仍能送达。
- 重复和乱序事件不造成重复持仓变化。
- 修正/撤销 revision 触发正确区间重算。
- 钱包超时后查询状态，确保不重复发送。

### 15.3 回放测试

- 使用固定区块范围回放，重复执行结果必须一致。
- 新旧解析器对同一原始数据的差异必须形成报告。
- 新旧计算版本同时运行并比较账户状态、盈亏和指标差异。

### 15.4 故障演练

- 第三方数据源中断和返回重复数据。
- 链节点延迟、链重组和区块回滚。
- 消息重复、延迟和消费积压。
- 数据库主从切换。
- 钱包不可用、签名超时和已发送但响应丢失。

## 16. 验收标准

验收阈值按各链、协议和数据供应商的 SLA 分别配置，以下为所有接入必须验证的能力：

### 交易日志获取

- 服务重启后可从检查点恢复。
- 同一区块范围重复采集不会产生重复当前事件。
- 标准账户事实可追溯到原始日志和解析器版本。
- 实时发布延迟可分段观测。

### 账户分析

- 固定事件集多次计算得到完全一致结果。
- 增量计算与从头重算结果一致。
- 修正和撤销事件能够正确更新持仓、盈亏和指标。
- 每个分析结果带计算版本、数据水位和质量状态。

### 数据审核

- 能发现预先注入的漏块、漏事件、重复、内容差异和链重组。
- 审核不阻塞实时事件发布。
- 修复请求能够闭环到回补、事件新 revision、账户重算和复验。

### 跟单

- 标准交易事实可在时效限制内完成规则判断。
- MCP、本地和第三方钱包通过统一接口替换。
- 钱包或网络超时时不会产生重复交易。
- 私钥不会出现在数据库普通字段、消息和应用日志中。

## 17. 多链与多协议扩展规范

### 17.1 链能力描述

每条链通过版本化的能力描述声明：

- 链标识、区块或等价排序位置、时间精度和最终性模型。
- 历史数据查询、实时订阅、游标恢复和对账数据源能力。
- 交易、撤销、修正、费用、资金变动和强制平仓等可观测事件类型。
- 钱包签名、交易序号、交易发送、状态查询和幂等能力。

通用服务根据能力描述启用相应功能，不假设所有链均具有 EVM 区块、日志索引、nonce 或确认数模型。

### 17.2 协议插件契约

每个协议插件至少实现：

```text
identify(raw_data) -> bool
parse(raw_data, context) -> AccountFact[]
validate(event, context) -> ValidationResult
buildTrade(intent, market_context) -> UnsignedTransaction
```

只读分析插件可不实现 `buildTrade`。插件必须声明版本、支持的交易类型、计算策略和所需市场数据。

### 17.3 标准化边界

- 标准账户事实保留通用事实，特定链和协议字段放入带命名空间的 `extension`；其中 `TRADE` payload 是标准交易记录。
- 任何影响持仓、盈亏或交易执行的字段不能只存在未定义的扩展中，必须升级通用 schema 或声明专用的版本化计算策略。
- 账户、资产和市场标识必须包含链或协议命名空间，防止跨链冲突。
- 跨链聚合分析必须先完成计价币换算和时间对齐，并保留汇率与市场价格来源。

### 17.4 接入验收

新链或协议不得绕过通用验收套件，必须通过：

- 历史数据回放和确定性测试。
- 断流恢复、重复事件和乱序事件测试。
- 事件完整审核和数据源对账测试。
- 持仓、盈亏与费用的协议样本测试。
- 如支持跟单，必须通过签名、发送、幂等、超时和状态恢复测试。

## 18. Roadmap 边界

本文档定义系统的整体目标架构，不规定具体实现顺序、分期、优先接入的链与协议或阶段交付范围。

上述内容由独立的 Roadmap 文档管理。Roadmap 可以调整实施优先级，但不得在未更新设计文档的情况下改变本文定义的服务边界、数据所有权和事件语义。
