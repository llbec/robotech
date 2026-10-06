# v0.5 验证记录

- 日期：2026-10-06
- 程序版本：0.5.0；新增 migration：0004_webhook_publishing.sql
- 环境：macOS 本机原生 Rust 程序、独立 PostgreSQL 18 测试库、标准库 Python/SQLite 接收器
- 结论：本机开发验证通过；本机没有 Docker，服务器容器部署及真实 Hyperliquid 候选投递仍待验收。

## 1. 开发验证

| 验证 | 结果 | 覆盖 |
| --- | --- | --- |
| 工作区测试 | 38 项通过 | 既有接口、解析、生命周期，以及新增候选策略、配置和响应分类；外部来源及数据库测试另行执行 |
| PostgreSQL 回归 | 36 项通过 | 原事实、库存、重放、水位、WS 会话、队列故障、续租并发和新增 outbox；各测试使用独立测试库 |
| v0.5 固定数据库及接收器套件 | 2 项通过 | HTTP 先入库、快照后无标记更新、抑制原因、原子回滚、固定正文、失效租约、500/429/401、确认丢失与去重；实际发布程序重启及 retry 命令 |
| 验收脚本测试 | 14 项通过 | PASS/FAIL/SKIP 判定、非法正文摘要、停服失败后恢复、独立资源清理 |
| 静态检查 | 通过 | cargo fmt、clippy 全工作区所有 target 且禁止 warning、shell 语法、Compose YAML 解析及 git diff 空白检查 |

可重复执行的开发命令：

```sh
cargo test --locked --offline --workspace
cargo clippy --locked --offline --workspace --all-targets -- -D warnings
PYTHONDONTWRITEBYTECODE=1 python3 -m unittest discover -s scripts/tests -v

# 使用独立 PostgreSQL 管理员测试连接，不要填写部署数据库。
cargo build --locked --offline -p trade-parser-publisher
ROBOTECH_TEST_DATABASE_URL='<独立测试 PostgreSQL URL>' \
  cargo test --locked --offline -p trade-log-adapters --test publishing_chain -- --ignored --nocapture
```

常规数据库开发测试会创建随机测试库；服务器 `--fixture-tests` 入口负责创建并清理自己的单个测试库、接收器和发布进程，不向部署库注入样本。

## 2. 实际进程验证

本机以独立数据库启动真实 `trade-parser-publisher`、`query-api` 和 Python 接收器，验证了：

1. 内部状态要求服务凭证，外部状态保持 data/meta 和 trace 契约。
2. 持久事件由发布程序投递并确认；发布进程停止时外部状态返回 503，网关健康仍正常。
3. 同配置重启后 activated_at 和 activation_epoch 保持，待发送事件继续交付。
4. 只中断该测试库的连接，状态返回 503；恢复后数据库和发布状态可访问，事件保留。
5. 显式加载 disabled 配置后状态为 DISABLED。
6. 发布角色可更新 outbox，不能更新标准事实或 collection_checkpoints。

本机样本事件、库、端口和 token 均为隔离测试资源。上述结果不登记为真实账户或服务器 Docker 验收通过。

## 3. 服务器待验收

执行顺序和参数见 [verify.md](verify.md)，部署及真实信号核对见 [v0.5 第 10 章](version0.5.md#10-交付部署与手动验证)。

```sh
sh scripts/init-v0.5.sh
docker compose up --build -d
sh scripts/verify-v0.5.sh --fixture-tests --lifecycle --database-fault --wait-seconds 180
```

发布默认关闭；真实投递前填写目标并启用，验收接收器需启动 verify-webhook profile。固定套件可以独立验证候选和故障行为；没有真实新成交、未启用发布或缺少接收端证据的项目必须保留 SKIP。

待记录：服务器代码版本、激活代次和时间、新鲜主动成交的 WS 原文及 snapshot_sequence、event_id、投递尝试和接收摘要、故障恢复结果。固定样本与真实来源分别登记。
