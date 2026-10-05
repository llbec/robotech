# v0.4 本地验收记录

验证日期：2026-10-06。验证对象：当前工作区代码、本机隔离 PostgreSQL 18.0、官方公开 HTTP/WSS，以及 query-api、trade-log-query、trade-collector 三个 release 程序。

## 1. 已完成验证

| 项目 | 结果 |
| --- | --- |
| workspace 常规测试 | 34 项通过；需要数据库和真实来源的 ignored 测试另行执行 |
| PostgreSQL 集成测试 | 32 项通过，使用独立测试数据库；含 8 项 v0.4 双通道测试 |
| 来源协议 | 本地真实 WS socket 验证订阅参数、确认、应用层 ping/pong、超限消息；快照、增量、未知模式、空 fills、错账户和畸形消息明确区分 |
| 重连调度 | 验证 5、10、20、40、60 秒基础退避、0～20% 抖动、最终 60 秒上限；稳定订阅且心跳正常满 60 秒才重置；滚动窗口最多 10 次尝试；worker 重建保留窗口和未结束退避 |
| 跨通道去重 | HTTP、WS 更新、重复推送与重连快照共享 fact_id；current 和 revision 不重复，来源观察分别保留；金额等价文本及无关原始字段不会误报冲突 |
| 历史兼容 | 保留原 payload/content_hash；旧事实首次比较时补充新业务指纹；旧 schema 数据迁移后身份、库存查询、重新解析及运行角色权限保留 |
| 事务一致性 | WS 真实业务冲突整批回滚；旧租约不能归档、提交或修改 WS 状态；HTTP 缺口更新失败时事实、水位和完成标记一起回滚 |
| 水位与缺口 | 收到较新 WS 成交不推进 HTTP 水位；HTTP 完整范围提交后才标记 HTTP_SCANNED；关闭 WS 后已有 HTTP 恢复范围继续扫描 |
| 自动调度与心跳 | 自动双通道断线后重新订阅、快照去重并补偿；没有成交时应用层 pong 持续更新；订阅超时、pong 缺失明确记录恢复需求 |
| 背压 | 阻塞事实提交，确定触发有界队列溢出；关闭 socket 后恢复已归档消息及 HTTP 范围，预置成交全部保存；等待数据库/提交能力恢复后才建立新连接 |
| 元数据与证据 | 完整原始 envelope 保留；缺失市场映射时记录待处理消息，追加元数据快照后解析；元数据及消息摘要可复核；原文篡改不能通过重新解析 |
| 重新解析 | WS 快照、增量、空消息及跨通道重复任务 reparse 为 SAME；旧 HTTP 与旧证据导入回归通过，重新解析不推进水位 |
| 状态接口 | 网关透传 monitoring_status/websocket/recovery；内部认证、trace、405 和故障映射保留；数据库不可读时状态 503、网关 health 200 |
| 官方 WSS | 实际 TLS 连接、userFills 订阅确认及应用层 ping/pong 通过，使用只读公开账户 |
| 实际 release 链路 | 三个原生进程使用分离数据库角色；官方 HTTP 至少两轮成功，WS 快照入库且 pong 正常，恢复目标完成；HTTP 与 WS 任务 reparse 均 SAME |
| 实际数据库停机恢复 | 短暂停止隔离 PostgreSQL，确认状态 503/health 200；恢复后建立新 WS 会话、HTTP 水位继续、旧事实保留 |
| 重启与退出 | 停止 collector 后状态 503、stored 200；重启不回退水位；三个程序 SIGTERM 后退出码 0，日志含 shutdown_completed |
| 静态与构建检查 | fmt、clippy（warnings 为错误）、三个 release binary 构建、初始化脚本语法、Compose YAML、文档 JSON/引用及差异空白检查通过 |

真实来源验收观察到订阅快照及正常 HTTP 扫描，不声称验收期间发生了新成交，也不声称官方无限历史完整。新成交及时入库、断线缺口补齐和队列溢出分别通过预置固定来源验证；LIVE 和 HTTP_SCANNED 均表示运行/扫描状态，coverage 仍为 SOURCE_HISTORY_NOT_VERIFIED。

## 2. 可重复执行的开发检查

```sh
cargo fmt --all -- --check
cargo clippy --locked --offline --workspace --all-targets -- -D warnings
cargo test --locked --offline --workspace
cargo build --locked --offline --release -p query-api -p trade-log-query -p trade-collector
sh -n scripts/init-v0.4.sh
```

使用可建库的隔离 PostgreSQL 测试管理员连接设置 ROBOTECH_TEST_DATABASE_URL，再执行：

```sh
cargo test --locked --offline -p trade-log-adapters \
  --test websocket_chain --test checkpoint --test collector_chain \
  --test postgres_persistence --test stored_query --test replay_import \
  -- --ignored --nocapture
```

需要网络时，可独立执行只读官方 WSS 协议验收：

```sh
cargo test --locked --offline -p hyperliquid --test websocket \
  official_wss -- --ignored --nocapture
```

测试连接不使用生产业务数据库。未新增仓库自动验收脚本。

## 3. 待服务器 Docker 手动验收

本机没有 Docker，未执行实际镜像构建、Compose 启动和容器数据卷重启。原生程序验证不能代替 Docker 部署验证。

填写 config/trade-collector.toml 的实际账户及 RFC3339 起点，保留已有凭证和数据卷，再执行：

```sh
docker compose stop query-api trade-log-query trade-collector
sh scripts/init-v0.4.sh
docker compose up --build -d
docker compose ps -a
docker compose logs --tail=50 trade-collector
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/version
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/watch-accounts
```

随后按 [v0.4 第 10.2 节](version0.4.md#102-手动验证步骤) 的“命令—观察内容—通过标准”，检查消息模式、跨通道去重、安静账户心跳、进程断线、数据库故障、重新解析以及关闭 WS 后的 HTTP 兼容。升级前备份数据库；不用 down -v 做普通升级或重启验证。

没有执行远程部署，没有修改 README、roadmap、详细设计或历史版本文档。
