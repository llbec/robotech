# 接口说明

本文是项目统一的接口说明文档，集中维护所有对外接口和内部接口。后续版本在本文件中更新接口定义，并注明新增、变更或废弃的适用版本，不另建按开发版本命名的接口文档。

当前对应程序版本：0.5.0；JSON schema_version：1。接口定义依据当前代码编写。

## 1. 地址和接口列表

默认网关地址：服务器本机 `http://127.0.0.1:8080`；公网地址示例 `http://server.example.com:8080`。`server.example.com` 是占位域名，使用时替换为实际服务器地址。公网是否可访问取决于实际部署及防火墙配置。

| 范围 | 方法 | 路径 | 用途 |
| --- | --- | --- | --- |
| 对外 | GET | [/api/v1/health](#41-网关健康接口) | 判断网关进程是否存活 |
| 对外 | GET | [/api/v1/version](#42-版本接口) | 查询网关程序版本 |
| 对外 | GET | [/api/v1/trade-events](#43-账户近期合约成交查询) | 实时查询并保存成交；[查询已保存成交](#44-已保存成交时间查询) |
| 对外 | GET | [/api/v1/watch-accounts](#45-自动更新状态) | 查看配置账户的自动采集状态、水位和错误 |
| 对外 | GET | [/api/v1/publishing-status](#46-webhook-发布状态) | 查看候选输出开关、队列、心跳和错误 |
| 输出 | POST | [配置的 webhook 地址](#47-webhook-候选事件) | 发送标准账户事实候选，接收方按 event_id 去重 |
| 内部 | GET | [/internal/v1/publishing-status](#36-内部发布状态) | 网关读取 publisher 状态，需要服务凭证 |
| 内部 | GET | [/internal/v1/collection-status](#35-内部自动采集状态) | 网关读取 collector 状态，需要服务凭证 |
| 内部 | GET | [/internal/v1/health](#33-内部健康接口) | 判断交易查询进程是否存活，需要服务凭证 |
| 内部 | POST | [/internal/v1/trade-queries](#32-内部成交查询) | 网关调用交易查询服务，需要服务凭证 |
| 内部 | POST | [/internal/v1/stored-trade-queries](#34-内部库存查询) | 查询数据库保存的成交，需要服务凭证 |

对外接口当前不要求客户端提供身份凭证。交易查询服务默认监听 8081，collector 默认监听 8082，publisher 默认监听 8083；Compose 不发布这些内部端口，网关分别通过 `http://trade-log-query:8081` 和 `http://trade-collector:8082` 调用。

对外 GET 接口支持 HEAD：响应不包含 JSON 正文。未定义路径返回 404；已定义路径使用不支持的方法返回 405。

## 2. 公共数据格式

### 2.1 成功响应

```json
{
  "data": {},
  "meta": {
    "trace_id": "trace_0123456789abcdef0123456789abcdef",
    "schema_version": 1
  }
}
```

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| data | object | 对应接口的业务数据，具体字段见后文 |
| meta.trace_id | string | 本次 HTTP 请求的追踪 ID，用于关联网关、内部服务日志及查询证据 |
| meta.schema_version | integer | 响应结构版本，当前为 1；与程序版本 0.5.0 分开管理 |

JSON 响应的 Content-Type 为 `application/json`。响应头 `x-trace-id` 与响应体中的 trace_id 相同。网关为每次外部请求生成新 trace_id，不沿用客户端传入的值。

### 2.2 错误响应

错误响应没有 data/meta 包装，三个字段位于顶层：

```json
{
  "code": "VALIDATION_ERROR",
  "message": "account must be 0x followed by 40 hexadecimal digits",
  "trace_id": "trace_0123456789abcdef0123456789abcdef"
}
```

| 字段 | 类型 | 含义 |
| --- | --- | --- |
| code | string | 稳定的错误分类，客户端按 HTTP 状态和此字段判断错误 |
| message | string | 本次错误的说明，不应依赖具体文案进行程序判断 |
| trace_id | string | 对应本次请求的追踪 ID |

时间字符串均为 UTC RFC 3339 格式，如 `2026-10-05T00:00:00.000Z`。价格、数量、金额和来源整数 ID 使用字符串传输，避免浮点数及大整数精度丢失。可选字段没有值时返回 JSON null，不等于 0 或 false。

## 3. 内部接口

内部接口用于项目内不同进程之间的通信。目前 query-api 负责对外 HTTP 请求，trade-log-query 负责手动成交查询和库存读取，trade-collector 独立执行单地址自动采集。通过内部接口明确这些进程的请求、响应、认证和错误边界，使业务查询过程可以独立运行和迭代，网关无需依赖交易查询进程的内部实现。

当客户端调用外部成交查询接口时，网关完成参数校验，再调用内部成交查询接口；交易查询服务完成采集和证据保存后，网关将结果返回客户端。source=stored 时网关改调用内部库存查询接口，仅读取本库，不进行来源采集。内部健康接口供运维或受信任的服务调用方检查交易查询进程是否存活，当前网关不会在每次成交查询前主动调用它。

内部接口位于容器网络，默认不向宿主机发布端口。所有内部接口均要求服务凭证，普通业务客户端通过第 4 章的外部接口访问。

### 3.1 通用请求头

| 请求头 | 必填 | 意义 |
| --- | --- | --- |
| Authorization | 是 | 格式 `Bearer <服务凭证>`，凭证由部署初始化生成，网关和交易查询服务读取同一文件 |
| Content-Type | POST 时是 | application/json |
| x-trace-id | 否 | 已认证调用方可传入 trace_ 加 32 位十六进制字符；符合格式时沿用，否则生成新 ID |

不要把服务凭证放入公网调用示例或业务 query 参数。凭证缺失或错误返回 401、code=AUTHENTICATION_REQUIRED；内部 health 同样需要认证。

### 3.2 内部成交查询

```http
POST /internal/v1/trade-queries
Authorization: Bearer <服务凭证>
Content-Type: application/json
```

```json
{
  "account": "0x010461c14e146ac35fe42271bdc1134ee31c703a",
  "limit": 10
}
```

account 必填、limit 可省略并默认 100；含义和校验规则与[第 4.3.1 节](#431-请求格式和参数)相同。limit 在 JSON 中必须为整数，不能写成字符串；未知字段被拒绝。请求体上限为 16 KiB，无效 JSON、错误 Content-Type 或超限请求由查询处理器返回 400 VALIDATION_ERROR。

成功响应与[外部成交查询](#432-成功响应结构)相同，错误格式与第 2.2 节相同。默认最多并行执行 4 个查询，超过立即返回 429；默认整体查询超时 40 秒，超时返回 503。并发和超时以 trade-log.toml 的实际配置为准，网关内部请求默认超时为 45 秒。

### 3.3 内部健康接口

```http
GET /internal/v1/health
```

没有业务参数，需要 Authorization 请求头。成功返回 200，data 格式：

```json
{"status":"ok","service":"trade-log-query","version":"0.5.0"}
```

status 表示内部进程可响应，service 为服务名，version 为程序构建版本；完整响应使用 data/meta 包装。该接口不主动检查上游来源。

### 3.4 内部库存查询

网关在外部请求 source=stored 时调用，只读数据库，不调用 Hyperliquid。内部认证、请求头、JSON 体上限与第 3.1、3.2 节相同。

```http
POST /internal/v1/stored-trade-queries
Authorization: Bearer <服务凭证>
Content-Type: application/json
```

```json
{
  "account": "0x010461c14e146ac35fe42271bdc1134ee31c703a",
  "limit": 10,
  "start_time": "2026-10-04T00:00:00.000Z",
  "end_time": "2026-10-05T00:00:00.000Z"
}
```

account 必填；limit 为整数，可省略但不可为 null。start_time、end_time、cursor 可省略或为 null，表示未提供。参数含义和返回数据见第 4.4 节；未知字段返回 400。数据库不可用或操作失败返回 503，不返回空库存冒充成功。

### 3.5 内部自动采集状态

新增于 v0.3。`GET /internal/v1/collection-status`，服务地址 `http://trade-collector:8082`。网关收到 watch-accounts 请求时调用；不会触发立即采集，不直接查询业务数据库。

请求没有正文或参数，必须携带内部 Bearer 服务凭证。成功 data 与第 4.5 节相同，meta 含内部请求 trace_id 和 schema_version=1。数据库读取默认时限 3 秒，失败或超时返回 503；未授权返回 401。能读取检查点而来源采集失败时返回 200，在业务状态中说明失败。

collector 的 `GET /internal/v1/health` 使用相同凭证要求，data 为 `{"status":"ok","service":"trade-collector","version":"0.5.0"}`，只表示进程可响应，不表示来源访问或自动采集成功。

### 3.6 内部发布状态

v0.5 新增。网关向 `http://trade-parser-publisher:8083/internal/v1/publishing-status` 发送 GET，使用第 3.1 节的服务凭证及 trace 头，无查询参数。返回第 4.6 节的 data/meta 信封。发布器也提供认证的 `/internal/v1/health`，仅判断进程存活；数据库不可读时状态接口返回 503。

## 4. 外部接口

外部接口是面向业务客户端、运维人员和调用方的统一 HTTP 入口，由 query-api 提供。定义外部接口是为了提供稳定的调用约定，使调用方通过账户等业务参数获取结果，无需了解内部进程地址、服务凭证或来源采集细节。

部署完成后，运维人员可调用健康接口确认网关可响应，调用版本接口确认运行版本；业务客户端在需要查看账户近期合约成交时调用成交查询接口。网关负责将成交查询转交内部服务，并统一返回结果、错误和请求追踪 ID。

外部接口当前不要求身份凭证。健康接口只反映网关进程存活；确认交易查询链路是否可用，应调用成交查询接口。以下分节说明各接口。

### 4.1 网关健康接口

```http
GET /api/v1/health
```

没有业务请求参数。

```sh
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/health
```

成功返回 200，示例：

```json
{
  "data": {
    "status": "ok",
    "service": "query-api",
    "started_at": "2026-10-05T00:00:00.000Z"
  },
  "meta": {
    "trace_id": "trace_0123456789abcdef0123456789abcdef",
    "schema_version": 1
  }
}
```

| data 字段 | 类型 | 含义 |
| --- | --- | --- |
| status | string | 当前为 ok，表示网关可以响应请求 |
| service | string | 服务名，当前为 query-api |
| started_at | string | 本次网关进程启动时间；进程重启后改变 |

health 不检查交易查询服务、Hyperliquid 或证据存储状态。验证业务链路需要调用成交查询接口。

### 4.2 版本接口

```http
GET /api/v1/version
```

没有业务请求参数。

```sh
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/version
```

成功返回 200；data 格式：

```json
{"service":"query-api","version":"0.5.0"}
```

| data 字段 | 类型 | 含义 |
| --- | --- | --- |
| service | string | 服务名 query-api |
| version | string | 当前运行程序的构建版本，v0.5 为 0.5.0 |

实际 HTTP 响应仍使用第 2 节的 data/meta 包装。

### 4.3 账户近期合约成交查询

#### 4.3.1 请求格式和参数

```http
GET /api/v1/trade-events?account=<账户地址>&limit=100
```

| 参数 | 位置 | 类型 | 必填 | 默认值 | 意义及限制 |
| --- | --- | --- | --- | --- | --- |
| account | URL query | string | 是 | 无 | 被查询的账户；必须以小写 `0x` 开头，后跟 40 个十六进制字符，总长度 42；地址字符可大小写混合，结果规范化为小写；不接受前后空格 |
| limit | URL query | integer | 否 | 100 | 最多展示多少条合约成交，范围 1–2000；不接受小数、负数或空字符串 |
| source | URL query | string | 否 | live | live 调用来源并保存；stored 查询数据库，详见第 4.4 节 |

source=live（或省略 source）使用本节的实时响应；不允许 start_time、end_time 或 cursor。

请求不带 JSON body。缺少 account、未知参数、重复的同名参数、非法类型或超出范围均返回 400。network 从服务配置读取，不能通过请求切换。实时查询不提供起止时间和游标；已保存记录的时间查询见第 4.4 节。不提供市场筛选或自定义排序参数。

示例账户为公开账户，实际近期成交会变化：

```sh
curl --noproxy '*' --max-time 45 -i --get 'http://127.0.0.1:8080/api/v1/trade-events' \
  --data-urlencode 'account=0x010461c14e146ac35fe42271bdc1134ee31c703a' \
  --data-urlencode 'limit=10'
```

省略 limit 时展示上限为 100。limit 只限制最终 HTTP 展示条数，不控制来源读取条数；完整标准结果仍会写入证据目录。

#### 4.3.2 成功响应结构

下面是无来源记录时的格式示例，ID 和时间仅用于展示：

```json
{
  "data": {
    "query_id": "query_0123456789abcdef0123456789abcdef",
    "account": "0x010461c14e146ac35fe42271bdc1134ee31c703a",
    "network": "mainnet",
    "queried_at": "2026-10-05T00:00:00.000Z",
    "query_scope": "RECENT_SOURCE_WINDOW",
    "coverage": "SOURCE_WINDOW_UNVERIFIED",
    "counts": {
      "source_records": 0,
      "duplicate_records": 0,
      "spot_records": 0,
      "unsupported_records": 0,
      "invalid_records": 0,
      "perpetual_records": 0,
      "returned_records": 0
    },
    "display_truncated": false,
    "observed_range": null,
    "trades": [],
    "warnings": ["NO_SOURCE_RECORDS"],
    "evidence_ref": "query_0123456789abcdef0123456789abcdef",
    "persistence": {"status":"COMMITTED","inserted_records":0,"existing_records":0}
  },
  "meta": {
    "trace_id": "trace_0123456789abcdef0123456789abcdef",
    "schema_version": 1
  }
}
```

| data 字段 | 类型 | 意义 |
| --- | --- | --- |
| query_id | string | 本次业务查询 ID，每次查询独立生成；用于定位证据 |
| account | string | 规范化为小写的查询账户 |
| network | string | 当前来源网络，mainnet 或 testnet，由服务配置决定 |
| queried_at | string | 本次来源采集开始时间，不是 HTTP 完成时间或成交发生时间 |
| query_scope | string | 当前固定为 RECENT_SOURCE_WINDOW，即官方来源返回的近期窗口 |
| coverage | string | 来源覆盖评价，取值见[第 4.3.6 节](#436-覆盖与警告)；不表示完整账户历史 |
| counts | object | 来源分类及标准成交数量，见下表 |
| display_truncated | boolean | 完整标准成交数是否超过请求 limit；true 表示 HTTP 只展示了前 limit 条 |
| observed_range | object/null | 完整标准成交集合中实际观察到的时间范围；无标准成交为 null |
| observed_range.first_at | string | 完整标准成交中的最早成交时间 |
| observed_range.last_at | string | 完整标准成交中的最晚成交时间 |
| trades | array | 展示的标准成交事实，结构见第 4.3.4、4.3.5 节 |
| warnings | array of string | 覆盖限制和解析诊断，可为空数组 |
| evidence_ref | string | 当前等于 query_id，是数据库采集任务及诊断证据标识；不是下载 URL |
| persistence | object | v0.2 保存回执，只有数据库提交成功才返回实时成功响应 |
| persistence.status | string | 固定 COMMITTED |
| persistence.inserted_records | integer | 本次完整标准集合中新入库的事实数，不受展示 limit 影响 |
| persistence.existing_records | integer | 本次完整标准集合中已经存在且内容相同的事实数；与 inserted 之和为 perpetual_records，不是原始数组内 duplicate_records |

按 occurred_at 从新到旧排序；时间相同时按来源 tid 数值降序，再按 fact_id 字典升序。observed_range 根据完整标准成交集合计算，可能比展示的 10 条覆盖更长，不代表已覆盖该区间的所有历史记录。

#### 4.3.3 counts：每种计数的意义

| 字段 | 类型 | 意义 |
| --- | --- | --- |
| source_records | integer | userFills 原始数组的元素总数，过滤和去重之前 |
| duplicate_records | integer | 同一市场和 tid 下，内容完全相同的重复元素数量 |
| spot_records | integer | 根据来源元数据识别为现货的记录数，不返回到 trades |
| unsupported_records | integer | 不在当前默认永续市场元数据中的未知或未支持市场记录数 |
| invalid_records | integer | 身份、数值、时间或其他必要字段无法可靠解析的记录数 |
| perpetual_records | integer | 去重后成功解析的全部默认永续成交数量，尚未应用展示 limit |
| returned_records | integer | 本次 HTTP 展示的成交数量，等于 trades.length |

分类计数满足：source_records = duplicate_records + spot_records + unsupported_records + invalid_records + perpetual_records。相同身份但内容不同会返回 409，不作为普通重复吞掉。

#### 4.3.4 trades[]：标准成交事实

| 字段 | 类型 | 意义及当前值 |
| --- | --- | --- |
| fact_id | string | 稳定事实身份，格式 hl_fill_v1_<SHA-256>；根据身份版本、network、protocol、规范化 account、原始 coin 和 tid 计算，不随 query_id 改变 |
| fact_type | string | 事实类型，当前 TRADE |
| revision | integer | 事实修订版本，当前 1；不是协议或 HTTP 版本 |
| change_type | string | 当前 UPSERT，表示以 fact_id 标识的事实写入或更新；接口本身不提供增量订阅 |
| confirmation_status | string | 当前 OBSERVED，表示从官方 API 观察到，不代表额外完成链上最终性确认 |
| chain_id | string | 来源网络标识，如 hyperliquid:mainnet |
| protocol | string | 当前 hyperliquid |
| account | string | 此成交所属的规范化账户 |
| account_key | string | 网络、协议和账户组成的复合键，如 hyperliquid:mainnet:hyperliquid:0x… |
| ordering_key | string | 20 位补零的来源毫秒时间戳与 20 位补零 tid，以冒号连接；供来源顺序定位 |
| sub_index | integer | 同一事实内的子事件序号，当前 0 |
| source | string | 当前 OFFICIAL_API |
| source_ref | string | 来源成交 tid，以十进制字符串保存 |
| raw_log_id | string | 该事实引用的 userFills 原始响应 ID，可定位对应 .body 和元数据文件 |
| occurred_at | string | 来源 time 转换后的 UTC 成交时间 |
| payload | object | 成交业务字段，见下节 |

同一成交跨查询 fact_id 保持一致；raw_log_id 属于本次采集，可能改变。不要把 query_id 或 raw_log_id 当作成交去重身份。

#### 4.3.5 payload：成交业务字段

| 字段 | 类型 | 意义及当前规则 |
| --- | --- | --- |
| market | string | 标准市场名，如 hyperliquid:BTC-USDC |
| instrument_type | string | 当前 PERPETUAL，永续合约 |
| base_asset | string | 来源 coin 对应的标的，如 BTC |
| quote_asset | string | 当前 USDC |
| action | string | 当前 TRADE |
| trigger_type | string | USER：来源 dir 被识别为普通开平仓；LIQUIDATION：强制清算相关；PROTOCOL：未能确认普通用户成交的其他来源方向 |
| copy_eligible | boolean | 普通已识别成交为 true，强制清算或未确认方向为 false；只是分类标志，不代表可以直接发出复制交易指令 |
| side | string | BUY 或 SELL；分别来自来源 side=B 或 A，不等同于持仓方向 |
| position_effect | string | OPEN、CLOSE、REVERSE、INCREASE、DECREASE，按成交前仓位和本次成交数量计算，见下表 |
| order_id | string/null | 来源 oid，以字符串保存；无有效整数 oid 时为 null |
| operation_id | string/null | 上层操作关联 ID，当前 null |
| price | string | 成交价格，来源 px 的精确十进制值 |
| quantity | string | 成交数量，来源 sz 的精确十进制值，必须大于 0 |
| notional | string | 精确 price × quantity，不使用二进制浮点数计算 |
| fee | string | 来源 fee；保留正负号，负值可能表示返佣，不能取绝对值 |
| fee_asset | string/null | 来源 feeToken；缺失或空值为 null，不根据 quote_asset 推定 |
| reported_realized_pnl | string/null | 来源 closedPnl 的报告值；并非本系统自行核算的最终收益 |
| reported_pnl_asset | string/null | 报告收益的资产单位，当前 null，表示未确认 |
| reported_pnl_includes_fee | boolean/null | 报告收益是否包含手续费，当前 null，表示未确认 |
| reported_pnl_includes_funding | boolean/null | 报告收益是否包含资金费，当前 null，表示未确认 |
| transaction_hash | string/null | 来源 hash 字符串；来源缺失则 null |
| extension | object | 来源辅助信息及原始记录，见下表 |

| position_effect | 意义 |
| --- | --- |
| OPEN | 成交前仓位为 0，本次建立仓位 |
| CLOSE | 成交后仓位恰好为 0 |
| REVERSE | 成交后仓位正负号与成交前相反，发生反转 |
| INCREASE | 仓位方向保持一致，绝对数量增加 |
| DECREASE | 仓位方向保持一致，绝对数量减少 |

| extension 字段 | 类型 | 意义 |
| --- | --- | --- |
| source_fill | object | 对应的原始 userFills 元素，供核对；来源字段可能随来源扩展，不作为本系统稳定业务字段契约 |
| source_indices | array of integer | 在原始 userFills 数组中的位置，从 0 开始；完全相同的重复记录会追加对应位置 |
| time_ms | integer | 来源成交时间，Unix 毫秒时间戳 |
| source_dir | string | 来源 dir 的原始值，如 Open Long |
| trigger_unconfirmed | boolean | true 表示既未识别为清算，也未识别为普通开平仓方向；此时 copy_eligible 为 false |

金额字段允许保留不同字符串表现形式，但应按精确十进制比较；不要用 JavaScript Number 或其他二进制浮点数进行财务计算。解析器不接受超出支持精度的输入或会产生精度损失的成交额。

#### 4.3.6 覆盖与警告

| coverage | 意义 |
| --- | --- |
| SOURCE_WINDOW_UNVERIFIED | 未发现明确的数量上限、无效记录或未支持市场限制，但仍没有证明来源窗口完整 |
| LIMITED | 来源元素达到 2000，或出现无效记录、未支持市场；查询覆盖存在明确限制 |

| warnings 值或格式 | 意义 |
| --- | --- |
| SOURCE_RECORD_LIMIT | 来源元素达到 2000，可能触及来源返回上限 |
| NO_SOURCE_RECORDS | 来源返回空数组；只表示本次没有来源记录 |
| NO_PERPETUAL_RECORDS | 没有标准永续成交，且来源包含被过滤的现货记录 |
| UNSUPPORTED_OR_UNKNOWN_MARKET:<index> | 原始数组该位置所属市场未知或当前未支持 |
| INVALID_RECORD:<index>:<reason> | 原始数组该位置无法可靠解析；reason 是诊断文本，不建议按具体文案分支 |

当前只返回默认永续 DEX 中、根据元数据确认的成交；现货被过滤，其他市场被披露为未支持。display_truncated 与 coverage 是不同维度：展示 limit 导致截断，不会单独把 coverage 改为 LIMITED。

成功的空 trades 不表示账户从未交易。来源请求失败、查询超时或存储失败会返回非 200 错误，不返回伪造的空列表。来源有记录但全部无法解析且没有现货或未支持市场时，返回 422。

### 4.4 已保存成交时间查询

沿用 GET /api/v1/trade-events，以 source=stored 选择数据库查询。此调用只读取已保存事实，不调用来源、不创建采集 query_id。

| 参数 | 类型 | 必填/默认 | 意义与校验 |
| --- | --- | --- | --- |
| account | string | 必填 | 地址规则同第 4.3.1 节 |
| source | string | 必须为 stored | 省略时仍执行 live，不是数据库查询 |
| limit | integer | 默认 100 | 每页展示上限，1–2000；下一页可调整 |
| start_time | string | 可选，无下界 | 包含开始时刻；带时区 RFC 3339，最多 3 位小数秒，转换为 UTC |
| end_time | string | 可选，无上界 | 不包含结束时刻；格式同 start_time |
| cursor | string | 可选，第一页无游标 | 直接使用上一页 next_cursor；最大 8192 字节，绑定账户、网络及时间条件 |

区间为 `[start_time,end_time)`；两者均存在时必须 start_time < end_time。无时区时间、负 Unix 时间、超出毫秒精度、未知 source、游标损坏或条件不匹配返回 400。分页期间新写入的事实不进入原快照；重新发起不带游标的请求才能看到新库存。

```sh
curl --noproxy '*' -i --get 'http://127.0.0.1:8080/api/v1/trade-events' \
  --data-urlencode 'account=0x010461c14e146ac35fe42271bdc1134ee31c703a' \
  --data-urlencode 'source=stored' \
  --data-urlencode 'start_time=2026-10-04T00:00:00.000Z' \
  --data-urlencode 'end_time=2026-10-05T00:00:00.000Z' \
  --data-urlencode 'limit=10'
```

成功响应仍为 data/meta；data 字段如下，与实时响应按 source 区分：

| data 字段 | 类型 | 意义 |
| --- | --- | --- |
| account/network | string | 规范化账户和配置网络 |
| query_scope | string | 固定 STORED_TIME_RANGE |
| coverage | string | 固定 STORED_RECORDS_ONLY，不承诺账户历史完整 |
| request_range | object | start_time/end_time 为标准 UTC 时间或 null；null 表示该方向无限制 |
| snapshot_seq | string | 此分页快照的已提交事实水位，以十进制字符串传输 |
| matched_records | integer | 本次条件及快照下的完整匹配数，不是本页条数 |
| returned_records | integer | 本页条数，等于 trades.length，不超过 limit |
| observed_range | object/null | 完整匹配集的 first_at/last_at，含义同第 4.3.2 节，无匹配为 null |
| trades | array | 持久化的标准事实；结构见第 4.3.4、4.3.5 节，保留首次来源引用 |
| has_more | boolean | 快照内是否有下一页 |
| next_cursor | string/null | 下一页游标，最后一页为 null |
| warnings | array of string | 包含 STORED_HISTORY_NOT_VERIFIED，说明保存历史尚未证明完整 |

示例：尚无保存记录时，data 返回如下（实际仍包含 meta）：

```json
{
  "account":"0x010461c14e146ac35fe42271bdc1134ee31c703a",
  "network":"mainnet",
  "query_scope":"STORED_TIME_RANGE",
  "coverage":"STORED_RECORDS_ONLY",
  "request_range":{"start_time":null,"end_time":null},
  "snapshot_seq":"0",
  "matched_records":0,
  "returned_records":0,
  "observed_range":null,
  "trades":[],
  "has_more":false,
  "next_cursor":null,
  "warnings":["STORED_HISTORY_NOT_VERIFIED"]
}
```

有下一页时，把 next_cursor 作为 URL cursor 参数传入，并保持 source、account、时间范围一致，limit 可调整。排序仍为成交时间降序、tid 数值降序、fact_id 升序。stored 不返回 live counts、采集 query_id 或 display_truncated。空库存不说明账户没有交易；数据库错误返回非 200。

### 4.5 自动更新状态

新增于 v0.3。`GET /api/v1/watch-accounts`，没有业务参数或请求正文。本版只读当前配置的一个账户；不提供添加、删除、暂停或恢复账户的管理接口。采集由独立 trade-collector 运行；v0.4 可同时启用 WebSocket 实时接收与 HTTP 补偿，读取状态不会执行来源查询。

```sh
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/watch-accounts
```

成功 data 为 `{ "items": [状态对象] }`，外层仍为第 2.1 节的 data/meta 信封，正常 items 只有一个对象：

```json
{
  "data": {
    "items": [{
      "account": "0x0000000000000000000000000000000000000001",
      "account_key": "hyperliquid:mainnet:hyperliquid:0x0000000000000000000000000000000000000001",
      "network": "mainnet",
      "status": "WAITING",
      "coverage": "SOURCE_HISTORY_NOT_VERIFIED",
      "initial_start_time": "2026-10-05T00:00:00.000Z",
      "scanned_through": "2026-10-05T00:01:00.000Z",
      "last_trade_at": null,
      "last_attempt_at": "2026-10-05T00:01:02.000Z",
      "last_success_at": "2026-10-05T00:01:03.000Z",
      "consecutive_failures": 0,
      "next_run_at": "2026-10-05T00:01:33.000Z",
      "pending_range": null,
      "last_query_id": "query_example",
      "last_success_query_id": "query_example",
      "last_error": null,
      "heartbeat_at": "2026-10-05T00:01:03.000Z",
      "lease_expires_at": "2026-10-05T00:04:03.000Z",
      "monitoring_status": "LIVE",
      "websocket": {
        "enabled": true,
        "status": "LIVE",
        "session_id": "00000000-0000-4000-8000-000000000001",
        "connected_at": "2026-10-05T00:00:01.000Z",
        "subscribed_at": "2026-10-05T00:00:02.000Z",
        "last_received_at": "2026-10-05T00:00:02.100Z",
        "last_committed_at": "2026-10-05T00:00:02.200Z",
        "last_trade_at": null,
        "last_pong_at": "2026-10-05T00:01:02.000Z",
        "reconnect_count": 0,
        "connection_count": 1,
        "next_retry_at": null,
        "last_error": null,
        "pending_messages": 0,
        "pending_bytes": 0,
        "metadata_stale": false
      },
      "recovery": {
        "status": "HTTP_SCANNED",
        "target_through": null,
        "open_gap_count": 0,
        "last_scanned_at": "2026-10-05T00:01:03.000Z",
        "last_error": null
      },
      "warnings": ["SOURCE_HISTORY_NOT_VERIFIED", "LATE_DATA_OUTSIDE_OVERLAP_NOT_VERIFIED"]
    }]
  },
  "meta": {"trace_id": "trace_0123456789abcdef0123456789abcdef", "schema_version": 1}
}
```

以上是字段格式示例，不代表实际采集证据。

| 字段 | 类型 | 意义 |
| --- | --- | --- |
| items | array | 当前配置中的采集账户，包含下列状态字段 |
| account/account_key/network | string | 规范地址、链协议账户身份、来源环境 |
| status | string | HTTP 调度状态：STARTING 初始化；RUNNING 本轮采集中；WAITING 成功后等待；RETRY_WAIT 退避或预算暂停；FAILED 无法继续或租约心跳过期；STOPPED 正常停止 |
| coverage | string | 固定 SOURCE_HISTORY_NOT_VERIFIED，不证明官方历史完整 |
| initial_start_time | string | 首次建立检查点的扫描起点，UTC 毫秒格式 |
| scanned_through | string/null | 已提交扫描结束边界，半开区间结束值；不是最新成交时间 |
| last_trade_at | string/null | HTTP 自动采集已提交的最新合约成交时间；无成交可为空；WS 最新成交另见 websocket.last_trade_at |
| last_attempt_at | string/null | 最近一次开始尝试的时间 |
| last_success_at | string/null | 最近一次 HTTP 事实与水位完整事务提交时间；HTTP 返回 200 或收到 WS 消息均不足以更新此字段 |
| consecutive_failures | integer | 连续失败轮数；完整成功后归零，预算暂停不计失败 |
| next_run_at | string/null | 下一次计划执行时间，包含有界退避 |
| pending_range | object/null | 未完成固定范围，含 start_time/end_time，开始包含、结束排除 |
| last_query_id | string/null | 最近尝试的采集任务标识 |
| last_success_query_id | string/null | 最近成功的自动采集任务标识，可用于 reparse |
| last_error | object/null | 最近未恢复错误，含 code/message/occurred_at；完整成功后清空 |
| heartbeat_at/lease_expires_at | string/null | 实例心跳与租约到期时间；数据库持久化值 |
| warnings | array of string | 来源历史、重叠区间外迟到数据、预算或心跳等限制说明 |

v0.4 增加以下字段，接口路径和 schema_version=1 保持不变。顶层水位、轮次失败计数、任务标识和调度时间仍表示 HTTP；WS 与 HTTP 的错误独立记录。

| 字段 | 类型 | 意义 |
| --- | --- | --- |
| monitoring_status | string | HTTP_ONLY 纯 HTTP；STARTING 初始化或尚未确认心跳；LIVE 双通道可用且恢复范围已扫描；RECOVERING 实时已订阅、HTTP 正在追赶；DEGRADED 某通道异常；FAILED 阻塞错误或租约过期；STOPPED 正常停止 |
| websocket.enabled | boolean | 是否启用来源实时订阅 |
| websocket.status | string | DISABLED、CONNECTING、SUBSCRIBING、LIVE、RECONNECT_WAIT、FAILED、STOPPED；LIVE 仅表示订阅已确认，需结合 pong 与 monitoring_status 判断健康 |
| websocket.session_id | string/null | 当前或最近连接的 UUID；不是来源游标 |
| websocket.connected_at/subscribed_at | string/null | 建立连接、收到订阅确认的 UTC 毫秒时间 |
| websocket.last_received_at | string/null | 最近归档账户数据消息的接收时间；不证明事实已提交，队列中尚未归档的消息不更新此持久化值 |
| websocket.last_committed_at | string/null | 最近数据消息完整事务提交时间；空快照也可以更新 |
| websocket.last_trade_at | string/null | 已提交 WS 合约成交的最大时间；不能推进或替代 HTTP scanned_through |
| websocket.last_pong_at | string/null | 最近一次应用层 ping 对应 pong 的确认时间；用于安静账户连接健康判断 |
| websocket.reconnect_count | integer | 首次成功连接之外重新建立连接的累计次数，跨进程重启保留；失败握手不计入此成功连接次数 |
| websocket.connection_count | integer | 成功建立连接的总次数，包含首次；跨重启保留 |
| websocket.next_retry_at | string/null | 下一次重连计划时间；退避、连接频率和恢复条件同时约束实际尝试 |
| websocket.last_error | object/null | 独立 WS 错误，含 code/message/occurred_at；订阅恢复后清除 |
| websocket.pending_messages/pending_bytes | integer | 有界内存队列的最近观测值；进程重启归零，非待补采成交数量；字节包含原始包络与提取的 fills |
| websocket.metadata_stale | boolean | 是否沿用旧市场映射；为 true 时 warnings 包含 METADATA_STALE |
| recovery.status | string | NOT_STARTED 恢复目标尚未固定；SCANNING HTTP 追赶中；HTTP_SCANNED 已完整扫描目标；BLOCKED 无法继续；DISABLED 未启用 WS 且无待补偿范围 |
| recovery.target_through | string/null | 重新订阅确认时固定的 HTTP 恢复目标；未固定或无待补偿范围时为 null |
| recovery.open_gap_count | integer | 尚未完成 HTTP 扫描的恢复范围数量，包含阻塞范围 |
| recovery.last_scanned_at | string/null | 最近恢复范围随 HTTP 水位事务完成扫描的时间 |
| recovery.last_error | object/null | 阻塞恢复的错误摘要 |

例如，WS 已保存较新的成交但 HTTP 尚未扫描中间范围时，websocket.status=LIVE，monitoring_status=RECOVERING，scanned_through 仍为原 HTTP 水位。只有 HTTP 范围及事实成功提交后，恢复范围才能标为 HTTP_SCANNED；此状态仍不证明来源历史完整。关闭 WS 不删除缺口，已有 HTTP 补偿继续执行。

WS 数据模式 SNAPSHOT、LIVE_UPDATE、UNKNOWN 保存在采集任务和原始证据中，不改变 stored 的标准成交字段。当前没有对外实时推送接口，stored 查询不会触发订阅或补偿。具体协作流程见 [v0.4 第 6 章](version0.4.md#6-处理流程)。

所有时间为 UTC、最大毫秒精度。无成交也可以成功推进扫描水位。原 live 查询和旧证据导入不会推进自动水位；stored 查询沿用 v0.2 行为。

未启用 collector、collector 不可访问、检查点尚未建立或数据库不可读时返回 503，不返回成功空列表。服务可读而采集失败时返回 200，status/last_error 表示业务失败。来源可能限流或截断，失败可见不等于自动补齐完整历史。

### 4.6 webhook 发布状态

v0.5 新增：`GET /api/v1/publishing-status`，无查询参数，从挂载的 collector 配置选择账户。用于确认发布是否启用、候选是否入队以及投递是否失败。进程不可达或数据库不可读返回 503，不能把不可用解释为空队列。

```sh
curl --noproxy '*' -sS http://127.0.0.1:8080/api/v1/publishing-status
```

返回公共 data/meta 信封，data 字段如下：

| 字段 | 类型 | 意义 |
| --- | --- | --- |
| enabled | boolean | 持久化业务开关；停止容器不等于停用候选生成 |
| status | string | DISABLED、STARTING、RUNNING、DEGRADED、BLOCKED 或 STOPPED |
| account / account_key | string | 配置账户及网络分区身份 |
| target_id | string | 目标 URL 的 SHA-256 身份，不返回地址、凭证或数据库连接串 |
| activation_epoch / activated_at | integer / RFC3339 string | 激活代次及起点；同配置重启保持，重新启用或改变策略创建新代次 |
| heartbeat_at | RFC3339 string/null | 发布循环心跳；超过 30 秒表示 STOPPED，仍不能证明接收端已消费 |
| policy | object | max_event_age_seconds、signal_ttl_seconds、clock_skew_tolerance_seconds、version |
| outbox | object | PENDING、SENDING、RETRY_WAIT、DELIVERED、BLOCKED 的计数；不存在的状态按 0 理解 |
| oldest_pending_at | RFC3339 string/null | 未送达队列的最老创建时间，包含 BLOCKED |
| last_published_at / last_delivered_at | RFC3339 string/null | 最近事件的首次准备发送时间及最近确认时间；重试不改首次发布时间 |
| last_error | string/null | 脱敏错误类别，或最近仍有错误的事件原因 |
| expired_pending | integer | 已过 expires_at、仍未送达的事件数 |
| suppression_counts | object | 各 SUPPRESSED 原因的累计数量；保存事实不等于发布候选 |

### 4.7 webhook 候选事件

发布器向配置 URL 发送 JSON POST，一个请求一个事件。默认要求 HTTPS，Docker 验收内网 HTTP 必须显式设置 allow_plain_http=true。凭证来自文件，不跟随重定向。

| 请求头 | 意义 |
| --- | --- |
| Content-Type: application/json | JSON 事件正文 |
| Authorization: Bearer `<凭证>` | 接收方验证的共享凭证；不进入事件或日志 |
| X-Robotech-Event-Id | 与正文 event_id 相同的稳定事件身份 |
| X-Robotech-Delivery-Attempt | 累计领取次数；重复发送时增加，可能包含领取后未发送的尝试 |

正文为标准 `AccountFactEnvelope`：

| 字段 | 类型 | 意义 |
| --- | --- | --- |
| schema_version | integer | 当前 1 |
| event_type | string | account.fact.v1 |
| event_id | string | `<fact_id>:1`，复用标准事实版本身份，全链路去重键 |
| partition_key | string | 原 account_key |
| occurred_at | RFC3339 string | 原成交时间 |
| received_at | RFC3339 string | 首次准入 WS 观察接收时间，不能用回补或重试时间替代 |
| stored_at | RFC3339 string | 事务内保存候选的时间，不是精确提交时刻 |
| published_at | RFC3339 string | 首次领取时持久化的发送准备时间；实际 webhook 请求必有此值 |
| expires_at | RFC3339 string | occurred_at 加固定 signal_ttl_seconds，默认 60 秒；重试不延期 |
| observation | object | transport=WEBSOCKET；query_id、session_id、message_sequence、raw_log_id、原始 message_mode、realtime_reason、publishing_policy_version |
| fact | object | 第 4.3.4～4.3.5 节的完整标准 AccountFact；金额继续使用十进制字符串 |

`observation.raw_log_id` 引用准入 WS 的来源 ID；`fact.raw_log_id` 可保留更早 HTTP 观察的来源 ID。realtime_reason 为 EXPLICIT_LIVE_UPDATE，或同会话成功快照之后字段确实缺失时的 POST_SNAPSHOT_UNFLAGGED。原 message_mode 不改写。快照、HTTP、激活前、过期、强制成交、元数据陈旧和无法确认模式的观察不生成候选。

接收方应先持久化 event_id 和原文，再返回 2xx；重复 event_id 成功返回且不重复触发下游决策。首次发送持久化最终字节及摘要，重试正文完全相同。过期事件仍可交付，接收方必须检查 expires_at，不能把候选直接当作下单授权。

2xx 确认送达；网络错误、408/425/429/5xx 有界退避；429/503 的 Retry-After 支持秒数或 HTTP 日期，超过 24 小时标为 BLOCKED；3xx 和其他 4xx 标为 BLOCKED。临时失败不自动丢弃。修复拒绝原因后执行：

```sh
docker compose exec -T trade-parser-publisher trade-parser-publisher retry --event-id '<实际 event_id>'
```

只允许重试本配置账户的 BLOCKED 事件，不修改目标、身份、正文或 TTL。接收成功后发送方确认失败可能重复投递，这是至少一次语义。

## 5. HTTP 状态及错误码

| HTTP 状态 | code | 意义 |
| --- | --- | --- |
| 400 | VALIDATION_ERROR | 参数、请求体或字段格式不合法 |
| 401 | AUTHENTICATION_REQUIRED | 内部接口缺少或使用错误服务凭证；对外接口当前不使用此认证 |
| 404 | RESOURCE_NOT_FOUND | 请求路径不存在 |
| 405 | METHOD_NOT_ALLOWED | 路径存在但不支持当前 HTTP 方法 |
| 409 | VERSION_CONFLICT | 相同来源市场和成交 tid 对应不同内容 |
| 422 | INCOMPLETE_DATA | 来源数据无法可靠解析，无法给出可用标准结果 |
| 429 | RATE_LIMITED | 内部并发容量已满，或来源限流在有限重试后仍未解除 |
| 500 | INTERNAL_INVARIANT_VIOLATION | 证据存储等内部处理失败 |
| 503 | DEPENDENCY_UNAVAILABLE | 查询未启用、内部服务不可用、来源请求失败、元数据异常、查询超时、数据库操作失败或自动采集状态不可用 |

自动状态中的 last_error 是业务诊断，不直接作为该状态请求的 HTTP 状态。WS 连接异常可记录 CONNECT_TIMEOUT、SUBSCRIBE_TIMEOUT、SUBSCRIPTION_REJECTED、HEARTBEAT_TIMEOUT、WS_IO_ERROR、DISCONNECT；数据及处理异常可记录 INVALID_WS_MESSAGE、WS_ACCOUNT_MISMATCH、MESSAGE_TOO_LARGE、QUEUE_OVERFLOW、COMMIT_TIMEOUT、METADATA_UNAVAILABLE、VERSION_CONFLICT、LEASE_LOST。数据库可读时仍以 HTTP 200 返回这些状态；数据库不可读或 collector 不可达才返回 503。

网关会对内部错误进行转换，内部认证失败等无法正常识别的内部响应会对外表现为 503。message 可能因错误发生位置不同而变化。

当前没有公开的证据下载接口。v0.2 原始字节、任务及完整事实保存于 PostgreSQL，文件仅为诊断副本；重新解析不依赖文件副本。文件证据在交易查询容器 `/var/lib/robotech/trade-log/<query_id>/`，包含 manifest.json、完整 result.json 和 requests/responses/metadata 子目录；通过部署运维命令查看。手动验证流程见 [v0.2 开发文档第 10.4 节](version0.2.md#104-手动验证步骤)。

## 6. 更新记录

所有接口的新增、变更和废弃均在本章记录，并同步更新正文。程序版本表示变更所属版本，接口路径版本及 schema_version 表示接口契约版本，两者分别维护。

| 日期 | 程序版本 | 类型 | 更新内容 |
| --- | --- | --- | --- |
| 2026-10-06 | 0.5.0 | 新增接口 | 发布状态内部及外部接口、标准候选 webhook 信封、投递头、时效和重试契约 |
| 2026-10-06 | 0.4.0 | 兼容扩展 | 自动状态增加 monitoring_status、websocket、recovery；HTTP 水位与 WS 接收、提交位置独立；stored 路径及 schema_version 不变 |
| 2026-10-06 | 0.3.0 | 新增接口 | 新增单地址自动采集状态、扫描水位和错误说明；collector 内部状态及健康接口 |
| 2026-10-05 | 0.2.0 | 新增及兼容扩展 | 实时查询增加持久化回执；新增 source=stored 时间范围及快照游标分页、内部库存查询；数据库故障明确返回 503 |
| 2026-10-05 | 0.1.0 | 文档调整 | 将接口说明集中到本文件；按内部、外部接口分章，增加用途、调用时机和接口列表页内跳转；接口行为未改变 |
| 2026-10-04 | 0.1.0 | 新增接口 | 新增外部 GET /api/v1/trade-events；新增内部 POST /internal/v1/trade-queries 与 GET /internal/v1/health；新增成交、覆盖、计数及证据字段说明 |
| 2026-10-04 | 0.0.0 | 初始接口 | 提供外部 GET /api/v1/health 与 GET /api/v1/version，以及统一响应、追踪 ID 和路由错误格式 |
