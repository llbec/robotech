# 版本开发记录

## V0.000

### 1. 版本信息

| 项目 | 内容 |
| --- | --- |
| 版本号 | V0.000 |
| 版本状态 | PLANNED |
| 记录日期 | 2026-09-11 |
| 开发状态 | 未开始 |
| 对应总体设计 | [详细设计](detailed-design.md) |

### 2. 版本目标

交付一个能够独立运行的 Nansen 数据能力验证程序，验证 Nansen API 是否能够提供完成 Hyperliquid 账户分析所需的全部输入数据，并完成以下闭环：

```text
读取运行参数
→ 按账户分析数据需求调用全部相关 Nansen API
→ 获取交易、持仓、账户快照、资金变化、费用和资金费等数据
→ 保存每个接口的完整原始响应和原始记录
→ 将可可靠映射的数据转换为标准账户事实
→ 保存标准账户事实和数据覆盖验证结果
→ 输出标准账户事实与覆盖报告
→ 正常退出
```

本版本的核心交付不是仅验证永续成交接口，而是回答以下问题：**仅依赖 Nansen API，是否能够取得完成 Hyperliquid 账户分析所需的全部数据。**

本版本必须以账户分析输入需求为基准逐项验证，不得因为 Nansen 暂无对应接口而缩小需求范围。能够获取的数据完成原始保存和标准化；无法获取、字段不足或历史范围不足的数据形成明确的覆盖缺口。覆盖缺口是本版本允许产出的验证结论，但不得用推测值补齐，也不得将“接口请求成功”等同于“账户分析数据完整”。

本版本仍不实现账户价值、盈亏、胜率等分析计算，也不实现完整系统持续运行能力。

### 3. 本版本范围

#### 3.1 包含范围

- 建立最小 Rust Cargo workspace 和交易日志服务目录。
- 提供一个单次执行的 Nansen 账户数据导入与覆盖验证程序。
- 接入所有与 Hyperliquid 账户分析输入有关的 Nansen API，而不是只接入永续成交接口。
- 获取指定地址在指定日期范围内可用的历史数据，并获取运行时账户与持仓快照。
- 验证并尽可能获取以下账户分析必需数据类别：
  - 永续成交：开仓、加仓、减仓、平仓和反向持仓。
  - 当前永续持仓：方向、数量、入场价格、标记价格、杠杆、保证金、清算价格和未实现盈亏。
  - 资金变化：充值、提现、地址间转账和内部账户划转。
  - 交易外费用：资金费、清算费用以及影响账户净值的其他费用。
  - 清算、奖励及其他会改变账户资产或持仓的事件。
  - 账户资产、保证金和账户价值快照，以及快照时间和计价信息。
- 保存 `Withdraw`、`SendAsset` 等与账户资金变化有关的记录；`ApproveAgent` 等账户管理操作至少保存为原始活动，用于证明活动覆盖范围，但不直接参与盈亏计算。
- 处理所有带分页的 Nansen 响应，直到 `is_last_page=true`。
- 保存每个接口的请求上下文、完整响应页和每条原始记录。
- 将能够可靠映射的数据转换为 `TRADE`、`TRANSFER`、`FEE`、`FUNDING`、`REWARD`、`LIQUIDATION_FEE` 或 `ACCOUNT_SNAPSHOT` 标准账户事实。
- 保存标准账户事实以及按数据类别、接口、时间范围统计的覆盖验证结果。
- 对完全相同的原始记录和标准事实进行幂等去重。
- 将本次获得的标准账户事实按发生时间顺序输出为 JSON Lines，并输出机器可读的覆盖报告。
- 提供基本配置校验、错误日志和退出码。
- 提供各数据类别的转换单元测试、接口 fixture 测试和一次完整覆盖验证流程测试。

#### 3.2 不包含范围

- 账户分析服务。
- 数据完整审核服务。
- 跟单服务。
- 账户价值、净盈亏、成交量、交易次数和胜率计算。
- Nansen 以外的数据源、链节点直连和链上区块采集；但 Nansen 的数据缺口必须记录，并作为后续版本选择补充数据源的依据。
- Nansen 实时订阅或持续轮询。
- 多数据源切换和来源对账。
- Kafka/Redpanda 消息发布。
- Outbox 发布流程。
- 完整的事实修订、撤销、最终性和链重组处理；如果 Nansen 返回相关状态，本版本仍须原样保存。
- 修复请求、历史回补任务和账户重算。
- HTTP 查询 API。
- Hyperliquid 以外的链和协议数据。
- 评分、地址选择和跟单信号执行。

### 4. 数据源

#### 4.1 Nansen 接口

| 项目 | 内容 |
| --- | --- |
| 产品入口 | `https://app.nansen.ai/api` |
| API 主机 | `https://api.nansen.ai` |
| 已确认接口 | `POST /api/v1/profiler/perp-trades`、`POST /api/v1/profiler/perp-positions` |
| 鉴权 | HTTP Header `apikey` |
| Content-Type | `application/json` |
| 目标数据类型 | 完成 Hyperliquid 账户分析所需的全部输入数据 |
| 调用方式 | 历史范围批量查询与运行时快照查询 |

官方参考：

- [Nansen API Authentication](https://docs.nansen.ai/getting-started/authentication)
- [Nansen Address Perp Trades](https://docs.nansen.ai/api/profiler/address-perp-trades)
- [Nansen Address Perp Positions](https://docs.nansen.ai/api/profiler/address-perp-positions)
- [Nansen Hyperliquid APIs](https://docs.nansen.ai/api/hyperliquid)

#### 4.2 接口发现与覆盖原则

V0.000 不预先假定通用 Nansen 地址接口支持 Hyperliquid Core。开发时必须依据 Nansen 官方接口文档、实际请求响应和固定 fixture 建立接口清单，并记录验证日期。

| 数据类别 | 账户分析用途 | V0.000 验证要求 |
| --- | --- | --- |
| 永续成交历史 | 还原持仓变化、成交量、交易次数和已实现盈亏 | 使用 `perp-trades` 验证完整分页、字段和时间覆盖 |
| 当前永续持仓 | 建立期末持仓与未实现盈亏快照 | 使用 `perp-positions` 验证仓位、保证金、价格和资金费字段 |
| 充值、提现、转账、内部划转 | 区分投资损益与本金流入流出 | 查找并验证 Nansen 官方支持的地址活动或资金流接口；不存在时记录 `UNAVAILABLE` |
| 独立资金费和费用 | 计算净盈亏和持仓成本 | 验证成交/持仓字段能否完整覆盖；不能逐笔还原时记录 `PARTIAL` |
| 清算和奖励 | 解释非普通成交导致的资产变化 | 验证是否存在事件类型、金额、资产、时间和交易引用 |
| 账户资产与价值快照 | 校验期末账户状态和账户价值 | 验证 Nansen 是否提供 Hyperliquid Core 账户级快照；不得用不支持该链的通用接口替代 |
| 管理活动 | 解释 `ApproveAgent` 等非资产活动 | 能获取则保存原始记录，不要求转换为账户分析事实 |

每个数据类别必须标记以下覆盖状态之一：

- `COMPLETE`：接口与字段足以支持对应账户分析计算，且时间/分页验证通过。
- `PARTIAL`：可以获取部分数据，但缺少关键字段、逐笔历史或必要语义。
- `UNAVAILABLE`：Nansen 没有可用接口，或当前订阅无权访问。
- `UNVERIFIED`：尚未完成真实请求或 fixture 验证。

只有所有账户分析必需数据类别均为 `COMPLETE`，才能得出“仅依赖 Nansen API 可以完成账户分析”的结论。

#### 4.3 请求参数

```json
{
  "address": "0x...",
  "date": {
    "from": "2026-09-01T00:00:00Z",
    "to": "2026-09-02T00:00:00Z"
  },
  "pagination": {
    "page": 1,
    "per_page": 100
  },
  "order_by": [
    {
      "field": "timestamp",
      "direction": "ASC"
    }
  ]
}
```

输入约束：

- `address` 为 42 字符十六进制 Hyperliquid 地址。
- `from < to`，使用 UTC ISO 8601。
- `page` 从 1 开始。
- `per_page` 固定为 100。
- 程序持续请求下一页，直到响应标记最后一页。

`from/to` 适用于支持历史范围的接口；当前持仓或账户快照接口按运行时点请求，并记录 `observed_at`。所有接口均遵循其各自的官方请求模型，不得强行复用永续成交参数。

#### 4.4 使用的来源字段

本版本保存所有接口的完整原始记录。永续成交标准化继续使用以下字段：

| Nansen 字段 | 用途 |
| --- | --- |
| `user` | 账户地址 |
| `timestamp` | 交易发生时间 |
| `block_number` | 来源位置证据 |
| `transaction_hash` | 交易引用 |
| `oid` | 订单标识 |
| `token_symbol` | 标的资产 |
| `side` | Long/Short 方向信息 |
| `action` | Open/Add/Close 等持仓动作 |
| `price` | 成交价格 |
| `size` | 成交数量 |
| `value_usd` | USD 名义价值 |
| `fee_usd` | USD 手续费 |
| `fee_token_symbol` | 手续费资产 |
| `closed_pnl` | 来源提供的已实现盈亏 |
| `start_position` | 成交前持仓 |
| `crossed` | 是否跨越持仓方向 |

来源响应字段缺失时，原始记录仍保存；无法满足对应标准账户事实必填字段的记录标记为转换失败，不输出猜测结果。

除永续成交外，各数据类别还必须验证并保存以下通用语义字段；Nansen 实际字段名由接口 fixture 固化：

- 事件或快照时间、观察时间及可用的区块/来源位置。
- 账户、交易哈希、订单 ID、事件 ID及其他可追溯引用。
- 资产、金额、方向、发送方、接收方和账户划转类型。
- 持仓方向、数量、入场价格、标记价格、杠杆、保证金、清算价格、已实现和未实现盈亏。
- 手续费、资金费、清算费用、奖励的金额及计价资产。
- 账户余额、可用余额、账户价值及其计价资产。

缺失的关键字段必须进入覆盖报告，不能通过默认值或业务猜测补齐。

### 5. 运行方式

本版本提供一个单次执行程序：

```text
trade-log-import \
  --address <hyperliquid-address> \
  --from <utc-time> \
  --to <utc-time>
```

运行配置：

| 配置 | 来源 | 要求 |
| --- | --- | --- |
| `NANSEN_API_KEY` | 环境变量 | 必填，不写入日志 |
| `DATABASE_URL` | 环境变量 | 必填，指向 PostgreSQL |
| `RUST_LOG` | 环境变量 | 默认 `info` |
| `--address` | CLI | 必填 |
| `--from` | CLI | 必填 |
| `--to` | CLI | 必填 |

部署环境中的环境变量由 HashiCorp Vault 注入。程序不读取项目内明文密钥文件。

### 6. 功能模块

```text
services/trade-log/
├── src/
│   ├── import_job/          # 单次导入流程协调
│   ├── raw_log/             # 原始页和原始记录模型
│   ├── normalization/       # Nansen 数据到标准账户事实
│   ├── coverage/            # 数据需求、接口能力和覆盖结果
│   ├── validation/          # 请求和标准事实校验
│   └── repository/          # 存储端口
├── adapters/
│   └── src/
│       ├── nansen/          # Nansen HTTP client
│       └── postgres/        # 原始记录和事实存储
├── bins/
│   └── trade-log-import/
└── migrations/
```

模块职责：

| 模块 | 职责 |
| --- | --- |
| `import_job` | 参数校验、多接口调用、分页循环、保存、转换、输出和结果统计 |
| `raw_log` | Nansen 各接口原始页、原始记录、快照和来源标识 |
| `normalization` | 字段映射、ID 生成和各类标准账户事实构造 |
| `coverage` | 定义账户分析数据需求，汇总接口、字段、时间和样本覆盖状态 |
| `validation` | 地址、时间、来源响应和标准事实必填字段校验 |
| `repository` | 原始数据和标准事实的存储端口 |
| `adapters/nansen` | 鉴权、请求、分页响应解析、超时和错误映射 |
| `adapters/postgres` | SQLx repository 实现和事务 |
| `trade-log-import` | composition root、配置、日志、退出码 |

### 7. 数据保存

#### 7.1 `nansen_import_runs`

| 字段 | 说明 |
| --- | --- |
| `id` | 导入运行 ID |
| `address` | 查询地址 |
| `from_time/to_time` | 查询范围 |
| `status` | RUNNING/COMPLETED/FAILED |
| `pages_fetched` | 已获取页数 |
| `endpoints_attempted/succeeded/failed` | 尝试、成功和失败的接口数 |
| `raw_count` | 原始记录数 |
| `normalized_count` | 成功转换数 |
| `failed_count` | 转换失败数 |
| `last_error` | 脱敏错误摘要 |
| `started_at/finished_at` | 运行时间 |

#### 7.2 `raw_logs`

保存：

- Nansen 来源信息。
- 请求地址和日期范围。
- 接口标识、请求类别、页码或快照标识及完整响应。
- 每条原始记录 JSON，包括成交、持仓、资金活动、费用、资金费、清算、奖励、账户快照和管理活动。
- 原始记录内容哈希。
- 获取时间和转换状态。

原始记录内容哈希建立唯一约束，重复执行相同查询不重复插入完全相同的记录。

#### 7.3 `nansen_coverage_results`

保存：

- 导入运行 ID、数据类别和 Nansen 接口。
- 覆盖状态：`COMPLETE/PARTIAL/UNAVAILABLE/UNVERIFIED`。
- 请求时间范围或快照观察时间。
- 页数、记录数、最早/最晚事件时间。
- 必需字段存在情况、缺失字段和验证证据。
- 脱敏后的失败原因或订阅权限信息。

#### 7.4 `account_fact_versions`

保存：

- `fact_id`。
- 固定 `revision=1`。
- `fact_type` 为本版本成功映射的标准账户事实类型。
- `account_key`。
- 发生时间和排序键。
- 标准化 payload。
- 对应原始记录 ID。
- `schema_version=1`。

V0.000 使用规范化原始记录的 SHA-256 生成 `fact_id`。该方式只保证完全相同记录的幂等，不实现来源事实的后续 revision 关联。

### 8. 标准化与验证输出

#### 8.1 固定字段

```text
event_type        = account.fact.v1
schema_version    = 1
fact_type         = TRADE / TRANSFER / FEE / FUNDING / REWARD /
                    LIQUIDATION_FEE / ACCOUNT_SNAPSHOT
revision          = 1
change_type       = UPSERT
confirmation      = OBSERVED
chain_id          = hyperliquid
protocol          = hyperliquid
source            = NANSEN
```

`instrument_type=PERPETUAL` 和 `quote_asset=USDC` 仅适用于能够确认该语义的 Hyperliquid 永续交易事实，不作为所有事实类型的固定字段。

#### 8.2 永续成交字段映射

| 标准字段 | 来源/规则 |
| --- | --- |
| `account` | `user` 规范化为小写 |
| `account_key` | `hyperliquid:hyperliquid:{account}` |
| `occurred_at` | `timestamp` 转 UTC |
| `transaction_hash` | `transaction_hash` |
| `order_id` | `oid` 转字符串 |
| `operation_id` | `transaction_hash + oid` |
| `base_asset` | `token_symbol` |
| `market` | `hyperliquid:{token_symbol}-USDC` |
| `price` | `price` 转十进制字符串 |
| `quantity` | `size` 转十进制字符串 |
| `notional` | `value_usd` 转十进制字符串 |
| `fee` | `fee_usd` 转十进制字符串 |
| `fee_asset` | `fee_token_symbol` |
| `reported_realized_pnl` | `closed_pnl` 转十进制字符串 |
| `side` | 根据 `side + action + start_position` 的固定映射表生成 |
| `position_effect` | 根据 `action + start_position + crossed` 的固定映射表生成 |

`side` 和 `position_effect` 的映射表必须使用 Nansen 固定样本验证。无法唯一映射时输出 `position_effect=UNKNOWN`；无法确认成交买卖方向时该记录转换失败，不猜测 `BUY/SELL`。

#### 8.3 其他事实映射要求

- 资金变化映射为 `TRANSFER`，必须保留资产、金额、方向、对手方、交易引用和发生时间；无法确认方向时转换失败。
- 独立手续费映射为 `FEE`，资金费映射为 `FUNDING`，不得把二者重复计入成交 payload 和独立事实。
- 清算相关费用映射为 `LIQUIDATION_FEE`；清算成交仍可映射为 `TRADE`，并标记协议触发语义。
- 奖励映射为 `REWARD`，必须具有来源类型、资产和金额。
- 持仓或账户状态映射为 `ACCOUNT_SNAPSHOT`，必须携带快照观察时间；当前快照不能伪装成历史事件。
- `ApproveAgent` 等不改变账户分析数值的管理活动只进入原始日志，不生成上述账户事实，除非总体 schema 后续增加对应类型。

#### 8.4 输出格式

- 标准输出使用 UTF-8 JSON Lines。
- 每行是一个完整 `AccountFactEnvelope`。
- 标准事实流只输出本次运行成功转换的事实。
- 输出按 `occurred_at`、`ordering_key`、`sub_index` 升序排列。
- 日志写入标准错误流，不能混入标准输出。
- 输出中的金额、价格和数量均为十进制字符串。
- 覆盖报告单独输出为 JSON，不混入事实 JSON Lines；报告必须列出每个必需数据类别的状态、接口、记录数、时间范围、字段缺口和证据。

### 9. 错误处理

| 场景 | 行为 |
| --- | --- |
| 缺少 API key/数据库配置 | 启动失败，退出码 2 |
| 地址或时间范围无效 | 参数失败，退出码 2 |
| Nansen 401/403/402 | 不重试，记录脱敏错误，运行失败 |
| Nansen 429 | 有界退避重试，超过次数后失败 |
| Nansen 5xx/网络超时 | 有界退避重试，超过次数后失败 |
| 单条记录转换失败 | 保存失败原因，继续处理其他记录 |
| 某数据类别无 Nansen 接口 | 记录 `UNAVAILABLE` 和验证证据，不伪造数据 |
| 接口可用但关键字段不足 | 保存原始数据，记录 `PARTIAL` 和缺失字段 |
| 接口因订阅权限不可用 | 保存脱敏状态和 HTTP 状态，区分能力缺失与权限不足 |
| PostgreSQL 写入失败 | 当前页事务回滚，运行失败 |
| 标准输出写入失败 | 运行失败，不删除已保存数据 |

API 响应和日志不得包含 `NANSEN_API_KEY`。错误日志包含运行 ID、页码、HTTP 状态和稳定错误码。

### 10. 验收标准

- 使用合法地址和时间范围可以完成一次运行并以退出码 0 结束。
- 程序按账户分析数据需求逐类调用和验证 Nansen 接口，不得只调用 `perp-trades`。
- 请求携带 `apikey`，并完整处理所有接口的分页或快照响应。
- 每个成功响应页、快照和原始记录均可在 PostgreSQL 查询。
- 永续成交和当前持仓必须完成真实接口验证；资金变化、独立费用/资金费、清算/奖励、账户快照必须形成有证据的覆盖结论。
- 对 `Withdraw`、`SendAsset` 等资金活动，必须证明能够从 Nansen 获取并映射，或在覆盖报告中明确标记缺口。
- `ApproveAgent` 等管理活动能获取时保存原始记录，并与数值型账户事实区分。
- 每个成功转换的原始记录可追溯到一个标准账户事实；无需转换的管理活动具有明确分类原因。
- 标准输出每行均能反序列化为 schema v1 的 `AccountFactEnvelope`。
- 输出顺序稳定，同一输入重复运行得到相同事实内容和 `fact_id`。
- 重复执行不会重复保存相同原始记录和标准事实。
- 单条坏数据不会阻塞同页其他数据转换。
- 覆盖报告包含所有必需数据类别，状态只能是 `COMPLETE/PARTIAL/UNAVAILABLE/UNVERIFIED`。
- 只有所有必需类别均为 `COMPLETE` 时，验收结论才能声明 Nansen 单一数据源足以支持完整账户分析。
- 任何 `PARTIAL/UNAVAILABLE/UNVERIFIED` 都必须列为 V0.000 已知数据缺口，并提出后续补充数据源或接口方案。
- API key 不出现在标准输出、日志、数据库普通字段和错误信息中。
- 本版本不启动账户分析、数据审核和跟单相关进程。

### 11. 测试记录要求

开发完成后在本节追加，不覆盖版本目标和范围：

| 项目 | 记录内容 |
| --- | --- |
| 单元测试 | 各事实映射和覆盖状态判定的用例数、通过数、失败数 |
| 集成测试 | 各 Nansen 接口 fixture、PostgreSQL、完整导入和覆盖报告结果 |
| 手工验证 | 执行命令、地址脱敏值、时间范围、快照时间、各接口和各类别统计 |
| 数据覆盖结论 | 每个必需数据类别的状态、证据、字段缺口和是否足以支持账户分析 |
| 已知问题 | 数据缺口、程序问题、影响和后续处理版本 |
| 构建信息 | Git revision、Rust toolchain 和 schema 版本 |

### 12. 变更记录

| 日期 | 状态 | 内容 |
| --- | --- | --- |
| 2026-09-07 | PLANNED | 建立 V0.000 版本目标、范围、数据源、模块和验收标准 |
| 2026-09-11 | PLANNED | 将范围从仅验证永续成交扩大为验证 Nansen 是否覆盖完成 Hyperliquid 账户分析所需的全部输入数据；增加持仓、资金活动、费用、资金费、清算、奖励、账户快照和覆盖报告要求 |
