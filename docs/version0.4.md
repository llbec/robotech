# v0.4 开发文档：单地址实时监控与断线恢复

- 文档版本：0.4
- 状态：已实现并完成本地验收；Docker 服务器验收待执行，详见 [验收记录](version0.4-acceptance.md)
- 更新日期：2026-10-06
- 设计依据：[概要设计](overview-design.md)、[详细设计](detailed-design.md)、[产品版本路线图](roadmap.md)
- 前置版本：[v0.3 开发文档](version0.3.md)
- 接口契约统一维护在 [接口说明](api.md)，实施本版时同步状态字段，不新增版本接口文档

## 1. 版本目标

在 v0.3 单地址定时采集基础上，通过官方 WebSocket 接收新成交并及时保存，HTTP 继续负责首次扫描、定期补偿以及断线后的范围扫描。程序重启或连接断开后能够恢复订阅，从持久化 HTTP 水位继续补采；无法确认或无法获取的范围必须可见。

保持 query-api、trade-log-query、trade-collector、PostgreSQL 四个常驻容器。WebSocket 和 HTTP 补偿在原 trade-collector 中协作，不新增业务服务、不重建事实存储。程序与应用镜像版本升级为 0.4.0，schema_version=1、稳定 fact_id 和既有接口路径继续复用。

实时监控表示来源主动推送后尽快入库，不承诺交易发生到本库提交的固定延迟，也不表示历史完整。用户仍通过 HTTP stored 接口查看保存结果；对外 WebSocket、候选信号、webhook 和跟单执行不属于本版。

## 2. 功能范围

| 能力 | 本版交付 | 边界 |
| --- | --- | --- |
| 单地址实时接收 | 订阅官方 userFills，保存合约成交 | 不扩展多地址管理、现货展示或其他账户事实 |
| 首条快照识别 | 区分快照、后续更新和无法确定的消息 | 快照不代表从配置起点开始的全部历史 |
| HTTP 持续补偿 | 保留 v0.3 的轮询、区间拆分和水位提交 | WebSocket 正常时也不关闭 HTTP 补偿 |
| 跨通道去重 | 同一成交由 HTTP、WebSocket 和重复推送获取时只保存一个 current 事实 | 真实金额、方向、手续费等冲突仍明确失败，不覆盖旧事实 |
| 断线恢复 | 心跳、重连、重新订阅、HTTP 追赶及缺口记录 | 超出来源可查询历史的成交不能保证补齐 |
| 状态查询 | 展示两个通道的状态、接收时间、入库时间和恢复范围 | 进程健康、连接成功、缺口扫描完成分别表达 |
| 可追溯与重放 | 保留完整 WebSocket 原始消息及逐笔来源关联，支持 reparse | 不将网络接收位置当成来源可恢复游标 |
| 升级兼容 | 增量 migration，复用原凭证、数据库及命名卷 | 不修改已发布 migration，不重置 v0.3 水位 |

本版不引入消息队列、outbox 发布、trade-parser-publisher 新进程、完整审核服务、事实修订、账户分析或签名下单。处理模式和时间字段为后续 v0.5 区分实时与历史准备基础，但不提前发布信号。

## 3. 使用方式

### 3.1 启用实时通道

沿用 config/trade-collector.toml 的 account、start_time 和所有 v0.3 collection 参数，新增可选 websocket 配置。没有该段或 enabled=false 时保持纯 HTTP 轮询行为；随附 v0.4 配置启用 WebSocket。

```toml
[websocket]
enabled = true
connect_timeout_seconds = 10
subscribe_timeout_seconds = 10
ping_interval_seconds = 20
pong_timeout_seconds = 10
reconnect_base_seconds = 5
reconnect_max_seconds = 60
reconnect_jitter_percent = 20
reconnect_reset_after_seconds = 60
max_message_bytes = 16777216
max_pending_messages = 256
max_pending_bytes = 33554432
metadata_refresh_seconds = 300
commit_timeout_seconds = 10
```

账户使用 v0.3 已配置的地址；不单独为 WebSocket 配置第二个账户。来源 URL 根据已有 network 选择，不把生产连接域名或端口重复写入配置。

start_time 继续只用于首次建立 HTTP 检查点。已有水位恢复后，修改起点不触发从头同步；启用 WebSocket 也不清空已有事实或未完成的 HTTP 范围。

### 3.2 查看实时与补偿状态

```text
GET /api/v1/watch-accounts
```

在已有 items 状态对象中增加 websocket、recovery 和 monitoring_status。原 scanned_through、last_success_at、consecutive_failures、next_run_at、pending_range 等字段仍表示 HTTP 自动扫描，不改成“最近收到 WebSocket 消息”。

连接成功但仍有恢复范围时应显示 RECOVERING；WebSocket 不可用而 HTTP 仍正常时显示 DEGRADED。用户可以分别判断“实时通道是否正常”“HTTP 是否继续补偿”“缺口是否仍未完成扫描”。

### 3.3 查询保存的成交

```text
GET /api/v1/trade-events?account=<ADDRESS>&source=stored&limit=100
```

复用 v0.2 的时间过滤、排序和快照分页。WebSocket 成交提交后，新查询可以看到；已有游标继续固定旧 snapshot_seq。stored 不连接 WebSocket，不触发 HTTP 补采。

本版不向业务客户端持续推送数据。用户若需要看到更新，可再次请求 stored；普通客户端不需要知道官方 WebSocket 地址、内部 token 或 collector 进程地址。

## 4. 数据来源与查询边界

### 4.1 官方 WebSocket 契约

mainnet 使用 `wss://api.hyperliquid.xyz/ws`，testnet 使用 `wss://api.hyperliquid-testnet.xyz/ws`。订阅消息为：

```json
{
  "method": "subscribe",
  "subscription": {
    "type": "userFills",
    "user": "0x010461c14e146ac35fe42271bdc1134ee31c703a",
    "aggregateByTime": false
  }
}
```

订阅确认 channel=subscriptionResponse，成交 channel=userFills；成交 data 包含 user、fills 和可选 isSnapshot。必须核对返回账户与配置账户一致，不能把其他账户或其他频道的数据写入本账户。[官方 WebSocket 说明](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/websocket)、[官方订阅说明](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/websocket/subscriptions)

isSnapshot=true 归类为 SNAPSHOT；明确 false 归类为 LIVE_UPDATE。字段缺失时归类 UNKNOWN，不自动将未知消息当成新实时成交。三种类别中的合法成交均可幂等保存，但来源模式分别记录；快照不计入“新增实时消息”验证。

客户端按间隔发送 `{"method":"ping"}`，等待 `{"channel":"pong"}`。官方说明空闲连接需要这种应用层心跳；本版不能只依赖 TCP 存活或 WebSocket 控制帧。[官方心跳说明](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/websocket/timeouts-and-heartbeats)

### 4.2 HTTP 补偿继续使用既有能力

userFillsByTime、meta、spotMeta 及范围 `[start_ms,end_ms)` 转换沿用 v0.3。达到单次来源记录上限时拆分范围，单毫秒仍满页则停止推进水位并明确失败。

官方时间成交接口当前最多每次 2,000 条，可查询最近 10,000 条；该限制包含现货等来源记录。[官方 HTTP Info endpoint 说明](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint#retrieve-a-users-fills-by-time)

WebSocket 快照和 HTTP 响应都不能证明无限历史完整。长期停机、超过来源可访问范围、窗口外迟到、数据冲突等情况下，必须保留覆盖警告。

### 4.3 两类位置不能混用

| 位置/时间 | 含义 | 可以做什么 |
| --- | --- | --- |
| scanned_through | HTTP 完整范围与事实一起提交后的结束边界 | 计算下一次 HTTP 扫描起点和断线补偿范围 |
| WebSocket last_trade_at | WebSocket 已提交成交中的最大时间 | 展示实时观察位置，不能推进 HTTP 水位 |
| WebSocket last_received_at | 最近已归档账户数据消息的接收时间 | 观察接收活跃度，不证明事实已提交；内存队列中的消息尚不更新此持久化值 |
| WebSocket last_committed_at | 最近完整保存 WebSocket 数据消息的时间 | 确认数据消息已提交，无成交快照也可有提交时间 |
| session_id/message_sequence | 本地连接身份与消息顺序 | 原始消息归档及重放；不是官方恢复游标 |

WebSocket 最新成交时间不能替代 HTTP scanned_through。收到时间很新的成交仍可能缺少更早成交，按最新时间重新开始补偿会跳过缺口。

## 5. 接口与持久化数据设计

### 5.1 接口兼容与状态字段

不新增外部路径。扩展现有 GET /api/v1/watch-accounts 和内部 GET /internal/v1/collection-status 的状态对象；内部凭证、错误信封和 trace 规则保持不变。

顶层 status 继续表示 v0.3 HTTP 调度状态，以兼容已有调用方；新增 monitoring_status 表示双通道监控状态。JSON schema_version 仍为 1，新字段为兼容扩展，读取旧数据时提供明确默认值。

| 新增字段 | 格式/含义 |
| --- | --- |
| monitoring_status | HTTP_ONLY、STARTING、LIVE、RECOVERING、DEGRADED、FAILED、STOPPED |
| websocket.enabled | boolean，是否启用实时通道 |
| websocket.status | DISABLED、CONNECTING、SUBSCRIBING、LIVE、RECONNECT_WAIT、FAILED、STOPPED |
| websocket.session_id | UUID string/null，当前或最近连接身份 |
| websocket.connected_at/subscribed_at | UTC 毫秒时间/null，连接与订阅确认分开 |
| websocket.last_received_at | UTC 毫秒时间/null，最近已归档账户数据消息的接收时间，不等于事实提交时间 |
| websocket.last_committed_at | UTC 毫秒时间/null，最近成功数据消息提交时间 |
| websocket.last_trade_at | UTC 毫秒时间/null，最近已提交成交位置，不替代 HTTP 水位 |
| websocket.last_pong_at | UTC 毫秒时间/null，用于安静账户的连接健康判断 |
| websocket.reconnect_count | 非负整数，首次成功连接之外重新建立连接的次数，跨重启保留；失败握手不计入 |
| websocket.connection_count | 非负整数，成功建立连接的总次数，包含首次 |
| websocket.next_retry_at | UTC 毫秒时间/null，下一次重连时间 |
| websocket.last_error | null 或 code/message/occurred_at 脱敏摘要 |
| websocket.pending_messages/pending_bytes | 非负整数，最近内存队列观测值；重启归零，不是持久化数据量；字节包含包络和提取 fills |
| websocket.metadata_stale | boolean，沿用旧映射时为 true，并增加 METADATA_STALE 警告 |
| recovery.status | NOT_STARTED、SCANNING、HTTP_SCANNED、BLOCKED；关闭 WebSocket 且无待补偿范围时 DISABLED，已有待补偿范围仍保留 SCANNING 或 BLOCKED，完成后转为 DISABLED |
| recovery.target_through | UTC 毫秒时间/null，当前需要追赶的固定结束边界 |
| recovery.open_gap_count | 非负整数，尚未完成 HTTP 扫描或已阻塞的恢复范围数量 |
| recovery.last_scanned_at | UTC 毫秒时间/null，最近恢复范围完成 HTTP 扫描的时间 |
| recovery.last_error | null 或脱敏错误摘要 |

HTTP 与 WebSocket 的失败次数和错误独立保存。WebSocket 连接恢复不清空 HTTP 错误；HTTP 成功也不伪造 WebSocket 正常。数据库不可读或 collector 不可达时状态接口仍返回 503；可读但某通道异常时返回 200，通过各状态说明。

monitoring_status=LIVE 要求订阅与心跳有效，恢复目标已被 HTTP scanned_through 覆盖，且两个通道没有阻塞错误。这里的 LIVE 只表示运行状态，coverage 仍为 SOURCE_HISTORY_NOT_VERIFIED，不证明历史完整。只有 HTTP 正常而 WebSocket 异常为 DEGRADED；正在追赶恢复目标为 RECOVERING；启用后尚未建立可用通道为 STARTING；不可恢复的业务冲突为 FAILED。

### 5.2 增量迁移与既有表

新增 `0003_websocket_monitoring.sql`，不修改 0001/0002 的 SQL 或摘要。

| 表 | 更新内容 |
| --- | --- |
| collection_checkpoints | 增加 websocket_state JSONB 和 recovery_state JSONB；复用原账户租约、epoch、HTTP position 和 pending_work |
| collection_jobs | 增加 transport=HTTP/WEBSOCKET，默认 HTTP；job_origin 继续 MANUAL/COLLECTOR；WS 任务 mode=REALTIME，记录 session_id/message_sequence/message_mode |
| raw_logs | 增加 transport、session_id、message_sequence、message_mode；HTTP 默认保留；WS kind=userFills、source_id=official_ws |
| account_fact_versions | 增加 semantic_hash_version、semantic_content_hash 可空列，用于跨通道标准业务比较 |
| account_facts_current/fact_observations/ingestion_state | 保持原事实身份、版本指针、观察关联和提交序号约定 |

WS 原始消息没有 HTTP 状态码，raw_logs.http_status 改为允许 null，并增加 transport 对应约束：HTTP 必须有状态，WEBSOCKET 必须为 null。不能用伪造的 HTTP 200 表示 WebSocket 数据。

raw_logs 新增 `(session_id,message_sequence)` 的非空唯一索引，保证本地同一消息归档幂等。该索引不能替代成交 fact_id 去重，不同连接重复发送同一成交仍产生新的来源观察。

每个账户数据消息建立一个独立采集任务和 query_id，空 fills 合法，元数据原始字节和 SHA-256 保存在任务 request.meta_snapshots，active_meta_snapshot 选择本次解析快照；需要刷新未知市场时追加快照，不覆写原证据。订阅确认、ping/pong 不生成成交任务，连接状态另行保存。格式错误的 userFills 数据消息保留原始字节和错误任务，不在解析失败后丢掉证据。

### 5.3 连接与缺口记录

新增 trade_log.collection_stream_sessions：session_id UUID 主键，checkpoint_key、lease_epoch、账户、network、连接状态、connected_at/subscribed_at/closed_at、close_reason、最后接收/提交和 pong 时间、接收/提交消息计数及 created_at。仅记录已实际发生的会话，不预先生成“成功连接”。

新增 trade_log.collection_gaps：gap_id UUID 主键，checkpoint_key、session_ids 引用数组、reason、start_ms、end_ms 可空、status、detected_at、scanned_at 及 last_error。reason 包含 STARTUP、DISCONNECT、HEARTBEAT_TIMEOUT、QUEUE_OVERFLOW、DATABASE_UNAVAILABLE、PROCESS_INTERRUPTED；相邻或重叠的未扫描范围合并处理，保留触发会话关联，避免每次重连建立无界重复任务。

| gap.status | 含义 |
| --- | --- |
| OPEN | 已发现中断，结束边界尚未固定或尚未启动补偿 |
| SCANNING | 已固定范围，HTTP 正在追赶 |
| HTTP_SCANNED | HTTP 完整扫描事务已覆盖结束边界；不是来源历史完整性证明 |
| BLOCKED | 来源饱和、冲突等使范围不能继续；保留原因和原水位 |

连接恢复不能直接把 gap 标为 HTTP_SCANNED。完成标记与 HTTP 水位提交同一事务更新；预算暂停、解析失败、仅保存原始响应均不满足完成条件。

### 5.4 跨通道内容比较

现有 content_hash 包含 extension.source_fill。HTTP 与 WebSocket 在可选原始字段、封装或字符串表示上可能不同，不能继续把所有原始 JSON 差异当作事实冲突。

新增版本化标准业务指纹 hl-trade-content-v2，使用标准事实身份、账户、类型、revision/change_type/confirmation_status、occurred_at、ordering_key/sub_index，以及 TradeFact 的市场、资产、动作、触发类型、可跟单标志、方向、持仓效果、订单/操作标识、价格、数量、成交额、手续费资产和金额、来源盈亏口径、交易哈希。金额先按既有定点类型规范化，null 不等同于 0；仅在业务含义一致时去重。

raw_log_id、传输通道、本地时间、source_indices、source_fill 及连接身份不进入该业务指纹；源语义差异必须映射到标准业务字段，不通过丢弃真实业务字段掩盖。保留 v0.3 content_hash 与原始 payload，不批量覆写历史事实。旧行的新指纹从已保存标准 payload 推导，首次比较时在事务内填入新增指纹列，不产生新 revision。

相同 fact_id 的标准金额、方向或其他业务字段不同仍返回 VERSION_CONFLICT，保留双方证据，不更新旧 current。source 保持 OFFICIAL_API；具体 HTTP/WS 通道和快照模式属于采集证据。reparse 分别验证原始字节摘要与版本化业务指纹，不能用新指纹跳过原始响应完整性检查。

### 5.5 WebSocket 提交原子性

接收完整消息后先保存原始字节，再解析 data.fills。适配器可提取 fills 数组供既有解析器使用，但不能只归档提取后的数组；source_index 对应原始 data.fills 的位置，关联路径保存在证据定位信息中。

一个数据消息的事实、观察关联、任务 COMPLETED 和 WS last_committed_at/last_trade_at 在同一事务提交，并校验同一账户租约及 epoch。重复消息增加合法观察，但不新增 current 事实；ingest_seq 仅新事实递增。

HTTP scanned_through 与 pending_work 不由 WS 提交修改。HTTP 提交也不能覆盖 websocket_state；两个事务更新各自字段，避免使用旧 JSON 整体覆盖对方的新状态。租约失效时不能提交事实或接收成功位置。

## 6. 处理流程

### 6.1 启动顺序

```text
配置/schema 校验 → 获取既有账户租约 → 恢复未完成原始消息和 HTTP pending_work
→ 初始化元数据 → 建立 WebSocket 并订阅 → 记录订阅边界与恢复目标
→ 实时消息持续归档及保存，同时 HTTP 从旧水位补偿
→ HTTP 水位达到恢复目标 → 更新运行状态，继续双通道运行
```

首次扫描和历史积压不能阻塞 WebSocket 接入。取得租约、准备元数据后即连接实时通道，历史 HTTP 范围单独有界推进。只要数据库允许提交，已订阅的新消息可以先入库；这不允许跳过旧 HTTP 范围。

账户只有一个采集租约，覆盖 WS 与 HTTP 两条执行路径。候补实例不得另行建立同账户订阅或写水位。query 进程的启动清理继续只处理 MANUAL 任务。

### 6.2 正常接收

按来源 channel 分类：匹配的订阅确认更新 SUBSCRIBING 状态；pong 更新心跳；userFills 核对账户、大小和结构后归档，解析并事务提交。未知业务频道不写入本账户；控制帧按协议处理，不当作成交消息。

元数据在启动加载并按配置刷新，解析任务固定引用一份已归档成功的 meta/spotMeta。刷新失败时，现有映射可继续处理已识别市场并标记过期；出现未知市场则保存原始消息、触发刷新和待解析处理，不按字符串猜测市场或静默丢成交。

队列同时限制消息数量和总字节。接收与 HTTP 补偿各使用有界执行资源；不得因 HTTP 批量积压让 socket 无法读 pong。接收、提交和来源成交时间分别记录，用于测量延迟。

### 6.3 断线、重连与补偿

连接关闭、订阅失败、pong 超时、读写错误或队列溢出时，停止该 WebSocket 连接并记录恢复需求。重连是重新建立 WebSocket 并订阅，不是 HTTP 补偿轮询。连续失败后的基础等待依次为 5、10、20、40、60 秒，此后保持 60 秒；每次增加基础等待的 0～20% 随机抖动，最终等待不超过 60 秒。首次启动可立即连接，失败或中断后的尝试才应用退避。

成功订阅后，连接连续稳定运行至少 60 秒且期间心跳正常，才重置本次连续重连退避；仅收到订阅确认或一次 pong 不重置，累计重连次数不清零。短暂恢复后再次断线继续原退避，避免反复建立连接。数据库或处理能力尚未恢复时暂停连接尝试，满足恢复条件后再按退避调度重连。

断线范围从可靠 HTTP scanned_through 回退 overlap 开始，没有水位时使用保存的 initial_start。重新订阅确认时固定恢复结束目标 T；HTTP 按 v0.3 的 safety_delay、最大窗口和请求预算逐步扫描，直到 committed scanned_through >= T。达到 T 前可接收新 WS 数据，但恢复状态仍未完成。

多次重连会扩展或合并未完成恢复范围，不丢弃原 pending_work。快照中的成交幂等保存，但收到快照本身不能完成恢复范围。

### 6.4 数据库故障与背压

数据库不可用或不能确认租约时停止提交，不积累无限内存消息；关闭 WS 连接并等待数据库恢复。队列满或单帧超限时同样关闭并记录 gap，不静默丢弃后还报告 LIVE。数据库不可用期间不反复建立 WS 连接；队列溢出后，须等提交任务恢复处理能力、队列低于容量限制且租约有效，再允许重连。持续单帧超限须先解决消息限制或处理能力问题，不能靠反复重连消除错误。

数据库恢复后读取旧检查点、已归档但未提交的消息和会话状态。原始证据仍可解析时优先恢复这些消息，再执行 HTTP 补偿；未能落盘的数据只能从来源可查询范围补取。

故障期间不能写 gap 时，在下次启动或恢复中根据未正常关闭的会话、过期租约和旧 HTTP 水位补建恢复范围。不得仅依赖内存里的 disconnected_at，否则进程被杀后缺口会消失。

### 6.5 退出与恢复

SIGTERM 停止新连接和新轮次，关闭 WebSocket，在既有关闭时限内提交已完整处理的任务或取消未完成任务。未提交的原始消息保留为待恢复，HTTP 水位不提前推进。

所有后台任务必须被等待或取消，不留下 detached worker。正常退出记录会话关闭原因并释放自有租约；强制退出依靠租约过期和未关闭会话恢复。数据库不可用时记录关闭失败日志，下一实例按可靠水位恢复，不伪造数据库 STOPPED。

## 7. 程序结构设计

### 7.1 进程与模块依赖

```text
query-api ──内部 HTTP──► trade-log-query ──► PostgreSQL
         └─内部 HTTP──► trade-collector ──► PostgreSQL
                              ├──────────► 官方 WebSocket userFills
                              └──────────► 官方 HTTP 范围补偿及元数据
```

实时业务放在 trade-log，连接适配在 protocols/implementations/hyperliquid，运行调度在 adapters，依赖装配在原 trade-collector binary。网关只读取状态，业务模块不依赖具体 WebSocket SDK。

### 7.2 复用与扩展原则

复用已有 acquisition/persistence 的 HTTP 端口，扩展 checkpoint 状态，新增 realtime 来源流端口、实时消息分类与提交意图，不复制 HTTP 收集流程和 Hyperliquid 金额解析。WebSocket 适配器负责 envelope 识别、账户验证和 fills 提取，既有 parser 负责标准事实。

租约管理从 v0.3 单轮调度扩展为两个 worker 共用的监督任务。HTTP 原子提交继续保留，只增加 gap 更新；WS 使用同一事实存储基础，仅更新自己的接收提交位置。

### 7.3 本版新增与更新文件

以 v0.3 为基线。以下为实际实现文件清单，明确新增、更新和复用。

#### 新增文件

| 文件 | 职责 |
| --- | --- |
| services/trade-log/src/realtime/mod.rs | 实时来源流端口、消息数据、配置与纯重连调度 |
| services/trade-log/src/collection/content.rs | 标准业务指纹规则，保持事实身份与原字节摘要独立 |
| protocols/implementations/hyperliquid/src/websocket.rs | WSS 连接、订阅、应用层心跳和包络适配 |
| services/trade-log/adapters/src/websocket_runtime.rs | 有界消息队列、重连与数据消息处理 |
| services/trade-log/adapters/src/postgres/realtime.rs | WS 原始消息、任务、会话、恢复记录及事务位置更新 |
| services/trade-log/migrations/0003_websocket_monitoring.sql | 状态列、传输证据、业务指纹、会话与 gap 表和约束 |
| services/trade-log/tests/realtime.rs | 重连退避、抖动、滚动频率窗口、稳定重置和旧配置默认值 |
| services/trade-log/adapters/tests/websocket_chain.rs | 固定 WS 来源、HTTP 重复、断线恢复、队列与事务测试 |
| protocols/implementations/hyperliquid/tests/websocket.rs | 订阅确认、账户匹配、ping/pong、消息限制和关闭协议 |
| tests/fixtures/v0.4/snapshot.json | 构造快照包络与固定成交 |
| tests/fixtures/v0.4/update.json | 与 HTTP 可比的实时更新包络 |
| tests/fixtures/v0.4/unknown.json | 缺少 isSnapshot 的保守分类样本 |
| tests/fixtures/v0.4/malformed.json | 可归档但不能解析的消息样本 |
| scripts/init-v0.4.sh | 保留已有凭证与角色，执行 0003 migration 和证据卷初始化 |

#### 更新文件

| 文件 | 变更点 |
| --- | --- |
| Cargo.toml | workspace 版本 0.4.0，统一必要 WebSocket 异步依赖 |
| Cargo.lock | 锁定新增依赖及所有 workspace 包版本 |
| compose.yaml | 应用镜像标签 0.4.0；四个常驻服务和原卷不变 |
| config/trade-collector.toml | 新增可选 websocket 段，随附配置启用，账户及起点保持 |
| services/trade-log/src/lib.rs | 导出 realtime 模块 |
| services/trade-log/src/checkpoint/mod.rs | 状态 DTO、会话/恢复位置及租约共用契约 |
| services/trade-log/src/collection/mod.rs | 导出版本化业务内容比较，复用原范围逻辑 |
| protocols/implementations/hyperliquid/Cargo.toml | 引入实际 WS 库及 TLS 能力，不引入完整交易 SDK |
| protocols/implementations/hyperliquid/src/lib.rs | 导出 WS 适配器 |
| protocols/implementations/hyperliquid/src/source.rs | HTTP 补偿与元数据共用滚动权重预算及来源限流等待 |
| protocols/implementations/hyperliquid/src/parser.rs | 同批重复记录按标准业务指纹比较，保留逐笔原始索引 |
| services/trade-log/Cargo.toml | 标准业务指纹所需定点金额与摘要依赖 |
| services/trade-log/adapters/src/lib.rs | 导出 WS runtime |
| services/trade-log/adapters/src/collector_runtime.rs | 双 worker 监督、统一租约和优雅关闭，保留 HTTP 轮询 |
| services/trade-log/adapters/src/postgres/mod.rs | 导出 realtime 适配器 |
| services/trade-log/adapters/src/postgres/checkpoint.rs | 读取双通道状态，恢复 gap 及独立字段更新 |
| services/trade-log/adapters/src/postgres/facts.rs | 版本化标准业务比较、共用事实事务和观察关联 |
| services/trade-log/adapters/src/postgres/replay.rs | 根据任务 transport 选择 HTTP/WS 重放；旧任务兼容 |
| services/trade-log/adapters/src/file_evidence.rs | 可选完整 WS 消息镜像，不伪造 HTTP 状态 |
| services/trade-log/bins/trade-collector/src/config.rs | WS 参数默认值、大小及关联时限校验 |
| services/trade-log/bins/trade-collector/src/bootstrap.rs | 装配来源流、HTTP 补偿及共享监督任务 |
| services/trade-log/bins/trade-collector/tests/startup.rs | 关闭 WS 的旧配置兼容与双通道退出验证 |
| services/trade-log/adapters/tests/common/mod.rs | 复用隔离数据库，增加固定 WS 与 HTTP 组合来源 |
| services/trade-log/adapters/tests/checkpoint.rs | 保留旧 schema 升级测试，按旧列保存历史事实后执行新版迁移 |
| services/trade-log/adapters/tests/collector_chain.rs | 验证网关透传双通道新增字段，保留认证与数据库故障映射 |

#### 直接复用的文件和约定

复用 account-facts 的事实 ID 与公开 payload、shared-types 定点金额、service-runtime 日志和生命周期、既有 stored 查询、网关 handler/client 以及 collector_http 路由。来源流端口集中在 realtime/mod.rs；原 acquisition 的 HTTP SourceReader、persistence 的 HTTP 提交端口、migration 的全表授权及版本校验直接复用。既有 replay_import 和网关测试继续执行；新增 WS 重放、状态和故障测试集中在 websocket_chain.rs。状态 DTO 增加字段即可传递新状态，不另建公网实时接口。

不要求修改 Dockerfile 的 build target 或引入新 binary；已有三个程序的 target 和构建缓存继续使用。0001/0002 migration、v0.3 的既有初始化脚本和历史版本文档保留，不把旧脚本偷偷改成另一版本入口。

本版同步 docs/api.md，并新增 docs/version0.4-acceptance.md 记录实际验收结果。不修改 README、roadmap、详细设计或历史版本文档。

### 7.4 固定验收脚本补充文件（2026-10-06）

| 文件 | 本次类型 | 职责 |
| --- | --- | --- |
| `scripts/verify-v0.4.sh` | 新增 | 本版本固定服务器验收入口 |
| `scripts/verify-version.py` | 复用 | 四个版本共用的检查实现，按版本启用能力 |
| `scripts/tests/test_verify_version.py` | 复用 | 验收脚本失败判定、跳过与恢复路径测试 |

这是本次补充交付清单；不改变原版本首次交付文件的分类。

## 8. 配置与异常处理

### 8.1 配置校验和调度周期

v0.3 的 interval_seconds=30 继续表示 HTTP 完整成功后等待 30 秒；WS 新消息不等待该轮询周期。HTTP 的 overlap=60 秒、safety_delay=2 秒、最大推进窗口=1 小时和预算=20 次继续有效，恢复积压按原规则推进。

WebSocket 重连默认基础等待 5 秒、指数增长至 60 秒，随机抖动上限 20%，稳定订阅并通过心跳 60 秒后重置退避。校验 reconnect_base_seconds > 0、reconnect_max_seconds >= reconnect_base_seconds、reconnect_jitter_percent 在 0～20 内、reconnect_reset_after_seconds >= 60。实际等待包含抖动且不超过 reconnect_max_seconds；这些默认时间是本项目选择，不是官方指定的重连周期。

Hyperliquid 官方限制按 IP 共享：最多同时保持 10 个 WebSocket 连接，每分钟最多建立 30 个新连接，每分钟发送至来源的 WS 消息总数最多 2000；HTTP REST 请求共享每分钟 1200 权重额度，userFillsByTime 等接口还按返回条数增加权重。[官方请求与连接限制](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/rate-limits-and-user-limits)

重连调度必须限制本采集进程滚动 60 秒内最多 10 次连接尝试，失败的握手也计入；超过后等待最早一次尝试退出窗口，再同时满足退避与恢复条件才发起下一次。此调度状态由同一进程监督任务持有，租约重获或 WS worker 重建不重置频率窗口及尚未结束的退避；进程重启后重新计时。此内部上限预留同一公网出口其他程序的余量，不能保证其他程序共同使用时仍不触及官方限额。部署时需统筹同出口连接数量及请求额度。HTTP 来源实例使用滚动 60 秒最多 1000 的内部权重预算；每个 fills 请求按最多 2000 条保守预留 120 权重，meta/spotMeta 各预留 20，并在收到 429 时共享来源冷却等待，避免元数据刷新绕过限流。该预算不协调其他进程或同出口其他应用。HTTP 的 30 秒等待不代表一次轮次只请求一次；恢复扫描、拆分及重试均须计入权重预算，限流时遵循来源可用的 Retry-After 或既有 HTTP 退避，不立即连续重试。

WS 心跳默认每 20 秒发送一次，发送后 10 秒内未收到对应 pong 则判为不可用。不得因账户没有成交而单独判为断线。ping_interval 与 pong_timeout 之和必须小于官方空闲关闭阈值；重连、连接及订阅时限均正数且有上界。

max_message_bytes 不超过既有 64 MiB 输入上界；队列条数和总字节都必须为正，max_pending_bytes 至少容纳一条最大消息。帧、队列和数据库提交均有界，不能用无界 channel。

元数据刷新默认 300 秒；commit_timeout 默认 10 秒，小于租约时限。租约续期由监督任务独立执行，不受 socket 接收或 HTTP 轮次等待阻塞。已有服务 token、collector 数据库 URL 和迁移 URL 分离规则继续沿用。

### 8.2 状态与异常

| 情况 | 处理与对外表现 |
| --- | --- |
| WS 未启用 | HTTP_ONLY，保持原 v0.3 行为 |
| 连接/订阅中 | CONNECTING/SUBSCRIBING，不能提前标为实时可用 |
| 已订阅但 HTTP 未追赶目标 | WS 可以保存新成交，总体 RECOVERING |
| 断线或心跳超时 | 记录 gap，RECONNECT_WAIT；HTTP 可用时总体 DEGRADED |
| 订阅错误或账户不匹配 | 归档必要证据，关闭连接并记录错误，不写错误账户事实 |
| 队列满或消息超限 | QUEUE_OVERFLOW/MESSAGE_TOO_LARGE，关闭连接，启动有界恢复 |
| 元数据刷新失败 | 已知市场可继续，标记 METADATA_STALE；未知市场消息待处理 |
| 同成交真实业务冲突 | VERSION_CONFLICT，失败任务和双方原始证据保留，不自动修订 |
| HTTP 毫秒饱和 | gap BLOCKED，HTTP 水位不推进；不能据 WS 新消息宣布缺口完成 |
| 数据库故障/租约失效 | 停止写入和实时接收，恢复后读取可靠检查点重建恢复范围 |
| 正常停止 | 会话 STOPPED，等待或取消全部 worker，保留已提交记录 |

日志记录 session_id、message_sequence、query_id、message_mode、lease_epoch、transport、接收/提交时间、入库计数和错误码。不打印凭证，不把完整来源 payload 直接塞入普通请求日志；来源原文保存在证据归档中。

## 9. 验证与验收

### 9.1 必须通过的开发验证

1. 保留 v0.0–v0.3 测试，包括 HTTP 水位、幂等、迁移、分页、旧证据导入和故障恢复；没有 websocket 段时兼容旧配置。
2. 固定 WS 服务验证订阅参数、确认、快照、更新、未知模式、空 fills、错账户、未知频道、畸形消息和超限帧。
3. 同一成交以 HTTP、WS 快照、重复实时推送、重连快照获取时 current 仅一条，观察关联分别保留；可选原始字段或等价金额文本不误报冲突，真实业务差异仍失败。
4. 在旧标准事实上推导新指纹，不改变旧 fact_id、revision、原 payload/content_hash；HTTP 和 WS reparse 均核对原字节。
5. 建立缺口后先收到较新 WS 成交，HTTP 水位不得跳过中间范围；补偿完成标记与 HTTP 提交原子性通过失败注入验证。
6. 强制连接关闭、订阅超时、pong 缺失和队列溢出后，重连、快照去重和范围补偿可恢复；缺口只能在 HTTP 扫描成功后更新状态。使用可控时钟验证 5、10、20、40、60 秒退避、0～20% 抖动及 60 秒最终上限；短暂成功不重置，稳定订阅且心跳正常满 60 秒才重置，累计次数保留；滚动连接尝试限额有效，数据库不可用或处理能力未恢复时不尝试重连。
7. socket 或 HTTP 请求阻塞时租约仍能续期；候补实例不重复订阅，失效令牌不能提交事实或位置。
8. 数据库故障期间不无限积压，恢复后处理已归档消息和未完成范围；进程强杀后通过会话记录和水位恢复。
9. HTTP 仍受来源限制；无法证明的旧历史、未知模式、单毫秒饱和等不输出虚假完整性。
10. 两通道同时运行时优雅退出、状态接口错误映射和既有 stored 快照分页一致性通过。

固定来源测试设置已知成交时间、明确接收顺序和断线区间；完整性结论依据预置预期集合，不依赖真实账户恰好发生新成交。正常固定来源下，收到一条合法新消息后应在 5 秒内完成入库，证明处理路径不等待 30 秒轮询；真实网络延迟只记录实测值，不承诺固定上限。

### 9.2 版本完成标准

Docker 环境能成功连接并订阅一个地址，安静账户靠 pong 保持连接。真实来源快照和任务可追溯；观察到的新消息及时入库，不依赖 live 请求。实际新消息未发生时该项明确待观察，不能把历史快照写成实时新增验收通过。

跨通道去重、标准内容兼容、断线与数据库故障恢复、有限队列和租约隔离均通过固定来源测试。状态能分别说明连接、提交与补偿进度；来源可查询范围内的预置断线成交全部补齐，超出范围的情况明确保留警告和阻塞原因。

## 10. 交付、部署与手动验证

### 10.1 升级与部署

沿用 Docker Compose、原四个常驻服务及 PostgreSQL 18 数据卷，不升级数据库主版本、不发布 collector 端口、不删除已有水位。新增初始化脚本生成缺失凭证、保留已有角色密码并执行 0003 migration。

迁移前备份数据库和文件证据。旧程序严格检查 schema，升级时先停止三个应用、保留 postgres，再初始化和启动新版。不用 down -v 进行升级或普通重启验证。

本版已提供 init-v0.4.sh 和 WS 配置。下面是服务器 Docker 手动验收步骤；本机原生进程验收不代替容器部署验收。

### 10.2 手动验证步骤

#### 第 1 步：启用、升级和检查状态

保留实际 account/start_time，在 config/trade-collector.toml 中设置 websocket.enabled=true：

```sh
docker compose stop query-api trade-log-query trade-collector
sh scripts/init-v0.4.sh
docker compose up --build -d
docker compose ps -a
docker compose logs --tail=50 trade-collector
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/version
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/watch-accounts
```

**看什么：** 四个常驻服务运行，版本 0.4.0；连接和订阅成功、session_id 存在，HTTP 水位继续保留并推进。恢复目标未追赶完时允许 RECOVERING，不要求立刻 LIVE。

**通过标准：** migration 与双通道装配成功，旧记录和水位保留，没有新部署一套数据库。

#### 第 2 步：区分订阅快照与实时更新

```sh
docker compose exec -T postgres psql -U robotech_admin -d robotech -c "SELECT query_id,transport,message_mode,status,created_at FROM trade_log.collection_jobs WHERE transport='WEBSOCKET' ORDER BY created_at DESC LIMIT 10;"
```

若观察期间账户产生新成交，再执行：

```sh
curl --noproxy '*' -i --get http://127.0.0.1:8080/api/v1/trade-events \
  --data-urlencode 'account=<配置中的实际账户>' \
  --data-urlencode 'source=stored' \
  --data-urlencode 'limit=10'
```

**看什么：** SNAPSHOT、LIVE_UPDATE 或 UNKNOWN 被明确分类；实际新消息归档后 stored 可见，last_received_at 与 last_committed_at 分开。记录来源成交、接收与提交时间。

**通过标准：** 数据消息成功入库且不是等待 HTTP 轮询才可见。没有新成交时只记录订阅和快照验证，实时新增部分待观察。

#### 第 3 步：跨通道去重

在 WS 保存某笔成交后等待至少一个成功 HTTP 重叠轮次，再执行：

```sh
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

用实际 fact_id 替换占位符查看证据通道：

```sh
docker compose exec -T postgres psql -U robotech_admin -d robotech -c "SELECT DISTINCT r.transport,r.message_mode,j.query_id FROM trade_log.fact_observations o JOIN trade_log.raw_logs r ON r.id=o.raw_log_id JOIN trade_log.collection_jobs j ON j.id=r.collection_job_id WHERE o.fact_id='<实际fact_id>';"
```

**看什么：** 查重 0 rows；同一 fact_id 可关联 HTTP 与 WS 观察，revision 仍为 1，没有因原始封装不同报业务冲突。

**通过标准：** 同成交只保存一份当前事实，双方原始证据均保留。

#### 第 4 步：安静账户与心跳

```sh
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/watch-accounts
docker compose logs --tail=40 trade-collector
```

等待超过一分钟再重复查看。

**看什么：** 没有新成交时 last_received_at 可不变，last_pong_at 持续更新，连接不会仅因账户安静而反复重连。HTTP last_success_at 可以继续更新。

**通过标准：** 连接健康与成交活跃度没有混淆。

#### 第 5 步：进程断线与恢复范围

先记下 session_id、HTTP scanned_through 和事实数量，再执行：

```sh
docker compose stop trade-collector
# 等待一个或多个原配置轮询间隔后启动。
docker compose start trade-collector
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/watch-accounts
docker compose exec -T postgres psql -U robotech_admin -d robotech -c "SELECT reason,start_ms,end_ms,status,detected_at,scanned_at FROM trade_log.collection_gaps ORDER BY detected_at DESC LIMIT 5;"
```

**看什么：** 新 session_id，HTTP 从旧水位继续；恢复范围先扫描，目标被已提交水位覆盖后才标记 HTTP_SCANNED。重连快照不重复事实。

**通过标准：** 程序恢复和范围扫描均可追踪。真实停机期间有没有新成交需要另行核对，不能只凭最后状态宣称全部补齐。

#### 第 6 步：数据库故障与恢复

仅在本项目独立验收环境执行：

```sh
docker compose stop postgres
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/watch-accounts
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/health
docker compose logs --tail=40 trade-collector
docker compose start postgres
```

恢复后等待重连和 HTTP 追赶，再重复状态、stored 与 gap 查询。

**看什么：** 状态接口 503、网关 health 200；WS 不无限积压，故障期不伪造提交时间或 HTTP 水位。恢复后旧记录保留，接续会话和范围可见。

**通过标准：** 两通道均从持久化状态恢复；补偿可查询范围而不承诺已超出来源范围的数据完整。

#### 第 7 步：重放 WS 任务

从第 2 步选择一个 COMPLETED WS query_id：

```sh
docker compose exec -T trade-log-query trade-log-query reparse --query-id query_实际值
```

**看什么：** comparison=SAME，原始 envelope、元数据与逐笔来源索引能够复算标准事实；不会更新 revision、水位或缺口状态。

**通过标准：** 新通道证据可追溯、可重放，旧 HTTP 任务仍兼容。

#### 第 8 步：关闭 WS 验证 HTTP 兼容

将 websocket.enabled 改为 false，保留账户和起点：

```sh
docker compose restart trade-collector
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/watch-accounts
```

等待一个成功 HTTP 轮次后再次查看。

**看什么：** HTTP_ONLY、websocket.status=DISABLED，原 scanned_through 与 last_success_at 继续更新；旧事实不变。

**通过标准：** 可退回纯 HTTP 模式，不需要删除检查点或另写采集程序。关闭 WS 不撤销已有 gap 历史；尚未完成的 HTTP 补偿仍应继续。

订阅超时、pong 缺失、错账户、队列溢出、失效租约和来源上限通过固定服务集成测试验证，不要求用户修改服务器防火墙或用真实账户碰巧制造。实际报告分别记录开发测试、真实来源观察和 Docker 部署结果。

#### 10.2.1 固定服务器验收脚本

新增 `scripts/verify-v0.4.sh`，共用 `scripts/verify-version.py`；版本入口固定，按版本启用相应能力。脚本测试位于 `scripts/tests/test_verify_version.py`。在项目根目录执行，需要 Python 3（标准库）、Docker 和 Compose，不需要 Rust、pip 或 jq。

```sh
sh scripts/verify-v0.4.sh
```

检查内容：库存、HTTP 推进及重放，并检查 WS 启用、pong 更新、缺口归零、WS 成功提交证据、跨通道观察关联和 WS 重放。UNKNOWN 允许正常保存，不当作错误；安静账户不要求出现新成交。

默认读取唯一配置账户；需明确选择时使用 `--account <实际地址>`。每段观察默认最多 90 秒，约每 15 秒报告等待进度；较长轮询间隔或扫描追赶可增加等待上限。

逐项输出 PASS、FAIL、SKIP 及汇总。退出码 0 表示所有已执行项目通过；1 表示检查失败；2 表示参数错误或缺少 Python。存在 SKIP 时明确提示非完整验收。无成交、库存不足以跨页、没有可重放任务或双通道共同事实，不伪报对应项目通过。接口失败、计数不符、重复事实或恢复超时均判 FAIL。

默认不停止服务、不修改配置，不构建镜像或迁移。SQL 使用只读事务；reparse 只读取证据并比较。基础健康检查仅覆盖基础契约，需要完整 v0.0 验收时另执行 `sh scripts/verify-v0.0.sh`。

仅在允许中断的独立验收环境执行：

```sh
sh scripts/verify-v0.4.sh --lifecycle
sh scripts/verify-v0.4.sh --database-fault --fault-seconds 15
```

`--lifecycle` 停止并启动 trade-collector，检查正常退出和恢复，以及适用版本的旧事实保留与水位不回退；会短暂中断对应能力。 `--database-fault` 停止并恢复本项目 PostgreSQL，检查业务 503、网关 health 200 和恢复结果，影响所有使用该数据库的服务。 停止测试在 finally 中尝试恢复服务；启动失败或被强制终止时需人工确认。前置检查失败时跳过停止测试。

部署参数可调整：

```sh
sh scripts/verify-v0.4.sh --base-url http://127.0.0.1:8080 --wait-seconds 180 --expected-version 0.4.0
```

不带 `--expected-version` 时允许在后续兼容版本上验收。URL、账户与本地 Compose 项目应指向同一部署。 `--live` 会主动访问官方来源，并保存证据和成交。

来源限流、业务冲突、饱和毫秒及真实历史完整性仍由受控开发测试或独立明细对账覆盖。 订阅失败、pong 超时、队列溢出及关闭 WS 不自动制造，按本章手动步骤或开发测试验证。 HTTP_SCANNED 不表示所有历史交易完整。

2026-10-06 开发验证：失败判定及恢复路径的 9 项标准库测试、shell 语法和帮助检查通过。v0.2 使用真实本机网关、查询服务及隔离 PostgreSQL 验证参数错误、stored 计数、去重 SQL 与 reparse=SAME；Docker 命令由临时替身衔接本机程序。开发机无 Docker，真实服务器 Docker 执行和停止测试仍待验收，不记为通过。


### 10.3 交付清单

交付 WS 来源适配器、双通道运行监督、标准业务指纹、原始消息与会话/gap 迁移、原子提交与恢复、状态 DTO 扩展、配置、初始化脚本、必要测试及统一接口文档更新。

手动验证与固定服务器验收脚本按本章执行。后续 v0.5 在已保存的事实和 message_mode 基础上输出候选信号，不能把本版快照或 HTTP 回补当作新增实时交易自动发布。

## 11. bug记录

每个 bug 按“描述、影响、确认记录、修改逻辑、修改记录摘要、如何验证、验证记录”七项记录。开发环境验证和服务器部署验证分别记录；未完成的验证不得记为通过。

### 11.1 BUG-0.4-001：租约续期等待导致采集任务暂停

#### 1 描述

HTTP、WebSocket 与租约续期由同一个异步监督任务驱动。旧实现进入心跳分支后直接等待数据库续期，不再轮询 HTTP 和 WebSocket。若采集事务正在持有检查点锁，续期等待该锁，而采集任务无法继续提交释放锁，可能一直等待到数据库语句超时。旧实现退出时顺序等待两条任务，也存在类似的相互等待风险。

#### 2 影响

可能造成 HTTP 轮次失败、WebSocket 提交超时、租约重新获取和会话重建，增加补偿延迟。故障期间的数据是否缺失，需要结合持久化证据及来源明细判断；成功恢复连接不等于历史完整。

#### 3 确认记录

2026-10-06 收到服务器排查结果：租约 UPDATE、HTTP 检查点 SELECT FOR UPDATE、WebSocket INSERT 多次在约 10 秒后超时；PostgreSQL 同时记录 `canceling statement due to statement timeout`。例如北京时间 2026-10-06 08:45:02，HTTP 报 `DEPENDENCY_UNAVAILABLE`，随后 WebSocket 停止。

故障后的 `pg_stat_activity` 快照显示连接为 idle，blockers 为空，只能说明检查时没有阻塞，不能排除此前的瞬时等待。代码检查确认旧监督循环会暂停采集任务；回归测试用延迟提交模拟检查点锁等待。服务器上每次超时是否均由此原因造成，尚未确认。

#### 4 修改逻辑

把心跳续期改为独立的异步 future，与 HTTP、WebSocket 一起参与 select。续期等待数据库时，监督循环继续轮询两条采集任务。任一任务结束、续期失败或收到停止信号后，取消本轮任务，并并发等待 HTTP 和 WebSocket 退出，再释放租约。新增 `collection_lease_renewal_failed` 日志，记录租约代数和错误码。

#### 5 修改记录摘要

| 日期 | 文件 | 类型 | 修改摘要 |
| --- | --- | --- | --- |
| 2026-10-06 | `services/trade-log/adapters/src/collector_runtime.rs` | 更新 | 并发驱动续期与采集；并发等待退出；记录续期失败 |
| 2026-10-06 | `services/trade-log/adapters/tests/websocket_chain.rs` | 更新 | 增加 `renewal_waiting_on_http_lock_does_not_stall_http_commit` 回归测试 |

#### 6 如何验证

开发回归测试在隔离 PostgreSQL 中安装测试触发器，让 HTTP 更新检查点时延迟 3 秒；租约为 6 秒、续期间隔为 2 秒，令续期遇到未完成的 HTTP 事务。通过标准：HTTP 在测试期限内成功提交，租约代数保持 1，程序能够正常停止。触发器只用于独立测试数据库，不在部署数据库安装。

部署修复后的采集程序：

```sh
docker compose up -d --build --no-deps trade-collector
docker compose logs -f --tail=50 trade-collector
```

另一个终端执行以下命令，等待超过一分钟后重复查询：

```sh
curl --noproxy '*' -sS http://127.0.0.1:8080/api/v1/watch-accounts
docker compose logs --no-color --since=10m trade-collector \
  | grep -E 'collection_completed|collection_failed|collection_lease_renewal_failed|websocket_stopped|statement timeout'
```

通过标准：HTTP scanned_through 和 last_success_at 继续推进；WebSocket 心跳正常；没有反复出现同类续期等待超时。继续观察原先出现间歇性错误的运行时段；短时间没有错误不能证明所有数据库问题已消除。

#### 7 验证记录

- 2026-10-06，开发环境：新增的并发续租回归测试通过；模拟慢提交期间 HTTP 成功完成，租约代数未增加。
- 2026-10-06，开发环境：HTTP 采集与 WebSocket 集成测试共 13 项通过，工作区常规测试通过，Clippy（`-D warnings`）及 diff 格式检查通过。
- 服务器部署验证：待执行，尚未收到修复后观察结果。

### 11.2 BUG-0.4-002：临时数据库错误导致恢复缺口永久 BLOCKED

#### 1 描述

旧数据库错误映射把 `DEPENDENCY_UNAVAILABLE` 标为不可重试。WebSocket 结束处理据此把恢复缺口标为 BLOCKED。后续 HTTP 成功提交仅关闭 OPEN、SCANNING 缺口，不能解除已有的数据库错误 BLOCKED 状态。

#### 2 影响

临时故障可能被当作不可恢复错误，造成恢复状态长期 BLOCKED、open_gap_count 无法归零。HTTP 水位已经越过缺口目标时，状态仍可能报告阻塞。这种状态不一致不能直接证明该范围没有成交数据，也不能作为历史完整的证据。

#### 3 确认记录

2026-10-06 服务器检查结果：北京时间 08:06:13.073 至 08:14:00.816 的缺口为 BLOCKED，错误码 DEPENDENCY_UNAVAILABLE；HTTP 水位已到 08:44:15.813，连续失败次数为 0，但 recovery.status 仍为 BLOCKED，open_gap_count 为 1。

代码检查确认：数据库错误默认 retryable=false；`end_stream` 将不可重试错误对应的缺口阻塞；`http_recovered` 的更新条件排除了 BLOCKED。

#### 4 修改逻辑

数据库操作错误仍映射为 DEPENDENCY_UNAVAILABLE，但设置 retryable=true，允许恢复重试。兼容已有数据：在 HTTP 事实及水位成功提交的同一事务内，允许解除错误码为 DEPENDENCY_UNAVAILABLE 的旧 BLOCKED 缺口；必须已固定 end_ms，且成功扫描水位达到该目标。

VERSION_CONFLICT 等真实数据错误的 BLOCKED 缺口不自动解除。未达到目标时不关闭缺口，不通过手工修改状态宣称补偿完成。HTTP_SCANNED 仍只表示完成来源可查询范围的扫描，不承诺来源历史完整。

#### 5 修改记录摘要

| 日期 | 文件 | 类型 | 修改摘要 |
| --- | --- | --- | --- |
| 2026-10-06 | `services/trade-log/adapters/src/postgres/mod.rs` | 更新 | 数据库操作错误设置为可重试 |
| 2026-10-06 | `services/trade-log/adapters/src/postgres/realtime.rs` | 更新 | HTTP 成功提交时兼容解除旧数据库错误 BLOCKED 缺口 |
| 2026-10-06 | `services/trade-log/adapters/tests/websocket_chain.rs` | 更新 | 增加 `successful_http_scan_clears_legacy_database_block_but_not_data_conflict` 回归测试 |

#### 6 如何验证

开发测试构造旧数据库错误 BLOCKED 缺口，分别验证：水位未到目标时保持 BLOCKED；到达目标后变为 HTTP_SCANNED，open_gap_count 归零；将错误改为 VERSION_CONFLICT 后，即使水位超过目标仍保持 BLOCKED。

按第 11.1 节部署后，等待一个成功 HTTP 轮次，再执行：

```sh
curl --noproxy '*' -sS http://127.0.0.1:8080/api/v1/watch-accounts
docker compose exec -T postgres psql -U robotech_admin -d robotech <<'SQL'
\pset pager off
SET TIME ZONE 'Asia/Shanghai';
SELECT checkpoint_key, reason, status,
       to_timestamp(start_ms / 1000.0) AS gap_start,
       to_timestamp(end_ms / 1000.0) AS gap_end,
       scanned_at, last_error
FROM trade_log.collection_gaps
ORDER BY detected_at DESC LIMIT 10;
SQL
```

通过标准：已有数据库错误缺口在成功扫描覆盖目标后变为 HTTP_SCANNED，scanned_at 有值、last_error 清空；没有其他待恢复缺口时 open_gap_count 为 0。真实业务冲突保持 BLOCKED。

独立验收环境可继续执行第 10.2 节第 6 步，验证新的数据库故障恢复后，两通道恢复、旧记录保留、HTTP 水位推进、临时数据库错误不永久阻塞缺口。不要在生产环境为验证而停止数据库。

#### 7 验证记录

- 2026-10-06，开发环境：旧数据库错误缺口回归测试通过，覆盖目标之前不关闭，达到目标后关闭，VERSION_CONFLICT 仍阻塞。
- 2026-10-06，开发环境：本次两项新增回归测试及既有 HTTP/WebSocket 集成测试共 13 项通过；工作区常规测试、Clippy 和 diff 格式检查通过。
- 服务器旧缺口恢复及部署后故障验证：待执行。此次修复不自动回补 HTTP 配置起点之前的历史记录。
