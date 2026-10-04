# v0.1 开发文档：Hyperliquid 近期合约成交查询

- 文档版本：0.1
- 状态：代码已实现，本地与真实来源验证通过；Docker 部署待服务器验收
- 更新日期：2026-10-04
- 设计依据：[概要设计](overview-design.md)、[详细设计](detailed-design.md)、[产品版本路线图](roadmap.md)
- 前置版本：[v0.0 程序基础](version0.0.md)

## 1. 版本目标

用户通过 HTTP 提交一个 Hyperliquid 实际交易账户地址，查看官方来源当前可返回的近期合约成交。结果可追溯到保留的来源响应，并明确区分无成交、筛选后为空、覆盖受限和查询失败。

沿用 v0.0 的 Docker Compose 启动方式、查询网关、配置加载、JSON 信封、trace ID、日志和有界退出。本版新增交易日志查询进程，协议请求与成交解析不写进网关 handler。

版本发布后仍通过 `docker compose up --build -d` 启动整个系统。健康和版本接口继续可用，程序及镜像版本升级为 `0.1.0`，HTTP 路径和响应 schema 保持 v1。

## 2. 功能范围

| 能力 | 本版边界 |
| --- | --- |
| 地址查询 | 一次提交一个主账户、子账户或 vault 地址；不自动展开关联账户，不使用 agent wallet 替代账户 |
| 近期成交 | 请求官方近期 fill 窗口，无需预先添加监控账户 |
| 合约展示 | 首版保证默认永续 DEX 的合约成交；其他市场需元数据确认，未支持或未知市场明确计数 |
| 来源保留 | 文件保留请求、原始响应、元数据及解析结果，事实关联来源记录位置 |
| HTTP 输出 | 统一 JSON、查询时间、来源数量、筛选数量、展示截断与覆盖警告 |
| 基础约束 | 参数、金额精度、市场识别、事实身份、去重及来源超时校验 |

不实现历史时间区间分页、定时采集、WebSocket、webhook、关注账户管理或现货成交展示。时间区间与存储查询在 v0.2 设计时细化，持续采集按路线图推进。

文件证据不等于成交数据库：本版不提供跨次查询合并、按 fact_id 的持久化查询、事实修订管理或采集水位。账户盈亏、胜率和跟单执行不属于本版。

## 3. 使用方式

### 3.1 对外接口

沿用详细设计的查询路径：

```text
GET /api/v1/trade-events?account=<ADDRESS>&limit=100
```

| 参数 | 规则 |
| --- | --- |
| `account` | 必填，0x 加 40 位十六进制字符；内部转小写，原输入保留在查询证据 |
| `limit` | 可选，默认 100，整数 1–2000；仅限制返回数量，不改变来源窗口 |

网络由服务配置指定，默认 mainnet；本版不允许请求携带来源 URL、网络覆盖或文件路径。未知查询参数返回 400，避免用户误认为 `start` 等尚未支持的条件已经生效。

请求示例，替换为实际账户后执行：

```sh
curl --get 'http://127.0.0.1:8080/api/v1/trade-events' \
  --data-urlencode 'account=<ADDRESS>' \
  --data-urlencode 'limit=100'
```

不支持的方法返回 JSON 405，Allow 包含 GET 和 HEAD；HEAD 沿用 HTTP 语义，可能执行相同查询但不返回响应体，不作为轻量健康检查。

### 3.2 操作方式

```sh
docker compose up --build -d
docker compose ps
docker compose logs -f query-api trade-log-query
```

对外仍访问网关 8080。应用监听配置以各自 TOML 为准；Compose 负责容器网络、端口发布和挂载，不重复覆盖应用监听值。

## 4. 数据来源与查询边界

### 4.1 官方请求

mainnet 来源为 `https://api.hyperliquid.xyz/info`；testnet 来源为 `https://api.hyperliquid-testnet.xyz/info`。公开查询不需要私钥或钱包签名。

近期成交请求体：

```json
{"type":"userFills","user":"<ADDRESS>","aggregateByTime":false}
```

来源接口同时可能返回现货成交。每次查询获取并保存 `meta` 与 `spotMeta`，用官方元数据确认市场类型，不依赖 `coin` 字符串外观。默认合约市场映射以 `meta.universe` 为准，现货按 `spotMeta` 的交易对与资产索引解析。

非默认永续 DEX、已下架且当前元数据无法确认的市场或其他新资产类型，不能直接归类为默认合约。首版以 `UNSUPPORTED_MARKET` 或 `UNKNOWN_MARKET` 披露；后续可通过 `perpDexs` 与对应元数据扩展同一解析器。

来源报告近期接口最多返回 2000 条，且按时间接口也有历史保留限制；因此不能由近期结果证明完整账户历史。`aggregateByTime=false` 保留 fill 粒度，不将部分成交合成订单。

### 4.2 覆盖与数量

- 本版查询范围是 `RECENT_SOURCE_WINDOW`，没有承诺固定天数。
- 返回最早和最晚成交时间只是实际记录区间，不是已验证完整覆盖起点。
- 请求成功且未发现具体异常，覆盖为 `SOURCE_WINDOW_UNVERIFIED`。
- 来源返回满 2000 条、存在未知或未支持市场、部分记录不可解析时，覆盖为 `LIMITED` 并列出原因。
- `limit` 截断通过 `display_truncated` 单独披露，不把显示 100 条描述为仅获取了 100 条。

无来源记录、来源只有现货、来源只有未支持市场必须分别说明。官方结果为空也不证明地址不存在或历史上没有交易。

按时间历史补偿在后续采集版本使用 `userFillsByTime`；本版不提前实现分页与历史完整性承诺。

官方依据：[Info endpoint](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint)、[市场与资产标识](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/asset-ids)、[永续元数据](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/info-endpoint/perpetuals)、[限流说明](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/rate-limits-and-user-limits)。核对日期：2026-10-04；实现前复核来源契约。

## 5. 接口与成交数据定义

### 5.1 查询结果

成功响应沿用 `data` 与 `meta`：

```json
{
  "data": {
    "query_id": "query_example",
    "account": "<NORMALIZED_ADDRESS>",
    "network": "mainnet",
    "queried_at": "2026-10-04T12:00:00.000Z",
    "query_scope": "RECENT_SOURCE_WINDOW",
    "coverage": "SOURCE_WINDOW_UNVERIFIED",
    "counts": {
      "source_records": 25,
      "duplicate_records": 0,
      "spot_records": 5,
      "unsupported_records": 0,
      "invalid_records": 0,
      "perpetual_records": 20,
      "returned_records": 20
    },
    "display_truncated": false,
    "observed_range": {"first_at": "2026-10-03T12:00:00.000Z", "last_at": "2026-10-04T11:00:00.000Z"},
    "trades": [],
    "warnings": [],
    "evidence_ref": "query_example"
  },
  "meta": {"trace_id": "trace_example", "schema_version": 1}
}
```

上例仅说明结构，实际 `trades` 数量必须与计数一致。无记录时 `observed_range=null`。范围按成功解析的合约全集计算，未受展示 limit 影响。`queried_at` 是查询开始时间，成交时间与接收时间分别保留。

`counts` 排他计数：原始记录扣除重复后，每条唯一记录只归入 spot、unsupported、invalid 或 perpetual 之一；来源身份冲突整体失败。unknown 市场计入 unsupported，并以警告区分原因。

`trades` 使用详细设计 `AccountFact` 的 TRADE 契约，不创建另一套只供 API 使用的业务成交模型；查询摘要 DTO 与协议原始 DTO 分离。`evidence_ref` 是服务内部引用，不暴露宿主机绝对路径，不提供公开证据文件下载。

### 5.2 来源字段与语义

| 来源字段 | 标准字段/用途 |
| --- | --- |
| time | UTC 毫秒成交时间，保留在来源扩展及标准时间字段 |
| coin | 通过元数据生成 market、instrument_type、base_asset、quote_asset |
| side | B → BUY，A → SELL；未知值报异常 |
| px、sz | price、quantity，正数定点十进制 |
| fee、feeToken | fee、fee_asset，保留返佣负值；缺失资产不默认为 USDC |
| tid | source_ref、稳定事实身份组成部分，整数不转浮点 |
| oid、hash | order_id、transaction_hash，不能单独识别一条 fill |
| dir、startPosition | 来源操作描述及仓位效果判断证据 |
| closedPnl | reported_realized_pnl，不等于本系统净收益 |
| crossed | 来源扩展，不单凭它判断主动交易或清算 |
| 清算/强制相关字段 | trigger_type 与来源证据，不能误标为普通可跟单成交 |

共享类型遵循详细设计：金额、价格和数量以十进制字符串输出；计算不用 f32/f64。notional 按 price × quantity 计算，检查精度范围和溢出，不静默舍入。

根据 startPosition、side 和 quantity 判断 OPEN、INCREASE、DECREASE、CLOSE、REVERSE；无法可靠判定则 UNKNOWN。买入不等于开多，卖出不等于开空。

已确认主动成交按详细设计标记 USER 和 copy_eligible=true；清算、强制事件标记不可跟单。无法确认来源触发语义时保守设置 copy_eligible=false 并记录原因。本版只查询，不发布信号。

### 5.3 身份与顺序

账户键包含网络、协议和规范化地址。稳定 fact_id 使用带版本的规范编码，对网络、协议、账户、来源市场和 tid 生成确定性摘要；具体编码与摘要算法在公共契约中固定，并提供固定向量测试。不能使用请求时间、数组位置或随机 ID 作为事实身份。

初次查询结果 revision=1、change_type=UPSERT、confirmation_status=OBSERVED；网络查询不冒充链最终性确认。排序使用成交时间降序，再按 tid 和 fact_id 确定同时间顺序，仅表达稳定展示顺序，不宣称链上全序。

缺必要身份字段的记录保留证据并标异常。同身份同内容去重并保留多个原始位置引用；同身份不同内容返回冲突错误，不能任意选一条。

## 6. 查询处理流程

1. 网关验证参数，生成 trace ID，调用内部交易日志查询接口。
2. 服务验证内部访问凭证及请求参数，分配 query_id，并取得有限并发许可。
3. 建立证据目录，获取官方元数据和近期成交响应；接收后先保存原始字节再解析。
4. 协议实现解析市场及来源 DTO，转为标准事实；交易日志业务执行校验、去重和冲突检查。
5. 筛选合约，计算排他数量与覆盖警告，排序并按 limit 返回。
6. 保存 manifest 和完整可用合约结果，响应网关；网关使用原来的 trace ID 输出外部响应。

同步查询不持有后台持续采集任务；超时与取消必须传播到内部调用及来源请求。健康接口不依赖来源在线，来源失败不能导致整个服务退出。

证据结构位于配置的数据目录：`<query_id>/manifest.json`、`requests/`、`responses/`、`metadata/`、`result.json`。manifest 记录网络、账户、查询及接收时间、工具/解析器/schema 版本、内容摘要、来源响应引用、记录位置、计数、覆盖和错误。

文件异步读写或使用受控 blocking 任务，临时文件写完后原子提交；失败证据保留。不使用用户提供的路径，不记录授权头。证据是诊断材料，由运营管理保留期和容量，本版不自动删除或充当数据库。

## 7. 程序结构设计

### 7.1 模块与进程

```text
crates/
  shared-types/                 # 网络、地址、时间、金额等值对象
  account-facts/                # 标准事实及身份契约
  protocol-api/                 # 本版实际使用的协议身份与能力
services/trade-log/
  src/{acquisition,raw_log,parsing,normalization,validation,query}/
  adapters/src/{file_evidence,internal_http}/
  bins/trade-log-query/          # 按需查询进程的装配与内部 HTTP 入口
protocols/implementations/hyperliquid/
  src/                          # 官方 HTTP、来源 DTO、元数据与解析
gateway/query-api/
  src/http/handlers/trade_events.rs
  src/clients/trade_log.rs
config/
  query-api.toml
  trade-log.toml
```

以上目录遵循详细设计的服务内聚结构；`trade-log-query` 是本版新增的查询 composition binary，不替代后续 trade-collector 和 trade-parser-publisher。获取与解析由同一查询流程按需调用，不写入事实数据库或发布事件；后续采集和发布进程复用同一业务模块与协议实现。

公共 crate 按实际契约创建，不提前实现资金费、完整账本或跟单业务。workspace 统一管理依赖和锁文件。

### 7.2 依赖方向

- 业务 crate 定义 SourceReader、ProtocolParser、RawEvidenceStore 等使用方端口，依赖公共契约，不依赖 Axum、reqwest、数据库或协议 SDK。
- Hyperliquid crate 实现来源读取与协议解析，来源 DTO 不向网关泄露。
- adapters 实现文件存储及内部 HTTP 边界。
- composition binary 装配业务、协议与适配器。
- 网关仅负责参数、访问控制、内部客户端和响应，不调用 Hyperliquid。

内部 `POST /internal/v1/trade-queries` 接收 account、limit，返回同一查询结果契约。内部接口不发布宿主机端口，且按详细设计要求服务身份认证：通过只读挂载的 token 文件接入最小服务凭证；缺失凭证启动失败，非法凭证返回 401，凭证不进入日志或镜像。

网关生成的 trace ID 通过受认证内部请求传播，内部错误转换为网关统一错误格式。内部服务只接受规定格式的 trace ID，否则自行生成；不直接信任外部客户端提供的 ID。

## 8. 配置与异常处理

### 8.1 新增配置

`query-api.toml` 保持原 server/logging；新增可选 `[trade_log]`，默认不启用，保持旧配置能启动基础接口。启用后必填内部 base_url、credential_file 和 request_timeout_seconds；未启用时成交接口返回 503，不返回假空结果。

`trade-log.toml` 包含 server/logging 及以下必填配置，config_version=1：

| 配置 | 默认部署值/校验 |
| --- | --- |
| hyperliquid.network | mainnet；允许 testnet，端点从网络确定，不接受请求覆盖 |
| hyperliquid.connect_timeout_seconds | 5，正整数 |
| hyperliquid.request_timeout_seconds | 10，覆盖响应读取 |
| query.timeout_seconds | 30，覆盖来源调用、重试、解析和证据保存 |
| query.max_concurrency | 4，正整数，超出立即返回 429，不无限排队 |
| query.max_response_bytes | 16777216，来源读取上限，元数据和成交各自检查 |
| evidence.directory | /var/lib/robotech/trade-log，启动时检查可写 |
| internal.credential_file | /run/secrets/trade-log-token，文件非空且不打印值 |

网关内部请求时限默认 35 秒，需大于业务 30 秒，避免把业务超时误判为无响应。服务启动校验本地配置；跨服务时限约束由部署配置及验收核对。

可重试网络错误、429 和 5xx 最多 2 次尝试，遵循 Retry-After 并受查询总时限约束；普通 4xx 与解析错误不重试。并发限制不能代替官方 IP 加权限流，遇到来源限流须明确返回，不持续压测公共接口。

新配置沿用文件加载与显式环境覆盖约定。本版环境覆盖仍仅支持 ROBOTECH_SERVER_HOST、ROBOTECH_SERVER_PORT、ROBOTECH_SHUTDOWN_TIMEOUT_SECONDS、ROBOTECH_LOG_LEVEL 和 ROBOTECH_LOG_FORMAT；业务新增配置由 TOML 设置；不再在 Compose 重复设置应用 host/port。

### 8.2 HTTP 行为

| 情况 | HTTP / code |
| --- | --- |
| 无效地址、limit 或未知参数 | 400 / VALIDATION_ERROR |
| 成功但无记录/只有现货 | 200，trades 为空，计数和原因明确 |
| 可用结果伴随未知市场/部分解析失败 | 200，coverage=LIMITED，warnings 列明原因 |
| 所有应解析记录不可用 | 422 / INCOMPLETE_DATA |
| 同身份不同内容 | 409 / VERSION_CONFLICT |
| 本服务并发限制或来源限流 | 429 / RATE_LIMITED |
| 来源网络失败、超时或元数据必需请求失败 | 503 / DEPENDENCY_UNAVAILABLE |
| 证据写入失败 | 500 / INTERNAL_INVARIANT_VIOLATION |
| 未启用/无法访问内部查询服务 | 503 / DEPENDENCY_UNAVAILABLE |

失败不能返回成功空数组。错误始终保留外部请求 trace ID，不把内部地址、文件路径、凭证或原始响应直接返回客户端。

## 9. 验证与验收

### 9.1 自动验证

保留 v0.0 全部测试，新增以下有实际语义的验证：

- 地址规范化、limit 边界、未知参数、非法输入不请求来源。
- 固定元数据与混合 fill 样本：合约、现货、未知市场、非默认 DEX、精度、费用负值和大整数 ID。
- 买入平空、卖出平多及反转，不能仅凭 side 判断仓位效果。
- 同订单多 fill 保留，重复去重，身份冲突报错，固定身份向量一致。
- 空来源、全现货、部分解析失败、全部解析失败、接口满额、展示截断，计数和覆盖一致。
- 来源超时、限流、非法 JSON、响应过大、证据写入失败及取消；使用本地 mock，不依赖公网稳定性。
- 原始字节与事实引用可核对，从证据重新解析结果一致。
- 网关与内部服务的错误、trace 传播、服务凭证验证及并发限制。

### 9.2 实际地址验收

选择至少一个有合约成交的实际账户；另以真实账户覆盖混合成交或无近期成交，无法找到的边界用固定样本补充并说明。保存查询时间、网络、地址、元数据及原始响应。

对真实 fill 按 tid、时间、市场、side、px、sz、fee、feeToken 核对；订单聚合界面只作辅助。报告明确核对条数和方式，不把一次实时请求与另一次已变化的来源窗口直接比较。

### 9.3 Docker 验收

两个容器正常启动；宿主机仅发布网关端口；网关能调用内部服务；数据卷重启后证据仍在；非 root 进程能写数据卷，根文件系统和配置挂载只读。来源暂时不可用时 health 仍可访问，成交接口返回明确错误。

修改网络/时限配置后重启生效；更新镜像版本后原 health/version 行为保持兼容；整个 Compose 停止和重启正常。公网或局域网验证根据实际部署地址执行，不在文档中硬编码生产服务器地址。

完成标准是自动测试、真实地址字段核对和双进程 Docker 查询链路均通过，并形成实际验收报告。本文件不表示这些项已经完成。

## 10. 交付与运行说明

### 10.1 容器部署

Compose 新增 `trade-log-query`，网关通过 `http://trade-log-query:8081` 访问。应用内部服务监听由 trade-log.toml 指定为 `0.0.0.0:8081`，不发布到宿主机。

两个进程使用 `gateway/query-api/Dockerfile` 的 `query-api` 和 `trade-log-query` 两个 target 分别构建镜像，共享一次 workspace 编译及运行基线；采用 v0.0 的多阶段构建、固定基础镜像摘要、非 root、exec ENTRYPOINT 与有界退出。加入公共 crate 后，更新构建 COPY 路径覆盖实际 workspace member。

交易日志服务只读挂载配置与内部凭证文件，并将命名数据卷挂载到证据目录；启动前由明确的部署初始化步骤设置卷目录 UID/GID 权限，不能依赖开发机 root 运行成功。卷初始化与 Compose 约定在实现时给出可执行命令。

内部凭证由部署者执行 `sh scripts/init-v0.1.sh` 生成并放在忽略的本地 secrets 目录，两个服务挂载同一凭证，不提交默认共享 token。脚本保留已有凭证，并运行仅在 init profile 中启用的一次性 evidence-init 容器，为命名数据卷设置 UID/GID 10001 和目录权限；正常运行的两个应用容器仍使用非 root。Compose 对文件存在性检查，缺少时明确启动失败。

启动仍使用 `docker compose up --build -d`；配置变化重启对应服务，Compose 挂载、镜像或网络变化重建容器。基础接口不依赖内部服务启动就绪；业务调用遇到内部服务未就绪返回 503，避免仅依赖容器启动顺序判断可用。

### 10.2 交付与迁移

交付公共契约、本版交易日志模块及协议实现、新进程与 Dockerfile、网关客户端和 handler、两份配置示例、更新后的 Compose、测试样本、运行与验收说明。

升级 package/镜像版本为 0.1.0；旧 query-api 配置仍可运行基础接口，新增成交能力需要显式启用并配置内部服务。v0.0 的 HTTP path/schema 不变；v0.1 成交响应按 schema_version=1 起步。

README 不增加版本入口。开发及部署说明集中在本文件和后续实际验收报告。

### 10.3 后续版本复用

v0.2 在交易日志服务接入原始数据与事实持久化，复用获取、市场映射、事实身份、解析和查询端口；v0.3 的采集入口复用同一协议与业务模块；v0.4 的 WebSocket 与 HTTP 数据使用相同事实身份；v0.5 的发布遵守存储后发布，不等待分析与审核。

本版文件证据适配器后续可以继续用于诊断或由正式归档适配器替换，业务查询不直接依赖文件系统。HTTP handler、标准事实及协议类型不会因增加数据库或消息系统而另起一套实现。


### 10.4 手动验证步骤

在服务器项目根目录，按下面顺序逐项操作。每一步先执行命令，再对照预期结果判断是否通过；不需要运行自动校验脚本。

#### 第 1 步：构建并启动

```sh
sh scripts/init-v0.1.sh
docker compose up --build -d
docker compose ps
```

**看什么：** `query-api` 和 `trade-log-query` 两个服务均显示 `Up` 或 `running`；query-api 发布 8080 端口，trade-log-query 没有宿主机端口映射。

**通过标准：** 构建成功，两个服务均运行，没有反复退出。初始化脚本需要 openssl；重复执行会保留已有凭证。

#### 第 2 步：确认启动日志

```sh
docker compose logs --tail=30 query-api trade-log-query
```

**看什么：** 两个服务分别出现 `server_started`，版本为 `0.1.0`。

**通过标准：** 无配置读取、凭证读取、端口绑定或证据目录写入错误。

#### 第 3 步：验证健康和版本接口

```sh
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/health
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/version
```

**看什么：** 两个请求均返回 `HTTP/1.1 200 OK`；health 的 `data.status` 为 `ok`，version 的 `data.version` 为 `0.1.0`。响应头有 `x-trace-id`，与响应体 `meta.trace_id` 相同。

**通过标准：** 状态码和上述字段全部符合。

#### 第 4 步：查询真实账户成交

```sh
curl --noproxy '*' --max-time 45 -i --get 'http://127.0.0.1:8080/api/v1/trade-events' \
  --data-urlencode 'account=0x010461c14e146ac35fe42271bdc1134ee31c703a' \
  --data-urlencode 'limit=10'
```

**看什么：** 返回 200；`data.account` 为查询账户；有 `query_id`、`counts`、`coverage`、`trades`、`warnings` 和 `evidence_ref`。`trades` 最多 10 条，每条含 `fact_id`、`occurred_at` 和 `payload`；payload 中有市场、方向、价格、数量和手续费。

**通过标准：** 真实成交正常返回，`counts.returned_records` 等于展示条数且不超过 10。如果 `counts.perpetual_records` 大于展示条数，`display_truncated` 应为 true。来源达到 2000 条时，coverage 应为 LIMITED，warnings 包含 SOURCE_RECORD_LIMIT。

记下返回的 **query_id**，第 7、8 步使用。这个公开账户的近期成交会变化；如果没有合约成交，应换一个有近期合约成交的公开账户验证。429、503 或超时不算本步骤通过。

#### 第 5 步：验证错误输入

分别执行：

```sh
curl --noproxy '*' -i 'http://127.0.0.1:8080/api/v1/trade-events?account=invalid'
curl --noproxy '*' -i 'http://127.0.0.1:8080/api/v1/trade-events?account=0x010461c14e146ac35fe42271bdc1134ee31c703a&limit=0'
curl --noproxy '*' -i 'http://127.0.0.1:8080/api/v1/trade-events?account=0x010461c14e146ac35fe42271bdc1134ee31c703a&limit=2001'
```

**看什么：** 三个请求均返回 400，响应体有 `code` 为 `VALIDATION_ERROR`，并有 trace_id。

**通过标准：** 每个非法请求都明确返回错误，没有返回 200 空列表。

#### 第 6 步：从自己的电脑验证公网访问

在自己的电脑执行：

```sh
curl --noproxy '*' -i http://120.77.207.116:8080/api/v1/health
curl --noproxy '*' --max-time 45 -i 'http://120.77.207.116:8080/api/v1/trade-events?account=0x010461c14e146ac35fe42271bdc1134ee31c703a&limit=10'
```

**通过标准：** health 返回 200 和 ok；成交查询返回 200，结构符合第 4 步。实时查询结果可能变化，不要求两次成交内容完全一致。如果服务器本机通过而此步失败，检查服务器安全组、系统防火墙和 Compose 的 8080 端口发布。

#### 第 7 步：检查查询证据已落盘

回到服务器，将下面的值替换为第 4 步实际返回的 query_id：

```sh
v01_query_id='query_替换为实际值'
docker compose exec -T trade-log-query ls -R "/var/lib/robotech/trade-log/$v01_query_id"
docker compose exec -T trade-log-query cat "/var/lib/robotech/trade-log/$v01_query_id/manifest.json"
mkdir -p var/manual-v0.1
docker compose cp "trade-log-query:/var/lib/robotech/trade-log/$v01_query_id" var/manual-v0.1/
```

**看什么：** 查询目录含 `manifest.json`、`result.json`、`requests/`、`responses/`、`metadata/`。manifest 的 status 为 COMPLETED，query_id 与 HTTP 返回相同，trace_id 与该次 HTTP 请求相同。

用编辑器打开复制出来的 `var/manual-v0.1/<query_id>/result.json`，与第 4 步结果对照：前面的成交相同；该文件保存完整标准成交，不能只保留 limit=10 的展示结果。例如 HTTP 显示 perpetual_records=2000，完整文件就应保留 2000 条成交。

**通过标准：** 文件齐全，状态和身份一致，完整结果未被展示 limit 截断。原始响应、请求及元数据都有对应文件。

#### 第 8 步：重启后检查证据仍在

```sh
docker compose restart trade-log-query
docker compose exec -T trade-log-query cat "/var/lib/robotech/trade-log/$v01_query_id/manifest.json"
```

**看什么：** 原 query_id 的 manifest 仍可读取，内容仍为 COMPLETED。

再执行第 4 步的查询命令。

**通过标准：** 历史证据仍存在，新查询仍返回 200。普通重启不要使用 `docker compose down -v`，它会删除证据卷。

#### 第 9 步：验证依赖停止时的行为

此步会短暂停止交易查询，请在允许暂停查询时执行。

```sh
docker compose stop trade-log-query
curl --noproxy '*' --max-time 45 -i 'http://127.0.0.1:8080/api/v1/trade-events?account=0x010461c14e146ac35fe42271bdc1134ee31c703a&limit=10'
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/health
docker compose start trade-log-query
```

**看什么：** 依赖停止期间，成交查询返回 503，`code` 为 DEPENDENCY_UNAVAILABLE；网关 health 仍返回 200 和 ok。

启动后等待日志出现 server_started，再执行第 4 步。

**通过标准：** 故障被明确报告，网关仍可响应，恢复后成交查询重新成功。

#### 第 10 步：验证正常停止并恢复服务

```sh
docker compose stop
docker compose ps -a
docker compose logs --tail=30 query-api trade-log-query
```

**看什么：** 两个常驻容器显示 `Exited (0)`；两个服务日志都有 `shutdown_started` 和 `shutdown_completed`。

```sh
docker compose start
```

再执行第 3 步检查健康和版本。

**通过标准：** 正常退出、恢复后基础接口正常。

完成后记录上述 10 步各自是否通过、失败时的状态码和日志。只有实际执行通过才能登记为服务器 Docker 验收通过；当前开发机的测试结果不能替代此记录。

自动检查命令和已执行结果见 [v0.1 验收报告](version0.1-acceptance.md)。
