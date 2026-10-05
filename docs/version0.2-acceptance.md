# v0.2 本地验收记录

验证日期：2026-10-05。验证对象：本版工作区代码、两个 release 程序及本机 PostgreSQL 18.0。

## 1. 已完成验证

| 验证项 | 结果 |
| --- | --- |
| workspace 测试 | 21 项通过；数据库测试另行执行，实时外部来源测试未执行 |
| PostgreSQL 集成测试 | 13 项通过，使用独立测试数据库 |
| 幂等与事务 | 重复和并发写入不重复事实；同身份内容冲突回滚整批，原始响应保留 |
| 角色权限与迁移 | 迁移可重复执行；迁移角色初始化成功；应用角色不能执行迁移或删除表记录 |
| 库存查询 | 时间边界、超大 tid 排序、固定快照分页、游标条件校验通过 |
| 导入与重放 | 固定旧证据样本重复导入通过，损坏摘要被拒绝，中断导入可恢复，reparse 不修改事实 |
| HTTP 链路 | 网关、内部接口与数据库完整链路通过；参数错误 400、未授权 401、数据库不可用 503 |
| 保存的真实来源记录 | v0.1 留存证据中 2,000 条成交导入成功，重新解析 SAME；release HTTP 库存查询匹配 2,000 条 |
| 数据库重启 | 两个 release 程序运行期间停止数据库，库存查询返回 503、网关健康接口仍为 200；数据库重启后恢复查询，记录保留 |
| 优雅退出 | 两个 release 程序收到 SIGTERM 后退出码 0，日志包含 shutdown_completed |
| 代码与部署文件检查 | fmt、clippy（warnings 为错误）、release 构建、初始化脚本语法、Compose YAML 和差异空白检查通过 |

旧证据的 normalized_account 缺失时从原请求推导；若提供了该字段但与请求不一致，拒绝导入。真实记录验证使用保存的来源字节，未将构造样本作为真实交易证据。

## 2. 可重复执行的开发检查

```sh
cargo fmt --all -- --check
cargo clippy --locked --offline --workspace --all-targets -- -D warnings
cargo test --locked --offline --workspace
cargo build --locked --offline --release -p query-api -p trade-log-query
sh -n scripts/init-v0.2.sh
```

数据库集成测试需要一个可连接且有建库权限的本机 PostgreSQL 管理员连接。设置 ROBOTECH_TEST_DATABASE_URL 后执行：

```sh
cargo test --locked --offline -p trade-log-adapters \
  --test postgres_persistence --test stored_query --test replay_import \
  -- --ignored --nocapture
```

数据库测试创建隔离测试库；不要用生产应用凭证作为测试管理员连接。以上是开发检查，服务器手动验收按 [v0.2 文档第 10.4 节](version0.2.md#104-手动验证步骤) 执行，逐步观察响应与数据库记录。

## 3. 尚需服务器验证

本机没有 Docker，未执行实际镜像构建、Compose 启动及容器卷重启验证。原生程序和数据库测试不替代容器验收。本次也未执行新的外部实时采集；升级后需按手动步骤验证来源访问、完整持久化回执以及容器重启后的库存。

```sh
sh scripts/init-v0.2.sh
docker compose up --build -d
```

本次没有远程部署，没有修改 README、roadmap 或详细设计说明，没有新增自动验收脚本。
