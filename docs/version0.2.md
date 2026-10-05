# v0.2 开发文档：成交持久化与时间区间查询

- 文档版本：0.2
- 状态：设计稿，尚未实现；以下新增接口、配置、命令均为本版交付要求
- 更新日期：2026-10-05
- 设计依据：[概要设计](overview-design.md)、[详细设计](detailed-design.md)、[产品版本路线图](roadmap.md)
- 前置版本：[v0.1 开发文档](version0.1.md)
- 接口契约统一维护在 [接口说明](api.md)，本文件描述实现与验收，不另建版本接口文档

## 1. 版本目标

将成功采集的成交保留为可跨次查询、按时间检索的标准事实。重复采集相同成交不会产生重复事实，程序和数据库重启后记录仍在；已保存的来源响应可以离线重新解析并核对结果。

v0.1 的文件证据按单次 query_id 保存，不能提供跨查询合并账本。v0.2 在这套获取、解析、身份及 HTTP 基础上接入正式存储，保留 query-api 与 trade-log-query 两个应用进程，Compose 新增 PostgreSQL 数据库容器。

用户仍可主动查询官方近期成交；也可以只查本地已保存的成交，不访问 Hyperliquid。保存记录不等于补齐完整历史，系统必须区分“库中有这些记录”和“该时间段完整”。程序与应用镜像版本升级为 0.2.0，已有接口路径及 schema_version=1 保持兼容。

## 2. 功能范围

| 能力 | 本版交付 | 边界 |
| --- | --- | --- |
| 采集后保存 | 实时查询成功前完成标准事实数据库事务 | 不自动定时采集 |
| 原始响应保存 | 原始字节、请求、来源状态、摘要、元数据分别持久化 | 不引入对象存储，先保存在 PostgreSQL |
| 成交去重 | 沿用 hl-fill-v1 稳定身份，唯一约束保障并发写入幂等 | 不按订单 oid 聚合多个 fill |
| 时间查询 | 查询本库已保存的指定账户、时间范围成交，支持游标分页 | 不将时间参数转换为来源历史回补 |
| 查询证据 | 记录每次采集、解析、保存状态及事实关联 | 不建设完整审核服务 |
| 重新解析 | 按 query_id 离线读取原始响应，复用协议解析器比较结果 | 不自动生成事实修订或覆盖冲突事实 |
| 升级迁移 | 导入 v0.1 文件证据，保留原 query_id 和来源引用 | 不信任旧 result.json 作为事实源 |

本版不新增定时任务、采集水位、WebSocket、消息系统、outbox 发布、webhook、账户分析、现货展示、跟单执行或事实修正流程。事实版本表按详细设计落实，当前只写 revision=1，后续版本沿用；不提前创建其他业务服务的表。

## 3. 使用方式

### 3.1 保留实时查询

```text
GET /api/v1/trade-events?account=<ADDRESS>&limit=100
```

省略 source 或使用 `source=live`：沿用 v0.1 来源查询语义，查询当前官方窗口，完成持久化后响应。相同调用再次执行仍会采集新原始响应，但已存在且相同的事实不重复写入。

limit 默认 100、范围 1–2000，仍只限制展示条数；成功解析出的完整事实集合全部入库。live 不允许 start_time、end_time 或 cursor，避免误以为已经执行历史范围采集。

### 3.2 查询本库已保存记录

```text
GET /api/v1/trade-events?account=<ADDRESS>&source=stored&start_time=<UTC>&end_time=<UTC>&limit=100
```

source=stored 时只调用交易日志服务的数据库查询端口，不调用来源、不生成新的采集 query_id，也不改写事实。start_time/end_time 均可省略，省略表示该方向无时间限制；两者都存在时必须 start_time < end_time。

区间为 `[start_time, end_time)`：包含开始时刻，不包含结束时刻。时间必须为带时区的 RFC 3339，最大毫秒精度；输入偏移换算为 UTC，输出统一 UTC。超出支持精度、负 Unix 时间或非法时间返回 400。

返回顺序为成交时间降序、tid 数值降序、fact_id 升序，与 live 一致。分页使用 next_cursor，不提供 offset。下一页携带原账户、时间范围、limit 和服务器返回的 cursor；游标固定分页快照，后续入库成交不插入正在浏览的分页结果。

空结果返回 200，说明“本库此查询条件下没有保存记录”；不能解释为账户没有交易，也不触发自动来源补采。

### 3.3 重新解析和导入

由运维在内部容器执行命令，不新增公网管理接口：

```sh
docker compose exec -T trade-log-query trade-log-query reparse --query-id query_实际值
docker compose exec -T trade-log-query trade-log-query import-evidence --directory /var/lib/robotech/trade-log
```

两个命令读取与服务相同的配置和数据库凭证。重新解析默认只核对；导入逐目录执行可恢复的事务，同一目录重复导入不重复入账。CLI 的详细退出码和输出字段见第 8 章。

## 4. 数据来源与查询边界

### 4.1 来源口径沿用

沿用 Hyperliquid userFills、meta、spotMeta，网络由配置固定；继续只支持元数据确认的默认永续市场。保持来源金额精度、费用符号、持仓效果、主动/强制分类及事实身份规则。

保存网络查询实际返回的原始字节，不以重新序列化后的 JSON 替代；保存状态码和 SHA-256，不记录授权头。一次响应及其重试分别保存，三类成功响应共同组成可重新解析的输入。

### 4.2 覆盖与库存分开

live 保留 SOURCE_WINDOW_UNVERIFIED、LIMITED 及 v0.1 警告。stored 的覆盖明确为 `STORED_RECORDS_ONLY`，只声明当前库存；不通过最早和最晚成交推断期间完整，没有记录的区间也不推断无交易。

stored 的 observed_range 按本次时间条件和固定分页快照下的完整匹配集计算，不能只用本页首尾生成。响应明确 request_range、snapshot_seq、matched_records、returned_records、has_more、next_cursor 和 warnings；默认警告 STORED_HISTORY_NOT_VERIFIED。

本版不证明连续覆盖、不补更早历史、不设置自动采集水位。v0.3 以后增加采集状态时复用保存接口，不改变本版事实身份。

## 5. 接口与持久化数据设计

### 5.1 接口增量

完整参数和字段维护在统一 [接口文档](api.md)，新增内容必须标为“v0.2 设计，尚未实现”。实现后才更新为已支持。

| 接口变化 | 设计约定 |
| --- | --- |
| 外部 GET /api/v1/trade-events | 增加 source，默认 live；stored 支持 start_time、end_time、cursor |
| 内部 POST /internal/v1/trade-queries | 保持原语义，成功前保存事实；不接收 stored 参数 |
| 内部 POST /internal/v1/stored-trade-queries | 接收账户、区间、limit、cursor；仅查询数据库 |
| 内部认证 | 新接口沿用 Bearer 服务凭证和 trace_id 传播 |
| 公共错误 | 沿用顶层 code/message/trace_id，不增加另一种错误包装 |

live 保留原有全部响应字段，新增 `persistence`：status=COMMITTED、inserted_records、existing_records。二者之和等于本次 perpetual_records；其中 existing 是此前已存在的事实数，不是原始数组内 duplicate_records。

stored 使用相同 data/meta 信封及 trades[] 标准事实，但 data 结构按读取场景定义：account、network、query_scope=STORED_TIME_RANGE、coverage=STORED_RECORDS_ONLY、request_range、snapshot_seq（十进制字符串）、matched_records、returned_records、observed_range、trades、has_more、next_cursor、warnings。没有采集 query_id、live 来源计数或 display_truncated；next_cursor 无下一页时为 null。

trades[] 保留持久化时的事实身份及首次来源 raw_log_id。每次 HTTP 的 trace_id 独立；重新读取不修改事实。调用方根据 source 区分响应，未指定 source 的已有客户端行为保持一致。

### 5.2 数据库归属与模型

使用 PostgreSQL，数据库 robotech、独立 schema trade_log；网关不持有数据库凭证，也不直接查表。时间列为 TIMESTAMPTZ，应用及运维会话输出使用 UTC，原始时间在来源 payload 中保留。[PostgreSQL 时间类型说明](https://www.postgresql.org/docs/18/datatype-datetime.html)

按详细设计落实以下表；组合字段在 SQL migration 中展开为独立列：

| 表 | 主要字段及约束 | 本版职责 |
| --- | --- | --- |
| collection_jobs | UUID id；UNIQUE query_id；chain_id/protocol/source_id；mode=HISTORICAL；status；request JSONB；trace_id；counts/coverage/error；created_at/updated_at | 保存一次按需采集或导入任务，HISTORICAL 不代表完整历史；已有 query_ 字符串为兼容外部标识 |
| raw_logs | UUID id；collection_job_id FK；chain_id/protocol/source_id；source_event_id；chain_position JSONB；ordering_key；payload JSONB 可空；body BYTEA；http_status；request JSONB；sha256；observed_at；parse_status/parser_version | 来源响应归档；非法 JSON 也保存 body，不能要求 payload 必填导致错误证据丢失 |
| account_fact_versions | PK(fact_id, revision)；UNIQUE event_id；account_key/fact_type；ordering_key/sub_index；occurred_at；完整标准事实 JSONB；raw_log_id UUID FK；parser_version；content_hash；created_at | 保存事实版本，v0.2 固定 revision=1；event_id 由 fact_id 和 revision 确定性生成 |
| account_facts_current | PK fact_id；current_revision；account_key/fact_type；ordering_key/sub_index；occurred_at；source_tid NUMERIC(20,0)；is_retracted；ingest_seq BIGINT UNIQUE；updated_at | 当前事实指针、查询索引及分页快照序列；复合 FK 指向事实版本 |
| fact_observations | PK(fact_id, revision, raw_log_id, source_index)；复合 FK 事实版本及 raw_logs FK | 同一事实对应多次来源响应与多个位置，不覆盖首次来源引用 |
| ingestion_state | 单行主键；committed_seq BIGINT | 为事实提交分配顺序，保证固定分页快照 |

raw_logs 的 source_event_id 使用兼容的原始响应字符串 ID，network/source_id 纳入唯一约束；UUID 是库内主键，不替换 API 的 raw_log_id 字符串。chain_position 保存 query_id、请求类型和 attempt；ordering_key 保存采集时间与响应身份。

创建索引 `(account_key, occurred_at DESC, source_tid DESC, fact_id)`，当前事实与版本按主键连接；查询只读取未撤回的 TRADE。NUMERIC tid 用于精确排序，大整数不可转浮点数。金额沿用标准 payload 的十进制字符串，本版不建立金额聚合列；未来独立数值列按详细设计采用 NUMERIC。

collection_jobs、raw_logs 的补充列和新增关联表属于本版具体化；不另建 trades_v02 等平行事实模型。raw_logs.payload 可空是对整体设计“保留错误原始响应”的必要补充，详细设计同步注明。

### 5.3 内容一致性和事务

fact_id 沿用 v0.1 固定编码及摘要算法。比较内容使用规范化事实业务字段和解析后的原始 source_fill，排除 raw_log_id、source_indices、采集时间等观测信息；JSON 对象键按固定规则规范化后计算 content_hash。价格精度、类型和原始内容变化必须可检测，不凭 fact_id 相同就忽略冲突。

同身份同内容：保持 revision=1，不改写既有事实，只追加新 observation。同身份不同内容：本次事实事务全部回滚，返回 409 VERSION_CONFLICT；保留原始响应和失败任务，不静默覆盖，不自动增加 revision。

网络采集不持有数据库长事务。原始响应接收后先以短事务归档，然后解析；完整事实集合、observation、任务 COMPLETED 和提交顺序在同一事务提交。发生解析/保存失败时任务标 FAILED；进程中断留下 RUNNING，重启时标 INTERRUPTED，不把未完成任务当成功。

事务提交后才返回 live 成功。提交结果不确定时返回明确失败，客户端重复查询依靠稳定身份恢复；不会为了给出 200 而只保存前 limit 条。

### 5.4 分页一致性

所有事实写入事务短暂锁定 ingestion_state 单行，先检查内容冲突，再为新事实分配递增 ingest_seq，最后一起提交。首个 stored 页面读取 committed_seq 作为 snapshot_seq，后续页面仅查询 ingest_seq <= snapshot_seq 的事实；不能只用数据库 sequence 的 last_value 充当已提交水位。

以成交排序三元组做 keyset 分页，cursor 编码版本、账户/网络/时间条件、snapshot_seq 及最后一条排序键。cursor 是有长度上限的 base64url JSON，不含 SQL；服务器校验格式、值域及条件匹配，用参数化查询。它不是权限令牌，也不承担鉴权。格式错误、条件不一致或未来水位返回 400。

limit 可在下一页调整；身份和时间范围必须保持一致。匹配数和页面数据在同一读事务下获得，使用 REPEATABLE READ，并按 snapshot_seq 过滤。[PostgreSQL 事务隔离说明](https://www.postgresql.org/docs/18/transaction-iso.html)

## 6. 处理流程

### 6.1 实时采集并保存

1. 网关校验参数、生成 trace_id，调用原内部查询接口。
2. 交易查询服务分配 query_id，建立采集任务，应用并发和总时限。
3. 请求来源，保存每次响应的字节、请求、元数据及摘要；数据库失败则中止，不伪造空结果。
4. 复用 Hyperliquid 解析器和交易日志校验，生成完整事实及覆盖说明。
5. 执行幂等事实事务，写版本、当前指针、来源关联和任务结果。
6. 数据库提交成功后应用展示 limit，返回原结果和 persistence 计数。

保留 v0.1 文件证据作为诊断副本；v0.2 数据库原始响应与任务状态是持久化权威。文件镜像失败记录日志和任务诊断，不把已提交事实伪装成未入库；原始数据库归档失败则不能返回成功。文件镜像目录格式尽量沿用，不需要它才能执行数据库重新解析。镜像任务受有界时限控制，数据库事实成功提交后镜像失败不会回滚；两者不能宣称跨介质原子事务。

### 6.2 本库查询

1. 校验 account、source、区间及 cursor，构造规范化过滤条件。
2. 内部服务开启读取事务，确定或校验分页 snapshot_seq。
3. 按相同条件查询匹配数、完整观察范围和 limit+1 条事实。
4. 返回本页、has_more 和下一页 cursor，披露 STORED_RECORDS_ONLY。

该流程独立于外部来源，不使用 live 的网络查询并发许可；使用受控数据库连接池和查询超时，不无限等待。

### 6.3 离线重新解析及导入

reparse 从数据库读取指定任务的成功 userFills/meta/spotMeta 响应，校验摘要，使用当前解析器重新解析并比较该任务的持久化事实引用，输出 SAME 或 DIFFERENT；失败任务输入不齐时输出 INCOMPLETE。默认不写事实，不发网络请求。

import-evidence 遍历本地 v0.1 查询目录；只导入 COMPLETED 且网络、请求、三个成功原始响应、SHA-256 均可核对的目录。重新解析 raw body，不以 result.json 代替来源；使用原 query_id、raw 字符串标识和时间，按照同一幂等写入流程保存。

一个目录失败不回滚已完成目录；输出每个目录及最终汇总，任一失败整体退出非零。同 query_id 同内容重复导入跳过；同 query_id 内容不同返回冲突。历史证据缺少元数据或损坏时明确拒绝，不自动补网络数据。

## 7. 程序结构设计

### 7.1 沿用目录和进程

```text
services/trade-log/
  src/
    acquisition/ raw_log/ parsing/ normalization/ validation/
    query/                       # 沿用实时查询业务
    persistence/                 # 事实事务与幂等端口
    stored_query/                # 已保存事实查询端口
    replay/                      # 重新解析与导入业务
  adapters/src/
    file_evidence.rs             # 诊断副本和旧证据读取
    postgres/                    # 原始响应、事实及查询适配器
    internal_http.rs             # 新增内部读取路由
  migrations/                    # 本服务数据库迁移
  bins/trade-log-query/           # 沿用进程，新增 migrate/reparse/import-evidence 子命令
gateway/query-api/
  src/clients/trade_log.rs        # 按 source 选择内部端口
  src/http/handlers/trade_events.rs
config/
  query-api.toml
  trade-log.toml
scripts/
  init-v0.2.sh                    # 凭证、数据目录和数据库初始化入口
```

复用 shared-types、account-facts、protocol-api、service-runtime 及 Hyperliquid 来源实现。不复制 v0.1 解析代码，不让 HTTP handler 执行 SQL；业务 crate 只定义使用方 trait，不依赖 SQLx、Axum 或 PostgreSQL。

### 7.2 数据库适配器

在 trade-log/adapters 引入 SQLx 或等价的异步 PostgreSQL 驱动，统一依赖与锁文件。DTO、业务事实和数据库 row 分离，映射发生在适配器内。迁移随本服务发布，不放到 gateway。

运行角色只能读写 trade_log 业务表；迁移角色负责 DDL；数据库管理密码不交给网关。事务与连接释放纳入原有生命周期和有界退出。v0.3 采集入口通过同一 persistence 端口写事实，v0.5 再接入 outbox，禁止另建一套去重逻辑。

### 7.3 本版计划新增与更新文件

以 v0.1 交付状态为基线。**v0.2 尚未实现，下面是计划文件清单，不代表文件已经创建或更新。** 新增模块采用与既有代码一致的 mod.rs 组织方式；实现时若文件拆分有变化，应同步修订本章，不能另起一套版本目录。

#### 计划新增文件

| 文件 | 计划职责或变更点 |
| --- | --- |
| `services/trade-log/src/persistence/mod.rs` | 新增幂等保存、事务结果及使用方存储端口 |
| `services/trade-log/src/stored_query/mod.rs` | 新增时间范围、分页快照、游标与库存查询业务端口 |
| `services/trade-log/src/replay/mod.rs` | 新增离线重新解析和旧证据导入业务 |
| `services/trade-log/adapters/src/postgres/mod.rs` | 导出 PostgreSQL 适配器并装配连接池 |
| `services/trade-log/adapters/src/postgres/raw_log.rs` | 数据库原始响应归档、任务状态与响应读取 |
| `services/trade-log/adapters/src/postgres/facts.rs` | 事实版本、当前指针、来源关联及幂等事务 |
| `services/trade-log/adapters/src/postgres/stored_query.rs` | 库存筛选、统计和固定快照 keyset 分页 SQL |
| `services/trade-log/adapters/src/postgres/migration.rs` | 服务 migration 执行及 schema 版本检查 |
| `services/trade-log/migrations/0001_trade_log_storage.sql` | 初始 trade_log schema、表、约束和索引；由迁移管理，不在 handler 建表 |
| `scripts/init-v0.2.sh` | 数据库凭证、角色、数据卷和 migration 初始化入口 |
| `services/trade-log/adapters/tests/postgres_persistence.rs` | 数据库幂等、并发、冲突、事务回滚及重启保留验证 |
| `services/trade-log/adapters/tests/stored_query.rs` | 时间边界、排序、游标和分页快照验证 |
| `services/trade-log/adapters/tests/replay_import.rs` | 重新解析、旧证据重复导入和损坏证据验证 |
| `tests/fixtures/v0.2/import-evidence/query_fixture/manifest.json` | 固定 COMPLETED 导入样本的任务说明 |
| `tests/fixtures/v0.2/import-evidence/query_fixture/requests/raw_query_fixture_userFills_1.json` | 固定成交来源请求样本 |
| `tests/fixtures/v0.2/import-evidence/query_fixture/requests/raw_query_fixture_meta_1.json` | 固定永续元数据请求样本 |
| `tests/fixtures/v0.2/import-evidence/query_fixture/requests/raw_query_fixture_spotMeta_1.json` | 固定现货元数据请求样本 |
| `tests/fixtures/v0.2/import-evidence/query_fixture/responses/raw_query_fixture_userFills_1.body` | 固定成交原始字节 |
| `tests/fixtures/v0.2/import-evidence/query_fixture/responses/raw_query_fixture_meta_1.body` | 固定永续元数据原始字节 |
| `tests/fixtures/v0.2/import-evidence/query_fixture/responses/raw_query_fixture_spotMeta_1.body` | 固定现货元数据原始字节 |
| `tests/fixtures/v0.2/import-evidence/query_fixture/metadata/raw_query_fixture_userFills_1.json` | 对应成交响应状态、身份和摘要 |
| `tests/fixtures/v0.2/import-evidence/query_fixture/metadata/raw_query_fixture_meta_1.json` | 对应永续元数据响应状态、身份和摘要 |
| `tests/fixtures/v0.2/import-evidence/query_fixture/metadata/raw_query_fixture_spotMeta_1.json` | 对应现货元数据响应状态、身份和摘要 |
| `tests/fixtures/v0.2/import-evidence/query_fixture/result.json` | 旧证据兼容样本的完整结果；导入不能将其当作事实源 |

导入样本是构造数据，query_fixture 及原始响应 ID 必须在 manifest、请求、元数据和结果中相互一致；原始字节的 SHA-256 按实际样本生成。样本不得使用部署凭证或实际服务器地址。

#### 计划更新文件

| 文件 | 计划职责或变更点 |
| --- | --- |
| `Cargo.toml` | 版本升级 0.2.0，统一新增数据库驱动、游标等实际依赖 |
| `Cargo.lock` | 锁定新增依赖与更新后的 workspace 包版本 |
| `compose.yaml` | 新增 postgres、trade-log-migrate，数据库卷及凭证挂载，保持原应用服务 |
| `gateway/query-api/Dockerfile` | 增加实际迁移/初始化构建资源和 Cargo 缓存，保留两个运行 target |
| `config/query-api.toml` | 网关内部请求默认时限调整为 45 秒 |
| `config/trade-log.toml` | 新增 database 配置，整体查询默认时限调整为 40 秒 |
| `gateway/query-api/src/clients/trade_log.rs` | 增加库存查询调用，按 live/stored 解析对应业务响应 |
| `gateway/query-api/src/http/handlers/trade_events.rs` | 新增 source、时间与游标参数及分支；保持原 live 请求兼容 |
| `services/trade-log/src/lib.rs` | 导出 persistence、stored_query、replay 新模块 |
| `services/trade-log/src/query/mod.rs` | 实时成功返回前保存完整事实，增加 persistence 结果 |
| `services/trade-log/src/raw_log/mod.rs` | 扩展可重新读取数据库原始响应的证据端口，保持业务不依赖数据库 |
| `services/trade-log/src/validation/mod.rs` | 增加新查询模式、时间精度、参数组合及游标条件校验 |
| `services/trade-log/adapters/Cargo.toml` | 注册数据库驱动及集成测试所需依赖 |
| `services/trade-log/adapters/src/lib.rs` | 导出 postgres 适配器 |
| `services/trade-log/adapters/src/file_evidence.rs` | 支持诊断镜像与旧证据读取；保留原目录格式 |
| `services/trade-log/adapters/src/internal_http.rs` | 增加库存内部路由和状态映射，实时成功返回持久化计数 |
| `services/trade-log/bins/trade-log-query/Cargo.toml` | 增加迁移和离线子命令实际使用的依赖 |
| `services/trade-log/bins/trade-log-query/src/bootstrap.rs` | 装配数据库、schema 检查和 migrate/reparse/import-evidence 子命令 |
| `services/trade-log/bins/trade-log-query/src/config.rs` | 新增数据库凭证文件、连接池和时限校验 |
| `services/trade-log/bins/trade-log-query/tests/startup.rs` | 新增数据库配置、迁移前启动失败及 CLI 验证 |
| `services/trade-log/adapters/tests/query_chain.rs` | 保留 v0.1 行为测试并验证完整保存及持久化回执 |
| `tests/integration-tests/query_api.rs` | 新增 source 分支、时间/cursor 参数和旧接口兼容验证 |

#### 计划直接复用的文件

`crates/account-facts/src/lib.rs` 及其稳定身份测试、`crates/shared-types/src/lib.rs`、`crates/protocol-api/src/lib.rs`、`crates/service-runtime/src/` 的既有文件，以及 `protocols/implementations/hyperliquid/src/lib.rs`、`parser.rs`、`source.rs` 均复用现有实现。本版存储功能不要求重新开发来源读取、事实身份、金额处理或协议解析。

`.gitignore` 和 `.dockerignore` 的既有 secrets/var 排除规则直接沿用，无须因数据库挂载新增版本专用规则。

网关的 `main.rs`、基础健康/版本 handler、公共响应和 trace 中间件不增加存储逻辑；`scripts/init-v0.1.sh` 保留，新版初始化另由新增的 `scripts/init-v0.2.sh` 提供。`rust-toolchain.toml` 继续固定已有工具链，部署构建减少检查组件下载由 Dockerfile 处理。

## 8. 配置与异常处理

### 8.1 计划新增配置

trade-log.toml 保留原配置，新增必填数据库配置；query-api 的内部 base_url/凭证仍沿用：

```toml
[database]
url_file = "/run/secrets/trade-log-database-url"
migration_url_file = "/run/secrets/trade-log-migration-url"
max_connections = 8
connect_timeout_seconds = 5
statement_timeout_seconds = 10
```

url_file 是运行角色连接，migration_url_file 仅 migrate 子命令必需，常驻应用不挂载该文件；服务启动不解析无须使用的迁移凭证。v0.2 默认 query.timeout_seconds 调整为 40，网关 trade_log.request_timeout_seconds 调整为 45，覆盖新增保存事务；连接字符串从文件加载，不把明文密码写入 TOML、日志、API 或镜像。事务操作受总时限约束，statement_timeout 在数据库会话设置。migration 不随每个 HTTP 请求执行；启动确认 schema 版本匹配，未迁移、权限不足或连接失败明确退出。

CLI 使用 `--config /etc/robotech/trade-log.toml` 默认路径，兼容现有无子命令启动方式；新增：

| 子命令 | 行为 | 退出码 |
| --- | --- | --- |
| migrate | 使用独立迁移凭证文件执行服务 migration，可重复运行 | 0 成功；2 配置错误；3 数据库/迁移错误 |
| reparse --query-id | 只读核对保存的来源和任务事实，不发网络请求 | 0 SAME；4 DIFFERENT/INCOMPLETE；2 参数错误；3 存储错误 |
| import-evidence --directory | 导入本地证据目录，可重复执行 | 0 全部成功或已存在；4 存在冲突/坏证据；2 参数错误；3 存储错误 |

重新解析输出 query_id、旧/新 parser_version、source_records、fact_records、comparison、差异 fact_id；导入汇总 imported/existing/failed。不输出数据库 URL 或凭证。

### 8.2 HTTP 和恢复规则

| 场景 | 响应及处理 |
| --- | --- |
| source/时间/cursor 参数错误或不适用组合 | 400 VALIDATION_ERROR |
| stored 无匹配事实 | 200 空 trades；matched_records=0，不访问来源 |
| 同身份不同内容、旧目录导入冲突 | 409 VERSION_CONFLICT；事实事务回滚 |
| live 来源失败/超时 | 沿用 v0.1 429/422/503 及原始证据保存 |
| 数据库不可达、连接池/查询等待超时 | 503 DEPENDENCY_UNAVAILABLE；不降级为“成功但未保存” |
| 原始证据损坏、不满足内部约束 | 500 INTERNAL_INVARIANT_VIOLATION；离线核对则非零退出 |
| 文件诊断副本失败，数据库已完整提交 | live 可成功，结构化日志记录镜像失败；不隐瞒数据库存储错误 |

health 仍只表示进程存活，数据库故障期间可响应；stored 查询失败不能冒充空库存。数据库恢复后新请求重新尝试，无后台无限重试。配置加载仍按既有显式环境覆盖约定，业务新增项只在 TOML/凭证文件维护。

## 9. 验证与验收

### 9.1 自动验证

保留 v0.0/v0.1 测试，增加实际 PostgreSQL 集成测试：首次插入、同内容重复、并发重复、内容冲突回滚、原始响应错误归档、事务失败、提交后响应丢失后的重复恢复、库故障、角色权限、迁移重复执行和旧证据重复导入。

查询测试覆盖开始包含/结束排除、偏移时区、毫秒精度、同时间大整数 tid 排序、空结果、游标条件错误、分页无重复遗漏、分页过程中插入更早成交、固定提交水位及连接超时。重新解析从已保存原始响应开始，验证不请求来源且不改写版本。

### 9.2 实际地址核对

至少一次真实账户 live 查询后：核对返回展示事实与对应原始 fill；完整 perpetual_records 与库中本次任务的事实关联数一致；inserted+existing 与完整事实数一致。同一记录多次查询不增加事实版本，但保留新的原始响应和 observation。

实时来源会变化，不能要求第二次插入数必为 0。重复采集不重复入账的确定性检查使用相同保存证据导入两次；再以真实重复查询验证共同 tid 的 fact_id 和 revision 保持不变。

### 9.3 完成标准

时间查询、稳定身份、原始字节与事实关联、事务及冲突、导入/重新解析、数据库重启保留和故障恢复全部验证通过。服务器完成第 10.4 节手动步骤后再编写实际验收报告；本文不代表功能已实现或验收通过。

## 10. 交付、部署与手动验证

### 10.1 Docker 部署和挂载

仍由 Docker Compose 启动：query-api、trade-log-query、postgres；evidence-init、trade-log-migrate 为一次性 init profile 服务，非常驻应用。仅网关发布 8080，内部查询 8081 和数据库 5432 不发布宿主机端口。

默认保持 Docker 命名卷，不擅自改为服务器目录绑定：

| 数据/文件 | 宿主机来源 | 容器目标 |
| --- | --- | --- |
| 网关/查询配置 | ./config 下 TOML，只读 bind | /etc/robotech 下对应文件 |
| 服务凭证 | ./secrets/trade-log-token，只读 bind | /run/secrets/trade-log-token |
| 应用数据库 URL | ./secrets/trade-log-database-url，只读 bind，仅查询服务 | /run/secrets/trade-log-database-url |
| 迁移数据库 URL | ./secrets/trade-log-migration-url，只读 bind，仅迁移服务 | /run/secrets/trade-log-migration-url |
| PostgreSQL 管理密码 | ./secrets/postgres-password，只读 bind，仅数据库/初始化服务 | /run/secrets/postgres-password |
| 原有文件证据 | trade-log-evidence 命名卷，保留 | /var/lib/robotech/trade-log |
| 数据库数据 | trade-log-postgres 命名卷 | /var/lib/postgresql |

计划采用 PostgreSQL 18 主版本，并配置 PGDATA=/var/lib/postgresql/18/docker（参照 [官方镜像目录约定](https://github.com/docker-library/docs/blob/master/postgres/README.md)）；发布时固定经过验证的小版本镜像及摘要，文档记录实际值。PostgreSQL 数据卷的权限按数据库镜像要求处理，不能用应用 UID 10001 覆盖数据库目录所有权。

数据库初始化明确创建 robotech_admin 管理角色、trade_log_migrator 迁移角色及 trade_log_app 运行角色。init-v0.2.sh 保留原服务 token，初始化数据库密码/角色及 URL 文件，不打印凭证；启动数据库、等待可用，再执行 migration，并为原证据卷保持 UID/GID 10001 写权限。既有数据库角色和密码重复初始化不更改；角色创建和迁移的权限分离。应用镜像仍使用共用多阶段 Dockerfile 的两个 target。

### 10.2 升级、备份和回退

升级前备份旧文件证据和数据库（若已存在）。通过 pg_dump 逻辑备份，不在数据库运行中直接复制数据目录作为一致备份。secrets 不入 Git 和构建上下文；不公开实际服务器地址。

先部署数据库和迁移，再更新 trade-log-query，最后更新网关。导入旧证据是明确操作，不自动把未知目录判为成功；失败目录保留供核对。本版 migration 为增量建表和索引，不删除旧证据或既有数据。

回退程序可恢复 v0.1 的实时文件查询，但不能读取新增数据库库存；保留数据库和卷，不自动执行反向删表。v0.2 运行配置要求数据库，缺失时必须报错；旧网关基础配置仍可运行 health/version。

构建交付应增加 Cargo registry/git/编译缓存，并避免生产构建安装检查工具组件，减少源码变化导致重复下载；必须保持 cargo --locked 和确定的工具链版本，不能用更换未验证依赖绕开网络问题。

### 10.3 交付清单

交付数据库 migration、PostgreSQL 适配器、幂等保存和库存查询端口、离线 CLI、网关 source 分支、配置和凭证初始化、Compose 服务及卷、构建缓存改进、测试、统一接口文档和实际验收报告。README 不新增版本入口。

### 10.4 手动验证步骤

**以下命令在 v0.2 实现完成后执行，当前版本尚不具备这些子命令或参数。** 在服务器仓库根目录，逐步执行并观察结果；不要求使用自动验收脚本。

#### 第 1 步：初始化、构建和检查容器

```sh
sh scripts/init-v0.2.sh
docker compose up --build -d
docker compose ps -a
docker compose logs --tail=40 query-api trade-log-query postgres
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/version
```

**看什么：** 三个常驻容器 running，两个应用出现 server_started，version 为 0.2.0；只有网关发布端口，迁移无错误。

**通过标准：** 数据库初始化与迁移成功，程序正确启动。命名卷没有删除或重新替换旧证据。

#### 第 2 步：采集真实账户并保存

```sh
curl --noproxy '*' --max-time 60 -i 'http://127.0.0.1:8080/api/v1/trade-events?account=0x010461c14e146ac35fe42271bdc1134ee31c703a&limit=10'
```

**看什么：** 200，persistence.status=COMMITTED；inserted_records+existing_records=counts.perpetual_records，展示不超过 10 条。记下 query_id、任一 fact_id/occurred_at 和完整成交数。

**通过标准：** 有真实合约成交，完整集合保存成功，不只保存展示的 10 条。无近期成交时更换实际有成交的公开账户。

#### 第 3 步：查询已保存的同一账户

```sh
curl --noproxy '*' -i 'http://127.0.0.1:8080/api/v1/trade-events?account=0x010461c14e146ac35fe42271bdc1134ee31c703a&source=stored&limit=10'
```

**看什么：** 200；coverage=STORED_RECORDS_ONLY，有 matched_records 和分页字段。第 2 步的事实可在库存中找到，fact_id、金额、时间和 revision=1 一致。已有库存较多时按成交时间定位或翻页。

**通过标准：** 查询来自本库，没有再次生成来源采集任务，没有声称完整历史。

#### 第 4 步：验证时间范围边界

把 T 替换为第 2 步记下的真实 occurred_at，T_NEXT 替换为 T 加 1 毫秒后的 UTC 时间：

```sh
curl --noproxy '*' -i --get 'http://127.0.0.1:8080/api/v1/trade-events' \
  --data-urlencode 'account=0x010461c14e146ac35fe42271bdc1134ee31c703a' \
  --data-urlencode 'source=stored' \
  --data-urlencode 'start_time=<T>' \
  --data-urlencode 'end_time=<T_NEXT>'
```

随后保持 start_time 省略，把 end_time 改为 T 再请求。

**看什么：** `[T,T_NEXT)` 包含该毫秒事实；结束为 T 的请求不包含 occurred_at=T 的事实；其他同毫秒成交允许同时出现。

**通过标准：** 开始包含、结束排除，返回所有记录都在请求范围内。

#### 第 5 步：验证分页

将第 3 步 limit 改为 1。库存多于一条时记下 next_cursor：

```sh
curl --noproxy '*' -i --get 'http://127.0.0.1:8080/api/v1/trade-events' \
  --data-urlencode 'account=0x010461c14e146ac35fe42271bdc1134ee31c703a' \
  --data-urlencode 'source=stored' \
  --data-urlencode 'limit=1' \
  --data-urlencode 'cursor=<上一页实际next_cursor>'
```

**看什么：** 下一页不重复上一页 fact_id，snapshot_seq 与第一页相同，顺序一致；最后一页 has_more=false、next_cursor=null。

**通过标准：** 分页没有重复；更换账户但沿用 cursor 应返回 400。

#### 第 6 步：查看数据库记录及重新解析

将 query_id 替换为第 2 步实际值：

```sh
docker compose exec -T postgres psql -U robotech_admin -d robotech -c "SELECT query_id,status FROM trade_log.collection_jobs ORDER BY created_at DESC LIMIT 5;"
docker compose exec -T postgres psql -U robotech_admin -d robotech -c "SELECT count(*) AS facts FROM trade_log.account_facts_current;"
docker compose exec -T postgres psql -U robotech_admin -d robotech -c "SELECT count(*) AS raw_responses FROM trade_log.raw_logs;"
docker compose exec -T trade-log-query trade-log-query reparse --query-id query_实际值
```

**看什么：** 采集任务 COMPLETED，原始响应和事实均有记录；重新解析输出 comparison=SAME、退出码 0。这些 SQL 是本版默认本机数据库管理员排查命令，不把管理凭证传给网关。

**通过标准：** 从已保存来源能复算同样的事实，数据库事实版本没有因为 reparse 增加。

#### 第 7 步：确定性验证重复导入

使用一个已有 v0.1 COMPLETED 证据目录；没有旧证据时，用发布附带的固定导入样本挂载到证据目录（样本明确标为构造数据），不要拿两次变化的实时窗口作为固定输入。

```sh
docker compose exec -T trade-log-query trade-log-query import-evidence --directory /var/lib/robotech/trade-log
docker compose exec -T postgres psql -U robotech_admin -d robotech -c "SELECT count(*) AS facts FROM trade_log.account_facts_current;"
docker compose exec -T trade-log-query trade-log-query import-evidence --directory /var/lib/robotech/trade-log
docker compose exec -T postgres psql -U robotech_admin -d robotech -c "SELECT count(*) AS facts FROM trade_log.account_facts_current;"
```

**看什么：** 第一次成功导入或已存在；第二次全部为已存在，事实总数不增加。同 query_id 不会重复建立任务。

**通过标准：** 相同输入重复执行不重复入账，来源仍可追溯。坏目录必须列出并非零退出，不能记作已通过。

#### 第 8 步：数据库重启后仍可查询

```sh
docker compose restart postgres trade-log-query
```

等待服务重新就绪，再执行第 3 步的 stored 请求，并读取第 6 步同一个任务。

**通过标准：** 原事实和原始响应仍在，query_id/fact_id 一致，stored 返回 200。不要用 down -v 验证普通重启。

#### 第 9 步：数据库故障与恢复

在允许短暂停止查询的验收环境执行：

```sh
docker compose stop postgres
curl --noproxy '*' --max-time 60 -i 'http://127.0.0.1:8080/api/v1/trade-events?account=0x010461c14e146ac35fe42271bdc1134ee31c703a&source=stored&limit=10'
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/health
docker compose start postgres
```

**看什么：** stored 返回 503，不返回成功空数组；网关 health 仍为 200。恢复就绪后重复 stored 查询成功。

**通过标准：** 故障明确，恢复后原记录仍能读取，未丢失数据。

#### 第 10 步：错误参数、旧接口及停止

```sh
curl --noproxy '*' -i 'http://127.0.0.1:8080/api/v1/trade-events?account=invalid&source=stored'
curl --noproxy '*' -i 'http://127.0.0.1:8080/api/v1/trade-events?account=0x010461c14e146ac35fe42271bdc1134ee31c703a&source=unknown'
curl --noproxy '*' -i 'http://127.0.0.1:8080/api/v1/trade-events?account=0x010461c14e146ac35fe42271bdc1134ee31c703a&source=stored&start_time=invalid'
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/health
docker compose stop
docker compose ps -a
docker compose logs --tail=30 query-api trade-log-query
docker compose start
```

**通过标准：** 三个非法请求均 400；旧基础接口正常；两个应用正常退出、日志含 shutdown_completed；重启后基础接口和 stored 查询恢复。另手动验证 start_time=end_time、limit=0 和损坏 cursor 都返回 400。

按以上各步记录实际结果、日期、执行人、代码版本及失败日志。其他冲突、并发和精度边界由固定样本集成测试验证，不能仅靠一次公网查询宣布全部通过。
