# 版本开发记录

## V0.000

### 1. 版本信息

| 项目 | 内容 |
| --- | --- |
| 版本号 | V0.000 |
| 版本状态 | PLANNED |
| 记录日期 | 2026-09-07 |
| 开发状态 | 未开始 |
| 对应总体设计 | [详细设计](detailed-design.md) |

### 2. 版本目标

交付一个能够独立运行的最小交易日志程序，完成以下闭环：

```text
读取运行参数
→ 调用 Nansen API 获取指定地址的 Hyperliquid 永续成交
→ 保存 Nansen 原始响应记录
→ 转换为标准交易事实
→ 保存标准交易事实
→ 以 JSON Lines 输出标准交易事实
→ 正常退出
```

本版本只验证第三方数据接入、原始数据保存、标准化转换和输出链路，不实现完整系统运行能力。

### 3. 本版本范围

#### 3.1 包含范围

- 建立最小 Rust Cargo workspace 和交易日志服务目录。
- 提供一个单次执行的交易日志导入程序。
- 接入 Nansen API。
- 获取指定 Hyperliquid 地址在指定日期范围内的永续成交。
- 处理 Nansen 分页响应，直到 `is_last_page=true`。
- 保存每页原始响应和每条原始成交记录。
- 将原始成交转换为 `AccountFactPayload::Trade`。
- 保存标准交易事实。
- 对完全相同的原始记录和标准事实进行幂等去重。
- 将本次获得的标准交易事实按发生时间顺序输出为 JSON Lines。
- 提供基本配置校验、错误日志和退出码。
- 提供转换单元测试和一次完整接入流程测试。

#### 3.2 不包含范围

- 账户分析服务。
- 数据完整审核服务。
- 跟单服务。
- 账户价值、净盈亏、成交量、交易次数和胜率计算。
- 链节点直连和链上区块采集。
- Nansen 实时订阅或持续轮询。
- 多数据源切换和来源对账。
- Kafka/Redpanda 消息发布。
- Outbox 发布流程。
- 事实修订、撤销、最终性和链重组处理。
- 修复请求、历史回补任务和账户重算。
- HTTP 查询 API。
- 现货交易和其他链、协议数据。
- 评分、地址选择和跟单信号执行。

### 4. 数据源

#### 4.1 Nansen 接口

| 项目 | 内容 |
| --- | --- |
| 产品入口 | `https://app.nansen.ai/api` |
| API 主机 | `https://api.nansen.ai` |
| 接口 | `POST /api/v1/profiler/perp-trades` |
| 鉴权 | HTTP Header `apikey` |
| Content-Type | `application/json` |
| 数据类型 | Hyperliquid 地址永续成交 |
| 调用方式 | 指定地址、日期范围和分页参数的批量查询 |

官方参考：

- [Nansen API Authentication](https://docs.nansen.ai/getting-started/authentication)
- [Nansen Address Perp Trades](https://docs.nansen.ai/api/profiler/address-perp-trades)

#### 4.2 请求参数

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

#### 4.3 使用的来源字段

本版本保存完整原始记录，标准化转换使用以下字段：

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

来源响应字段缺失时，原始记录仍保存；无法满足标准交易事实必填字段的记录标记为转换失败，不输出猜测结果。

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
│   ├── normalization/       # Nansen 成交到 TradeFact
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
| `import_job` | 参数校验、分页循环、保存、转换、输出和结果统计 |
| `raw_log` | Nansen 原始页、原始成交和来源标识 |
| `normalization` | 字段映射、ID 生成和标准交易事实构造 |
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
| `raw_count` | 原始记录数 |
| `normalized_count` | 成功转换数 |
| `failed_count` | 转换失败数 |
| `last_error` | 脱敏错误摘要 |
| `started_at/finished_at` | 运行时间 |

#### 7.2 `raw_logs`

保存：

- Nansen 来源信息。
- 请求地址和日期范围。
- 页码及完整页响应。
- 每条原始成交 JSON。
- 原始记录内容哈希。
- 获取时间和转换状态。

原始记录内容哈希建立唯一约束，重复执行相同查询不重复插入完全相同的记录。

#### 7.3 `account_fact_versions`

保存：

- `fact_id`。
- 固定 `revision=1`。
- `fact_type=TRADE`。
- `account_key`。
- 发生时间和排序键。
- 标准化 payload。
- 对应原始记录 ID。
- `schema_version=1`。

V0.000 使用规范化原始记录的 SHA-256 生成 `fact_id`。该方式只保证完全相同记录的幂等，不实现来源事实的后续 revision 关联。

### 8. 标准化输出

#### 8.1 固定字段

```text
event_type        = account.fact.v1
schema_version    = 1
fact_type         = TRADE
revision          = 1
change_type       = UPSERT
confirmation      = OBSERVED
chain_id          = hyperliquid
protocol          = hyperliquid
instrument_type   = PERPETUAL
quote_asset       = USDC
source            = NANSEN
```

#### 8.2 字段映射

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

#### 8.3 输出格式

- 标准输出使用 UTF-8 JSON Lines。
- 每行是一个完整 `AccountFactEnvelope`。
- 只输出本次运行成功转换的事实。
- 输出按 `occurred_at`、`ordering_key`、`sub_index` 升序排列。
- 日志写入标准错误流，不能混入标准输出。
- 输出中的金额、价格和数量均为十进制字符串。

### 9. 错误处理

| 场景 | 行为 |
| --- | --- |
| 缺少 API key/数据库配置 | 启动失败，退出码 2 |
| 地址或时间范围无效 | 参数失败，退出码 2 |
| Nansen 401/403/402 | 不重试，记录脱敏错误，运行失败 |
| Nansen 429 | 有界退避重试，超过次数后失败 |
| Nansen 5xx/网络超时 | 有界退避重试，超过次数后失败 |
| 单条记录转换失败 | 保存失败原因，继续处理其他记录 |
| PostgreSQL 写入失败 | 当前页事务回滚，运行失败 |
| 标准输出写入失败 | 运行失败，不删除已保存数据 |

API 响应和日志不得包含 `NANSEN_API_KEY`。错误日志包含运行 ID、页码、HTTP 状态和稳定错误码。

### 10. 验收标准

- 使用合法地址和时间范围可以完成一次运行并以退出码 0 结束。
- 请求携带 `apikey`，并完整处理所有分页。
- 每个成功响应页和原始成交均可在 PostgreSQL 查询。
- 每个成功转换的原始成交可追溯到一个标准交易事实。
- 标准输出每行均能反序列化为 schema v1 的 `AccountFactEnvelope`。
- 输出顺序稳定，同一输入重复运行得到相同事实内容和 `fact_id`。
- 重复执行不会重复保存相同原始记录和标准事实。
- 单条坏数据不会阻塞同页其他数据转换。
- API key 不出现在标准输出、日志、数据库普通字段和错误信息中。
- 本版本不启动账户分析、数据审核和跟单相关进程。

### 11. 测试记录要求

开发完成后在本节追加，不覆盖版本目标和范围：

| 项目 | 记录内容 |
| --- | --- |
| 单元测试 | 用例数、通过数、失败数 |
| 集成测试 | Nansen fixture、PostgreSQL 和完整导入结果 |
| 手工验证 | 执行命令、地址脱敏值、时间范围和统计 |
| 已知问题 | 问题、影响和后续处理版本 |
| 构建信息 | Git revision、Rust toolchain 和 schema 版本 |

### 12. 变更记录

| 日期 | 状态 | 内容 |
| --- | --- | --- |
| 2026-09-07 | PLANNED | 建立 V0.000 版本目标、范围、数据源、模块和验收标准 |

