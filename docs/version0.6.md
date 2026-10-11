# v0.6 开发文档：多地址监控

- 文档版本：0.6
- 状态：已实现，本机开发验证通过；Docker 和真实来源多账户验收待服务器执行
- 更新日期：2026-10-11
- 设计依据：[概要设计](overview-design.md)、[详细设计](detailed-design.md)、[版本路线图](roadmap.md)
- 前置版本：[v0.5 开发文档](version0.5.md)
- 接口契约见 [统一接口说明](api.md)，不另建版本接口文档

## 1. 版本目标

在 v0.5 单账户 HTTP 补偿、WebSocket 实时采集和单 webhook 发布的基础上，同时监控多个 Hyperliquid 账户。通过接口添加、暂停和恢复监控，分别查询采集水位、连接状态、缺口和发布情况。

复用已有事实身份、原始证据、账户租约、采集水位、候选策略和 outbox，不为每个账户创建一套程序或数据库。保持现有多进程部署：一个 trade-collector 管理多个账户任务，一个 trade-parser-publisher 处理多个账户的发布控制和队列。

应用版本升级为 0.6.0。本版完成账户管理和任务隔离，不扩展现货、账户分析、跟单执行或多接收方订阅服务。

## 2. 功能范围

| 能力 | 本版交付 | 边界 |
| --- | --- | --- |
| 添加账户 | 提交地址和采集起点，持久保存监控配置 | 只支持配置网络上的 Hyperliquid 永续成交 |
| 暂停账户 | 停止该账户的 HTTP 和 WS 新采集，关闭新候选准入 | 保存事实、水位、证据和既有 outbox，不删除账户 |
| 恢复账户 | 从原水位继续补偿并重新建立实时会话 | 不批量发布暂停期间的历史成交 |
| 分账户状态 | 列表、详情、采集状态和发布状态 | 一个账户失败不使全部列表不可读 |
| 故障隔离 | 独立租约、任务、队列、退避和错误 | 数据库、进程和出口 IP 仍是共享依赖 |
| 容量限制 | 默认最多 5 个活动账户，可配置但不超过 10 | 暂停账户不占活动额度；登记总量默认最多 100 |
| 可靠发布 | 多账户沿用一个固定 webhook 目标 | 账户间公平领取，事件身份和正文不变 |
| 运维验收 | 固定脚本、隔离样本、多账户手动步骤 | 固定样本与真实来源验收分别报告 |

不提供账户硬删除、修改已有账户的历史起点、动态修改网络、按账户配置不同 webhook 或多租户权限系统。详细设计中的 DELETE 停止关注，本版实现为持久暂停而非删除；另提供明确的 resume 操作。

## 3. 使用方式

### 3.1 启动和管理

沿用 Docker Compose 和原数据卷。升级后先迁移原单账户，再通过 HTTP 添加其他账户；正常添加、暂停、恢复不要求重启容器或编辑 TOML。

账户清单以 PostgreSQL 为运行时唯一来源。TOML 保留来源、采集默认策略、资源上限和首次导入配置，不能在每次启动时覆盖接口管理结果。

### 3.2 管理接口的访问

读接口沿用既有访问方式。新增和变更监控会产生持续外部请求和存储开销，写接口必须使用独立管理凭证，网关验证后使用内部服务凭证转发；管理凭证不传给来源或写入日志。

管理凭证以只读文件挂载，不放在 URL、示例正文或仓库中。直接访问公网时由部署反向代理提供 HTTPS；本机手动验收使用回环地址。

### 3.3 数据和信号查询

成交查询继续使用 account 参数。发布状态支持指定账户；未指定时返回多账户汇总，不能静默选择第一个账户。现有单账户调用方在升级前需适配本文第 5 节的列表结构。

webhook 继续接收 v0.5 标准事件，从事件中的账户身份区分来源，不增加按账户的接收地址或新的事件身份算法。

## 4. 来源与候选边界

### 4.1 来源限制与本版容量

截至文档编写时，官方按出口 IP 限制 REST 总权重 1200/分钟、WS 连接 10 个、新连接 30 次/分钟、用户订阅涉及的独立用户 10 个。userFillsByTime 等查询还按返回条数增加权重。本版按官方公布限制设计，不因实测可能更宽松而放大容量。来源依据：[官方限制](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/rate-limits-and-user-limits)。

每个活动账户复用既有独立 WS 连接，默认 5 个活动账户，为共享 IP 上的其他调用保留余量。max_active_accounts 允许 1–10，但上限不是来源容量保证；同 IP 的其他部署、人工 live 查询和外部程序也会占用配额。

collector 与 trade-log-query 的所有官方 HTTP 请求共用 PostgreSQL 中的权重预算，默认 600/分钟，硬配置上限 1200；预算是本部署的保守调度值，不声称管理同 IP 的外部程序。请求基础权重和返回条数附加权重按官方规则计算；发送前按请求可能返回的最大条数预留，成功后依据实际条数结算，失败或响应不确定时保留预留消耗。预算不足等待下个可用窗口，日志记录 source_budget_wait 及账户、传输类型，不推进水位；等待达到轮次上限时状态记录 ROUND_TIMEOUT。

预算使用数据库时间、跨进程原子预留和滚动 60 秒计量，不能由进程重启清空。官方规则可能变化，权重表应集中定义并有固定测试；实现时重新核对官方说明。

### 4.2 连接与请求调度

WS 新连接共享滚动一分钟预算，默认 15 次，最大 30 次；首次连接和重连均计数，耗尽时等待。账户仍沿用 v0.4 的 5 秒起始、最大 60 秒有界退避及抖动；成功后的重置条件不变，避免多个账户同时重连造成连接风暴。

HTTP 每个账户最多一个采集轮次，全局默认最多 2 个并发来源请求，轮次到期按账户轮转，补偿大窗口必须分批让出调度机会。共享限流等待不计入实际外部请求超时，但不能无限持有数据库事务或阻塞租约续租。

### 4.3 暂停和恢复的信号边界

暂停时立即在数据库中禁止该账户新候选；worker 随后退出。暂停前已经生成的 outbox 继续按原目标交付，可能晚于暂停响应送达。暂停监控不等于撤回已经产生的信号，也不表示停止下游跟单。

恢复时生成新的账户监控代次和发布激活边界，重新订阅并建立本次快照边界。HTTP 补偿、初始快照、暂停前成交及恢复激活前的迟到成交只保存事实；只有恢复后符合 v0.5 全部条件的实时主动成交可产生候选。实时发布不等待完整历史补偿，缺口状态单独披露。

官方订阅协议提供订阅确认及初始快照标记；本版沿用 v0.5 对显式标记和快照后缺省标记的处理，不把新账户或恢复账户的快照当作实时信号。来源依据：[官方订阅说明](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/websocket/subscriptions)。

## 5. 数据与接口定义

### 5.1 监控账户模型

新增 watch_accounts，唯一键为 account_key，ID 使用稳定 UUID。地址按现有规则校验和归一化；网络和协议参与账户身份。

| 字段 | 含义 |
| --- | --- |
| id、account、account_key | 管理 ID、规范地址及既有统一账户身份 |
| start_time | 首次登记时请求的历史起点，已有水位不回退 |
| desired_state | 用户期望状态 ACTIVE 或 PAUSED |
| observed_state | STARTING、RUNNING、DEGRADED、PAUSED；带观测时间，PAUSED 表示旧租约已撤销，旧连接可能仍在关闭 |
| generation | 每次有效暂停或恢复递增的管理代次 |
| created_at、updated_at、paused_at、resumed_at | 管理时间，UTC RFC3339 |
| last_error、last_error_at | 脱敏的账户级错误；不替代原采集错误字段 |

observed_state 根据当前代次的 worker 心跳和租约计算；进程退出后不得永久展示 RUNNING。ACTIVE 账户即使故障或等待预算，仍占用活动额度，不能借故障突破容量限制。

已有 checkpoints、WS sessions、publishing_control 继续按 account_key 隔离。增加监控 generation 到任务提交校验；WS 会话与创建时的 generation 绑定。旧任务不得在暂停或恢复之后提交新事实、推进水位或创建候选；来源原始文件若已经写出，可保留为未完成证据，不伪造完成任务。

### 5.2 外部接口

所有结果沿用 data/meta 信封、trace_id 和统一错误结构。列表按 created_at、id 稳定排序，默认 limit=20、最大 100，使用不透明 cursor 分页。

| 方法与路径 | 请求 | 成功响应 |
| --- | --- | --- |
| POST /api/v1/watch-accounts | JSON account、start_time | 201，返回账户对象；仅表示登记成功，不表示已追平 |
| GET /api/v1/watch-accounts | 可选 state、limit、cursor | 200，data.items、next_cursor、total、active_count、max_active_accounts |
| GET /api/v1/watch-accounts/{id} | 管理 ID | 200，账户对象及 collection、publishing 状态 |
| DELETE /api/v1/watch-accounts/{id} | 无正文 | 202，期望 PAUSED，返回当前停止进度 |
| POST /api/v1/watch-accounts/{id}/resume | 无正文 | 202，期望 ACTIVE，返回当前启动进度 |
| GET /api/v1/publishing-status | 可选 account | 指定账户返回原账户级结构；不指定返回 items 和队列汇总 |

POST start_time 必须是非未来的有效 RFC3339 时间，时区明确，金额和事实模型不变。重复添加相同规范账户返回 409 WATCH_ACCOUNT_EXISTS，并提供 existing_id，不隐式恢复、不覆盖起点。暂停已暂停账户、恢复已活动账户保持幂等，不递增 generation，不重新激活发布。

管理操作对同一账户加行锁，容量检查与新增/恢复处于同一事务；并发请求不得超过上限。数据库提交成功即为写操作成功；调度器随后执行，调用方通过详情轮询最终状态。

### 5.3 内部接口和权限

collector 内部增加对应账户列表、详情和管理路由，使用既有内部服务认证。网关只做外部认证、参数转换和转发，不直接访问业务表。publisher 的账户级状态按 account 查询，不再从 collector TOML 固定选择一个地址。

collector 角色可管理 watch_accounts、账户租约及共享预算，publisher 对账户注册表只读，继续只写发布控制、队列和尝试。trade-log-query 增加预算预留权限，不获得监控账户管理权限。migration 在原角色基础上增量授权，不重新初始化库。

### 5.4 覆盖与状态

每个账户分别返回原水位、最后成功时间、失败次数、WS 会话状态、待补偿区间和 coverage 状态。恢复后保留未确认覆盖区间；若来源窗口已不足，标记覆盖不足，不能用最新水位证明历史完整。

列表查询不逐账户同步调用来源接口；从数据库读取持久状态与心跳。共享数据库不可用返回 503；单账户 WS 或来源失败时列表仍返回 200，并明确其 DEGRADED 和错误原因。

## 6. 处理流程

### 6.1 首次迁移和启动

1. 执行增量 migration，新增账户注册表、监控代次和共享预算结构，保留原事实与数据卷。
2. 原单账户以现有 checkpoint 为依据登记；不覆盖水位、发布目标、activation_epoch 或已有正文。
3. 发布器改为扫描账户注册表；原活动账户首次迁移沿用原激活边界，普通进程重启不重新激活。
4. collector 周期读取 ACTIVE 清单，按公平顺序启动缺少的账户 worker，并获取原有独占账户租约。
5. 全局发布开关开启时，为新增或恢复账户建立发布控制，再允许当前代次 WS 生成候选；控制未就绪时保存事实并披露未准入原因。

### 6.2 添加账户

认证、参数校验、去重和容量检查后写入账户记录，初始 generation=1。起点用于首次 HTTP 扫描，独立 WS 同时建立实时观察；两条通道继续去重。订阅或 HTTP 失败只退避该账户，管理登记不会因为暂时来源失败回滚。

### 6.3 暂停账户

在同一事务中锁定账户，设置 PAUSED、递增 generation 并关闭候选准入；事实提交和候选判定取得同一账户锁并校验 generation，保证暂停边界有明确先后顺序。暂停事务之前已提交的事实和候选保留，之后旧 generation 的提交必须拒绝。

暂停事务立即撤销旧租约，接口返回 202 和 PAUSED；调度器随后取消请求并关闭连接。PAUSED 表示已无有效提交权限，不承诺物理连接已关闭。再次恢复不得同时持有两个有效 worker，旧任务始终受 generation 和 lease_epoch 双重约束。

### 6.4 恢复账户

原子检查活动额度，设为 ACTIVE，递增 generation；保留原水位和缺口，创建新的激活边界，不重置已存在事件。worker 获取新租约并重新建立 WS 会话，HTTP 从已确认水位减 overlap 继续补偿；没有旧水位时使用原 start_time。

### 6.5 多账户发布

发布器扫描所有登记账户，包括已暂停账户的既有 outbox。每账户最多一个在途发送，全局默认并发 2，按账户轮转，不让持续失败或大量积压账户占满领取机会。网络请求不持有数据库事务；现有领取租约、确认 fencing、固定正文和至少一次语义不变。

发布控制的有效启用条件为全局 publishing.enabled 且账户 ACTIVE。关闭全局发布沿用 v0.5 语义；账户暂停禁止新候选，但已入队事件继续交付。全局开关、账户代次和激活时间更新必须幂等，publisher 定期同步不能把管理暂停重新打开。

### 6.6 故障隔离和退出

账户任务返回错误或 panic 时由 supervisor 记录、回收并按规则重建，不能导致 collector 主进程退出。每账户复用原有有界 WS 队列，另设进程级待处理字节上限；超限时停止受影响连接并记录补偿需求，不能无限增长内存。

共享 DB 故障影响所有账户，明确标为共享依赖故障；恢复后各账户独立续采。进程退出先停止管理变更和新领取，再并行取消 worker、等待在途事务、释放租约；有界退出，超时由租约恢复，不删水位或队列。

## 7. 程序结构设计

### 7.1 进程与职责

| 容器 | 本版职责 |
| --- | --- |
| query-api | 管理接口认证、监控列表/详情/操作转发、多账户发布状态 |
| trade-log-query | 既有 live/stored/证据/重放，共享来源 HTTP 预算 |
| trade-collector | 监控注册表、调度器、多账户 worker 和暂停 fencing |
| trade-parser-publisher | 多账户控制同步、公平领取和同一 webhook 投递 |
| postgres | 原库增量表、账户状态、共享预算及已有事实和 outbox |

不新增常驻容器，不按账户复制 Compose 服务。业务状态和端口位于 trade-log，数据库与网络实现位于 adapters，binary 只负责配置和装配。

### 7.2 新增文件

| 文件 | 职责 |
| --- | --- |
| services/trade-log/src/watch/mod.rs | 账户管理模型、校验、状态和管理端口 |
| services/trade-log/src/acquisition/budget.rs | 来源预算和权重策略端口 |
| services/trade-log/adapters/src/postgres/watch.rs | 注册、容量、代次及暂停事务 |
| services/trade-log/adapters/src/postgres/source_budget.rs | 跨进程滚动预算预留和结算 |
| services/trade-log/adapters/src/source_budget.rs | SourceReader/RealtimeSource 预算包装、共享 HTTP 在途槽位和取消释放 |
| services/trade-log/adapters/src/collector_supervisor.rs | 多账户调度、任务回收及资源上限 |
| services/trade-log/migrations/0005_multi_account_monitoring.sql | 增量表、generation、索引和角色权限 |
| services/trade-log/tests/watch.rs | 管理状态、参数及幂等测试 |
| services/trade-log/adapters/tests/multi_account_chain.rs | 隔离、容量竞争、暂停恢复和迁移测试 |
| scripts/init-v0.6.sh | 增量迁移、管理凭证和旧账户导入 |
| scripts/verify-v0.6.sh | 服务器验收入口 |
| scripts/verify-v0.6.py | 多账户状态、持续更新、重启和停库验证 |
| scripts/verify-v0.6-fixtures.py | 复用隔离资源方式运行多账户固定测试 |
| scripts/run-v06-fixtures.sh | 精简测试镜像中的固定样本入口 |
| scripts/tests/test_verify_v06_fixtures.py | 验收异常时只清理本次测试资源 |
| docs/version0.6.md | 本开发文档 |
| docs/version0.6-acceptance.md | 开发验证及服务器待验记录 |

### 7.3 更新文件

| 文件 | 修改 |
| --- | --- |
| Cargo.toml、Cargo.lock | 版本 0.6.0 和必要依赖 |
| services/trade-log/src/lib.rs、src/acquisition/mod.rs | 导出监控配置与预算权重策略 |
| services/trade-log/adapters/src/lib.rs、src/postgres/mod.rs、src/postgres/migration.rs | 模块导出、迁移及权限 |
| services/trade-log/adapters/src/collector_runtime.rs、src/websocket_runtime.rs | 账户级取消、预算等待及 generation 传递 |
| services/trade-log/adapters/src/postgres/checkpoint.rs、src/postgres/realtime.rs | 提交 fencing、账户状态及暂停原子边界 |
| services/trade-log/adapters/src/collector_http.rs | 管理、列表、详情内部接口 |
| services/trade-log/adapters/src/publisher_runtime.rs、src/publisher_http.rs、src/postgres/publishing.rs | 多账户公平发布、控制同步及状态 |
| services/trade-log/bins/trade-collector/src/config.rs、src/bootstrap.rs | 默认策略、容量及 supervisor 装配 |
| services/trade-log/bins/trade-parser-publisher/src/config.rs、src/bootstrap.rs | 移除固定账户选择，装配多账户控制和调度 |
| services/trade-log/bins/trade-log-query/src/bootstrap.rs | live 来源请求预算装配 |
| gateway/query-api/Cargo.toml、src/config.rs、src/bootstrap.rs、src/state.rs | 管理凭证、恒定时间校验依赖及客户端装配 |
| gateway/query-api/src/clients/collector.rs、src/clients/publisher.rs | 管理调用、错误码及多账户状态 |
| gateway/query-api/src/http/handlers/watch_accounts.rs、src/http/handlers/publishing.rs、src/http/router.rs | 外部管理及状态路由 |
| gateway/query-api/src/http/error.rs | 支持服务透传错误码 |
| config/trade-collector.toml、config/trade-parser-publisher.toml、config/query-api.toml | 新配置、旧字段兼容和凭证引用 |
| compose.yaml、gateway/query-api/Dockerfile | 版本、管理凭证、测试 target；保留原数据卷 |
| scripts/verify-version.py | 旧版本脚本读取 collection 嵌套，显式查询所选账户发布状态 |
| services/trade-log/adapters/tests/publishing_chain.rs 及既有启动/集成测试 | 覆盖多账户后的原行为回归 |
| docs/api.md、docs/verify.md | 实现时更新统一契约和简洁执行命令 |

以上为本版文件归属。既有测试中 Postgres 构造增加共享内存资源字段；collector 启动测试支持空账户清单，旧集成测试适配新列表结构。旧版本验收默认选择首个登记账户，必要时显式 --account；v0.6 验收读取全部账户。

## 8. 配置与异常处理

### 8.1 配置示例

以下为 collector 新增配置；既有 server、logging、database、来源、采集默认策略和 websocket 配置继续保留。

```toml
[monitoring]
max_active_accounts = 5
max_registered_accounts = 100
reconcile_interval_seconds = 2
worker_stale_after_seconds = 30
max_concurrent_http_requests = 2
max_total_pending_bytes = 67108864

[source_budget]
http_weight_per_minute = 600
websocket_connections_per_minute = 15
```

共享预算参数在数据库保存一份，由初始化/配置同步写入，query 服务读取同一份参数，不各自维护独立配额。活动账户上限不超过 10，总登记上限不小于活动上限；并发数、内存和心跳阈值均校验合理范围，数据库连接池须支持续租与预算等待期间的其他任务。

旧 collection.account/start_time 仅用于一次性导入原账户，完成后标记导入记录；重启不重新添加或恢复该账户。新部署可以没有导入账户，以空清单启动后通过接口登记。旧字段保留一个兼容周期，显式废弃，不静默改变含义。

publisher 保留全局目标、策略和 enabled，新增 max_inflight=2；顶层 network 默认为示例中的 mainnet，须与 collector/query 的来源网络一致。不再从 collector 配置选择固定账户；旧 collection_config_path 仅作为未设置 network 时读取网络的兼容路径。管理凭证配置位于 query-api 的独立 management 段，默认路径 /run/secrets/watch-management-token。

### 8.2 错误及处理

| 情况 | 状态与行为 |
| --- | --- |
| 管理凭证缺失/错误 | 401 AUTHENTICATION_REQUIRED，不修改状态 |
| 地址、起点、ID 或分页参数非法 | 400 INVALID_ARGUMENT |
| 账户不存在 | 404 WATCH_ACCOUNT_NOT_FOUND |
| 重复登记 | 409 WATCH_ACCOUNT_EXISTS，返回既有 ID |
| 活动或登记容量已满 | 409 WATCH_ACCOUNT_LIMIT_REACHED，不创建半完成记录 |
| 参数未知或试图修改原起点 | 400，拒绝未支持字段 |
| 来源 429 或共享预算不足 | 保留水位和任务，等待；预算等待记录 source_budget_wait，来源错误保留原错误原因，不变为管理接口容量错误 |
| 单账户异常 | 账户 DEGRADED，其他账户继续运行 |
| 数据库故障 | 依赖该库的接口 503 DEPENDENCY_UNAVAILABLE；网关基础 health 存活 |
| 旧代次提交 | 拒绝提交，终止旧 worker；不算来源数据已保存 |
| 旧会话归档恢复 | 可恢复事实，候选判定记录 MONITORING_GENERATION_UNCONFIRMED；不冒充新信号 |
| webhook 失败 | 分账户队列和错误可见，不阻塞其他账户采集或可发送事件 |

## 9. 验证与验收

### 9.1 固定开发验证

1. 地址归一化去重、起点校验、未知字段、认证、分页、重复操作幂等。
2. 并发添加或恢复争抢最后一个额度，恰有一个成功；故障中的 ACTIVE 账户仍占额度。
3. 两个账户独立水位、租约和会话；同 tid 的不同账户按原身份规则保存，不串账。
4. 账户 A 超时、WS 异常、队列溢出或 worker panic，B 在固定来源样本下持续提交和续租。
5. 暂停与事实提交并发：事务边界前的提交保留，边界后的旧代次不能入库、推进或生成候选。
6. 立即暂停再恢复、反复重启和租约接管，不产生双有效 worker，不沿用旧会话快照边界。
7. 恢复 HTTP 补偿和快照无候选；新激活之后的合格实时成交有候选；覆盖不足明确展示。
8. 暂停前 outbox 继续送达，正文和时效不变；A 大量积压或 500 不饿死 B。
9. HTTP 与 live 共享预算、返回条数结算、未知响应保守扣费、滚动窗口、WS 重连风暴和数据库故障。
10. 旧账户迁移不回退水位、不重置激活、不改写事件；迁移和初始化重复执行不重复登记。
11. 列表结构及错误透传正确；日志带 account_key、generation、lease_epoch 和 trace，不含凭证。
12. 原 v0.0–v0.5 能力回归；固定测试不依赖部署 TOML 的地址和开关，也不依赖精简镜像未提供的工具。

### 9.2 完成标准

至少两个真实账户同时监控并分别可查询；其中一个暂停和恢复时另一个保持采集，暂停区间在来源可查询范围内补偿，补偿不生成历史候选。真实来源无成交不能记为实时采集或发布通过。

固定测试必须确定覆盖单账户故障隔离、容量竞争、暂停提交竞争、公平发布和限流；PASS、FAIL、SKIP 分开报告。共享数据库故障不能作为“单账户异常不影响其他账户”的反例或通过依据，必须单独验收恢复。

## 10. 交付、部署与手动验证

### 10.1 部署与升级

备份数据库和凭证，在维护窗口停止旧采集及发布写入，执行 init-v0.6 增量迁移及原账户导入，再升级全部应用镜像。沿用 PostgreSQL 18、证据卷和 webhook 验收卷；初始化脚本不得 down --volumes，不重建旧凭证。

升级前检查原账户水位和发布目标，升级后核对原账户仍为活动状态且数据保留。原本已关闭发布时保持关闭。回退前停止新版本写入并备份；旧程序不保证兼容多账户 schema，不提供破坏性自动降级。

初始化会生成缺失的管理凭证，并复用原初始化流程运行当前迁移；保留旧凭证和数据卷。执行：

```sh
sh scripts/init-v0.6.sh
docker compose up --build -d
```

### 10.2 手动验证步骤

在服务器项目根目录执行。预先准备两个不同的有效账户及 RFC3339 起点，地址不填入本文或仓库。以下管理请求以宿主机 secrets/watch-management-token 读取凭证，示例中的地址、时间和 ID 需替换为实际值。

#### 第 1 步：检查升级后的原账户

```sh
curl --noproxy '*' -sS http://127.0.0.1:8080/api/v1/watch-accounts
curl --noproxy '*' -sS http://127.0.0.1:8080/api/v1/publishing-status
```

通过标准：原账户、事实、水位、队列和激活边界保留，列表显示当前活动数量和上限。

#### 第 2 步：添加第二个账户

```sh
curl --noproxy '*' -sS -X POST http://127.0.0.1:8080/api/v1/watch-accounts \
  -H "Authorization: Bearer $(cat secrets/watch-management-token)" \
  -H 'Content-Type: application/json' \
  -d '{"account":"<第二个实际账户>","start_time":"<RFC3339起点>"}'
```

记下返回 id。通过标准：返回 201，列表中两个账户分别进入运行或明确的异常状态；再次添加返回 409，不覆盖原起点。

#### 第 3 步：分别查询采集和发布

```sh
curl --noproxy '*' -sS 'http://127.0.0.1:8080/api/v1/watch-accounts/<实际ID>'
curl --noproxy '*' -sS 'http://127.0.0.1:8080/api/v1/publishing-status?account=<实际账户>'
docker compose logs --tail=100 trade-collector trade-parser-publisher
```

通过标准：各账户水位、会话、错误和队列独立，日志可定位账户；成交查询指定 account 后不混入另一个账户的事实。

#### 第 4 步：暂停一个账户

```sh
curl --noproxy '*' -sS -X DELETE \
  'http://127.0.0.1:8080/api/v1/watch-accounts/<实际ID>' \
  -H "Authorization: Bearer $(cat secrets/watch-management-token)"
```

轮询详情直到 PAUSED，记录暂停前水位和时间。通过标准：该账户停止新采集和候选，其他账户水位继续推进；已有 outbox 可继续投递。重复暂停不改变代次。

#### 第 5 步：恢复和补偿

```sh
curl --noproxy '*' -sS -X POST \
  'http://127.0.0.1:8080/api/v1/watch-accounts/<实际ID>/resume' \
  -H "Authorization: Bearer $(cat secrets/watch-management-token)"
```

通过标准：从旧水位继续、建立新会话和激活边界；HTTP 补偿只保存事实，不输出暂停期间历史信号；独立来源对账或固定样本证明补偿结果，不能只凭水位推进断言无缺口。

#### 第 6 步：重启和共享数据库故障

记录两个账户 ID、水位、管理状态和待发送事件，再分别重启 collector、publisher；允许中断的验收环境中再停库并恢复。

```sh
docker compose restart trade-collector trade-parser-publisher
# 下面两条仅在允许中断时执行，期间先检查业务 503 和网关 health。
docker compose stop postgres
docker compose start postgres
```

通过标准：重启保留清单、暂停状态和活动账户激活边界；数据库恢复后各账户续采，outbox 保留；不重复导入、不超额启动、不把回补当新信号。

### 10.3 固定服务器验收脚本

```sh
# 当前部署：检查清单、分账户状态、水位和队列
sh scripts/verify-v0.6.sh --wait-seconds 180

# 独立测试库和固定双账户样本：隔离故障、容量、暂停恢复及公平发布
sh scripts/verify-v0.6.sh --fixture-tests --wait-seconds 180

# 允许中断的环境：重启和共享数据库故障
sh scripts/verify-v0.6.sh --lifecycle --database-fault --wait-seconds 180
```

默认从数据库读取账户列表，不要求重复传 account。--fixture-tests 复用已有隔离执行方式，只清理本次创建的库、容器和凭证，不添加或暂停部署账户。真实列表不足两个活动账户时，多账户实测记为 SKIP，并给出原因；固定样本通过不能替代真实来源验收。

### 10.4 交付和 bug 记录

验证记录见 [version0.6-acceptance.md](version0.6-acceptance.md)，区分本机开发测试和服务器待验，不将固定样本通过记作真实来源通过。发现 bug 时按既有七项格式记录：描述、影响、确认记录、修改逻辑、修改记录摘要、如何验证、验证记录。
