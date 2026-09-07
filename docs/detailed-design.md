# 链上账户投资分析服务详细设计

## 1. 文档概述

### 1.1 编写目的

本文定义链上账户投资分析服务的完整详细设计，包括系统边界、服务与功能模块、接口契约、数据库模型、核心流程、Rust 程序结构和部署方式，作为编码、测试、部署及验收依据。

本文描述系统整体目标，不规定具体实施顺序、阶段范围或优先接入的链与协议；相关内容由 Roadmap 管理。

### 1.2 项目背景

系统根据链上地址的交易和账户事实，还原地址在各链及协议中的账户状态，计算账户价值、净盈亏、成交量、交易次数和胜率，为人工判断地址投资表现提供数据基础，并为后续策略评分和自动选择跟单地址保留扩展能力。

交易数据既可以直接从链或协议获取，也可以来自第三方数据服务。实时数据需要尽快整理并发布，供账户分析和跟单使用；完整性复核在主流程之外独立运行，发现问题后推动回补、修正和重算。

系统的主要交付是：

1. 账户分析数据。
2. 实时标准交易事实；同一事实在跟单场景中直接作为跟单信号。

### 1.3 适用范围

本文适用于：

- 交易日志服务。
- 账户分析服务。
- 数据完整审核服务。
- 跟单服务。
- 查询接口、消息契约和数据库设计。
- 多链、多协议和现货、合约等交易类型扩展。
- Rust workspace、模块边界和部署进程。

评分规则、自动地址选择算法和具体实施 Roadmap 不在本文中固定，仅定义扩展边界。

### 1.4 术语与缩略语

| 术语 | 含义 |
| --- | --- |
| 账户/地址 | 某条链、某个协议内被跟踪的链上账户 |
| `account_key` | `chain_id + protocol + normalized_account` 组成的协议账户唯一键 |
| 原始日志 | 数据源返回且尚未标准化的交易、资金或账户数据 |
| 标准账户事实 | 账户分析的统一输入模型，包含交易、资金、费用和快照等事实 |
| 标准交易事实 | 标准账户事实中的 `TRADE` 类型，同时作为跟单信号 |
| `fact_id` | 一个业务事实跨修订保持稳定的唯一 ID |
| `revision` | 同一事实的版本号，从 1 开始递增 |
| 数据水位 | 已处理到账户事实流中的确定位置 |
| 完整持仓周期 | 合约净持仓从零变为非零，最终重新回到零的完整过程 |
| 现货卖出操作 | 合并同一订单或协议操作中多个 fill 后的一次逻辑卖出 |
| Outbox | 业务写入与待发布消息在同一事务保存的可靠发布模式 |
| MPC | Multi-Party Computation，多方安全计算；用于多方协同管理密钥材料并完成签名 |
| KMS/HSM | 密钥管理服务/硬件安全模块 |
| SLA/SLO | 服务协议/服务目标 |

### 1.5 参考资料

- [项目说明](../readme.md)
- [概要设计](overview-design.md)
- [8 维权重参考](8维权重.md)
- 各接入链、协议、数据供应商和钱包的官方接口文档

## 2. 系统概述

### 2.1 系统目标

- 持续获取指定账户的交易及账户事实。
- 将异构数据整理为统一、可追溯、可修订的标准账户事实。
- 还原现货库存、合约持仓、账户状态和交易过程。
- 计算账户价值、净盈亏、成交量、交易次数和胜率。
- 快速发布标准交易事实，支持跟单规则判断和交易执行。
- 在不阻塞实时主流程的情况下复核数据完整性并形成修复闭环。
- 支持计算规则、链、协议、数据源和钱包实现持续扩展。

### 2.2 系统功能概述

| 服务 | 定位 | 核心职责 |
| --- | --- | --- |
| 交易日志服务 | 前置服务 | 数据获取、原始日志存储、解析标准化、轻量校验、事实发布和修复 |
| 账户分析服务 | 核心服务 | 事实账本、持仓与库存、账户估值、盈亏和指标计算、重算与查询 |
| 数据完整审核服务 | 独立保护服务 | 采集健康、范围完整、来源对账、业务一致性、异常和修复闭环 |
| 跟单服务 | 扩展服务 | 订阅交易事实、规则判断、交易构建、钱包签名、发送和状态追踪 |

### 2.3 系统边界

系统负责：

- 指定账户数据的采集、标准化、存储和发布。
- 账户事实版本、顺序和修正语义。
- 账户分析结果的计算、版本化和查询。
- 数据完整性审核和修复编排。
- 跟单配置、决策、签名调用、交易发送和执行记录。

系统不负责：

- 改变链或协议的最终结算结果。
- 保证第三方数据源本身绝对正确。
- 在标准交易事实中给出策略评价或风险提示。
- 将评分规则固化为不可替换的唯一算法。
- 管理外部 MPC 钱包、第三方托管钱包或其他远程签名服务内部的密钥材料。

### 2.4 用户角色

| 角色 | 权限范围 |
| --- | --- |
| 分析数据用户 | 查询账户状态、持仓、交易过程、盈亏和指标 |
| 跟单用户 | 管理自己的跟单配置、钱包引用并查询执行记录 |
| 运营人员 | 管理关注账户、查看采集状态和一般告警 |
| 审核人员 | 创建审核任务、查看异常、确认需要人工处理的修复 |
| 系统管理员 | 管理服务配置、协议能力、权限和计算版本发布 |
| 内部服务身份 | 按最小权限访问本服务数据库、消息主题和内部接口 |

### 2.5 运行环境

| 能力 | 基准实现 |
| --- | --- |
| 操作系统 | Linux x86_64 |
| 容器编排 | Kubernetes |
| 编程语言 | Rust stable |
| 异步运行时 | Tokio |
| Web 框架 | Axum |
| 服务接口 | HTTP/JSON |
| 事件总线 | Redpanda |
| 业务数据库 | PostgreSQL |
| 数据访问 | SQLx |
| 原始日志归档 | MinIO S3 对象存储 |
| 缓存 | Redis |
| 可观测性采集 | OpenTelemetry Collector |
| 指标 | Prometheus |
| 日志 | Loki |
| 链路追踪 | Tempo |
| 可视化与告警 | Grafana |
| 密钥与凭证 | HashiCorp Vault |

以上组件构成统一运行基线。服务实现、部署配置和验收环境均以该技术路线为准。

## 3. 系统架构设计

### 3.1 总体架构

```text
链节点 / 协议接口 / 第三方数据源
                  │
                  ▼
           交易日志服务
       原始日志 → 标准账户事实 → account.fact.v1
                  │                    │
                  │                    ├── TRADE → 跟单服务 → 钱包 → 链/协议
                  │                    │
                  │                    └── 全部事实 → 账户分析服务 → 查询结果
                  │
                  └── 数据完整审核服务 → 修复请求 → 交易日志服务
                                               └→ 新 revision → 分析重算
```

实时事实发布不等待完整性审核。交易日志服务负责主流程内的基本格式、来源、位置和幂等检查；审核服务异步发现漏数、错数、重复、重组和业务不一致。

### 3.2 技术架构

系统采用服务目录内聚、功能模块优先、端口适配器隔离的代码架构：

- 四个业务服务分别拥有自己的业务模块、适配器、运行入口和 migration。
- 跨服务模型集中在公共契约 crate。
- 服务模块通过使用方定义的 trait 访问数据库、消息系统、市场数据、协议和钱包。
- 业务计算保持同步、确定性和无副作用；I/O 和事务由服务协调模块执行。
- 所有消息按至少一次投递设计，生产与消费双方均实现幂等。
- 账户分析结果为可重建派生数据，事实与结果均版本化。

### 3.3 系统部署架构

| 进程 | 所属服务 | 职责 |
| --- | --- | --- |
| `trade-collector` | 交易日志 | 历史拉取、实时订阅、检查点和原始日志保存 |
| `trade-parser-publisher` | 交易日志 | 解析、标准化、事实版本、最终性和 outbox 发布 |
| `account-analyzer` | 账户分析 | 消费账户事实并执行增量计算 |
| `account-recalculator` | 账户分析 | 历史回放、修正重算和计算版本切换 |
| `data-auditor` | 数据审核 | 各类审核、异常管理和修复编排 |
| `copy-trader` | 跟单 | 规则判断、交易构建、钱包调用和执行追踪 |
| `query-api` | 查询网关 | 认证、路由和各服务查询结果组合 |

进程可以独立扩缩容。实时采集、实时分析和跟单使用独立资源池；重算与审核不得抢占实时处理资源。

### 3.4 系统组件划分

#### 交易日志服务

- acquisition：链上获取和第三方订阅。
- checkpoint：游标、水位和恢复。
- raw_log：原始日志及归档。
- parsing：协议解析。
- normalization：标准账户事实构造。
- validation：同步轻量校验。
- publishing：事实版本、outbox 和消息发布。
- finality：确认、撤销和替代事实。
- repair：回补和重新解析。

#### 账户分析服务

- fact_ledger：事实入账、顺序、revision 和水位。
- derivative_position/episode：合约持仓和完整持仓周期。
- spot_inventory/disposal：现货库存和逻辑卖出。
- capital_flow、fee_and_funding：资金、费用和资金费。
- account_state、valuation：账户状态和账户价值。
- pnl、performance_metrics：净盈亏和基础指标。
- calculation_engine：计算器注册、依赖和调度。
- recalculation、result_store、query：重算、结果和查询。

#### 数据完整审核服务

- collection_health：采集健康。
- range_completeness：范围连续性。
- source_reconciliation：多来源对账。
- business_consistency：账户事实业务一致性。
- findings：异常与证据。
- repair_request、repair_verification：修复请求和复验。

#### 跟单服务

- subscription：交易事实过滤和消费。
- config、rule_engine、sizing：配置、规则和数量。
- trade_intent、transaction：交易意图和交易构建。
- wallet、submission：签名、发送和状态查询。
- execution_record：决策与执行状态。

### 3.5 组件依赖关系

```text
运行入口
├── 服务功能模块
├── 服务 adapters
└── 协议 implementations

服务 adapters ──实现──→ 功能模块定义的 ports
协议 implementations ──实现──→ 采集、结算、估值或交易构建 ports
服务功能模块 ──依赖──→ shared-types / account-facts / internal-events
```

禁止：

- 业务模块依赖 PostgreSQL、Kafka、HTTP 或具体 SDK。
- 一个服务导入另一个服务的数据库实现。
- handler、consumer、SQL mapper 或 `main` 实现业务规则。
- 计算器直接读写数据库或访问网络。

### 3.6 外部系统依赖

| 外部系统 | 用途 | 失败处理 |
| --- | --- | --- |
| 链节点/协议接口 | 历史数据、实时数据、交易状态 | 重试、检查点恢复、备用源和审核 |
| 第三方数据服务 | 交易和账户数据订阅 | 主流程按配置信任，异步对账 |
| 市场数据服务 | 账户估值和统一计价 | 使用版本化价格证据，缺失时降级 |
| Kafka/Redpanda | 账户事实和内部事件 | outbox、消费幂等、重试和死信 |
| PostgreSQL | 各服务业务状态 | 事务、约束、备份和迁移 |
| 对象存储 | 原始日志长期归档 | 重试并保留数据库状态 |
| 本地签名器/外部钱包服务 | 本地 KMS/HSM 签名，或调用 MPC、托管及其他远程签名服务 | 能力检测、幂等、超时后查询 |
| 链 RPC/提交接口 | 广播和回执查询 | 未知状态先查询，不直接重发 |

## 4. 功能模块详细设计

### 4.1 交易日志服务

#### 4.1.1 模块概述

交易日志服务是事实生产方，负责从链、协议或第三方获取数据，保存原始内容，转换为标准账户事实并及时发布。该服务是原始日志、事实版本和 current 指针的唯一写入方。

#### 4.1.2 功能说明

- 支持历史范围拉取和实时订阅。
- 保存可恢复检查点和数据源定位信息。
- 幂等保存原始日志。
- 按链和协议版本解析、标准化。
- 执行发布前轻量校验。
- 同事务保存事实版本与 outbox。
- 处理确认、重组、撤销和供应商修正。
- 接收审核服务修复请求并执行回补。

#### 4.1.3 业务流程

```text
Source Adapter → acquisition → checkpoint + raw_log
raw_log → parsing → normalization → validation
validation → fact version + current pointer + outbox
outbox worker → account.fact.v1
```

#### 4.1.4 输入与输出

输入：

- 关注账户和起始范围。
- 链节点、协议接口或第三方数据。
- 修复请求。
- 链最终性和源侧修正信息。

输出：

- 原始日志和采集状态。
- `AccountFactEnvelope`。
- `account.fact.v1` 消息。
- 修复执行结果和新版事实。

#### 4.1.5 核心处理逻辑

1. 根据采集任务读取适配器能力。
2. 拉取或订阅原始数据，并以稳定来源 ID 幂等保存。
3. 使用 `(chain_id, protocol, parser_version)` 选择解析器。
4. 将一条原始记录转换为一个或多个账户事实，以 `sub_index` 区分。
5. 校验必填字段、标识、数值、来源位置和 schema 版本。
6. 在同一事务内写入事实版本、更新 current 指针并写 outbox。
7. outbox worker 发布后记录结果；失败按有界退避重试。
8. 修正沿用 `fact_id` 并递增 revision；替代事实使用新 `fact_id`。

#### 4.1.6 状态及状态转换

采集任务：

```text
PENDING → RUNNING → COMPLETED
                 ├→ RETRY_WAIT → RUNNING
                 ├→ PAUSED
                 └→ FAILED
```

事实确认：

```text
OBSERVED → CONFIRMED
    └────→ ORPHANED/RETRACTED
```

Outbox：

```text
PENDING → PUBLISHING → PUBLISHED
                      └→ RETRY_WAIT → PUBLISHING
```

#### 4.1.7 异常处理

| 异常 | 处理 |
| --- | --- |
| 数据源超时 | 指数退避，达到阈值后切换备用源并告警 |
| 订阅断开 | 从持久化游标恢复；游标失效时范围回拉 |
| 重复原始日志 | 唯一约束忽略，不重复解析 |
| 解析失败 | 保存失败原因和解析器版本，进入隔离/死信 |
| 未知协议版本 | 隔离并告警，不生成猜测事实 |
| 发布失败 | 保留 outbox 并重试，不回滚已保存事实 |
| 链重组 | 发布更高 revision 的撤销或替代事实 |

#### 4.1.8 权限控制

- 采集进程只写采集任务、检查点和原始日志。
- 解析发布进程写事实版本、current 和 outbox。
- 修复接口仅允许审核服务和授权运维身份调用。
- 查询身份不得修改事实。
- 数据源凭证通过密钥引用加载，不写入日志或数据库明文字段。

### 4.2 账户分析服务

#### 4.2.1 模块概述

账户分析服务消费全部标准账户事实，建立版本化事实账本，分别计算合约和现货状态，再生成账户价值、净盈亏、成交量、交易次数和胜率。评分和地址自动选择作为扩展计算器预留。

#### 4.2.2 功能说明

- 事实幂等消费、排序和 revision 管理。
- 合约持仓状态和完整持仓周期归并。
- 现货移动加权平均成本库存和卖出操作归并。
- 充值、提现、转账、手续费、资金费和清算费用处理。
- 当前账户价值和统一 USD 计价。
- 正式净盈亏和基础表现指标。
- 迟到、修正、撤销和规则升级后的重算。
- 版本化结果、数据水位、覆盖范围和质量状态查询。

#### 4.2.3 业务流程

```text
AccountFact Consumer → fact_ledger → calculation_engine
                                      ├→ derivatives
                                      ├→ spot
                                      ├→ account_state/valuation
                                      └→ pnl/performance_metrics
                                                ↓
                                      result_store → query
```

#### 4.2.4 输入与输出

输入：

- 当前有效的交易、资金、费用和账户快照事实。
- 版本化市场价格、汇率和资产元数据。
- 计算方案、计算器版本及协议结算策略。
- 重算请求。

输出：

- 当前合约持仓和现货库存。
- 合约持仓周期和现货卖出操作。
- `account_value`。
- `realized_net_pnl`、`unrealized_pnl`、`total_net_pnl`、`estimated_pnl`。
- USD 成交量和成交笔数。
- 分类交易次数、胜负数、持平数和胜率。
- 计算版本、水位、覆盖范围、估算标记和数据质量。

#### 4.2.5 核心处理逻辑

##### 合约

- 按 `(account_key, market)` 归并。
- 净持仓从 0 变为非 0 时开始持仓周期，回到 0 时结束。
- 加仓和减仓属于同一周期。
- 穿越 0 反向时，结束原周期并以剩余数量开始新周期。
- 手续费、资金费和清算费用归属对应周期。
- 已结束周期计入交易次数和胜率，进行中周期不计入。

##### 现货

- 按 `(account_key, base_asset)` 维护数量、总成本和移动平均成本。
- 买入增加数量和成本；卖出按数量分摊成本并产生已实现盈亏。
- 同一订单的多个 fill 合并为一个逻辑卖出操作。
- 交易对中的计价资产不重复统计为另一笔卖出。
- 转入增加库存；没有可靠成本时使用转入时市场价格并标记 `ESTIMATED`。
- 转出按数量同比例移出库存成本，不产生交易盈亏。

##### 账户价值与盈亏

- 账户价值是数据水位对应时点的状态，不按持仓周期计算。
- 协议提供并通过校验的账户价值优先作为状态输入，否则由协议估值策略计算。
- 每个合约周期或现货卖出操作只有一个正式 `realized_net_pnl`。
- `calculation_method` 为 `PROTOCOL_REPORTED`、`RECONSTRUCTED` 或 `ESTIMATED`。
- `total_net_pnl = realized_net_pnl + unrealized_pnl`。
- `estimated_pnl` 是总盈亏的一部分，不重复相加。
- 成交量保留原始计价，并使用版本化价格折算 USD；无法换算的成交不进入总值并降低质量状态。

##### 指标

- 合约 `trade_count` 为已结束完整持仓周期数。
- 现货 `trade_count` 为产生已实现盈亏的逻辑卖出操作数。
- `win_rate = winning_trade_count / (winning_trade_count + losing_trade_count)`。
- 持平交易计入 `trade_count` 和 `breakeven_trade_count`，不进入胜率分母。
- 没有盈利或亏损交易时 `win_rate=null`。
- 现货和合约的交易次数、胜率分别输出，不生成混合总值。

##### 历史与估算

- 默认统计系统当前可验证的全部历史，不使用默认时间窗口。
- 结果包含 `tracking_started_at`、`coverage_start`、`calculated_to`、`history_complete`。
- 观察起点已有合约持仓以当时可靠价格建立 `ESTIMATED_OPEN` 周期。
- 观察起点已有现货资产建立估算库存批次。
- 历史补齐后触发重算并替换估算成本。

#### 4.2.6 状态及状态转换

合约周期：

```text
OPEN → INCREASE/DECREASE → CLOSED
  ├→ REVERSED（原周期 CLOSED，新周期 OPEN）
  └→ LIQUIDATED（按实际剩余仓位决定 CLOSED 或继续 OPEN）
```

计算结果：

```text
BUILDING → VALIDATING → ACTIVE
    └───────────────→ INVALID
ACTIVE → SUPERSEDED
```

重算任务：

```text
PENDING → RUNNING → COMPLETED
                 └→ RETRY_WAIT → RUNNING
                 └→ FAILED
```

#### 4.2.7 异常处理

| 异常 | 处理 |
| --- | --- |
| 重复事实 | `(event_id)` 消费幂等，不重复改变状态 |
| 旧 revision | 忽略并记录指标 |
| revision 跳号 | 暂停该事实应用并请求缺失版本 |
| 迟到/修正/撤销 | 从受影响位置前快照创建重算任务 |
| 缺少市场价格 | 输出 `UNAVAILABLE` 或降级，不猜测 |
| 无前态减仓 | 标记业务不一致并通知审核服务 |
| 计算不变量失败 | 不切换结果版本，保留旧结果并告警 |
| 重算期间查询 | 返回旧 ACTIVE 版本并标记 `RUNNING` |

#### 4.2.8 权限控制

- 分析 consumer 仅消费账户事实并写分析账本和增量结果。
- 重算进程只写重算任务和新结果版本，不修改事实源。
- 计算方案发布、默认版本切换需要管理员权限和审计日志。
- 外部查询按账户数据授权范围过滤。
- 分析服务不持有钱包签名权限。

### 4.3 数据完整审核服务

#### 4.3.1 模块概述

审核服务独立于实时主流程，持续检查采集健康、范围完整、来源差异、最终性和账户事实业务一致性。审核服务只发现和记录问题，并通过修复请求驱动交易日志服务处理。

#### 4.3.2 功能说明

- 分钟级采集健康检查。
- 范围连续性和分页完整性检查。
- 第三方与权威/备用来源抽样或全量对账。
- 链重组和最终性变化检查。
- 无前态减仓、余额异常和无法闭合序列检查。
- 异常证据、严重程度和状态管理。
- 自动或人工修复请求及修复后复验。

#### 4.3.3 业务流程

```text
调度审核 → 读取检查点/事实/来源 → 执行检查 → findings
findings → 自动或人工决策 → RepairRequested
交易日志服务修复 → 新 revision → 账户重算 → 审核复验关闭
```

#### 4.3.4 输入与输出

输入：

- 采集任务、检查点和原始日志只读数据。
- 标准账户事实和分析结果只读数据。
- 权威来源或备用来源数据。
- 修复完成通知。

输出：

- 审核任务和统计。
- 异常、证据和数据质量状态。
- `RepairRequested`。
- 告警和复验结果。

#### 4.3.5 核心处理逻辑

| 检查 | 内容 |
| --- | --- |
| 采集健康 | 心跳、游标推进、延迟和错误率 |
| 范围完整 | 区块/位置连续性、分页结束、数量异常 |
| 来源对账 | 缺失、重复、字段差异和最终性差异 |
| 业务一致 | 仓位、库存、余额和费用关联不变量 |
| 修复复验 | 缺口是否补齐、新 revision 是否发布、分析是否重算 |

审核结果按链与来源范围、账户与事实范围维护 `UNKNOWN`、`CHECKING`、`HEALTHY`、`DEGRADED`、`INVALID`、`REPAIRING` 状态。

#### 4.3.6 状态及状态转换

```text
AuditJob: PENDING → RUNNING → COMPLETED/FAILED
Finding:  OPEN → REPAIRING → VERIFYING → CLOSED
                 └────────────────────→ OPEN
```

#### 4.3.7 异常处理

- 数据源不可用时保留任务进度并退避重试。
- 对账来源均不可用时保持 `UNKNOWN`，不标记健康。
- 无法自动修复时保留证据并升级人工处理。
- 重复修复请求通过幂等键合并。
- 审核失败不得阻塞账户事实实时发布。

#### 4.3.8 权限控制

- 审核服务只读交易日志和分析数据。
- 审核服务只能写审核库并发布修复请求。
- 人工关闭、忽略或重新打开异常需要审核权限并记录审计日志。
- 审核服务无权修改事实、分析结果和跟单记录。

### 4.4 跟单服务

#### 4.4.1 模块概述

跟单服务消费 `fact_type=TRADE` 的标准账户事实，根据用户配置执行规则判断，通过后构建交易，调用本地签名器或外部钱包服务完成签名，并发送至链或协议。外部钱包服务包括 MPC 钱包、第三方托管钱包和其他远程签名服务。

#### 4.4.2 功能说明

- 交易事实订阅、过滤、时效和幂等检查。
- 跟单配置和版本管理。
- 市场、方向、动作、金额、敞口和滑点规则。
- 固定数量、固定金额或资金比例 sizing。
- 协议交易构建。
- 钱包能力查询、签名、签名并发送。
- 交易提交、状态查询、回执和失败记录。

#### 4.4.3 业务流程

```text
TRADE fact → subscription → rule_engine → sizing → trade_intent
trade_intent → transaction → wallet → submission → execution_record
```

#### 4.4.4 输入与输出

输入：

- 标准交易事实。
- 跟单配置及版本。
- 钱包引用和能力。
- 市场报价、手续费和链状态。

输出：

- 规则通过或拒绝结果。
- 待签名交易和内容哈希。
- 提交结果、交易引用和回执。
- 完整跟单执行记录。

#### 4.4.5 核心处理逻辑

1. 使用 `fact_id + revision + copy_config_id` 幂等创建执行项。
2. 检查事件年龄和确认状态。
3. 匹配市场、方向和动作范围。
4. 计算目标数量并校验单笔和总敞口。
5. 获取报价并校验滑点。
6. 生成 `TradeIntent`，调用协议 `TransactionBuilder`。
7. 根据钱包能力调用 `sign`、`submit` 或 `signAndSubmit`。
8. 保存提交引用并查询最终状态。
9. 超时或响应丢失时先查询，不直接重签或重发。

#### 4.4.6 状态及状态转换

```text
RECEIVED → RULE_REJECTED
         → APPROVED → SIGNING → SIGNED → SUBMITTING → SUBMITTED → CONFIRMED
                              └→ FAILED              └→ FAILED
         → EXPIRED
```

`RULE_REJECTED`、`CONFIRMED`、`FAILED` 和 `EXPIRED` 为终态。未知提交状态不是失败终态。

#### 4.4.7 异常处理

| 异常 | 处理 |
| --- | --- |
| 重复事实 | 返回已有执行项 |
| 事件过期 | 标记 `EXPIRED`，不创建交易 |
| 规则拒绝 | 记录拒绝原因，不重试 |
| 用户拒签 | 标记失败，不重试 |
| 钱包/网络超时 | 查询幂等键或交易引用后决定 |
| nonce/序号冲突 | 重新读取状态、构建交易并复查时效 |
| 报价过期 | 重新报价并再次检查滑点 |
| 历史修正 | 默认不补单；仅显式补偿配置允许处理 |

#### 4.4.8 权限控制

- 用户只能管理自己的跟单配置和钱包引用。
- 只有钱包适配进程拥有签名权限。
- 钱包适配器不得返回或记录私钥和助记词。
- 本地私钥保存在 KMS/HSM 或专用密钥服务，业务数据库只保存引用。
- 签名请求限制链、钱包、目标协议和金额范围。
- 配置变更、钱包绑定、签名和发送均写审计日志。

## 5. 接口详细设计

### 5.1 接口设计规范

- 外部路径使用 `/api/v1`。
- 内部路径使用 `/internal/v1`，仅服务身份可访问。
- 时间统一为 UTC ISO 8601，精度至少毫秒。
- 金额、价格和数量在 JSON 中使用十进制字符串。
- 列表接口使用游标分页，不使用大 offset。
- 写接口接受 `Idempotency-Key`。
- 请求和响应使用明确 schema 版本。
- 错误响应统一包含 `code`、`message`、`trace_id` 和可选 `details`。
- API DTO 不复用数据库 row 或第三方 SDK 类型。

通用响应：

```json
{
  "data": {},
  "meta": {
    "trace_id": "trace_...",
    "schema_version": 1
  }
}
```

账户分析响应额外包含：

```json
{
  "data": {},
  "calculation_version": "pnl-v1/metrics-v1",
  "data_watermark": "source-position",
  "tracking_started_at": "2026-09-01T00:00:00Z",
  "coverage_start": "2026-09-01T00:00:00Z",
  "calculated_to": "2026-09-04T01:02:03Z",
  "history_complete": false,
  "data_quality": "HEALTHY",
  "recalculation_status": "IDLE",
  "updated_at": "2026-09-04T01:02:03Z"
}
```

### 5.2 内部接口

| 方法与路径 | 调用方 | 用途 |
| --- | --- | --- |
| `POST /internal/v1/repair-requests` | 审核服务 | 提交幂等修复请求 |
| `GET /internal/v1/repair-requests/{id}` | 审核服务 | 查询修复执行状态 |
| `POST /internal/v1/recalculations` | 交易日志/审核/管理端 | 创建账户重算任务 |
| `GET /internal/v1/recalculations/{id}` | 内部服务 | 查询重算状态 |
| `GET /internal/v1/collection-ranges` | 审核服务 | 读取采集范围和检查点 |
| `GET /internal/v1/account-facts` | 审核服务 | 读取事实版本和证据 |

### 5.3 外部接口

外部 API 分为交易日志查询、账户分析查询、审核管理和跟单管理。所有外部接口通过查询网关进入，由所属服务完成授权和业务处理。

### 5.4 HTTP/REST API

#### 交易日志

| 方法与路径 | 用途 |
| --- | --- |
| `POST /api/v1/watch-accounts` | 添加关注账户和起始范围 |
| `DELETE /api/v1/watch-accounts/{id}` | 停止关注账户 |
| `GET /api/v1/trade-events` | 查询标准交易事实 |
| `GET /api/v1/trade-events/{fact_id}` | 查询当前版本及版本历史 |
| `GET /api/v1/collection-status` | 查询采集水位、延迟和来源健康 |

#### 账户分析

| 方法与路径 | 用途 |
| --- | --- |
| `GET /api/v1/accounts/{account_key}/overview` | 账户价值、盈亏和基础指标 |
| `GET /api/v1/accounts/{account_key}/positions` | 当前合约持仓和现货库存 |
| `GET /api/v1/accounts/{account_key}/trades` | 合约周期和现货卖出操作 |
| `GET /api/v1/accounts/{account_key}/analysis-results` | 按计算器、版本和统计范围读取结果 |
| `GET /api/v1/accounts` | 按同口径指标筛选和排序账户 |
| `GET /api/v1/calculation-profiles` | 查询可用计算方案和版本 |

`overview.data`：

```json
{
  "account_key": "namespace:network:protocol:account",
  "valuation_currency": "USD",
  "account_value": "10000.00",
  "realized_net_pnl": "800.00",
  "unrealized_pnl": "200.00",
  "total_net_pnl": "1000.00",
  "estimated_pnl": "50.00",
  "volume_usd": "250000.00",
  "fill_count": 120,
  "estimated_trade_count": 1,
  "derivative_metrics": {
    "trade_count": 20,
    "winning_trade_count": 12,
    "losing_trade_count": 7,
    "breakeven_trade_count": 1,
    "win_rate": "0.6315789474"
  },
  "spot_metrics": {
    "trade_count": 5,
    "winning_trade_count": 3,
    "losing_trade_count": 2,
    "breakeven_trade_count": 0,
    "win_rate": "0.6000000000"
  }
}
```

#### 数据审核

| 方法与路径 | 用途 |
| --- | --- |
| `POST /api/v1/audits` | 创建专项审核 |
| `GET /api/v1/audits/{id}` | 查询审核任务 |
| `GET /api/v1/findings` | 查询异常 |
| `POST /api/v1/findings/{id}/repair` | 确认修复 |
| `GET /api/v1/data-quality` | 查询范围或账户质量 |

#### 跟单

| 方法与路径 | 用途 |
| --- | --- |
| `POST /api/v1/copy-configs` | 创建跟单配置 |
| `PUT /api/v1/copy-configs/{id}` | 创建新配置版本 |
| `POST /api/v1/copy-configs/{id}/enable` | 启用配置 |
| `POST /api/v1/copy-configs/{id}/disable` | 停用配置 |
| `GET /api/v1/copy-executions` | 查询执行列表 |
| `GET /api/v1/copy-executions/{id}` | 查询执行详情 |

### 5.5 RPC 接口

基准实现不强制使用 RPC。内部 HTTP 在吞吐或类型约束不能满足要求时可替换为 gRPC，但必须保持：

- 相同业务语义和幂等键。
- 相同鉴权范围。
- protobuf schema 独立版本管理。
- transport 类型不得进入业务模型。

### 5.6 消息接口

| Topic | 分区键 | 生产方 | 消费方 | 保留策略 |
| --- | --- | --- | --- | --- |
| `account.fact.v1` | `account_key` | 交易日志 | 账户分析、审核；跟单过滤 `TRADE` | 不少于最长重放周期 |
| `trade.repair.requested.v1` | `chain_id` | 审核 | 交易日志 | 30 天以上 |
| `account.recalculation.requested.v1` | `account_key` | 交易日志/审核 | 账户分析 | 30 天以上 |

消费规则：

- 每个消费者组维护独立 offset。
- 业务写入和消费幂等记录同事务提交后确认消息。
- 连续失败进入死信主题，保留原消息、错误码和次数。
- 失败消息必须重放或由授权人员关闭，不得静默跳过。

### 5.7 接口数据结构

#### 事件信封

```json
{
  "event_id": "evt_...",
  "event_type": "account.fact.v1",
  "schema_version": 1,
  "occurred_at": "2026-09-04T01:02:03.456Z",
  "observed_at": "2026-09-04T01:02:04.120Z",
  "published_at": "2026-09-04T01:02:04.180Z",
  "producer": "trade-log-service",
  "trace_id": "trace_...",
  "fact": {}
}
```

#### 事实公共结构

```json
{
  "fact_id": "stable-fact-id",
  "fact_type": "TRADE",
  "revision": 1,
  "change_type": "UPSERT",
  "confirmation_status": "OBSERVED",
  "chain_id": "namespace:network",
  "protocol": "protocol-name",
  "account": "chain-account",
  "account_key": "namespace:network:protocol-name:chain-account",
  "ordering_key": "source-defined-position",
  "sub_index": 0,
  "source": "THIRD_PARTY",
  "source_ref": "provider-event-id",
  "raw_log_id": "raw_...",
  "payload": {}
}
```

`fact_type` 支持 `TRADE`、`TRANSFER`、`FEE`、`FUNDING`、`REWARD`、`LIQUIDATION_FEE` 和 `ACCOUNT_SNAPSHOT`。

#### 交易 payload

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
  "extension": {}
}
```

`trigger_type` 为 `USER`、`PROTOCOL` 或 `LIQUIDATION`；`position_effect` 为 `OPEN`、`INCREASE`、`DECREASE`、`CLOSE`、`REVERSE`、`NONE` 或 `UNKNOWN`。

#### 资金与费用 payload

```json
{
  "asset": "USDC",
  "amount": "12.34000000",
  "direction": "CREDIT",
  "funding_type": "FUNDING",
  "related_market": "protocol:ETH-USD",
  "related_position_ref": null,
  "related_trade_fact_id": null,
  "transaction_hash": "source-transaction-id",
  "extension": {}
}
```

#### 账户快照 payload

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

### 5.8 接口错误码

| 错误码 | HTTP | 含义 |
| --- | ---: | --- |
| `VALIDATION_ERROR` | 400 | 参数或事实字段无效 |
| `AUTHENTICATION_REQUIRED` | 401 | 未认证 |
| `PERMISSION_DENIED` | 403 | 权限不足 |
| `RESOURCE_NOT_FOUND` | 404 | 资源不存在 |
| `VERSION_CONFLICT` | 409 | revision、配置或状态版本冲突 |
| `IDEMPOTENCY_CONFLICT` | 409 | 相同幂等键对应不同请求 |
| `UNSUPPORTED_PROTOCOL` | 422 | 不支持的协议或语义 |
| `INCOMPLETE_DATA` | 422 | 输入不足，无法可靠计算 |
| `RATE_LIMITED` | 429 | 超过限流 |
| `DEPENDENCY_UNAVAILABLE` | 503 | 外部依赖暂不可用 |
| `UNKNOWN_SUBMISSION` | 503 | 交易提交结果未知，需要查询 |
| `INTERNAL_INVARIANT_VIOLATION` | 500 | 内部状态不变量失败 |

### 5.9 接口鉴权

- 外部 API 使用用户身份令牌。
- 内部 API 使用双向 TLS 或短期服务凭证。
- 授权维度包括用户、操作、账户范围、服务身份和管理员权限。
- 钱包操作需要跟单配置所有权和钱包绑定权限。
- 修复、重算、计算版本发布和人工异常处理使用独立权限。
- 所有管理和签名操作记录审计日志。

### 5.10 接口调用时序

#### 实时分析与跟单

```text
数据源 → trade-log: 原始数据
trade-log → PostgreSQL: raw log + fact revision + outbox
trade-log → account.fact.v1: AccountFact
account.fact.v1 → account-analysis: 全部事实
account-analysis → analysis DB: 账本、计算结果和水位
account.fact.v1 → copy-trading: 仅 TRADE
copy-trading → wallet: 签名/发送
copy-trading → chain/protocol: 提交或查询
```

#### 审核修复

```text
data-audit → trade-log/read source: 对账
data-audit → trade-log: RepairRequested
trade-log → account.fact.v1: 新 revision
account-analysis → account-analysis: 重算
data-audit → data-audit: 复验并关闭 finding
```

## 6. 数据库详细设计

### 6.1 数据库总体设计

- 四个业务服务分别拥有独立数据库 schema 或独立实例。
- 禁止跨服务直接写表。
- 交易日志库保存事实源和 outbox。
- 分析库保存可重建账本、状态、过程和结果。
- 审核库保存任务、异常、证据和修复状态。
- 跟单库保存配置、钱包引用和执行状态。
- 金额字段使用 `NUMERIC(78,30)` 或按资产精度收窄。
- 所有时间字段使用 `TIMESTAMPTZ`。
- 主键默认 UUID/ULID；稳定链上 ID 使用 TEXT。
- JSONB 只承载原始 payload、版本化 payload 和扩展字段。

### 6.2 ER 模型

```text
trade-log:
collection_jobs ──< raw_logs ──< account_fact_versions >── account_facts_current
collection_jobs ── collection_checkpoints
account_fact_versions ── outbox_events

account-analysis:
account_ledger_events ──> derivative_episodes
account_ledger_events ──> spot_inventories ──< spot_disposals
account_ledger_events ──> account_snapshots
calculation_results ──> account_metrics
recalculation_jobs ──> calculation_results

data-audit:
audit_jobs ──< audit_findings ──< repair_requests

copy-trading:
copy_configs ──< copy_config_versions ──< copy_executions
wallet_refs ──< copy_executions ──< wallet_requests ──< transaction_receipts
```

跨服务关系使用稳定 ID，不建立跨数据库外键。

### 6.3 数据表设计

#### 6.3.1 表结构

##### 交易日志服务

###### `collection_jobs`

| 字段 | 类型 | 约束/说明 |
| --- | --- | --- |
| `id` | UUID | PK |
| `chain_id/protocol/source_id` | TEXT | NOT NULL |
| `mode` | TEXT | `HISTORICAL/REALTIME/REPAIR` |
| `range_start/range_end` | JSONB | 可空 |
| `cursor` | TEXT | 可空 |
| `status` | TEXT | NOT NULL |
| `attempts` | INT | 默认 0 |
| `heartbeat_at` | TIMESTAMPTZ | 可空 |
| `last_error` | TEXT | 脱敏摘要 |
| `created_at/updated_at` | TIMESTAMPTZ | NOT NULL |

索引：`(status, updated_at)`、`(chain_id, protocol, source_id)`。

###### `collection_checkpoints`

| 字段 | 类型 | 约束/说明 |
| --- | --- | --- |
| `chain_id/source_id/partition_key` | TEXT | 联合 PK |
| `position` | JSONB | NOT NULL |
| `cursor` | TEXT | 可空 |
| `finality_proof` | JSONB | 可空 |
| `updated_at` | TIMESTAMPTZ | NOT NULL |

###### `raw_logs`

| 字段 | 类型 | 约束/说明 |
| --- | --- | --- |
| `id` | UUID | PK |
| `chain_id/protocol/source_id` | TEXT | NOT NULL |
| `source_event_id` | TEXT | NOT NULL |
| `chain_position` | JSONB | NOT NULL |
| `ordering_key` | TEXT | NOT NULL |
| `transaction_ref` | TEXT | 可空 |
| `payload` | JSONB | NOT NULL |
| `observed_at` | TIMESTAMPTZ | NOT NULL |
| `parse_status/parser_version` | TEXT | NOT NULL |

唯一约束：`(chain_id, source_id, source_event_id)`；索引：`(chain_id, ordering_key)`、`(parse_status, observed_at)`。

###### `account_fact_versions`

| 字段 | 类型 | 约束/说明 |
| --- | --- | --- |
| `fact_id/revision` | TEXT/INT | 联合 PK |
| `event_id` | TEXT | UNIQUE, NOT NULL |
| `fact_type/account_key` | TEXT | NOT NULL |
| `ordering_key/sub_index` | TEXT/INT | NOT NULL |
| `change_type/confirmation_status` | TEXT | NOT NULL |
| `occurred_at` | TIMESTAMPTZ | NOT NULL |
| `payload` | JSONB | NOT NULL |
| `raw_log_id` | UUID | FK `raw_logs(id)` |
| `parser_version` | TEXT | NOT NULL |
| `created_at` | TIMESTAMPTZ | NOT NULL |

索引：`(account_key, ordering_key, sub_index)`、`(fact_type, occurred_at)`。

###### `account_facts_current`

| 字段 | 类型 | 约束/说明 |
| --- | --- | --- |
| `fact_id` | TEXT | PK |
| `current_revision` | INT | NOT NULL |
| `account_key/fact_type` | TEXT | NOT NULL |
| `ordering_key/sub_index` | TEXT/INT | NOT NULL |
| `is_retracted` | BOOL | NOT NULL |
| `updated_at` | TIMESTAMPTZ | NOT NULL |

外键 `(fact_id,current_revision)` 指向事实版本。revision 更新使用乐观锁，新版本必须为当前版本加一。

###### `outbox_events`

| 字段 | 类型 | 约束/说明 |
| --- | --- | --- |
| `event_id` | TEXT | PK |
| `topic/partition_key` | TEXT | NOT NULL |
| `payload` | JSONB | NOT NULL |
| `status` | TEXT | NOT NULL |
| `attempts` | INT | 默认 0 |
| `next_retry_at/published_at` | TIMESTAMPTZ | 可空 |
| `created_at` | TIMESTAMPTZ | NOT NULL |

索引：`(status, next_retry_at, created_at)`。

##### 账户分析服务

###### `account_ledger_events`

| 字段 | 类型 | 约束/说明 |
| --- | --- | --- |
| `fact_id/revision` | TEXT/INT | 联合 PK |
| `event_id` | TEXT | UNIQUE, 消费幂等 |
| `account_key/fact_type` | TEXT | NOT NULL |
| `ordering_key/sub_index` | TEXT/INT | NOT NULL |
| `change_type/is_current` | TEXT/BOOL | NOT NULL |
| `occurred_at` | TIMESTAMPTZ | NOT NULL |
| `payload` | JSONB | NOT NULL |
| `consumed_at` | TIMESTAMPTZ | NOT NULL |

索引：`(account_key, is_current, ordering_key, sub_index)`。

###### `positions_current`

| 字段 | 类型 | 约束/说明 |
| --- | --- | --- |
| `account_key/market` | TEXT | 联合业务键 |
| `side` | TEXT | LONG/SHORT |
| `quantity/average_cost` | NUMERIC | NOT NULL |
| `unrealized_pnl` | NUMERIC | 可空 |
| `episode_id` | UUID | 可空 |
| `is_estimated` | BOOL | NOT NULL |
| `calculation_version/data_watermark` | TEXT | 联合 PK 组成部分 |
| `updated_at` | TIMESTAMPTZ | NOT NULL |

###### `derivative_episodes`

| 字段 | 类型 | 约束/说明 |
| --- | --- | --- |
| `episode_id` | UUID | PK |
| `account_key/market/side` | TEXT | NOT NULL |
| `status` | TEXT | OPEN/CLOSED |
| `opened_at/closed_at` | TIMESTAMPTZ | `closed_at` 可空 |
| `open_quantity/closed_quantity` | NUMERIC | NOT NULL |
| `realized_pnl/fees/funding/net_pnl` | NUMERIC | NOT NULL |
| `calculation_method` | TEXT | NOT NULL |
| `is_estimated` | BOOL | NOT NULL |
| `calculation_version/data_watermark` | TEXT | NOT NULL |

索引：`(account_key, market, status, opened_at)`。

###### `spot_inventories`

| 字段 | 类型 | 约束/说明 |
| --- | --- | --- |
| `account_key/asset/calculation_version` | TEXT | 联合 PK |
| `quantity/total_cost/average_cost` | NUMERIC | NOT NULL |
| `estimated_quantity` | NUMERIC | NOT NULL |
| `data_watermark` | TEXT | NOT NULL |
| `updated_at` | TIMESTAMPTZ | NOT NULL |

###### `spot_disposals`

| 字段 | 类型 | 约束/说明 |
| --- | --- | --- |
| `disposal_id` | UUID | PK |
| `account_key/market` | TEXT | NOT NULL |
| `order_id/operation_id` | TEXT | 至少一个可归并标识 |
| `quantity/proceeds/allocated_cost/fees/net_pnl` | NUMERIC | NOT NULL |
| `calculation_method` | TEXT | NOT NULL |
| `is_estimated` | BOOL | NOT NULL |
| `calculation_version/data_watermark` | TEXT | NOT NULL |

唯一约束：`(account_key, calculation_version, operation_id)`；协议只有 order ID 时使用规范化 operation ID。

###### `account_snapshots`

保存账户余额、持仓、估值、价格证据、事实水位、计算版本、质量状态和创建时间。主键 `(account_key, calculation_version, data_watermark)`。

###### `calculation_results`

| 字段 | 类型 | 约束/说明 |
| --- | --- | --- |
| `result_id` | UUID | PK |
| `account_key/calculator_id/calculator_version` | TEXT | NOT NULL |
| `data_watermark` | TEXT | NOT NULL |
| `output_schema_version` | INT | NOT NULL |
| `payload/dependency_versions` | JSONB | NOT NULL |
| `input_hash` | TEXT | NOT NULL |
| `quality_status` | TEXT | NOT NULL |
| `created_at` | TIMESTAMPTZ | NOT NULL |

唯一约束：`(account_key, calculator_id, calculator_version, data_watermark, input_hash)`。

###### `account_metrics`

保存当前查询投影：账户价值、已实现/未实现/总盈亏、估算盈亏、USD 成交量、fill 数、合约指标、现货指标、覆盖范围、质量状态、计算版本和水位。主键 `(account_key, calculation_version, data_watermark)`。

###### `recalculation_jobs`

保存任务 ID、账户、最早影响位置、目标计算版本、状态、租约、进度、尝试次数、错误和时间。索引 `(status, next_retry_at)`；同一账户的待处理任务合并到最早位置。

##### 数据完整审核服务

###### `audit_jobs`

保存审核类型、链、协议、来源、范围、状态、租约、进度、统计、错误和时间。索引 `(status, scheduled_at)`。

###### `audit_findings`

保存 finding ID、审核任务、异常类型、影响范围、严重级别、证据 JSON、状态、首次/最近发现时间和关闭信息。唯一键由异常类型和规范化影响范围生成。

###### `repair_requests`

保存请求 ID、finding ID、目标服务、修复范围、预期操作、幂等键、状态、响应、复验状态和时间。`idempotency_key` 唯一。

##### 跟单服务

###### `copy_configs` 与 `copy_config_versions`

`copy_configs` 保存配置 ID、用户、当前版本、启用状态和时间。`copy_config_versions` 保存来源账户、钱包引用、市场/方向/动作、sizing、事件年龄、确认状态、单笔上限、总敞口、滑点和创建时间；主键 `(config_id,version)`。

###### `wallet_refs`

保存钱包引用、用户、链、地址、`provider_type`、`custody_model`、`signing_technology`、外部资源 ID、状态和非敏感配置。`provider_type` 区分本地签名器和外部钱包服务；`custody_model` 描述本地、共同管理或第三方托管；`signing_technology` 可标记 MPC、HSM 或供应商定义方式。敏感凭证只保存 KMS/HSM 或供应商引用。

###### `copy_executions`

保存执行 ID、幂等键、事实 ID/revision、配置及版本、决策、拒绝原因、目标参数、钱包引用、状态和时间。`idempotency_key` 唯一。

###### `wallet_requests`

保存请求 ID、执行 ID、动作、内容哈希、幂等键、尝试次数、交易引用、状态和脱敏错误。`idempotency_key` 唯一。

###### `transaction_receipts`

保存回执 ID、执行 ID、交易引用、确认状态、链位置、实际数量、价格、费用、失败原因和时间。交易引用建立唯一或条件唯一索引。

#### 6.3.2 字段定义

各表字段、类型、空值要求和业务说明已在 6.3.1 对应表结构中定义。组合字段写法如 `chain_id/protocol/source_id` 表示三个独立数据库字段，不表示单个复合字段；编码前的 SQL migration 必须展开为独立列。

#### 6.3.3 主键与外键

- 服务内稳定父子关系使用外键；跨服务只保存稳定 ID。
- 事实版本以 `(fact_id, revision)` 为主键，current 表通过复合外键指向版本表。
- 计算过程、审核任务和跟单执行使用 UUID 主键。
- 用户、钱包、配置、执行和回执只在跟单服务数据库内建立外键。
- 删除事实和计算历史不使用级联删除；业务失效通过 revision、状态或 current 指针表达。

#### 6.3.4 索引设计

- 高频 current 查询使用单独投影或 current 指针，不在版本表执行全表 `MAX(revision)`。
- 账户顺序读取索引统一以 `account_key, ordering_key, sub_index` 开头。
- worker 任务表使用 `(status, next_retry_at)` 或 `(status, scheduled_at)` 索引。
- 账户交易查询使用 `(account_key, market, opened_at/occurred_at)` 索引。
- JSONB 只对已确定的查询条件创建表达式或 GIN 索引，不为整个 payload 默认建立索引。

#### 6.3.5 约束设计

- 状态字段使用 CHECK 约束或 PostgreSQL enum，并随 migration 管理。
- 数量必须非负，方向由独立字段表达。
- revision 必须大于零并连续递增。
- `total_net_pnl` 如持久化，必须由写入服务校验等于已实现与未实现盈亏之和。
- JSON payload 必须在应用写入前通过对应 schema 校验。

### 6.4 数据字典

| 字段 | 定义 |
| --- | --- |
| `fact_id` | 同一事实跨 revision 稳定的 ID |
| `event_id` | 一次消息版本的唯一 ID |
| `ordering_key` | 协议账户内可稳定比较的位置 |
| `sub_index` | 同一来源位置拆分事实的顺序 |
| `data_watermark` | 计算已覆盖的确定事实位置 |
| `calculation_version` | 计算器组合及规则版本 |
| `input_hash` | 计算输入和依赖版本的摘要 |
| `is_estimated` | 是否使用估算起始成本或估算数据 |
| `history_complete` | 是否覆盖协议可验证范围内的完整历史 |
| `quality_status` | `UNKNOWN/HEALTHY/DEGRADED/INVALID` 等质量状态 |

### 6.5 数据关系

- 一个原始日志可以解析为多个标准账户事实。
- 一个事实拥有多个 revision，但只有一个 current 版本。
- 一个账户事实可以影响一个持仓周期、库存批次或账户快照。
- 一个账户可拥有多个计算版本和水位结果，但查询只指向一个 ACTIVE 版本。
- 一个交易事实可命中多个跟单配置，每个配置产生独立执行项。
- 一个审核异常可产生多个修复尝试，但每个幂等修复请求只有一条记录。

### 6.6 数据初始化

- 创建各服务 schema、角色和最小权限。
- 初始化事件 schema 版本、默认计算器描述和默认计算方案。
- 初始化支持的链、协议能力和资产精度元数据。
- 不初始化具体用户跟单配置或钱包凭证。
- 初始化数据必须通过版本化 migration 完成，不在应用启动时静默修改表结构。

### 6.7 数据迁移与版本管理

- migration 归属服务目录，按单调版本执行。
- 迁移必须向前兼容：先部署兼容新旧结构的读取方，再切换写入方。
- 删除列、改名或收紧约束分多次发布完成。
- 事件 schema 删除、改名或语义变化必须升级主版本和 topic。
- 计算规则语义变化必须发布新计算器版本，旧结果不可原地修改。
- 发布前执行 migration 前向测试、数据回填演练和回滚方案检查。

## 7. 核心流程设计

### 7.1 核心业务流程

#### 交易事实主流程

```text
数据获取
→ 原始日志幂等保存
→ 协议解析
→ AccountFact 标准化
→ 轻量校验
→ 事实版本 + current + outbox 原子提交
→ 发布 account.fact.v1
```

#### 账户分析主流程

```text
消费 AccountFact
→ 保存账本和消费幂等
→ 判断增量或重算
→ 现货/合约分类计算
→ 账户状态与估值
→ 盈亏与指标
→ 保存不可变结果
→ 切换当前查询投影
```

#### 跟单主流程

```text
消费 TRADE
→ 时效与幂等
→ 规则判断
→ sizing
→ 构建交易
→ 钱包签名
→ 发送/查询
→ 保存最终状态
```

### 7.2 系统交互时序

#### 正常实时处理

```text
Source        TradeLog        EventBus       Analyzer       CopyTrader       Wallet
  | raw/log      |                |               |               |               |
  |------------->|                |               |               |               |
  |              | fact+outbox    |               |               |               |
  |              |--------------->|               |               |               |
  |              |                | AccountFact   |               |               |
  |              |                |-------------->|               |               |
  |              |                | TRADE         |               |               |
  |              |                |------------------------------>|               |
  |              |                |               | calculate     |               |
  |              |                |               |-----> DB      | rules/build   |
  |              |                |               |               |-------------->|
  |              |                |               |               |<--------------|
```

#### 修复与重算

```text
Auditor → TradeLog: RepairRequested
TradeLog → EventBus: corrected AccountFact revision
EventBus → Analyzer: corrected revision
Analyzer → Recalculator: recalculation job
Recalculator → ResultStore: validated new result version
Auditor → Auditor: verify and close finding
```

### 7.3 异步处理流程

- outbox publisher：领取待发布记录、发布、确认或退避。
- account analyzer：按 `account_key` 分区消费事实。
- recalculator：按账户租约执行回放和版本切换。
- data auditor：按链、来源和范围执行检查。
- copy trader：按来源账户或配置分片执行规则和交易。
- dead-letter processor：隔离失败消息并提供重放入口。

所有异步 worker 必须具有：

- 持久化进度或消息 offset。
- 稳定幂等键。
- 有界重试和退避。
- 心跳、租约或消费者组所有权。
- 可取消和优雅退出。
- 延迟、积压和失败指标。

### 7.4 定时任务流程

| 任务 | 调度粒度 | 处理内容 |
| --- | --- | --- |
| 采集健康检查 | 分钟级 | 心跳、水位、延迟和错误率 |
| 范围完整审核 | 小时级或链级配置 | 连续性、分页和数量异常 |
| 来源对账 | 抽样持续、全量按配置 | 第三方与权威/备用源比较 |
| 链最终性更新 | 随链状态或周期执行 | 确认、重组和撤销 |
| 账户业务一致性 | 计算后及周期执行 | 持仓、库存、余额不变量 |
| 过期任务回收 | 分钟级 | 释放失效租约并恢复任务 |
| 历史数据归档 | 日级或容量触发 | 原始日志和旧结果归档 |

调度频率属于运行配置，不能写死在业务模块中。

### 7.5 状态机设计

#### 事实版本

```text
OBSERVED → CONFIRMED
    └────→ ORPHANED/RETRACTED
```

#### 计算版本

```text
BUILDING → VALIDATING → ACTIVE → SUPERSEDED
    └───────────────→ INVALID
```

#### 修复异常

```text
OPEN → REPAIRING → VERIFYING → CLOSED
          └───────────────→ OPEN
```

#### 跟单执行

```text
RECEIVED → RULE_REJECTED
         → APPROVED → SIGNING → SIGNED → SUBMITTING → SUBMITTED → CONFIRMED
                              └→ FAILED              └→ FAILED
         → EXPIRED
```

状态迁移使用乐观锁或条件更新，非法迁移返回 `VERSION_CONFLICT` 或内部不变量错误。

### 7.6 异常流程

#### 数据迟到或修正

1. 保存事实新 revision。
2. 定位最早受影响顺序。
3. 创建或合并账户重算任务。
4. 从前一个有效快照回放。
5. 验证事件数量、持仓、余额和结果水位。
6. 原子切换当前计算版本。

#### 外部依赖超时

1. 根据错误类型判断是否可重试。
2. 可重试错误进入有界退避。
3. 涉及钱包发送时先查询幂等键或交易引用。
4. 超过阈值后进入失败或人工处理状态并告警。

#### 不完整数据

- 可估算的起始成本使用 `ESTIMATED` 并明确披露。
- 缺少必要价格或无法闭合状态时输出 `UNAVAILABLE`，不使用零值。
- 审核服务创建 finding，修复后触发重算。

## 8. 程序结构设计

### 8.1 工程目录结构

```text
robotech/
├── Cargo.toml
├── Cargo.lock
├── rust-toolchain.toml
├── crates/
│   ├── shared-types/
│   ├── account-facts/
│   ├── internal-events/
│   └── protocol-api/
├── services/
│   ├── trade-log/
│   │   ├── Cargo.toml
│   │   ├── src/{功能模块}/
│   │   ├── adapters/{Cargo.toml,src/}
│   │   ├── bins/
│   │   │   ├── trade-collector/
│   │   │   └── trade-parser-publisher/
│   │   └── migrations/
│   ├── account-analysis/
│   │   ├── Cargo.toml
│   │   ├── src/{功能模块}/
│   │   ├── adapters/{Cargo.toml,src/}
│   │   ├── bins/
│   │   │   ├── account-analyzer/
│   │   │   └── account-recalculator/
│   │   └── migrations/
│   ├── data-audit/
│   │   ├── Cargo.toml
│   │   ├── src/{功能模块}/
│   │   ├── adapters/{Cargo.toml,src/}
│   │   ├── bins/data-auditor/
│   │   └── migrations/
│   └── copy-trading/
│       ├── Cargo.toml
│       ├── src/{功能模块}/
│       ├── adapters/{Cargo.toml,src/}
│       ├── bins/copy-trader/
│       └── migrations/
├── protocols/implementations/{protocol}/
├── gateway/query-api/
├── config/
├── tests/{contract-tests,integration-tests,replay-tests,fixtures}/
└── docs/
```

每个服务目录内聚业务功能、适配器、运行入口和 migration。服务业务 crate、adapters crate 和各 composition binary crate 均为 workspace member。

### 8.2 Package/Module 划分

#### 公共 package

| Package | 职责 |
| --- | --- |
| `shared-types` | 定点金额、价格、数量、链、协议、账户、市场和版本值对象 |
| `account-facts` | 标准账户事实 Rust 类型和 schema |
| `internal-events` | 修复、重算和审核通知 |
| `protocol-api` | 协议描述、能力枚举和注册键 |

#### 服务 package

| Package | 主要 module |
| --- | --- |
| `trade-log` | acquisition、checkpoint、raw_log、parsing、normalization、validation、publishing、finality、repair、query |
| `account-analysis` | fact_ledger、derivatives、spot、capital_flow、fees、state、valuation、pnl、metrics、calculation、recalculation、result、query |
| `data-audit` | health、completeness、reconciliation、consistency、findings、repair、verification、query |
| `copy-trading` | subscription、config、rules、sizing、intent、transaction、wallet、submission、execution、query |

功能目录默认为 Rust module。跨服务复用、独立版本、重型 SDK 或必须编译隔离时才拆为 crate。

### 8.3 核心类/结构体设计

```rust
pub struct AccountFactEnvelope { /* 消息信封与 AccountFact */ }
pub struct AccountFact { /* fact_id、revision、账户、顺序和 payload */ }
pub enum AccountFactPayload { Trade(TradeFact), Funding(FundingFact), AccountSnapshot(AccountSnapshotFact) }

pub struct AccountLedger { /* 当前事实视图与水位 */ }
pub struct DerivativePosition { /* 方向、数量、成本和当前周期 */ }
pub struct DerivativeEpisode { /* 开平过程、费用、资金费和净盈亏 */ }
pub struct SpotInventory { /* 数量、总成本和移动平均成本 */ }
pub struct SpotDisposal { /* 卖出收入、分摊成本和净盈亏 */ }
pub struct AccountValuation { /* 账户价值和价格证据 */ }
pub struct PnlResult { /* 正式盈亏、方法和估算标记 */ }
pub struct AccountMetrics { /* 总值及现货/合约分类指标 */ }

pub struct AuditJob { /* 审核范围、状态和进度 */ }
pub struct Finding { /* 异常类型、证据和状态 */ }
pub struct RepairRequest { /* 幂等修复范围和动作 */ }

pub struct CopyConfig { /* 当前配置和版本 */ }
pub struct TradeIntent { /* 规则通过后的目标交易 */ }
pub struct CopyExecution { /* 决策与执行状态机 */ }
pub struct WalletRef { /* 非敏感钱包引用 */ }
```

所有金额类型封装定点十进制值和资产精度，不使用 `f32/f64`。

### 8.4 核心方法/函数设计

#### 计算器

```rust
pub trait AccountCalculator: Send + Sync {
    fn descriptor(&self) -> CalculatorDescriptor;
    fn validate(&self, input: &CalculationInput) -> ValidationResult;
    fn calculate(&self, input: &CalculationInput)
        -> Result<CalculationOutput, CalculationError>;
}
```

#### 采集与解析端口

```rust
pub trait SourceReader { /* fetch、subscribe、finality、health */ }
pub trait ProtocolParser { /* RawLog -> Vec<AccountFact> */ }
pub trait FactRepository { /* 事实版本与 current */ }
pub trait FactPublisher { /* AccountFactEnvelope 发布 */ }
```

#### 分析端口

```rust
pub trait LedgerRepository { /* 事实版本、revision 和顺序读取 */ }
pub trait CalculationResultRepository { /* 不可变结果和 current 指针 */ }
pub trait MarketDataReader { /* 版本化价格和汇率 */ }
pub trait MarketDataRecorder { /* 保存价格证据 */ }
pub trait SettlementPolicy { /* 唯一正式结算盈亏 */ }
pub trait ValuationPolicy { /* 协议账户估值 */ }
```

#### 跟单端口

```rust
pub trait TransactionBuilder { /* TradeIntent -> UnsignedTransaction */ }
pub trait WalletGateway { /* capabilities、sign、submit、query */ }
pub trait ChainSubmissionReader { /* 交易状态和回执 */ }
```

trait 由使用方模块定义，具体协议和 adapters 实现。计算器不执行 I/O。

### 8.5 公共组件设计

- 定点十进制金额与精度校验。
- ID、时间、版本和数据水位类型。
- AccountFact schema 编解码和兼容性校验。
- outbox 与消费幂等基础实现。
- 有界重试、退避、租约和任务心跳。
- UTC 时间和 trace 上下文传播。
- 敏感字段脱敏和审计日志接口。
- HTTP 错误信封、游标分页和认证中间件。

公共组件不得包含具体账户分析、审核或跟单规则。

### 8.6 配置项设计

| 类别 | 内容 |
| --- | --- |
| 进程 | 监听地址、并发、队列容量、超时、优雅退出 |
| 数据库 | DSN 密钥引用、连接池、事务超时 |
| 消息 | broker、topic、consumer group、批量、重试和死信 |
| 数据源 | 链/协议、端点、凭证引用、游标和限流 |
| 计算 | 计算器版本、计价币、精度、价格来源和参数 |
| 审核 | 审核类型、频率、范围、阈值和备用来源 |
| 钱包 | 本地签名器或外部钱包服务类型、MPC/托管方式、能力、KMS/供应商引用和允许范围 |

配置加载顺序为默认文件、环境覆盖、密钥引用解析和完整校验。动态跟单配置存放数据库并版本化，不属于静态配置。

### 8.7 依赖库设计

依赖类别限定为：

- Tokio：异步运行时。
- Serde：契约序列化。
- 定点十进制库：金额和价格。
- HTTP 框架：接口和中间件。
- SQLx 或等价 PostgreSQL client：数据库适配器。
- Kafka/Redpanda client：消息适配器。
- OpenTelemetry：可观测性。
- 各协议和钱包 SDK：仅位于协议实现或对应 adapters crate。

依赖版本在 workspace 统一管理并提交锁文件。业务 crate 禁止直接依赖数据库、消息和协议 SDK。生产依赖执行漏洞、许可证和来源审计。

## 9. 部署设计

### 9.1 部署单元

七个进程分别构建部署：

- `trade-collector`
- `trade-parser-publisher`
- `account-analyzer`
- `account-recalculator`
- `data-auditor`
- `copy-trader`
- `query-api`

每个发布物包含 Git revision、Rust toolchain、事件 schema 和协议模块版本。

### 9.2 数据与权限隔离

- 四个服务使用独立数据库 schema、数据库角色和 migration。
- 进程只获得职责需要的表权限。
- 内部服务使用 mTLS 或短期服务凭证。
- 只有钱包适配进程具有签名权限。
- 密钥通过 secret manager 注入引用，不进入镜像和普通配置。

### 9.3 扩缩容

- 采集按链、协议和来源分片。
- 解析按原始日志 ID 分片。
- 账户分析按 `account_key` 分区并保持账户内串行。
- 重算按账户租约并行，设置独立资源配额。
- 审核按链、来源和范围分片。
- 跟单按来源账户或用户配置分片，并保持执行幂等。

### 9.4 高可用与恢复

- 所有无状态 API 和 worker 可以运行多个副本。
- Kafka consumer group、数据库租约和唯一约束负责工作归属。
- PostgreSQL 配置备份、时间点恢复和高可用复制。
- 原始日志归档支持重新解析和事实重建。
- 服务重启从检查点、offset、任务进度或交易状态恢复。
- 未知钱包提交状态必须查询确认，不能通过重启触发重复发送。

### 9.5 可观测性与告警

交易日志：

- 采集水位差、接收率、解析失败、死信和 outbox 年龄。
- `observed_at - occurred_at` 和 `published_at - observed_at`。

账户分析：

- 消费积压、计算水位、计算耗时、重算数量和缺失价格。

数据审核：

- 未审核范围、来源差异率、finding 数量、修复时长和复验失败。

跟单：

- 发布到规则、签名、提交的分段延迟。
- 拒绝、过期、签名失败、链上失败和未知提交数量。

指标禁止使用账户地址和交易哈希等高基数字段作为标签。

### 9.6 发布与迁移

1. 执行向前兼容 migration。
2. 部署能够同时读取新旧结构的消费者。
3. 部署新生产者并切换 schema 版本。
4. 观察积压、错误率和新旧结果差异。
5. 切换默认计算版本或查询投影。
6. 稳定后停止旧写入路径；结构清理在后续 migration 完成。

计算规则、解析器和事件 schema 均可灰度。旧事实和计算结果不原地修改。

### 9.7 安全要求

- 镜像使用最小运行用户和只读文件系统。
- 网络策略限制数据库、broker、钱包和链端点访问。
- 私钥和助记词不得出现在日志、消息、错误、数据库普通字段和备份明文中。
- 签名请求限制链、钱包、目标协议和金额。
- 配置变更、修复、重算版本切换、签名和发送记录不可篡改审计日志。

### 9.8 部署验收

- 各进程可独立启动、停止和扩缩容。
- 服务重启能够从持久化位置恢复。
- 重复消息不会产生重复状态和重复交易。
- 审核与重算压力不影响实时发布和跟单 SLO。
- 增量分析与完整回放结果一致。
- 修复能够闭环到新 revision、账户重算和复验。
- 钱包超时和响应丢失不会导致重复广播。
- 监控能够定位数据源、事实、账户、任务和执行 ID。
