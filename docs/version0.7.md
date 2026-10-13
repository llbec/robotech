# v0.7 开发文档：现货成交与候选信号

- 文档版本：0.7
- 状态：已实现，本机开发验证通过；Docker 及真实来源验收待服务器执行
- 编写日期：2026-10-13
- 前置版本：[v0.6 开发文档](version0.6.md)
- 设计依据：[概要设计](overview-design.md)、[详细设计](detailed-design.md)、[路线图](roadmap.md)
- 已更新 [统一接口说明](api.md)，不另建版本接口文档

## 1. 版本目标

在多账户 HTTP 补偿、WebSocket 采集、标准事实存储和 webhook 发布基础上，识别并保存 Hyperliquid 现货成交，支持现货与永续分别查询，对符合条件的主动现货买卖生成候选信号。

沿用 query-api、trade-log-query、trade-collector、trade-parser-publisher 和 PostgreSQL，不新增现货专用进程或第二套存储、去重、租约、发布队列。应用版本升级为 0.7.0。

本版交付成交事实和候选，不计算现货库存成本、收益、胜率，不提交跟单订单。现货成本与盈亏仍归属 v1.5，跟单执行仍归属 v3.x。

## 2. 功能范围

| 能力 | 本版交付 | 边界 |
| --- | --- | --- |
| 现货解析 | 从原有 fills 通道识别现货买卖，建立基础资产和计价资产映射 | 不凭 coin 字符串猜测资产，不把未知市场当作现货 |
| 现货持久化 | HTTP、WS、补录共享事实身份和版本存储 | 重复观察不重复事实，不重建永续身份 |
| 分类查询 | live/stored 支持 PERPETUAL、SPOT、ALL 筛选 | 默认 PERPETUAL，现有调用不突然混入现货 |
| 手续费 | 保存精确金额和实际币种，保留来源字段 | 不把基础资产手续费强行换算为 USDC，不计算净收益 |
| 实时输出 | 主动现货 BUY/SELL 可形成 webhook 候选 | 历史 HTTP、快照、旧归档重解析不产生候选 |
| 升级补录 | 持久记录现货历史补录任务及进度 | 不倒退原 HTTP 水位；受已保存证据及来源历史范围限制 |
| 可观测性 | 分类型计数、能力启用时间、补录状态和异常原因 | 不用扫描完成证明来源历史完整 |
| 验收 | 固定样本、数据库集成、服务器脚本及真实现货对账 | 无实际成交时标记 SKIP，不能记作实时发布通过 |

不增加按账户选择市场订阅、不增加第二个 webhook 目标，不提供账户余额、充值提现、资产划转、现货订单簿或预测市场支持。监控清单仍通过管理接口维护，不恢复 TOML 中的 account、start_time。

## 3. 使用方式

### 3.1 查询方式

继续按功能使用 `/api/v1/trade-events/live` 和 `/api/v1/trade-events/stored`，不恢复 source 参数。instrument_type 是同一查询功能内的市场类型筛选条件，不用于选择实时或库存功能。

| 场景 | 调用方式 | 结果 |
| --- | --- | --- |
| 查询近期现货 | GET /api/v1/trade-events/live?account=…&instrument_type=SPOT | 查询来源、保存全部有效事实，展示现货 |
| 查询已保存现货 | GET /api/v1/trade-events/stored?account=…&instrument_type=SPOT | 仅读数据库，支持时间区间和分页 |
| 查询合约 | instrument_type=PERPETUAL 或省略 | 保持既有默认展示范围 |
| 查询全部成交 | instrument_type=ALL | 展示两类成交，沿用确定排序 |

展示筛选不能影响采集入库范围。请求 SPOT 不得丢弃同一来源响应中的永续事实；请求 PERPETUAL 也不得丢弃已支持的现货事实。

### 3.2 监控与发布

原账户自动支持两类成交采集，无需重新添加账户。添加、列表、详情、暂停、恢复和分账户状态接口沿用 v0.6 功能路径；暂停同时停止该账户的两类新采集和候选准入。

发布仍受 publishing.enabled 总开关控制。新增 publishing.spot_enabled，默认 false；明确开启后，才允许符合实时条件的现货生成候选。开启现货不重置永续发布激活时间和既有队列。

### 3.3 历史范围

升级只扩展后续解析能力，不保证旧水位之前的现货已入库。旧版虽然保存过含现货的原始响应，但此前过滤的现货不能仅凭水位被认定为“已保存”。通过独立补录任务恢复可用历史，并披露未覆盖范围。

## 4. 来源与候选边界

### 4.1 来源与市场映射

复用 userFills、userFillsByTime 和 userFills WS。HTTP 保持 aggregateByTime=false；不另开现货专用 WS 连接。请求继续经过 v0.6 共享预算，补录任务也不能绕过并发和权重限制。

spotMeta 提供 tokens 和 universe；universe 的 tokens 关联基础、计价资产索引，index 标识交易对。来源存在可读名称及 @index 表达，因此必须建立经元数据验证的映射，不能只按斜杠或 @ 前缀判断。依据：[官方现货元数据](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint/spot)。

资产显示名不等于稳定身份。事实扩展字段保存 tokenId、token index、pair index 和来源名称；基础/计价显示名由对应 token 元数据取得，不假定计价资产一定是 USDC。同名 token 或名称变更不能合并为同一资产。

fills 的字段、历史窗口和订阅包络按 [官方 Info 接口](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint) 与 [官方 WS 订阅](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/websocket/subscriptions) 核对。实现和验收保留来源响应及其元数据，公开文档示例不作为真实账户成交证据。

### 4.2 映射失败与刷新

收到无法识别的市场时，先在现有缓存策略内刷新元数据并重新解析。仍无法定位、索引重复、基础/计价 token 缺失或永续/现货识别冲突时，保存证据和具体错误，不伪造资产、不发布候选。

可确定与本版无关的市场记为 unsupported；属于支持范围但元数据不完整时应阻塞对应提交范围，不能作为普通过滤后静默推进水位。已保存且仍可验证的旧映射允许用于事实落库，记录 metadata_stale；候选仍拒绝过期元数据。

### 4.3 主动成交和费用

现货 side=B 表示买入基础资产，side=A 表示卖出基础资产；标准 side 使用 BUY/SELL。数量、价格必须为正，notional=price×quantity，保持十进制字符串，不能使用浮点计算。

主动成交要求支持的普通现货成交语义、已成交数量、可信市场映射及无强制事件标记。强平、系统处置、未知 dir 或相互矛盾的来源字段不得仅因存在 side 而准入。保留来源 dir、liquidation 等原字段供核对。

fee 保留有符号金额，允许返佣负数。fee_asset 优先解析来源 feeToken 并验证身份；缺失时不默认 USDC，事实记录 fee_asset=null 和警告，候选拒绝手续费币种未确认的现货。无法消除同名歧义时也拒绝候选。手续费不从原始成交 quantity 或 notional 中扣减，避免双算。

### 4.4 实时候选条件

沿用 v0.5 的订阅确认、快照边界、监控代次、时钟偏差、成交年龄、TTL、事实首次准入和 outbox 原子事务；按 instrument_type 分支判断可跟单操作。

| 类型 | 允许的操作 | 额外要求 |
| --- | --- | --- |
| PERPETUAL | 原 OPEN/INCREASE/DECREASE/CLOSE/REVERSE | 永续策略不变 |
| SPOT | 主动 BUY/SELL | spot_enabled=true、资产及手续费币种确认、现货能力激活后发生和接收 |

HTTP 查询、定时补偿、手动补录、历史重解析、WS 初始快照和旧代次消息一律不生成现货实时候选。恢复补录后不追发过去成交。候选仍不表示下单授权。

## 5. 数据与接口定义

### 5.1 标准事实与身份

沿用 AccountFact、TradeFact 和 AccountFactEnvelope，不改变现有永续 fact_id、revision、partition_key、event_id 或已持久化 webhook 正文。新增现货采用独立的规范市场身份输入，避免同账户同 tid 的永续和现货碰撞。

| 字段 | 现货规则 |
| --- | --- |
| fact_type / instrument_type | TRADE / SPOT |
| 身份输入 | 沿用 hl-fill-v1 摘要框架；market 输入为 spot:<pair index>，网络和账户仍参与 |
| payload.market | hyperliquid:spot:<pair index>，稳定标识；交易对名称另存扩展字段 |
| base_asset / quote_asset | 元数据确认的显示名；tokenId 和 index 同时保存 |
| action / side | BUY 或 SELL；不伪造开仓、平仓或反转 |
| position_effect | NOT_APPLICABLE |
| trigger_type / copy_eligible | 根据来源语义和完整映射判断 USER/UNKNOWN、true/false；不能将全部现货标成可跟单 |
| price / quantity / notional | 十进制字符串；notional 为计价资产金额 |
| fee / fee_asset | 来源有符号费用和实际币种；未知币种为 null 并警告 |
| reported_realized_pnl 及其口径字段 | null；不使用现货 closedPnl=0 冒充已计算盈亏 |
| extension | 来源 fill、索引、pair index/name、base/quote tokenId/index、费用资产身份、解析版本及元数据证据引用 |

可读名称和 @index 指向同一交易对时必须得到同一身份；同身份、同语义重复只增加 observation，不新建事实。市场分类或资产身份不确定时，不先生成临时 fact_id 再静默改写。

### 5.2 查询请求参数

以下是新增参数，其余 account、limit、时间及 cursor 沿用 api.md。完整参数表已同步到统一文档。

| 接口 | 参数 | 位置 | 类型 | 必填 | 默认值 | 校验 |
| --- | --- | --- | --- | --- | --- | --- |
| GET /api/v1/trade-events/live | instrument_type | URL query | string | 否 | PERPETUAL | PERPETUAL、SPOT、ALL，区分大小写 |
| GET /api/v1/trade-events/stored | instrument_type | URL query | string | 否 | PERPETUAL | 同上；下一页必须保持不变 |
| POST /internal/v1/trade-queries | instrument_type | JSON body | string | 否 | PERPETUAL | 与外部一致，由网关透传 |
| POST /internal/v1/stored-trade-queries | instrument_type | JSON body | string | 否 | PERPETUAL | 与外部一致；游标绑定此条件 |

非法枚举、重复参数、未知参数返回 JSON 400 VALIDATION_ERROR。现货筛选不得改变 stored 的时间边界、快照或排序。旧游标解码为 PERPETUAL；新游标增加 instrument_type 并升级游标版本；跨类型复用返回 400。

现有旧外部路径保持兼容转发，默认范围保持 PERPETUAL，新示例只使用功能路径。

### 5.3 计数与覆盖

在响应 data.counts 中新增计数版本 counts_version=2，计数覆盖完整来源响应，不随展示筛选变化。source_records 拆为 duplicate_records、perpetual_records、spot_records、unsupported_records、invalid_records；后四项互斥，perpetual_records 和 spot_records 只包含成功标准化的去重事实。

旧版 spot_records 表示被识别但过滤的现货；新版表示成功标准化现货，必须通过 counts_version 明确区别。历史 collection_jobs.result 和证据文件保留原计数，不原地改写。新增 matched_records 表示筛选命中数，returned_records 表示展示数；display_truncated 按匹配集合计算，持久化计数覆盖两类事实。

每项金额使用自身 quote_asset，不把不同计价资产的 notional 直接相加。live 仍是来源当前窗口，stored 仍只代表库存，水位仍是已扫描边界，均不证明完整历史。

### 5.4 现货发布激活与补录任务

新增按账户和 instrument_type 保存的能力激活记录，至少包含 enabled、activation_epoch、activated_at、policy_version。现货首次启用和重新启用建立自身边界；不调用全局激活重置来代替现货开关，不影响永续已激活状态。管理恢复同时更新适用的监控代次与能力边界。

新增现货补录任务，保存 job_id、account_key、固定 start/end、已扫描位置、阶段、状态、重试时间、计数、最后错误、来源覆盖说明及创建/完成时间。每账户最多一个运行或等待中的补录任务，重复提交相同范围返回同一任务；冲突范围拒绝，不并发争抢同账户补录游标。

发布详情增加按类型的能力状态及抑制计数；账户详情增加 spot_backfill 状态。汇总保持 items/outbox，不新建通过参数切换不同业务功能的入口。数据库迁移只建结构，不批量重解析全部旧日志。

## 6. 处理流程

### 6.1 HTTP 与 WS

1. 获取并归档原始 fills 和实际使用的永续/现货元数据。
2. 将来源市场规范化后做身份去重，分别解析两类成交；错误与证据关联。
3. 在既有代次和租约 fencing 下保存完整事实、观测及扫描位置。
4. 只有符合实时条件的 WS 事实进入按类型候选策略。
5. 候选和 outbox 与事实同事务提交，沿用至少一次投递和不可变正文。
6. 查询在完整结果上按 instrument_type 筛选，再应用 limit；stored 在数据库内过滤后分页。

不创建独立的日常“现货水位”。同一来源范围必须完成两类支持成交的处理后才可推进原水位；现货解析失败不能被合约成功掩盖。

### 6.2 升级历史补录

补录有独立持久进度，不复用或倒退原采集水位。通过 collector CLI 提交固定半开时间范围任务，常驻 collector 在共享预算内分批处理，不由临时 CLI 持有长期采集循环。

先扫描已完成旧任务的原始响应及其对应元数据，仅追加过去未入库的现货事实；不得改写旧永续结果、历史计数或已投递正文。历史任务重新解析得到新内容时输出带解析版本的独立执行报告，不把旧任务回执覆盖为“新版本一致”。

随后对指定范围执行官方 HTTP 扫描；不能把保存过一个原始包视为整个时间段已覆盖。来源分页饱和、范围不可继续或元数据不足时保留进度并标记 BLOCKED/覆盖不足，不能直接成功跳过。补录期间实时采集优先，按批让出共享资源。

账户暂停同时暂停其补录执行；恢复继续原固定范围。进程或数据库故障后由任务进度续跑，重试不重复事实，任何阶段均不生成历史候选。

## 7. 程序结构设计

### 7.1 归属与复用

协议的市场映射、字段解析归属 protocols/implementations/hyperliquid；标准事实与策略归属现有业务 crate；持久化、查询、补录调度和投递归属 trade-log adapters；HTTP handler 只做参数、认证和转发。保持 workspace 和多进程 Docker 结构。

以下为实施文件清单，新增和更新分开列示。开发验证与服务器验收分别记录，未验证项目不记为通过。

### 7.2 新增文件

| 文件 | 用途 |
| --- | --- |
| protocols/implementations/hyperliquid/src/spot.rs | 现货交易对、token 与别名映射及字段解析 |
| services/trade-log/src/instrument.rs | 市场类型枚举、默认值和筛选规则 |
| services/trade-log/adapters/src/postgres/spot_activation.rs | 独立现货激活和幂等控制同步 |
| services/trade-log/adapters/src/postgres/spot_backfill.rs | 补录任务、检查点、租约和结果持久化 |
| services/trade-log/adapters/src/spot_backfill_runtime.rs | 分批重解析与 HTTP 补录，复用预算及暂停约束 |
| services/trade-log/migrations/0006_spot_trades.sql | 能力激活、补录表及按类型查询索引 |
| protocols/implementations/hyperliquid/tests/spot.rs | 映射、费用、精度、未知市场及身份测试 |
| services/trade-log/adapters/tests/spot_chain.rs | 两类事实、查询、候选、迁移和补录集成测试 |
| tests/fixtures/hyperliquid/spot/ | 可核对的 fills、完整元数据、费用及错误样本 |
| scripts/init-v0.7.sh | 增量迁移，沿用数据卷及凭证 |
| scripts/verify-v0.7.sh、scripts/verify-v0.7.py | 固定服务器验收入口及检查 |
| scripts/verify-v0.7-fixtures.py、scripts/run-v07-fixtures.sh | 独立测试库固定样本及镜像执行入口 |
| scripts/tests/test_verify_v07_fixtures.py | 资源隔离和失败清理测试 |
| docs/version0.7.md | 本开发文档 |
| docs/version0.7-acceptance.md | 实现阶段新增，分别记录开发验证及服务器验收 |

### 7.3 更新文件

| 文件 | 修改内容 |
| --- | --- |
| Cargo.toml、Cargo.lock | 版本 0.7.0，按需要调整依赖 |
| services/trade-log/adapters/src/file_evidence.rs | 新证据记录解析版本；复用原标准事实结构 |
| protocols/implementations/hyperliquid/src/lib.rs、src/parser.rs | 接入现货解析，统一市场去重与计数 |
| services/trade-log/src/lib.rs、src/normalization/mod.rs、src/query/mod.rs、src/validation/mod.rs、src/stored_query/mod.rs | 类型筛选、计数版本、游标兼容 |
| services/trade-log/src/collection/mod.rs、src/collection/content.rs、src/publishing/mod.rs | 合并计数、现货语义身份和能力状态 |
| services/trade-log/src/publishing/policy.rs | 按成交类型判断准入，保留原永续条件 |
| services/trade-log/adapters/src/lib.rs、src/postgres/mod.rs、src/postgres/migration.rs | 模块与迁移注册、角色最小权限 |
| services/trade-log/adapters/src/postgres/facts.rs、src/postgres/stored_query.rs、src/postgres/replay.rs | 现货持久化、类型筛选和不可覆盖历史回执 |
| services/trade-log/adapters/src/postgres/publishing.rs、src/postgres/watch.rs | 分类型能力边界、发布及补录状态 |
| services/trade-log/adapters/src/publisher_runtime.rs | 按类型同步控制，不重置永续激活 |
| services/trade-log/bins/trade-collector/src/bootstrap.rs | 补录提交/状态命令和后台执行装配 |
| services/trade-log/bins/trade-collector/Cargo.toml | 补录 CLI 使用已有 workspace UUID 依赖 |
| services/trade-log/bins/trade-parser-publisher/src/config.rs、src/bootstrap.rs | spot_enabled 校验及分类型策略 |
| gateway/query-api/src/http/handlers/trade_events.rs | 外部筛选参数，复用原客户端转发 |
| config/trade-parser-publisher.toml | spot_enabled 及明确中文注释 |
| compose.yaml、gateway/query-api/Dockerfile | 版本和验收 target，沿用原挂载及数据卷 |
| scripts/verify-version.py、scripts/tests/test_verify_version.py | 兼容计数版本和旧默认 PERPETUAL，保持旧版本验收 |
| protocols/implementations/hyperliquid/tests/parsing.rs、services/trade-log/adapters/tests/ 及既有请求校验测试 | 请求字段默认值、完整现货样本及原能力回归 |
| tests/fixtures/hyperliquid/spot-meta.json | 固定元数据补齐 tokenId，不作为真实来源记录 |
| docs/api.md、docs/verify.md、docs/test-clean.md | 统一参数表、执行命令、现货补录及能力记录的清理说明 |

README、roadmap 及其他设计文件不增加入口。只有发现与当前设计存在实质冲突时再明确修订对应设计条款，不借版本实施改动无关内容。

## 8. 配置与异常处理

### 8.1 配置

保留 collector 中的来源、通用采集策略、监控容量和共享预算；不增加账户列表，也不新增重复的现货来源 URL、端口或 HTTP 周期。

| 配置项 | 类型 | 默认值 | 意义及约束 |
| --- | --- | --- | --- |
| publishing.spot_enabled | boolean | false | 仅开启现货候选准入；必须同时 publishing.enabled=true，具备确认元数据、实时会话及有效时效 |

新增字段在代码中有默认值，旧配置可继续解析。配置每项保留明确中文注释。修改后重启 publisher 生效；重启保持未变更类型的激活边界和已保存 outbox。

### 8.2 异常

| 异常 | 行为 |
| --- | --- |
| 非法筛选或游标类型不匹配 | 400，明确校验错误，不能改为 ALL 或成功空数组 |
| 市场身份冲突或重复事实内容不同 | 沿用 409 冲突，事务不推进，保存证据 |
| 支持市场元数据缺失 | 刷新后仍失败则阻塞对应范围；不得猜资产或静默跳过 |
| 手续费币种缺失/歧义 | 可保存事实和警告，不生成现货候选 |
| 来源限流或共享预算不足 | 有界等待/退避，保留任务进度；补录不抢占全部容量 |
| 数据库不可用 | 业务 503，保留原有故障恢复语义 |
| 补录范围来源不可完整查询 | 保留扫描证据，报告覆盖不足或 BLOCKED；不标“完整历史” |
| 接收方失败或响应丢失 | 原文重试、event_id 去重、TTL 到期处理不变 |

## 9. 验证与验收

### 9.1 固定验证

1. 标准名称与 @index 指向同一现货身份；同名不同 tokenId 不合并；未知索引不猜测。
2. 正常买卖、基础/计价资产手续费、负手续费、缺失币种、非 USDC 计价和高精度金额核对。
3. 同响应两类成交分别展示但完整入库；SPOT/PERPETUAL/ALL 的数量、排序、limit 和计数守恒。
4. HTTP/WS、多个来源包、名称别名重复不重复事实；同 tid 的不同市场不碰撞。
5. stored 类型过滤在 limit 前执行，快照稳定、跨页不重复，跨类型游标拒绝；旧游标按 PERPETUAL 兼容。
6. 主动实时现货在明确启用后产生候选；历史、快照、旧代次、旧激活、未知费用和强制事件不产生候选。
7. 默认关闭现货发布；开关、重启和暂停恢复不重置无关永续激活，不修改旧 outbox 字节。
8. 旧数据库迁移不改永续事实、旧任务回执、水位和历史计数；已有原始证据的现货可独立补录。
9. 补录暂停、断库、重启、重复范围和来源饱和恢复；当前采集继续推进，不被历史任务饿死。
10. 旧版本脚本识别 counts_version，默认查询仍只展示合约；固定测试资源独立创建并在失败后清理。

### 9.2 服务器验收

新增 verify-v0.7 脚本，读取已登记账户，不要求再次传 account。至少选择一个确有现货成交的账户；缺少样本时标 SKIP。固定样本与真实来源分开报告，不向部署账户伪造事实，不为制造成交而下单。

检查两类筛选、数据库事实身份、证据关联、金额与手续费、重启保留和既有能力回归。真实 WS 候选需匹配接收器 event_id 与实际成交；仅连接正常或 stored 有记录不能作为发布通过。

## 10. 交付、部署与手动验证

以下命令适用于 v0.7，需先升级应用及执行迁移。验证情况见 [验收记录](version0.7-acceptance.md)。

### 10.1 升级

备份数据库和凭证，执行增量迁移；不删除卷、不重建原账户、不自动全量回补。现货发布默认关闭，旧接收方需支持 SPOT 字段后再开启。

```sh
sh scripts/init-v0.7.sh
docker compose up --build -d
```

### 10.2 查询及对账

在服务器项目根目录设置一个确有现货成交的实际账户，按来源记录选择时间，不将示例地址当作已验证活跃账户。

```sh
ACCOUNT='<实际钱包地址>'

# 近期现货：保存完整来源结果，展示 SPOT。
curl --noproxy '*' -sS --get http://127.0.0.1:8080/api/v1/trade-events/live \
  --data-urlencode "account=$ACCOUNT" --data-urlencode 'instrument_type=SPOT' \
  --data-urlencode 'limit=10'

# 已保存现货：只读数据库。
curl --noproxy '*' -sS --get http://127.0.0.1:8080/api/v1/trade-events/stored \
  --data-urlencode "account=$ACCOUNT" --data-urlencode 'instrument_type=SPOT' \
  --data-urlencode 'limit=10'

# 合约仍可独立查询。
curl --noproxy '*' -sS --get http://127.0.0.1:8080/api/v1/trade-events/stored \
  --data-urlencode "account=$ACCOUNT" --data-urlencode 'instrument_type=PERPETUAL'
```

通过标准：返回类型符合筛选，资产、价格、数量、费用及币种与原始响应一致；两次查询相同成交的 fact_id 一致，不重复入账。

### 10.3 历史现货补录

补录只写事实，不生成历史候选。以下 CLI 为新增补录命令；起点和终点使用带时区 RFC3339，结束不包含，且不得晚于提交时刻。

| 命令 | 参数 | 类型 | 必填 | 默认值 | 含义及校验 |
| --- | --- | --- | --- | --- | --- |
| backfill-spot | --account | string | 是 | 无 | 已登记的钱包地址，非登记 UUID；网络取 collector 配置 |
| backfill-spot | --start-time | RFC3339 string | 是 | 无 | 包含的起点，时区明确、非负 Unix 时间、最多毫秒精度 |
| backfill-spot | --end-time | RFC3339 string | 是 | 无 | 不包含的终点，大于起点且不晚于提交时刻 |
| backfill-spot-status | --job-id | UUID string | 是 | 无 | 提交补录命令返回的任务 ID，不是钱包或监控登记 ID |
| 两者 | --config | string | 否 | /etc/robotech/trade-collector.toml | 容器内配置文件路径；使用原数据库角色及共享预算 |

提交时账户尚未登记则拒绝。已暂停账户允许登记补录任务但不执行，恢复后继续；不得因为提交补录而隐式恢复监控。

```sh
docker compose exec -T trade-collector trade-collector backfill-spot \
  --account "$ACCOUNT" --start-time '<RFC3339开始时间>' --end-time '<RFC3339结束时间>'

# 从提交结果复制 job_id。
docker compose exec -T trade-collector trade-collector backfill-spot-status \
  --job-id '<实际补录任务UUID>'
```

记录补录前后永续事实数、现货事实数、原 HTTP 水位、激活时间和 outbox；核对现货新增、永续身份不变、无历史候选。重复提交相同范围复用任务；若任务为 BLOCKED，修复原因后重新提交相同范围会重新进入 PENDING。尚在运行、BLOCKED 或覆盖不足不能记作历史完整。

### 10.4 实时发布与恢复

1. 按 v0.5 接收器流程启动 webhook-receiver，配置对应 URL、令牌及测试用 HTTP 开关。
2. 保持原发布总开关，显式设置 spot_enabled=true，重启 publisher；检查该账户现货激活时间，永续激活保持不变。
3. 等待账户实际现货交易，通过分账户 publishing-status 和 `/receipts` 核对标准 SPOT 事实、费用、event_id、时间及 TTL。
4. 暂停/恢复该账户，检查另一账户继续采集；暂停期间历史补偿没有现货候选。
5. 重启应用和短暂停库后核对事实、补录任务、能力边界及队列保留；故障测试必须恢复数据库。

### 10.5 固定验收脚本

```sh
sh scripts/verify-v0.7.sh --fixture-tests --lifecycle --database-fault --wait-seconds 180
```

命令会短暂中断服务。参数沿用现有脚本含义，已加入 docs/verify.md；可额外加 --live 查询官方近期现货，默认只读部署库存。PASS、FAIL、SKIP 单独列出，未验证不算通过。验证记录写入 version0.7-acceptance.md；bug 沿用描述、影响、确认记录、修改逻辑、修改记录摘要、如何验证、验证记录七项格式。
