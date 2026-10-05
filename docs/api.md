# 接口说明

本文是项目统一的接口说明文档，集中维护所有对外接口和内部接口。后续版本在本文件中更新接口定义，并注明新增、变更或废弃的适用版本，不另建按开发版本命名的接口文档。

当前对应程序版本：0.1.0；JSON schema_version：1。接口定义依据当前代码编写。

## 1. 地址和接口列表

默认网关地址：服务器本机 `http://127.0.0.1:8080`；公网地址示例 `http://server.example.com:8080`。`server.example.com` 是占位域名，使用时替换为实际服务器地址。公网是否可访问取决于实际部署及防火墙配置。

| 范围 | 方法 | 路径 | 用途 |
| --- | --- | --- | --- |
| 对外 | GET | [/api/v1/health](#41-网关健康接口) | 判断网关进程是否存活 |
| 对外 | GET | [/api/v1/version](#42-版本接口) | 查询网关程序版本 |
| 对外 | GET | [/api/v1/trade-events](#43-账户近期合约成交查询) | 查询账户近期合约成交 |
| 内部 | GET | [/internal/v1/health](#33-内部健康接口) | 判断交易查询进程是否存活，需要服务凭证 |
| 内部 | POST | [/internal/v1/trade-queries](#32-内部成交查询) | 网关调用交易查询服务，需要服务凭证 |

对外接口当前不要求客户端提供身份凭证。内部服务默认监听 8081，Compose 不发布该端口；网关在容器网络中使用 `http://trade-log-query:8081`。

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
| meta.schema_version | integer | 响应结构版本，当前为 1；与程序版本 0.1.0 分开管理 |

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

内部接口用于项目内不同进程之间的通信。目前 query-api 负责对外 HTTP 请求，trade-log-query 负责来源采集、成交解析和证据保存。通过内部接口明确两者的请求、响应、认证和错误边界，使业务查询过程可以独立运行和迭代，网关无需依赖交易查询进程的内部实现。

当客户端调用外部成交查询接口时，网关完成参数校验，再调用内部成交查询接口；交易查询服务完成采集和证据保存后，网关将结果返回客户端。内部健康接口供运维或受信任的服务调用方检查交易查询进程是否存活，当前网关不会在每次成交查询前主动调用它。

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

成功响应与[外部成交查询](#432-成功响应结构)相同，错误格式与第 2.2 节相同。默认最多并行执行 4 个查询，超过立即返回 429；默认整体查询超时 30 秒，超时返回 503。并发和超时以 trade-log.toml 的实际配置为准，网关内部请求默认超时为 35 秒。

### 3.3 内部健康接口

```http
GET /internal/v1/health
```

没有业务参数，需要 Authorization 请求头。成功返回 200，data 格式：

```json
{"status":"ok","service":"trade-log-query","version":"0.1.0"}
```

status 表示内部进程可响应，service 为服务名，version 为程序构建版本；完整响应使用 data/meta 包装。该接口不主动检查上游来源。

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
{"service":"query-api","version":"0.1.0"}
```

| data 字段 | 类型 | 含义 |
| --- | --- | --- |
| service | string | 服务名 query-api |
| version | string | 当前运行程序的构建版本，v0.1 为 0.1.0 |

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

请求不带 JSON body。缺少 account、未知参数、重复的同名参数、非法类型或超出范围均返回 400。network 从服务配置读取，不能通过请求切换；不提供起止时间、分页游标、市场筛选或排序参数。

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
    "evidence_ref": "query_0123456789abcdef0123456789abcdef"
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
| evidence_ref | string | 当前等于 query_id，是证据目录标识；不是下载 URL |

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
| 503 | DEPENDENCY_UNAVAILABLE | 查询未启用、内部服务不可用、来源请求失败、元数据异常或查询超时 |

网关会对内部错误进行转换，内部认证失败等无法正常识别的内部响应会对外表现为 503。message 可能因错误发生位置不同而变化。

当前没有公开的证据下载接口。证据在交易查询容器 `/var/lib/robotech/trade-log/<query_id>/`，包含 manifest.json、完整 result.json 和 requests/responses/metadata 子目录；通过部署运维命令查看。手动验证流程见 [v0.1 开发文档第 10.4 节](version0.1.md#104-手动验证步骤)。

## 6. 更新记录

所有接口的新增、变更和废弃均在本章记录，并同步更新正文。程序版本表示变更所属版本，接口路径版本及 schema_version 表示接口契约版本，两者分别维护。

| 日期 | 程序版本 | 类型 | 更新内容 |
| --- | --- | --- | --- |
| 2026-10-05 | 0.1.0 | 文档调整 | 将接口说明集中到本文件；按内部、外部接口分章，增加用途、调用时机和接口列表页内跳转；接口行为未改变 |
| 2026-10-04 | 0.1.0 | 新增接口 | 新增外部 GET /api/v1/trade-events；新增内部 POST /internal/v1/trade-queries 与 GET /internal/v1/health；新增成交、覆盖、计数及证据字段说明 |
| 2026-10-04 | 0.0.0 | 初始接口 | 提供外部 GET /api/v1/health 与 GET /api/v1/version，以及统一响应、追踪 ID 和路由错误格式 |
