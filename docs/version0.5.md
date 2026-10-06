# v0.5 开发文档：候选信号与 webhook 输出

- 文档版本：0.5
- 状态：已实现；本机开发验证通过，Docker 部署及真实来源投递待服务器验收
- 更新日期：2026-10-06
- 设计依据：[概要设计](overview-design.md)、[详细设计](detailed-design.md)、[版本路线图](roadmap.md)
- 前置版本：[v0.4 开发文档](version0.4.md)，保留其续租并发及数据库故障恢复修复
- 本版实现时，在 [统一接口说明](api.md) 中补充 webhook 契约和发布状态，不另建版本接口说明

## 1. 版本目标

在 v0.4 已保存的标准成交事实基础上，把实时观察到的主动永续成交作为候选信号，向一个配置的 webhook 接收服务输出。候选信号沿用标准账户事实，不新建另一套成交身份或金额模型；表示可进入下游跟单规则判断，不代表推荐、下单或保证成交。

事实、来源观察、候选判定和 outbox 在同一数据库事务中提交，独立发布进程随后投递。接收方故障不阻塞正常采集；发布不等待 HTTP 历史补偿、分析或完整性审核。系统记录成交时间、实时接收时间、事件保存时间、首次发送时间及成功确认时间，供接收方判断时效。

程序与应用镜像升级为 0.5.0。沿用原 PostgreSQL、事实身份、HTTP 水位、证据卷和外部查询接口，增加 `trade-parser-publisher` 容器。首版发布目标为 webhook，后续消息系统通过发布端口扩展，不重写采集和标准事实。

## 2. 功能范围

| 能力 | 本版交付 | 范围 |
| --- | --- | --- |
| 候选判定 | 根据来源模式、会话边界、时效及 copy_eligible 判定 | 单配置账户、默认永续市场，主动已成交操作 |
| webhook 输出 | 一个目标，每个请求发送一个标准事实事件 | JSON POST，配置文件指定地址和凭证文件 |
| 可靠发布 | 原子 outbox、领取租约、有界退避、重启恢复 | 至少一次投递语义，允许重复；永久拒绝须处理 |
| 时效标识 | 保存各阶段时间及 expires_at | 过时事件可继续交付，下游明确忽略过期候选 |
| 可追溯 | 关联事实版本、WS 消息、候选判定及投递尝试 | 可定位未生成、未发送、失败或重复的原因 |
| 状态查询 | 发布进程和队列状态 | 通过 query-api 查看，不开放任意修改队列接口 |
| 验收 | 固定接收器、服务器脚本及故障测试 | 不依赖活跃账户碰巧产生指定异常 |

本版不执行跟单，不计算跟随金额，不连接签名器，不增加分析或审核服务，也不实现多账户管理、现货候选、事实修订、Kafka 或 Redpanda。HTTP 回补和订阅快照继续保存，不作为本版 webhook 候选。未来完整事实流可复用事件信封和发布端口；本版 webhook 接收的是标准事实中的候选子集，不是全部账户事实。

## 3. 使用方式

### 3.1 配置并启动

新增 `config/trade-parser-publisher.toml`，配置一个 webhook 地址及凭证文件。账户和网络通过只读挂载的 `trade-collector.toml` 读取，不在发布配置中重复填写账户，不要求用户在验收命令中再次输入账户。

默认 `publishing.enabled=false`。填写目标并显式启用后，发布进程首次激活时以数据库时间保存该账户的 `activated_at`；既有历史事实不自动生成候选。重启不重置激活时间，停用和重新启用会创建新的激活代次。

```sh
sh scripts/init-v0.5.sh
docker compose up --build -d
curl --noproxy '*' -sS http://127.0.0.1:8080/api/v1/publishing-status
```

脚本执行增量 migration 并保留旧卷及凭证；发布默认关闭，需配置后启用。

### 3.2 查看结果

原 `source=stored` 查看成交事实；新发布状态接口查看队列和错误；接收端按 event_id 查看是否接收。stored 有成交不等于对应交易会发布，发布资格独立判定并保存原因。

发布进程停止期间，已启用的候选仍可由采集事务进入 outbox；进程恢复后继续交付。显式关闭 publishing 时，停止新候选入队和新发送领取，已有队列保留；已经在途的请求按关闭机制收尾。启停配置由发布控制记录持久化，不把进程存活等同于业务开关。

### 3.3 接收方要求

接收方校验凭证，把 event_id 与事件原文持久化后返回 2xx。重复 event_id 返回成功，不重复触发后续决策。接收成功只表示已接收事件，不表示已经跟单。

接收方必须检查 expires_at，并按自身规则决定接受、忽略或告警。事件已经过期时仍可保存并返回成功，但不能仅因 payload.copy_eligible=true 自动执行交易。

## 4. 来源与候选边界

### 4.1 来源契约及 UNKNOWN 的处理

官方说明 userFills 先发送快照，快照标记 isSnapshot=true；后续流式更新标记 false，而类型定义中的 isSnapshot 为可选布尔值。服务器实际记录中，初始快照为 SNAPSHOT，后续大量消息缺少该字段并保存为 UNKNOWN。[官方订阅说明](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/websocket/subscriptions)

不能把全部 UNKNOWN 改成实时消息，也不能要求所有实际增量必须携带 false。本版保留原始 message_mode，在候选判定中引入有版本的会话规则：

| 原始模式及上下文 | 有效判定 | 能否成为候选 |
| --- | --- | --- |
| SNAPSHOT | 快照 | 否，保存事实并记录 SNAPSHOT |
| LIVE_UPDATE，订阅已确认 | 明确增量 | 满足其余条件时可以 |
| UNKNOWN，字段确实缺失，同会话已有成功提交的初始快照，序号在该快照之后 | 快照之后的无标记更新 | 满足其余条件时可以，记录 POST_SNAPSHOT_UNFLAGGED |
| UNKNOWN，字段为 null、字符串或其他非法类型 | 模式无法确认 | 否，记录 MODE_UNCONFIRMED |
| UNKNOWN，尚无成功提交的初始快照 | 模式无法确认 | 否，等待后续快照建立边界，已有消息不补发 |
| 订阅确认之前收到的数据 | 会话边界无法确认 | 否，记录 SESSION_UNCONFIRMED |
| HTTP、旧证据导入、reparse | 历史或复算 | 否 |

“快照之后的无标记更新”是本系统依据官方顺序约定和实际来源行为采用的适配规则，不是官方对缺省字段作出的明确保证。规则必须由固定样本和实际原始消息验证；原始证据和判定原因保留，不修改 v0.4 已保存的消息模式。若来源不再满足该顺序约定，模式无法确认的成交继续入库，但不生成候选。

每次重新连接都重新建立会话边界，不能继承上一个会话的快照确认。只有快照任务成功提交后，才持久化本会话的 snapshot_sequence；解析失败或收到消息但未完成提交不算确认。

### 4.2 候选准入条件

必须同时满足：

1. publishing 控制记录已启用，目标已配置，事件所属账户和网络与配置一致。
2. 事实为 TRADE、PERPETUAL、UPSERT、revision=1，trigger_type=USER、copy_eligible=true，持仓效果为已确认的 OPEN、INCREASE、DECREASE、CLOSE 或 REVERSE。
3. 来源是已确认订阅中的 WS 增量，或第 4.1 节允许的快照后无标记更新；市场映射可靠，metadata_stale=false。
4. 交易发生时间及 WS received_at 均不早于当前激活时间。received_at 与 occurred_at 的间隔不超过 max_event_age_seconds，首次判定时交易相对数据库时间也未超过该上限。
5. 成交时间不超过数据库当前时间加 clock_skew_tolerance_seconds；未来超限、缺少时间或来源身份无法确认均不准入。
6. 标准事实、来源观察及必要证据能够成功提交，同身份业务字段不存在冲突。

默认 max_event_age_seconds=30，clock_skew_tolerance_seconds=5，signal_ttl_seconds=60。expires_at=occurred_at+signal_ttl_seconds；TTL 至少为准入年龄上限。时间阈值是本版产品默认值，不是官方接口限制。

历史补偿不能通过更新时间或重新解析转成候选。恢复旧 WS 任务时，按原 received_at、原会话证据和当前判定时间检查；已经过期的任务只保存事实。系统不等待恢复缺口归零才能发布新鲜且已确认的实时成交。

### 4.3 重复与 HTTP/WS 竞争

HTTP 可能先保存同一成交，WS 随后才完成提交。因此“事实本次是否新增”不作为唯一发布条件：符合条件的首次 WS 观察可以为已有标准事实创建 outbox。使用标准事实版本的 event_id 去重，只产生一条待发布事件。

HTTP 先入库、WS 先入库、WS 重复消息和快照重放都必须保持相同 fact_id、revision 和 event_id。同一事实的首次准入观察固定事件接收时间与来源引用；后续重复不更新旧事件或重复创建 outbox。

## 5. 数据与接口定义

### 5.1 webhook 请求

```http
POST <配置的 webhook 路径>
Content-Type: application/json
Authorization: Bearer <凭证文件内容>
X-Robotech-Event-Id: <event_id>
X-Robotech-Delivery-Attempt: <attempt_number>
```

请求体采用 `AccountFactEnvelope`：

```json
{
  "schema_version": 1,
  "event_type": "account.fact.v1",
  "event_id": "<fact_id>:1",
  "partition_key": "<account_key>",
  "occurred_at": "2026-10-06T03:00:00.000Z",
  "received_at": "2026-10-06T03:00:00.120Z",
  "stored_at": "2026-10-06T03:00:00.150Z",
  "published_at": "2026-10-06T03:00:00.200Z",
  "expires_at": "2026-10-06T03:01:00.000Z",
  "observation": {
    "transport": "WEBSOCKET",
    "query_id": "query_example",
    "session_id": "<session_uuid>",
    "message_sequence": 2,
    "raw_log_id": "<ws_raw_log_id>",
    "message_mode": "UNKNOWN",
    "realtime_reason": "POST_SNAPSHOT_UNFLAGGED",
    "publishing_policy_version": "candidate-v1"
  },
  "fact": {}
}
```

fact 必须是完整的既有 AccountFact，示例空对象仅表示省略展开，实际发送不允许为空。fact.payload.copy_eligible 必须为 true；价格、数量、手续费等沿用精确十进制字符串。fact.raw_log_id 可能引用首次 HTTP 观察，observation.raw_log_id 独立指向促成本次候选的 WS 证据，不覆盖原事实。

| 字段 | 定义 |
| --- | --- |
| event_id | 复用 account_fact_versions.event_id，稳定对应 fact_id 与 revision |
| occurred_at | 原始成交时间 |
| received_at | 首次准入 WS 观察的接收时间 |
| stored_at | 同一事务记录的事件保存时间，不冒充精确的事务提交时刻 |
| published_at | 第一次准备发送的时间，首次领取后持久化；重试不改写 |
| expires_at | 候选时效边界，重试不延长 |
| observation | 候选判定所用原始观察和规则版本 |
| fact | 完整标准账户事实，不另造跟单成交模型 |

首次发送前持久化包含 published_at 的最终 wire_body 及 SHA-256，后续重试发送完全相同的请求体。投递次数在请求头和尝试表中变化，不改变事件身份、时效或事实内容。成功确认时间 delivered_at 保存在发送方投递记录中。

### 5.2 数据库增量结构

新增 `0004_webhook_publishing.sql`，保留原 0001～0003 migration 和所有事实、水位、角色凭证。

| 表或字段 | 关键内容 |
| --- | --- |
| publishing_control | account_key、target_id、enabled、activation_epoch、activated_at、策略阈值、策略版本和更新时间 |
| publication_decisions | fact_id/revision、raw_log_id/source_index、session_id/message_sequence、activation_epoch、result、reason、decided_at；同一来源观察唯一 |
| outbox_events | event_id 主键、fact_id/revision 外键、topic、partition_key、target_id、payload、wire_body、body_sha256、状态、attempts、next_retry_at、published_at、delivered_at、lease_owner/epoch/expires_at、last_error、created_at/updated_at |
| delivery_attempts | event_id、attempt_number 联合唯一、领取及结束时间、HTTP 状态、结果、脱敏错误、发送租约代数 |
| collection_stream_sessions.snapshot_sequence | 成功提交初始快照的消息序号，可空 |

outbox 状态：PENDING、SENDING、RETRY_WAIT、DELIVERED、BLOCKED。建立状态与下次重试时间索引、账户与创建时间索引；同一 event_id 只允许一条 outbox。凭证原文不进入数据库或事件。

publication_decisions 保存 ELIGIBLE、ALREADY_ENQUEUED 或 SUPPRESSED，并记录 DISABLED、SNAPSHOT、MODE_UNCONFIRMED、SESSION_UNCONFIRMED、BEFORE_ACTIVATION、STALE_EVENT、FUTURE_EVENT、NOT_COPY_ELIGIBLE、METADATA_STALE 等原因。未能解析成事实的消息沿用采集失败记录，不伪造事实级判定。

旧事实和旧 UNKNOWN 消息不回填 outbox。原始 parse_status、HTTP 水位及 WS 消息提交位置仍按 v0.4 规则保存。

### 5.3 发布状态接口

新增内部 `GET /internal/v1/publishing-status` 和内部健康接口，由发布进程提供，沿用服务凭证认证。网关新增 `GET /api/v1/publishing-status`，返回既有 data/meta 信封。

状态包含 enabled、status、account/account_key、target_id、activation_epoch/activated_at、heartbeat_at、策略阈值、各 outbox 状态数量、最老待发送事件时间、最后发送及确认时间、最后错误、过期未送达数量与抑制原因计数。状态枚举：DISABLED、STARTING、RUNNING、DEGRADED、BLOCKED、STOPPED。

目标地址只返回脱敏摘要，不返回凭证、内部连接串或事件中的敏感请求头。发布进程不可访问时接口返回 503，不能以空队列代替不可用；health 和 stored 继续保持各自原有语义。状态只说明发布链路，不证明接收方执行交易。

## 6. 处理流程

### 6.1 初始化与激活

发布进程校验配置、凭证和数据库版本，从 collector 配置读取账户及网络，保存发布控制记录。首次启用或关闭后重新启用时，生成新的 activation_epoch 和 activated_at；同配置重启保持原代次。策略变更须记录新策略版本与激活代次，不能追溯改变已生成事件。

本版只有一个目标。target_id 固定关联配置地址；已有目标队列存在时，地址变化启动失败并显示 TARGET_CONFIGURATION_CHANGED，不能把旧队列静默发送到另一接收方。凭证可在同目标下更新，重新启动加载后重试。

关闭 publishing 的配置必须由发布进程成功加载并持久化，状态确认 DISABLED 后才表示新候选停止入队。单纯停止发布容器表示发送暂停，采集仍可继续入队。

### 6.2 事实及 outbox 原子提交

WS 原始消息先按既有流程归档，使用固定元数据解析及验证。事务内检查采集租约和会话，保存或复用标准事实、保存观察、更新 WS 提交位置，并计算 publication_decisions；准入时 INSERT outbox ON CONFLICT DO NOTHING。同一事务提交失败时，不留下成功事实对应的半完成候选判定。

若是快照成功提交，同时记录 snapshot_sequence；其自身全部抑制。若 HTTP 已保存该事实，仍比较业务内容后使用原事实版本和事件身份；冲突则整体回滚，沿用 VERSION_CONFLICT，不为冲突成交发布。

发布控制开关必须在同事务中读取并锁定。锁顺序沿用 ingestion_state → collection_job → checkpoint → publishing_control → 事实与 outbox；发布领取和确认不反向获取采集锁。外部 HTTP 请求在数据库事务之外执行，避免重现 v0.4 的锁等待与任务暂停问题。

### 6.3 领取和发送

发布器默认一次领取一条，最大并发 1，按 created_at、event_id 确定领取顺序。来源乱序时，该顺序不等同于成交时间顺序；接收方仍需使用 fact.ordering_key 处理来源顺序。用短事务及 FOR UPDATE SKIP LOCKED 领取到期 PENDING/RETRY_WAIT 或租约已过期的 SENDING；递增 attempts，记录本次 attempt 和租约，持久化首次 wire_body 后提交。领取租约默认 30 秒，单次请求超时 10 秒。

随后发送 webhook，不持有数据库锁。收到结果后，以 event_id、lease_owner 和 lease_epoch 条件更新状态。进程崩溃或发送成功后确认写库失败，租约到期后再次发送相同事件；接收方按 event_id 去重。

attempts 表示已领取的发送尝试，不表示成功送达次数。进程可能在领取之后、发送之前退出，尝试记录须区分未确认、失败和确认成功，不能把未知结果当成未送达或已送达。

### 6.4 响应与重试

| 结果 | 处理 |
| --- | --- |
| 2xx | 标记 DELIVERED，记录 delivered_at；接收方应已持久化或可靠接收 |
| 网络错误、请求超时、408、425、429、5xx | RETRY_WAIT；按 5、10、20、40……秒有界退避，最大 300 秒，附加不突破上限的小幅抖动 |
| 429/503 含 Retry-After | 解析秒数或 HTTP 日期，最早重试时间取退避与来源要求中较晚者；支持最长 24 小时的等待；更长有效要求标为 BLOCKED，不截短后提前请求 |
| 其他 4xx | BLOCKED，保存脱敏原因，修正目标或接收契约后显式重试 |
| 3xx | 不跟随跳转，BLOCKED，避免把凭证转交给另一目标 |
| 确认写库失败 | 保留可重新领取状态；重复发送可能发生 |

临时失败不设自动丢弃次数上限，重复次数持续累计。到达 expires_at 后仍可交付，接收方按固定时效判断忽略；发布器不能把过期事件改成新鲜事件。队列积压和最老事件年龄必须可查询。至少一次语义不等于所有永久拒绝的事件已经送达，BLOCKED 必须明确可见。

提供运维 CLI `trade-parser-publisher retry --event-id <event_id>`，仅允许把已处理原因的 BLOCKED 事件重新安排发送，保留 event_id、正文、首次发布时间及总尝试次数。不提供删除事实、改写金额或更换投递目标的操作。

### 6.5 故障和退出

接收方离线时采集和 stored 正常，outbox 持久积压并重试。数据库离线时不能入库或确认投递，沿用 v0.4 采集恢复，发布器在恢复后重新领取；既有待发送事件保留。正常停止先停止领取，再在关闭上限内等待在途请求及确认；超时保留未确认记录供后续重试。

实时采集、续租和发布重试各自持续得到轮询，不在持锁事务中等待另一个任务。投递状态和 HTTP 扫描水位彼此独立，webhook 成功不推进采集水位。

## 7. 程序结构设计

### 7.1 进程与职责

| 常驻容器 | 本版职责 |
| --- | --- |
| query-api | 既有接口及发布状态网关 |
| trade-log-query | live/stored、原始证据、重放和 migration |
| trade-collector | 原 HTTP/WS 采集及事务内候选判定、outbox 入队 |
| trade-parser-publisher | 发布控制、outbox 领取、webhook、投递确认和状态 |
| postgres | 沿用原库和数据卷 |

trade-parser-publisher 使用详细设计中的进程名。本版先承担可靠发布职责，既有解析和事实提交仍复用 collector 的实现；不在新进程重复解析并保存第二套事实。候选策略属于 trade-log，数据库和网络实现放在 adapters，后续解析进程分离时沿用这些端口。

### 7.2 新增文件

| 文件 | 职责 |
| --- | --- |
| `config/trade-parser-publisher.toml` | 发布进程、目标和策略配置 |
| `services/trade-log/bins/trade-parser-publisher/Cargo.toml` | 新 binary crate |
| `services/trade-log/bins/trade-parser-publisher/src/main.rs` | CLI 入口 |
| `services/trade-log/bins/trade-parser-publisher/src/lib.rs` | 服务模块导出 |
| `services/trade-log/bins/trade-parser-publisher/src/config.rs` | 配置校验、collector 账户选择 |
| `services/trade-log/bins/trade-parser-publisher/src/bootstrap.rs` | 启停、发布循环和 retry 命令 |
| `services/trade-log/src/publishing/mod.rs` | 候选策略、发布端口、状态 DTO |
| `services/trade-log/src/publishing/policy.rs` | 模式、激活和时效判定 |
| `services/trade-log/adapters/src/postgres/publishing.rs` | 控制、判定、outbox 和尝试事务 |
| `services/trade-log/adapters/src/publisher_runtime.rs` | 领取、发送、确认和退避 |
| `services/trade-log/adapters/src/webhook.rs` | 单目标 HTTP 投递适配器 |
| `services/trade-log/adapters/src/publisher_http.rs` | 内部健康、认证及状态接口 |
| `gateway/query-api/src/clients/publisher.rs` | 内部发布状态客户端 |
| `gateway/query-api/src/http/handlers/publishing.rs` | 外部发布状态 handler |
| `services/trade-log/migrations/0004_webhook_publishing.sql` | 本版增量表和权限 |
| `services/trade-log/tests/publishing.rs` | 候选策略固定样本测试 |
| `services/trade-log/adapters/tests/publishing_chain.rs` | 原子事务、跨通道竞争与故障验证 |
| `services/trade-log/bins/trade-parser-publisher/tests/startup.rs` | 配置、CLI 和生命周期测试 |
| `tests/fixtures/v0.5/` | 快照、无标记更新、主动与强制成交、过期和重复样本 |
| `scripts/init-v0.5.sh` | 保留旧凭证、增量角色及 migration |
| `scripts/verify-v0.5.sh` | 本版固定服务器验收入口 |
| `scripts/webhook-receiver.py` | 独立验收接收器，持久化去重及可控响应 |
| `scripts/verify-v0.5-fixtures.py`、`scripts/run-v05-fixtures.sh` | 独立测试数据库、接收器和测试镜像的执行及清理 |
| `scripts/tests/test_verify_v05_fixtures.py` | 验证失败时仅清理本次创建资源 |
| `docs/version0.5-acceptance.md` | 本机验证和服务器待验项目 |

### 7.3 更新文件

| 文件 | 修改 |
| --- | --- |
| `Cargo.toml`、`Cargo.lock` | 注册新进程、应用版本 0.5.0 及必要依赖 |
| `crates/account-facts/src/lib.rs` | 增加标准事件信封，既有事实身份和字段不变 |
| `services/trade-log/src/lib.rs` | 导出 publishing 模块 |
| `services/trade-log/adapters/src/lib.rs` | 导出发布适配器与运行模块 |
| `services/trade-log/adapters/Cargo.toml` | 将 webhook 所需 reqwest 纳入运行依赖 |
| `services/trade-log/adapters/src/postgres/facts.rs` | 候选判定和 outbox 原子提交，兼容 HTTP 先入库 |
| `services/trade-log/adapters/src/postgres/mod.rs`、`migration.rs` | 导出发布数据库模块，增量 schema 校验和角色权限 |
| `services/trade-log/adapters/tests/common/mod.rs` | 支持验收器指定独立数据库，避免创建未清理的嵌套测试库 |
| `gateway/query-api/src/config.rs`、`src/bootstrap.rs` | 发布状态客户端配置与装配 |
| `gateway/query-api/src/state.rs` | 保存发布状态客户端 |
| `gateway/query-api/src/clients/mod.rs`、`src/http/handlers/mod.rs`、`src/http/router.rs` | 模块导出及路由注册 |
| `config/query-api.toml` | 新内部发布状态服务地址及凭证文件路径 |
| `gateway/query-api/Dockerfile` | 增加发布进程 build target |
| `compose.yaml` | 新发布容器、凭证和 collector 配置只读挂载；验收接收器 profile |
| `scripts/verify-version.py`、`scripts/tests/test_verify_version.py` | v0.5 检查及错误判定测试，旧版本入口保留 |
| `docs/api.md` | 统一 webhook、状态、字段和错误契约 |
| `docs/verify.md` | 追加 v0.5 执行命令，默认从配置读取账户 |
| `docs/version0.5.md` | 同步实际交付文件、配置和验证状态 |

以上为实际新增和更新文件。原 realtime.rs、websocket_runtime.rs 直接复用；候选上下文从归档、会话和 checkpoint 在提交事务中读取，不改写消息模式。旧版本验收入口、原数据库和证据卷直接复用。

## 8. 配置与异常处理

### 8.1 配置示例

```toml
config_version = 1

[server]
host = "0.0.0.0"
port = 8083
shutdown_timeout_seconds = 25

[logging]
level = "info"
format = "json"

[internal]
credential_file = "/run/secrets/trade-log-token"

[database]
url_file = "/run/secrets/trade-log-publisher-database-url"
max_connections = 4
connect_timeout_seconds = 5
statement_timeout_seconds = 10

[publishing]
enabled = false
collection_config_path = "/etc/robotech/trade-collector.toml"
webhook_url = ""
credential_file = "/run/secrets/webhook-token"
allow_plain_http = false
request_timeout_seconds = 10
max_response_bytes = 65536
lease_seconds = 30
poll_interval_ms = 500
retry_base_seconds = 5
retry_max_seconds = 300
max_event_age_seconds = 30
signal_ttl_seconds = 60
clock_skew_tolerance_seconds = 5
```

启用时必须提供有效目标和非空凭证。默认要求 HTTPS；验收 Docker 内网接收器可显式 allow_plain_http=true。目标地址不接受用户名密码、查询参数或片段；不跟随重定向。凭证以只读单文件挂载，不提交仓库、不写事件或日志。

领取租约须大于请求超时与数据库 statement_timeout 之和；关闭宽限须覆盖这两项，默认 25 秒，Compose stop_grace_period=30s。TTL 不小于候选准入年龄。关闭 publishing 时允许目标为空，但状态接口及数据库版本校验仍应正常。发布角色只读事实与原始证据，写发布控制、队列和尝试；collector 获得候选判定及 outbox 入队权限，不能由发布角色改写事实或采集水位。

### 8.2 异常可见性

| 情况 | 行为 |
| --- | --- |
| publishing 未启用 | 保存成交，不创建新候选，状态 DISABLED |
| UNKNOWN 缺少可靠会话边界 | 保存并记录 MODE_UNCONFIRMED，不能冒充信号成功 |
| 快照、历史、过旧或强制成交 | 保存事实，候选抑制原因可查询 |
| outbox 写入失败 | 整体事务失败，WS 任务恢复后重试，不伪造已发布 |
| webhook 临时失败 | 队列保留并退避，采集继续 |
| webhook 永久拒绝或目标变化 | BLOCKED 或启动失败，明确原因 |
| 数据库故障 | 采集和发送确认按持久化状态恢复 |
| 接收成功但确认丢失 | 相同 event_id 重发，由接收方去重 |
| 发布进程不可用 | 发布状态接口 503，网关健康及库存按原职责工作 |

## 9. 验证与验收

### 9.1 固定开发验证

必须覆盖：

1. HTTP、导入、快照和历史 UNKNOWN 无候选；明确 false 和快照后缺省标记按规则准入；非法类型不准入。
2. 新会话不能继承旧快照边界；订阅确认前、激活前、过期、未来超限、陈旧元数据和强制成交均有原因。
3. HTTP 先保存、WS 先保存、同消息重放及不同消息重复，只有一个 outbox；标准事实版本和原证据不变。
4. 事实、观察、判定、outbox 与提交位置整体提交或回滚；发送失败不回滚已保存事实。
5. 无数据库锁覆盖外部请求；续租等待时采集继续；失效领取租约不能覆盖新确认。
6. 2xx、429/Retry-After、5xx、超时、连接失败、401、3xx 分类正确；凭证不进入日志。
7. 领取后崩溃、接收后确认失败、重启重新领取，事件身份、正文摘要、时效不改变；接收端只保存一份业务事件。
8. 过期重试不延长 expires_at；显式关闭与重启、重新启用的激活边界正确；旧目标队列不静默改投。
9. 增量 migration、旧事实和水位保持、发布角色权限、网关状态与基础接口兼容。
10. 固定验收接收器能制造失败、确认丢失和去重，脚本区分 FAIL、未启用、条件不足和缺少对账依据。

### 9.2 完成标准

真实来源的新鲜主动成交能够追溯到 WS 原始消息、候选判定、标准事实和 webhook 接收记录；快照与回补不产生候选。接收端暂时故障、数据库故障及进程重启后队列恢复，重复交付不重复消费。实际没有新成交时，真实成交验收记为待验证，固定样本通过不能替代来源观察。

所有报告区分开发测试、固定样本验收及真实来源验收。本机开发验证及实际进程验证见 [v0.5 验证记录](version0.5-acceptance.md)。固定样本通过不替代服务器 Docker 或真实来源投递验收。

## 10. 交付、部署与手动验证

### 10.1 部署与升级

新增发布进程 Docker target 和 Compose service，仅在容器网络监听 8083；正常部署仍只发布 query-api 的宿主机端口。沿用原 PostgreSQL 18 和原数据卷，不升级数据库主版本、不重置水位。

init-v0.5 保留已有数据库及内部服务凭证，生成缺失的发布角色和 webhook 凭证，执行 0004。升级前备份数据库；迁移和镜像更新在维护窗口顺序执行，避免旧进程继续提交不含候选判定的事实。默认关闭发布，配置完整后启用并确认激活时间。回退应用前关闭发布并处理队列；旧程序不保证兼容新增 migration，不自动执行破坏性降级。

验收接收器放在独立 `verify-webhook` profile，复用标准库 Python，独立保存接收及投递尝试记录；测试控制端口仅绑定宿主机回环地址 18080。不要求另开任务或使用外部 webhook 平台。

### 10.2 手动验证步骤

以下命令在服务器项目根目录执行。先在独立验收环境验证，不把固定样本注入生产账户。

#### 第 1 步：启动接收器并配置

```sh
sh scripts/init-v0.5.sh
docker compose --profile verify-webhook up -d webhook-receiver
curl --noproxy '*' -sS http://127.0.0.1:18080/health
```

发布配置启用 publishing，webhook_url 填 `http://webhook-receiver:8080/events`，allow_plain_http=true；接收器和发布器挂载同一个 webhook-token。初始化并修改发布配置后执行 `docker compose up --build -d`。

**通过标准：** 发布状态 RUNNING，激活时间和账户正确，接收器健康；未迁移旧事实为候选。

#### 第 2 步：核对真实信号和抑制原因

```sh
curl --noproxy '*' -sS http://127.0.0.1:8080/api/v1/publishing-status
curl --noproxy '*' -sS http://127.0.0.1:18080/receipts
docker compose logs --tail=50 trade-collector trade-parser-publisher
docker compose exec -T postgres psql -U robotech_admin -d robotech <<'SQL'
SELECT event_id,status,attempts,published_at,delivered_at,last_error
FROM trade_log.outbox_events ORDER BY created_at DESC LIMIT 10;
SELECT result,reason,count(*)
FROM trade_log.publication_decisions
WHERE decided_at >= now()-interval '10 minutes'
GROUP BY result,reason ORDER BY result,reason;
SQL
```

选择新鲜主动成交，核对 fact_id、市场、金额、时间和 WS 来源引用。初始快照只保存事实；后续无标记更新须具备成功快照边界及 POST_SNAPSHOT_UNFLAGGED 原因。通过 stored 的事实和候选并不是相同总数。

**通过标准：** 新候选最终 DELIVERED，接收记录一致；快照、HTTP、过期及强制交易无候选。无真实新成交时记录条件不足。

#### 第 3 步：制造接收方失败及恢复

```sh
curl --noproxy '*' -sS -X POST http://127.0.0.1:18080/control \
  -H 'Content-Type: application/json' -d '{"response_status":500}'
# 等待一个新候选，查看 RETRY_WAIT 和递增的 attempts。
curl --noproxy '*' -sS http://127.0.0.1:8080/api/v1/publishing-status
curl --noproxy '*' -sS -X POST http://127.0.0.1:18080/control \
  -H 'Content-Type: application/json' -d '{"response_status":200}'
```

**通过标准：** 失败期间采集水位继续推进，队列保留；恢复后同一事件送达，原正文、身份和 expires_at 不变。过期事件仍可记录，下游不会把它当作新鲜信号。

#### 第 4 步：确认丢失及重复消费

```sh
curl --noproxy '*' -sS -X POST http://127.0.0.1:18080/control \
  -H 'Content-Type: application/json' -d '{"drop_response_after_store_once":true}'
```

**通过标准：** 接收器持久化一条新事件后断开响应，发送方再次投递相同事件；delivery_attempts 有多次尝试，接收器有多次请求但唯一业务事件只有一条。

#### 第 5 步：发布进程重启

先制造 500 留下待发送事件，记下 event_id：

```sh
curl --noproxy '*' -sS -X POST http://127.0.0.1:18080/control \
  -H 'Content-Type: application/json' -d '{"response_status":500}'
# 等待一条候选进入 RETRY_WAIT 后再重启。
docker compose stop trade-parser-publisher
docker compose start trade-parser-publisher
curl --noproxy '*' -sS -X POST http://127.0.0.1:18080/control \
  -H 'Content-Type: application/json' -d '{"response_status":200}'
# 等待到期重试后检查原事件。
curl --noproxy '*' -sS http://127.0.0.1:8080/api/v1/publishing-status
```

**通过标准：** 原队列和尝试记录保留，激活时间不变，领取租约过期后继续发送；重启不批量发布旧成交。

#### 第 6 步：数据库故障及历史补偿

记下 HTTP 水位及至少一个待发送 event_id，在独立验收环境执行：

```sh
docker compose stop postgres
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/publishing-status
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/watch-accounts
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/health
docker compose start postgres
# 等待数据库、采集和发布恢复后检查。
curl --noproxy '*' -sS http://127.0.0.1:8080/api/v1/publishing-status
curl --noproxy '*' -sS http://127.0.0.1:8080/api/v1/watch-accounts
curl --noproxy '*' -sS http://127.0.0.1:18080/receipts
```

**通过标准：** stored/状态不可用时明确返回 503，网关 health 存活；恢复后事实、水位及 outbox 保留。HTTP 补回的数据不生成实时候选；未确认投递允许重发，候选时效不被延长。

#### 第 7 步：永久拒绝

验收接收器设置 401，等待一条候选进入 BLOCKED，随后恢复 200：

```sh
curl --noproxy '*' -sS -X POST http://127.0.0.1:18080/control \
  -H 'Content-Type: application/json' -d '{"response_status":401}'
# 等待一条候选进入 BLOCKED，记录实际 event_id。
curl --noproxy '*' -sS -X POST http://127.0.0.1:18080/control \
  -H 'Content-Type: application/json' -d '{"response_status":200}'
docker compose exec -T trade-parser-publisher \
  trade-parser-publisher retry --event-id '<实际 event_id>'
```

**通过标准：** 401 不无限重试，修正后显式重试成功；总尝试次数保留，事件正文不改写。

### 10.3 固定服务器验收脚本

验收入口默认从配置读取账户：

```sh
sh scripts/verify-v0.5.sh --wait-seconds 180
```

默认检查配置、发布状态、候选决策、outbox 约束和已有接收记录，不自动变更真实接收器。需要独立接收器的确定性测试时执行：

```sh
sh scripts/verify-v0.5.sh --fixture-tests --wait-seconds 180
```

fixture-tests 必须使用独立测试数据库、独立接收器和固定样本，覆盖候选/抑制、重复、500、401、429、确认丢失及重启；结束后清理仅由本次测试创建的资源，不触碰部署库、水位和原卷。它不是向生产 collector 塞入模拟成交。

在允许中断的部署验收环境，再启用：

```sh
sh scripts/verify-v0.5.sh --lifecycle --database-fault --wait-seconds 180
```

统一 PASS、FAIL、SKIP 和退出码；SKIP 区分未启用、条件不足、未启用故障参数和缺少接收对账依据。执行命令见 [verify.md](verify.md)。fixture-tests 首次构建 verify-v05 测试镜像，需下载 Python 镜像并编译测试程序；使用当前 PostgreSQL 容器中新建的随机测试库和测试容器内独立接收器，结束后仅删除本次创建的库、容器和临时凭证。正常发布镜像不编译测试程序。

### 10.4 交付记录

已建立 [v0.5 验证记录](version0.5-acceptance.md)，区分本机开发验证与服务器待验证项。服务器验收时补充代码版本、激活时间、来源观察、事件及尝试 ID 和失败恢复过程。实际 bug 按 v0.4 的七项格式记录：描述、影响、确认记录、修改逻辑、修改记录摘要、如何验证、验证记录。
