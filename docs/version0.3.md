# v0.3 开发文档：单地址自动更新与采集水位

- 文档版本：0.3
- 状态：待实现，本文件中的新增接口、配置和文件为实现约定，尚未交付
- 更新日期：2026-10-05
- 设计依据：[概要设计](overview-design.md)、[详细设计](detailed-design.md)、[产品版本路线图](roadmap.md)
- 前置版本：[v0.2 开发文档](version0.2.md)
- 接口契约统一维护在 [接口说明](api.md)，实施本版时同步新增接口及参数说明

## 1. 版本目标

配置一个 Hyperliquid 账户后，程序自动定时获取成交并保存。用户无需持续请求 live 接口，即可通过已有 stored 接口查看自动更新的合约成交，同时查询采集运行状态、最后成功时间、最后错误和持久化水位。采集进程重启后从水位继续，不重复入账。

沿用 query-api、trade-log-query、PostgreSQL，新增详细设计中既定的 trade-collector 独立进程和 Docker 容器。采集业务仍属于 services/trade-log，复用现有来源、解析器、标准事实、数据库及幂等逻辑。应用与镜像版本升级为 0.3.0，保留 schema_version=1 及既有接口路径。

“自动更新成功”表示本轮来源请求、解析、完整事实保存和水位提交成功，不表示已经证明全部历史完整。轮询可能存在延迟，不能当作实时跟单信号；WebSocket、信号输出、多地址管理分别按 roadmap 后续版本交付。

## 2. 功能范围

| 能力 | 本版交付 | 边界 |
| --- | --- | --- |
| 单地址自动采集 | 配置一个账户，按间隔自动拉取时间范围成交 | 不新增账户添加、删除、暂停、恢复管理接口 |
| 增量与重叠查询 | 使用持久化扫描水位及重叠窗口查询 | 不依赖手动 live 查询触发，不以最新成交时间代替扫描水位 |
| 水位保存 | 成交保存与水位推进在同一事务提交 | 任一失败均不推进水位 |
| 重启恢复 | 读取原水位、未完成范围和调度状态继续运行 | 配置起点不覆盖已有水位 |
| 状态查询 | 查看运行状态、成功时间、失败次数、当前范围和水位 | 服务健康不代表最近采集成功 |
| 来源上限处理 | 区间拆分、有界请求、不能完整处理时明确停留水位并记录错误 | 不保证恢复官方已不可访问的旧历史 |
| 幂等与证据 | 每轮请求保留原始字节及关联，重复成交不重复保存 | 原始响应与观察记录增加属于正常采集留痕 |
| 部署升级 | Docker 新增采集进程，增量 migration | 不删除现有事实、证据、凭证或命名卷 |

本版不增加消息队列、outbox 发布、webhook、现货展示、账户收益分析或跟单执行。采集进程在本版同步完成来源读取和解析保存；未来拆出 trade-parser-publisher 时沿用业务端口和原始响应身份，不改写事实模型。

## 3. 使用方式

### 3.1 配置自动更新账户

新增 config/trade-collector.toml。使用实际交易账户地址，不使用 agent wallet；地址验证沿用既有规则并规范化为小写。默认 network=mainnet，与已有存储保持一致。

首次启动必须明确 start_time，表示希望从何时开始扫描；配置示例使用占位符，用户须填写实际值。已有同账户检查点时从数据库恢复，不根据每次启动的当前时间重新建立起点。

```toml
[collection]
account = "<实际账户地址>"
start_time = "<带时区的RFC3339起点，最大毫秒精度>"
interval_seconds = 30
overlap_seconds = 60
safety_delay_seconds = 2
max_window_seconds = 3600
max_requests_per_round = 20
round_timeout_seconds = 120
lease_seconds = 180
retry_base_seconds = 5
retry_max_seconds = 300
```

本版只配置一个地址，修改账户并重启时建立或读取该账户自己的检查点，旧账户事实和水位保留。相同账户修改 start_time 不重置水位；需要更早历史时使用后续明确的回补能力。修改 network 必须与已有查询服务、来源环境一致。

### 3.2 查看自动更新状态

```text
GET /api/v1/watch-accounts
```

本版仅提供只读列表，正常最多返回一个配置中的账户，不读取所有历史检查点充当正在监控的账户。新增 collector 内部调用后由网关返回统一响应。账户没有成交时，成功时间和扫描水位仍应更新，last_trade_at 可以为空。

### 3.3 查看保存的成交

```text
GET /api/v1/trade-events?account=<ADDRESS>&source=stored&limit=100
```

继续使用 v0.2 的时间范围和快照游标。自动采集完成后，新成交可通过新的 stored 请求看到；正在翻页的旧快照仍按原 snapshot_seq 返回。stored 不触发来源访问或立即轮询。

省略 source 的 live 查询继续保持 v0.2 行为，手动查询和后台采集可同时保存同一成交，依靠稳定身份去重。手动查询不会推进自动采集的扫描水位，也不会修改最后自动成功时间。

## 4. 数据来源与查询边界

### 4.1 来源请求

自动采集使用官方 /info 的 userFillsByTime，aggregateByTime=false，同时保存 meta、spotMeta 以复用现有合约/现货识别。手动 live 查询继续使用 userFills。

官方文档说明：时间范围成交单次最多 2,000 条、可访问最近 10,000 条；startTime/endTime 均为包含边界的毫秒时间。上述上限针对来源成交，含本版不展示的现货，不能根据解析后的合约条数判断是否截断。[官方 Info endpoint 说明](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint#retrieve-a-users-fills-by-time)

本版内部范围统一为 `[start_ms,end_ms)`，转换为官方 startTime=start_ms、endTime=end_ms-1。空区间不发请求。时间使用整数毫秒、UTC；不允许浮点换算或自动丢弃亚毫秒精度。

### 4.2 扫描范围与水位

- scanned_through 表示已成功提交的扫描结束边界，不表示最新一笔成交的时间，也不表示历史完整性证明。
- last_trade_at 表示自动采集已提交的最新合约成交时间，可以为空；重复和乱序返回不使它倒退。
- 已有水位 W 时，下一轮起点取 max(initial_start,W-overlap)；没有水位时使用明确的 initial_start。
- 本轮结束 E 取 min(now-safety_delay,W+max_window)，首次以 initial_start 代替 W。now 在一轮开始时固定，不随分页变化。
- 结束不大于起点时等待下次调度，不执行空任务。now 早于已有水位时记录 CLOCK_BEHIND_WATERMARK，不回退水位或伪造成功。

无成交的正常响应也可以将 scanned_through 推进到 E；否则安静账户永远无法推进。迟到数据通过重叠区间和既有事实身份吸收，超出重叠窗口的迟到成交可能漏采，状态必须保留覆盖限制。

### 4.3 截断、分页与覆盖限制

不直接取最后返回时间加 1 毫秒翻页，避免漏掉同一毫秒内的其他成交。返回少于上限时可进入本轮解析保存；达到 2,000 条时将当前范围按整数毫秒中点拆成两个不重叠区间，再分别查询，原始满页同样归档。

拆分到单个毫秒仍达到上限时，无法通过时间接口证明该毫秒已取全：本轮返回 SOURCE_WINDOW_SATURATED，保留原水位和原始响应，不越过该毫秒。请求预算用尽时保存剩余子区间及已完成页引用，进入 RETRY_WAIT，下一轮继续；尚未完成完整范围不得提前提交水位。

官方只开放有限的最近历史。即使某个范围返回空或少量数据，也无法仅凭这一响应确认更早记录没有被来源淘汰。因此覆盖字段固定 SOURCE_HISTORY_NOT_VERIFIED；若检测到截断、无法推进或历史限制线索，增加明确警告并保留缺口说明，不输出“完整历史”或“已补齐”。停机时间较长后能够恢复运行，不等于能够恢复全部缺失成交。

## 5. 接口与持久化数据设计

### 5.1 外部与内部接口

| 类型 | 接口 | 用途 |
| --- | --- | --- |
| 外部，新增 | GET /api/v1/watch-accounts | 查看当前配置账户和自动采集状态 |
| 内部，新增 | GET /internal/v1/collection-status | 网关向 collector 读取配置账户与数据库检查点 |
| 内部，新增 | GET /internal/v1/health | collector 进程存活状态，需内部凭证 |
| 外部，复用 | GET /api/v1/trade-events | live/stored 成交查询 |
| 外部，复用 | GET /api/v1/health、GET /api/v1/version | 网关基础接口 |

新增接口沿用 data/meta 成功信封和 code/message/trace_id 错误格式；内部使用既有服务 token。collector 不发布宿主机端口，网关通过 Compose 服务名访问。暂不实现详细设计中的 POST/DELETE watch-accounts，避免在单配置地址版本提前交付多账户生命周期。

状态查询成功时 data 为 `{items:[状态对象]}`，meta.schema_version=1；collector 未启用或不可访问返回 503，不能把未启用解释为正常空列表。数据库不可读同样返回 503；能读取检查点而采集自身失败时返回 200，并在状态对象中说明失败。

### 5.2 状态字段

| 字段 | 格式/含义 |
| --- | --- |
| account、account_key、network | 配置账户及规范身份，与事实模型一致 |
| status | STARTING、RUNNING、WAITING、RETRY_WAIT、FAILED、STOPPED |
| coverage | SOURCE_HISTORY_NOT_VERIFIED |
| initial_start_time | 首次建立检查点时保存的起点，后续配置不覆盖 |
| scanned_through | UTC 毫秒时间或 null，已提交的扫描结束边界 |
| last_trade_at | 自动采集保存的最新合约成交时间或 null |
| last_attempt_at | 最近开始尝试的时间或 null |
| last_success_at | 最近完整提交事实与水位的时间或 null |
| consecutive_failures | 连续失败轮数；完整成功后归零 |
| next_run_at | 下一次调度时间或 null，包含退避 |
| pending_range | 当前未完成的起止范围或 null，端点为 UTC 时间 |
| last_query_id | 最近尝试的采集任务 ID 或 null，可用于重新解析 |
| last_success_query_id | 最近成功任务 ID 或 null |
| last_error | null 或 `{code,message,occurred_at}`，脱敏错误摘要 |
| heartbeat_at、lease_expires_at | 运行实例心跳和租约到期时间或 null |
| warnings | 来源范围、迟到数据和历史覆盖等说明 |

last_error 在下一次完整成功后清空；失败任务与原始响应仍保留在 collection_jobs/raw_logs 中。无成交成功不能伪造 last_trade_at。状态读取时若持久化状态仍 RUNNING 但租约已经过期，对外显示 FAILED 并增加 COLLECTOR_HEARTBEAT_EXPIRED，不把旧心跳当作仍正常采集。

### 5.3 检查点表

新增 migration `0002_collection_checkpoints.sql`，新增详细设计既定的 trade_log.collection_checkpoints。主键沿用 `(chain_id,source_id,partition_key)`，partition_key 使用规范 account_key；一个网络、来源和账户共享一个自动采集检查点。

| 字段 | 类型/要求 |
| --- | --- |
| chain_id、source_id、partition_key | TEXT，联合主键；source_id=official_http |
| position | JSONB NOT NULL，含 initial_start_ms、scanned_through_ms、last_trade_at_ms，使用非负整数毫秒或明确 null |
| cursor | TEXT 可空，本版官方接口无不透明游标，保持 null |
| finality_proof | JSONB 可空，本版不构造链最终性证明 |
| status | TEXT NOT NULL，状态枚举约束 |
| pending_work | JSONB 可空，包含固定范围、剩余子区间、已完成页和原始响应引用 |
| last_attempt_at、last_success_at | TIMESTAMPTZ 可空 |
| last_query_id、last_success_query_id | TEXT 可空，外键关联 collection_jobs.query_id |
| consecutive_failures | INT NOT NULL DEFAULT 0，非负 |
| next_run_at | TIMESTAMPTZ 可空 |
| last_error | JSONB 可空，脱敏摘要 |
| lease_owner | UUID 可空，每次进程启动生成实例 ID |
| lease_epoch | BIGINT NOT NULL DEFAULT 0，接管时递增的隔离令牌 |
| heartbeat_at、lease_expires_at | TIMESTAMPTZ 可空 |
| created_at、updated_at | TIMESTAMPTZ NOT NULL |

不另建 v03_trades 或单独的成交表，不迁移既有 fact_id。事实仍写 v0.2 的版本/current/observation 表，ingest_seq 规则保持不变。

### 5.4 采集任务、分页证据和事务

collection_jobs 增加 job_origin（MANUAL/COLLECTOR，旧记录默认 MANUAL）、checkpoint 关联身份和 lease_epoch。每个完整扫描范围对应一个采集任务；可拆分为多个来源页，但不能建立并行的成交身份体系。

raw_logs 增加 page_no INT NOT NULL DEFAULT 0。将原 `(collection_job_id,kind,attempt)` 唯一约束扩展为 `(collection_job_id,kind,page_no,attempt)`；旧页保持 page_no=0。既有原始 ID 和既有摘要不修改；新的页身份纳入 page_no，精确请求时间范围保存在 request 中。该变更兼容旧单页任务和旧证据导入。

完整范围所有页归档并解析成功后，在一个数据库事务内完成：事实幂等保存、观察关联、任务 COMPLETED、检查点水位推进、last_success_at 和计数更新、清除 pending_work。沿用现有事实冲突整批回滚规则。必须将 v0.2 事实事务能力扩展为可参与同一个检查点事务，不能先提交事实、再另发 SQL 推水位。

原始响应在来源读取后归档，即使后续失败也保留。请求预算耗尽只是未完成调度，保存 pending_work 并等待，不记为完整成功。不可重试的解析/身份冲突将本轮 FAILED，保留旧水位和任务证据；修复后继续原范围，不能跳过问题记录。

reparse 扩展为读取全部成功来源页，按各页 raw_log_id/source_index 解析并合并去重，比较已提交事实；手动单页任务保持兼容。旧 import-evidence 不推进 collector 检查点。

## 6. 处理流程

### 6.1 启动与恢复

```text
加载配置并校验 → 数据库 schema 检查 → 读取或建立账户检查点
→ 获取账户租约 → 恢复 pending_work 或根据已提交水位建立范围
→ 定时来源读取 → 原始响应归档 → 解析、校验、去重
→ 事实与检查点同事务提交 → 更新成功时间 → 等待下次调度
```

已有 RUNNING 任务不能由另一个查询进程启动时统一标记 INTERRUPTED。v0.2 的启动清理需要按 job_origin 区分：query 服务只清理自己的手动任务；collector 仅在租约到期后恢复对应采集任务，并验证实例及 epoch。

首次配置的起点晚于当前可查询时间时明确配置错误；已有检查点恢复时不以配置起点重置。数据库不可用、schema 不匹配或凭证文件缺失时按现有启动约定退出，不建立仅在内存中的“成功水位”。

### 6.2 调度与失败重试

一个账户同一时间只处理一个范围。正常间隔从一次完整成功之后计算，不使用可能累积重叠任务的无界定时 tick。轮询时限覆盖来源请求、重试、解析和数据库写入。

可重试的网络、429、5xx、数据库短时故障使用有界指数退避和小幅随机抖动，遵守可解析的 Retry-After；只持久化 next_run_at，不阻塞状态接口。连续失败次数递增，成功归零。来源限流不允许启动更多并行请求补偿。

不可重试的身份冲突、协议数据不可解析和毫秒饱和进入 FAILED，状态接口仍可读；进程保持存活供排查，但不自动跳过范围。修改配置或修复原因并重启后，在旧水位基础上重新尝试；不需要删除检查点。

### 6.3 租约与并发

实例获取租约使用数据库条件更新或插入冲突处理，仅在租约未持有或已过期时接管，同时递增 lease_epoch。时间判断使用数据库 now()，续租周期不超过 lease_seconds/3，不在 HTTP 请求期间持有数据库行锁。

提交事实和水位时检查 lease_owner、epoch、租约有效期；令牌失效则整个提交回滚。旧实例恢复连接后不能覆盖新实例水位。重复启动 collector 时没有获得租约的实例等待并报告 LEASE_HELD，不进行来源采集；不能依靠“Compose 通常只有一个副本”保证唯一执行。

数据库断线时无法确认租约和事务提交结果，不推进内存水位；重连后重新读取检查点、任务和租约。若事务已提交但响应丢失，按已提交状态继续；若未提交，则重跑，唯一约束保证事实不重复。

### 6.4 停止

收到 SIGTERM 后停止领取新轮次，在 shutdown_timeout 内完成或取消当前工作。取消不提交半轮水位；保留已有原始响应和可恢复工作。正常释放自有租约并记录 STOPPED；强制退出后依靠租约过期接管。

网关、查询程序、collector 分别优雅退出。collector 退出不删除监控账户的历史事实。重启后的第一次执行优先恢复尚未完成的范围，再进入正常间隔。

## 7. 程序结构设计

### 7.1 进程与依赖

```text
query-api ──内部 HTTP──► trade-log-query ──► PostgreSQL
         └─内部 HTTP──► trade-collector ──► PostgreSQL
                              └──────────► Hyperliquid 官方 HTTP
```

业务端口定义在 trade-log crate，Tokio 调度、HTTP、SQLx 和日志适配放在 adapters 或 binary 装配层。现有 QueryService 的 live 路径继续复用；后台流程通过 acquisition、checkpoint、parsing 和 persistence 协作，不将循环任务写进网关 handler，也不复制协议解析器。

### 7.2 核心模块

| 模块 | 职责 |
| --- | --- |
| acquisition | 增加时间范围来源请求及页身份，保留现有单页 SourceReader 能力 |
| checkpoint | 检查点、扫描范围、租约令牌、状态及存储端口 |
| collection | 一轮采集编排、区间拆分、解析合并、提交意图，不依赖 Tokio/SQLx |
| persistence | 扩展完整事实与检查点的事务提交端口 |
| adapters/postgres | 检查点租约、任务页保存和原有事实事务的组合实现 |
| adapters/collector_runtime | 有界调度、超时、续租、退避和生命周期 |
| bins/trade-collector | 配置、依赖装配、内部状态 HTTP 和进程入口 |

### 7.3 本版新增与更新文件

以 v0.2 为基线，以下是实现文件清单；新增与更新明确区分。实现期间如必须调整文件归属，先保持设计职责，再同步本节，不把设计中的未实现项标为已交付。

#### 新增文件

| 文件 | 职责 |
| --- | --- |
| services/trade-log/src/checkpoint/mod.rs | 检查点 DTO、扫描范围、租约与存储端口 |
| services/trade-log/src/collection/mod.rs | 自动采集单轮业务及来源满页拆分规则 |
| services/trade-log/adapters/src/postgres/checkpoint.rs | 水位读取、租约获取续期、状态保存和事务提交检查 |
| services/trade-log/adapters/src/collector_runtime.rs | 独立轮询、续租、退避、超时和停止 |
| services/trade-log/adapters/src/collector_http.rs | 内部 health/collection-status 路由及凭证校验 |
| services/trade-log/bins/trade-collector/Cargo.toml | 新 binary crate，依赖既有服务及运行组件 |
| services/trade-log/bins/trade-collector/src/main.rs | 程序入口与退出码 |
| services/trade-log/bins/trade-collector/src/lib.rs | 导出配置和装配以支持测试 |
| services/trade-log/bins/trade-collector/src/config.rs | 单账户采集配置和关联约束校验 |
| services/trade-log/bins/trade-collector/src/bootstrap.rs | collector 装配与优雅关闭 |
| services/trade-log/bins/trade-collector/tests/startup.rs | 配置错误、schema、租约和退出验证 |
| services/trade-log/migrations/0002_collection_checkpoints.sql | 检查点、job_origin、页身份和约束增量迁移 |
| gateway/query-api/src/clients/collector.rs | collector 状态内部调用，不读取数据库 |
| gateway/query-api/src/http/handlers/watch_accounts.rs | 只读状态接口及统一响应 |
| config/trade-collector.toml | collector 运行及单账户配置 |
| scripts/init-v0.3.sh | 保留既有初始化行为，迁移新 schema 并初始化采集凭证 |
| services/trade-log/tests/collection.rs | 范围计算、同毫秒边界、满页拆分及重叠测试 |
| services/trade-log/adapters/tests/checkpoint.rs | 原子水位、租约隔离、无成交成功和重启恢复测试 |
| services/trade-log/adapters/tests/collector_chain.rs | 固定来源下调度、持久化和状态 HTTP 链路 |
| tests/fixtures/v0.3/collection-cases.json | 确定性的空范围、多页、迟到、冲突和故障测试场景 |

#### 更新文件

| 文件 | 变更点 |
| --- | --- |
| Cargo.toml | 注册 trade-collector，workspace 版本 0.3.0；仅增加必要依赖 |
| Cargo.lock | 锁定版本及实际依赖 |
| compose.yaml | 新增 collector target、配置与凭证挂载，继续使用 postgres 和原卷 |
| gateway/query-api/Dockerfile | 构建新增 binary，增加 trade-collector target，保留 Cargo 缓存 |
| config/query-api.toml | 新增 collector 内部地址和请求时限 |
| services/trade-log/src/lib.rs | 导出 checkpoint 和 collection |
| services/trade-log/src/acquisition/mod.rs | 增加时间范围来源能力，不破坏现有单页调用 |
| services/trade-log/src/persistence/mod.rs | 统一完整事实及检查点原子提交契约 |
| services/trade-log/adapters/src/lib.rs | 导出 collector runtime/http 适配器 |
| services/trade-log/adapters/src/postgres/mod.rs | 装配检查点；按任务来源区分启动恢复 |
| services/trade-log/adapters/src/postgres/raw_log.rs | 新页身份、完整请求范围及 collector 任务归档 |
| services/trade-log/adapters/src/postgres/facts.rs | 复用事实保存，支持在同事务更新检查点并检查租约 |
| services/trade-log/adapters/src/postgres/migration.rs | 新 schema 版本及 collector 角色授权 |
| services/trade-log/adapters/src/postgres/replay.rs | 多页重新解析，保持旧导入行为 |
| crates/protocol-api/src/lib.rs | 注册 userFillsByTime 请求类别，不改变原枚举值的序列化含义 |
| protocols/implementations/hyperliquid/src/source.rs | 实现时间范围请求及原始响应读取，复用 HTTP 限制 |
| gateway/query-api/src/config.rs | 新 collector 配置及凭证校验 |
| gateway/query-api/src/state.rs | 增加可选 collector 客户端 |
| gateway/query-api/src/bootstrap.rs | 装配 collector 客户端，旧基础配置仍可运行 |
| gateway/query-api/src/clients/mod.rs | 导出 collector 客户端 |
| gateway/query-api/src/http/handlers/mod.rs | 导出 watch_accounts handler |
| gateway/query-api/src/http/router.rs | 注册只读监控状态路由 |
| services/trade-log/adapters/tests/common/mod.rs | 复用隔离数据库，扩展分范围固定来源 |
| services/trade-log/adapters/tests/postgres_persistence.rs | 验证新 migration 与角色权限仍兼容 |
| services/trade-log/adapters/tests/replay_import.rs | 多页重放及旧证据回归 |
| tests/integration-tests/query_api.rs | 新状态接口、失败映射及旧基础接口回归 |

文档只新增本文件；实施后同步 docs/api.md，完成真实验收后新增 docs/version0.3-acceptance.md。本次不修改 README、roadmap 或详细设计说明，不在其中添加版本入口。

#### 直接复用的文件与约定

复用 account-facts、shared-types、service-runtime、现有 Hyperliquid parser、trade-log normalization/validation、stored_query 及 file_evidence。复用 0001 migration 而不修改已发布 SQL 的摘要；新增 0002 承接变更。复用已有 token 和证据卷，不另建版本目录或版本成交 ID。

collector 解析合并要保留每页原始引用和页内 source_index；不能为方便合并把所有事实指向同一个 raw_log_id。existing_records/inserted_records 的业务计数沿用 v0.2 语义。

## 8. 配置与异常处理

### 8.1 配置与权限

collector 配置沿用 config_version、server、logging、hyperliquid、evidence、internal、database 结构，新增第 3.1 节 collection。内部监听默认 0.0.0.0:8082，不对宿主机发布；网关 collector.base_url=http://trade-collector:8082，状态请求默认时限 5 秒。

租约必须大于 round_timeout，默认 180 秒对 120 秒；shutdown_timeout 默认 10 秒，Compose stop_grace_period 默认 15 秒。interval、overlap、window、预算、重试均为正整数，safety_delay 允许 0；overlap 小于 max_window，retry_base 不大于 retry_max。配置未知字段、非法地址、非法起点、时限矛盾或缺失凭证启动失败。

新增 secrets/trade-log-collector-database-url，使用 trade_log_collector 运行角色，权限由 migration 授予业务 schema 的 SELECT/INSERT/UPDATE 及迁移记录 SELECT，无建表、删除和数据库管理权限。collector 不挂载迁移 URL。网关不挂载任何数据库凭证。

init-v0.3.sh 复用已有 v0.2 数据库管理员凭证和应用/迁移凭证，不重置密码或重建卷；仅在缺少 collector 凭证时生成，再创建对应角色并执行新版 migration。迁移通过 existing trade-log-query migrate 执行，新版本 schema 校验后再启动所有应用。

### 8.2 运行状态与异常

| 情况 | 状态与处理 |
| --- | --- |
| 首次初始化/恢复 | STARTING，未取得水位前不声称成功 |
| 正在采集 | RUNNING，last_attempt_at 更新 |
| 完整提交后等待 | WAITING，last_success_at 和水位更新 |
| 来源 429/5xx/网络故障 | RETRY_WAIT，记录错误、失败次数及 next_run_at，不推进水位 |
| 有界预算耗尽 | RETRY_WAIT，保留 pending_work，warning=ROUND_BUDGET_EXHAUSTED，不计完整成功或来源失败 |
| 来源毫秒满页 | FAILED，SOURCE_WINDOW_SATURATED，保留证据和旧水位 |
| 解析失败/事实冲突 | FAILED，沿用 INCOMPLETE_DATA/VERSION_CONFLICT，禁止覆盖已有事实 |
| 租约由其他实例持有 | STARTING，warning=LEASE_HELD，不读取来源 |
| 租约失效 | 停止本轮提交并重新获取检查点，不允许旧实例写新水位 |
| 数据库不可用 | 无成功响应或成功水位，状态接口 503；调度重连后恢复 |
| collector 进程不可达 | 网关状态接口 503；原 trade-events 和网关 health 可继续独立响应 |
| 正常退出 | STOPPED，释放自有租约，保留记录 |

日志使用既有结构化格式，记录 service、account_key、query_id、范围、页号、lease_epoch、耗时、入库计数和错误码；不打印 token 或数据库 URL。last_success_at 必须来自成功事务提交结果，而不是请求开始、HTTP 200 或仅原始响应归档。

## 9. 验证与验收

### 9.1 必须通过的开发验证

1. 保留 v0.0–v0.2 的基础 HTTP、live/stored、证据导入、去重、时间范围和分页测试。
2. 固定来源验证首次范围、重叠范围、空成交水位推进、迟到和乱序事实，避免“无成交”掩盖后台没有执行。
3. 在原始归档、解析、事实事务、水位更新和提交后响应丢失处注入失败；未提交时水位不变，已提交时恢复读取正确水位。
4. 满页拆分及 1 毫秒饱和、预算暂停恢复、所有页事实引用与重新解析一致；计数按来源记录判断。
5. 并发 collector、租约续期失败和旧实例恢复不能重复领取或覆盖水位。
6. 重启读取旧水位和 pending_work；query 进程重启不会中断正在运行的 collector 任务。
7. collector 状态鉴权、网关 503 映射、服务健康与采集失败区分，以及数据库连接恢复。
8. 增量迁移前后旧事实数量和身份一致、角色权限正确；旧 v0.1 文件证据仍可导入。

定时测试使用可控时钟和固定来源，避免依靠真实等待或真实账户恰好发生交易才能验证。PostgreSQL 测试使用隔离数据库，真实来源验证单独记录，不用构造数据冒充真实成交。

### 9.2 版本完成标准

Docker 启动后不调用 live，collector 能自动生成至少两个成功采集任务，状态可查询，原始响应和事实可追溯。已存在成交反复采集不重复入账；安静账户仍有新的成功时间和水位。

重启恢复、水位与事实原子性、限流退避、故障可见、恢复后继续采集以及来源满页边界均验证通过。实际新成交若验收期间未发生，不能声称已经验证“新成交被自动发现”；该部分使用固定来源测试证明机制，真实持续更新待观察结果记录。

## 10. 交付、部署与手动验证

### 10.1 部署与升级

常驻容器为 query-api、trade-log-query、trade-collector、postgres；evidence-init、trade-log-migrate 为一次性初始化服务。仍由 Docker Compose 启动，仅网关发布 8080。collector 挂载自己的配置、共享内部 token、专用数据库 URL 和现有证据卷。

PostgreSQL 沿用 v0.2 的镜像、数据库及 trade-log-postgres 命名卷，本版不升级数据库主版本。证据仍使用 trade-log-evidence；raw ID 含唯一任务及页身份，两应用不会以相同任务名覆盖文件。

升级前备份数据库和旧证据。先停止旧应用，执行新版初始化和 0002 migration，再启动新版应用。不要在旧 query 程序仍运行时迁移，因为旧程序严格检查 schema 版本。停止应用命令不包含 postgres，不删除卷。

以下命令及步骤在实现交付后可执行；当前文档不会创建容器、凭证或数据库。

### 10.2 手动验证步骤

#### 第 1 步：填写账户、升级并启动

在仓库根目录填写 config/trade-collector.toml 的实际账户和起点，建议首次验收选取近期确有合约成交的账户和较近起点。

```sh
docker compose stop query-api trade-log-query
sh scripts/init-v0.3.sh
docker compose up --build -d
docker compose ps -a
docker compose logs --tail=40 query-api trade-log-query trade-collector
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/version
```

**看什么：** 四个常驻容器 running，migration 成功，version=0.3.0，collector 有启动和采集日志。

**通过标准：** 正确加载账户、读取数据库并开始调度；旧库事实和卷保留。

#### 第 2 步：不调用 live，查看自动更新

```sh
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/watch-accounts
docker compose exec -T postgres psql -U robotech_admin -d robotech -c "SELECT query_id,status,created_at FROM trade_log.collection_jobs WHERE job_origin='COLLECTOR' ORDER BY created_at DESC LIMIT 5;"
```

等待至少两个配置轮询间隔，再重复上述命令。

**看什么：** 配置账户正确，last_success_at 和 scanned_through 更新；有至少两个 COMPLETED 自动采集任务。无成交时 last_trade_at 允许为空。

**通过标准：** 自动更新不依赖用户主动查询，状态来自真实任务与检查点。

#### 第 3 步：查看自动保存的成交

```sh
curl --noproxy '*' -i --get http://127.0.0.1:8080/api/v1/trade-events \
  --data-urlencode 'account=<配置中的实际账户>' \
  --data-urlencode 'source=stored' \
  --data-urlencode 'limit=10'
```

**看什么：** 自动任务保存的事实可以查询；若该账户在观察期间有新成交，下一次成功轮询后可看到相同 fact_id 和来源信息。

**通过标准：** 无需 live 查询即可入库；实际新成交部分与官方记录核对。库中已有旧成交不能单独证明本轮采集生效。

#### 第 4 步：验证重复采集不重复事实

```sh
docker compose exec -T postgres psql -U robotech_admin -d robotech -c "SELECT count(*) AS facts FROM trade_log.account_facts_current;"
docker compose exec -T postgres psql -U robotech_admin -d robotech <<'SQL'
SELECT c.account_key, v.payload->'payload'->>'market' AS market,
       c.source_tid, count(*) AS records
FROM trade_log.account_facts_current c
JOIN trade_log.account_fact_versions v
  ON v.fact_id=c.fact_id AND v.revision=c.current_revision
WHERE c.fact_type='TRADE'
GROUP BY c.account_key,v.payload->'payload'->>'market',c.source_tid
HAVING count(*) > 1;
SQL
```

等待多个重叠轮次，再执行一次。必要时按配置 account_key 限定统计，避免其他手动采集账户影响总数。

**看什么：** 查重返回 0 rows；无新成交时事实数不增加，有新成交时只增加新的稳定身份。自动任务和 raw_logs 增加是正常证据留存。

**通过标准：** 重叠轮询未重复入账。

#### 第 5 步：验证应用与采集进程重启恢复

先记下状态中的 scanned_through、last_success_at 和最近成功 query_id：

```sh
docker compose restart trade-collector trade-log-query
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/watch-accounts
```

等待租约释放或到期，以及下一轮成功后重复状态和 stored 查询。

**看什么：** 原水位不被配置起点覆盖，继续生成新成功任务；旧事实不丢失、不重复。query 服务没有把 collector 活跃任务错误标为 INTERRUPTED。

**通过标准：** 重启从数据库继续，而不是从头重复建立监控。

#### 第 6 步：验证数据库故障和恢复

仅在本项目独立验收环境执行。先记下水位与成功时间：

```sh
docker compose stop postgres
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/watch-accounts
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/health
docker compose logs --tail=30 trade-collector
docker compose start postgres
```

恢复后等待有界退避及租约恢复，再查看状态。

**看什么：** 数据库不可读时状态接口 503，网关 health 200；collector 不记录虚假成功。恢复后旧记录仍在，水位继续推进，失败次数在完整成功后归零。

**通过标准：** 故障不丢事实、不误推进水位，恢复不需要删除检查点。

#### 第 7 步：停止 collector，验证职责独立

```sh
docker compose stop trade-collector
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/watch-accounts
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/health
docker compose start trade-collector
```

collector 停止期间再执行第 3 步的 stored 查询。

**看什么：** watch-accounts 返回 503，health 和可用数据库下 stored 仍为 200；collector 恢复后继续原账户。

**通过标准：** 采集故障不阻塞已有库存查询。

#### 第 8 步：重新解析自动采集任务

用第 2 步的成功任务 ID 替换下方占位符：

```sh
docker compose exec -T trade-log-query trade-log-query reparse --query-id query_实际值
```

**看什么：** comparison=SAME、退出码 0；多页任务保留每页来源引用，不因 reparse 修改水位或事实版本。

**通过标准：** 自动采集的原始响应能够复算相同事实。

来源满页、毫秒饱和、事务中断和租约竞争由固定来源集成测试验收，不要求用户靠真实账户碰巧发生这些情况；实际验收报告分别记录开发测试、真实来源观察和 Docker 部署结果，不把未执行项写为通过。

### 10.3 交付清单

交付 collector 程序及镜像 target、单地址配置、检查点 migration、范围采集与调度、原子水位提交、租约恢复、只读状态接口、新初始化脚本、必要测试和统一接口文档更新。验收完成后提交独立验收记录，写明数据来源、观察时间、水位变化和未完成项。

本版不新增自动验收脚本，手动验证按本章执行。后续版本继续沿用 trade-log 的 acquisition/checkpoint/persistence 和事实身份；新增 WebSocket 或信号发布时扩展适配器与提交能力，避免重建采集和存储基础。
