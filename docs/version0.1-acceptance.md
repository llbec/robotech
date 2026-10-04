# v0.1 验收报告

- 日期：2026-10-04（Asia/Shanghai）
- 状态：代码实现、本地自动测试、真实来源核对及双进程 release 验证通过；Docker 运行待服务器验证
- 开发环境：macOS aarch64，Rust/Cargo 1.96.1
- 基准提交：`c556dc0d7067b4393a0d7871c9e5e869bfcda8c6`；本版改动尚未提交
- 程序版本：query-api 与 trade-log-query 均为 0.1.0

## 实现范围

新增 GET `/api/v1/trade-events`，请求经网关进入受服务凭证保护的交易日志查询进程，读取官方 userFills、meta、spotMeta，并按默认永续市场元数据解析合约事实。保留原始字节、请求、元数据、内容摘要与完整标准结果。

查询结果包括来源覆盖、各类计数、展示截断及稳定事实身份。价格、数量和手续费采用定点数，整数 ID 保留精度。无成交与来源错误分开处理。未支持市场明确披露，不冒充完整账户历史。

v0.0 的配置、日志和退出逻辑抽取到 service-runtime，网关原有基础接口及旧格式配置保持兼容。

## 自动检查

| 检查 | 结果 |
| --- | --- |
| cargo test --locked --workspace | 20 项自动测试通过，1 项联网验收默认 ignored |
| cargo build --locked --release -p query-api -p trade-log-query | 两个 release 程序构建通过 |
| cargo fmt --all -- --check | 通过 |
| cargo clippy --locked --workspace --all-targets -- -D warnings | 通过 |
| release --version | 两个程序均返回 0.1.0 |
| sh -n scripts/init-v0.1.sh | 通过 |
| Compose YAML 解析及服务/端口/profile 检查 | 通过；不能替代 Docker Compose 实际运行 |
| git diff --check | 通过 |

测试覆盖既有 v0.0 行为、配置兼容、账户与 limit 校验、混合市场、开平仓和反转、强制成交、金额精度和溢出、大整数 ID、稳定身份向量、重复与冲突、完整证据与展示截断、回放一致性、401、超时、429/5xx 重试、400 不重试、来源非法 JSON、响应大小限制、并发许可及证据写入失败。

固定来源样本由本项目构造，位于 tests/fixtures/hyperliquid，不冒充实际账户历史。

## 真实来源验收

官方 vaultDetails 查询返回 HLP 主地址 `0xdfc24b077bc1425ad1dea75bcb6f8158e10df303` 的公开 childAddresses。主地址本次 userFills 返回空；选用其中子账户 `0x010461c14e146ac35fe42271bdc1134ee31c703a` 验证实际成交。

2026-10-04T12:02:49.536Z 至 12:02:53.006Z，联网验收通过真实官方来源、受认证的内部 HTTP 服务和网关完成完整链路：

- 来源返回 2000 条，解析为 2000 条合约成交，无重复、无无效或未知市场记录。
- API 返回 100 条，display_truncated 为 true；完整标准结果保留 2000 条。
- 返回的 100 条逐条按原始元素位置核对 tid、时间、coin、side、价格、数量、手续费与费用资产，全部一致。
- coverage 为 LIMITED，warnings 包含 SOURCE_RECORD_LIMIT，没有宣称完整历史。
- 本地证据：`var/acceptance/v0.1/query_1b00249be0d74936ab4ba1520e0dfa4a/`，按忽略规则不提交。

联网验证命令：

```sh
cargo test --locked -p trade-log-adapters --test query_chain live_public_account_acceptance -- --ignored --nocapture
```

可通过 ROBOTECH_ACCEPTANCE_ACCOUNT 指定其他公开实际交易账户。该验收要求有可用近期合约成交；实时结果变化，不作为普通离线测试前置条件。

## 双进程实际运行

使用两个 release 可执行文件分别加载临时 TOML，监听本地临时端口，通过临时生成的凭证通信：

- 网关 health 返回 ok，version 返回 0.1.0。
- 未授权的内部 health 请求返回 401。
- 经网关查询同一公开子账户，官方返回 2000 条，API 返回 10 条；coverage=LIMITED，无解析错误。
- 证据：`var/acceptance/v0.1-process/query_b7f1edfb227f477e923cd915b31db4dc/`。
- 两个进程均在 SIGTERM 后以退出码 0 结束，无遗留监听进程。

未连接生产服务器或修改其运行配置；没有将本地程序验证描述为容器验证。

## Docker 部署待验证

开发机未安装 Docker。已交付 Compose 的两个常驻服务和仅在 init profile 使用的一次性数据卷初始化服务，以及共用 Dockerfile 的两个镜像 target。

服务器部署入口：

```sh
sh scripts/init-v0.1.sh
docker compose up --build -d
docker compose ps
docker compose logs --tail=30 query-api trade-log-query
```

初始化脚本生成本地凭证且保留已有凭证，为证据卷设置 UID/GID 10001。secrets 与 var 被 Git 和 Docker 构建上下文忽略。启动前初始化是本版新增的部署步骤。

待验证镜像构建、凭证挂载可读、非 root 数据卷写入、内部服务无宿主机端口发布、证据重启保留、正常停止和部署后的真实账户 HTTP 查询。以下步骤用于重复执行这些检查。Docker 验收完成前不标记完整发布验收通过。

## 手动验收步骤

按 [开发文档第 10.4 节](version0.1.md#104-手动验证步骤) 顺序执行。每一步均列出命令、观察项和通过标准，包括真实查询、错误输入、公网访问、证据落盘、重启保留、依赖故障与恢复、正常停止。

验收时记录执行时间、代码版本、执行人、各步骤结果，以及失败时的状态码和日志。全部实际执行通过后，再将服务器 Docker 验收状态更新为通过。

默认永续 DEX 之外的市场会明确标为未支持；历史区间、数据库和实时监控留到后续版本。README 未改动。
