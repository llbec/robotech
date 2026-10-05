# v0.3 本地验收记录

验证日期：2026-10-06。验证对象：当前工作区代码、本机 PostgreSQL 18.0，以及 query-api、trade-log-query、trade-collector 三个 release 程序。

## 1. 已完成验证

| 项目 | 结果 |
| --- | --- |
| workspace 常规测试 | 26 项通过；真实外部来源的既有 ignored 测试未在该命令中执行 |
| PostgreSQL 集成测试 | 24 项通过，使用独立测试数据库 |
| 水位与幂等 | 无成交成功推进扫描水位；重叠采集和手动查询不重复成交；配置起点不覆盖已有水位 |
| 事务一致性 | 事实冲突整批回滚；在事实写入之后令检查点更新失败，事实、水位、成功任务一起回滚；重试提交成功；已提交事务重试不重复入账 |
| 来源边界 | 时间转换为包含边界的官方毫秒请求；2,001 条固定成交拆分完整保存；单毫秒满页明确失败且水位不变 |
| 有界工作与恢复 | 请求预算暂停后保留剩余区间；新实例从同一任务继续；多页重新解析 SAME |
| 租约隔离 | 同账户第二实例不能领取活跃租约；旧令牌不能归档或提交；查询进程启动不误中断自动采集任务 |
| 自动调度 | 固定来源 429 的错误、连续失败与退避可见；Retry-After 写入下一次调度；恢复成功清空错误与失败计数，连续生成成功轮次 |
| 取消与重启 | 采集请求中取消时保留 pending_work，释放租约；新实例从保存的元数据和任务继续 |
| HTTP 状态接口 | 内部未授权 401；状态外部响应与 trace 对应；不支持的方法 405；数据库故障状态 503、网关 health 200 |
| v0.2 升级与权限 | 单独创建 0001 schema 并保存成交，再执行 0002；旧事实、身份及重放保留；collector 运行角色能采集但不能建表或删除检查点 |
| 实际 release 程序 | 三个原生进程使用分离的应用、迁移、采集角色连接；官方 HTTP 自动采集至少两轮成功，对自动任务 reparse 输出 SAME |
| 职责独立 | 停止 collector 时状态接口 503、stored 查询仍 200；重启 collector 后原水位不回退 |
| 优雅退出 | 三个 release 程序收到 SIGTERM 后退出码 0，日志包含 shutdown_completed |
| 构建和静态检查 | fmt、clippy（warnings 为错误）、三个 release binary 构建、脚本语法、Compose YAML、文档引用及差异空白检查通过 |

真实来源验证使用官方公开账户接口，不使用构造样本冒充真实成交。成功轮询与重新解析证明来源访问及自动保存链路可运行；没有据此声称验收期间出现了新的成交或官方历史已完整。新成交发现、满页、事务中断等机制分别通过固定来源测试验证。

## 2. 可重复执行的开发检查

```sh
cargo fmt --all -- --check
cargo clippy --locked --offline --workspace --all-targets -- -D warnings
cargo test --locked --offline --workspace
cargo build --locked --offline --release -p query-api -p trade-log-query -p trade-collector
sh -n scripts/init-v0.3.sh
```

数据库测试需要可建库的 PostgreSQL 测试管理员连接。设置 ROBOTECH_TEST_DATABASE_URL 后执行：

```sh
cargo test --locked --offline -p trade-log-adapters \
  --test checkpoint --test collector_chain --test postgres_persistence \
  --test stored_query --test replay_import -- --ignored --nocapture
```

测试创建隔离数据库，不使用生产业务连接。未新增仓库自动验收脚本。

## 3. 待服务器手动验收

本机没有 Docker，未实际执行镜像构建、Compose 启动及容器数据卷重启。本机原生程序结果不代替 Docker 服务器验收，也不证明长期运行中的来源覆盖完整。

先填写 config/trade-collector.toml 中的账户和 RFC3339 起点，再在仓库根目录执行：

```sh
docker compose stop query-api trade-log-query
sh scripts/init-v0.3.sh
docker compose up --build -d
curl --noproxy '*' -i http://127.0.0.1:8080/api/v1/watch-accounts
```

逐步查看响应、水位和数据库记录，按 [v0.3 文档第 10.2 节](version0.3.md#102-手动验证步骤) 完成验收。没有执行远程部署，没有修改 README、roadmap 或详细设计说明。
