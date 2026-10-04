# v0.0 验收报告

- 日期：2026-10-04（Asia/Shanghai）
- 状态：本地程序验证通过，Docker 发布验收待完成
- 基准提交：`99a77e50895c9ee99d2ec969a49c10d1c380bb11`，本次实现尚未提交
- 验证环境：macOS，aarch64；Rust/Cargo 1.96.1
- 实现入口：`gateway/query-api`，程序版本 `0.0.0`

## 已完成验证

| 检查 | 结果 |
| --- | --- |
| `cargo build --locked -p query-api` | 通过 |
| `cargo build --locked --release -p query-api` | 通过 |
| release 程序的 `--help` 与 `--version` | 通过，版本为 `query-api 0.0.0` |
| `cargo test --locked --workspace` | 7 项集成测试通过 |
| `cargo fmt --all -- --check` | 通过 |
| `cargo clippy --locked --workspace --all-targets -- -D warnings` | 通过，无警告 |
| Compose YAML 语法解析 | 使用 Ruby YAML 通过；不能替代 Compose 语义检查 |
| `git diff --check` | 通过 |

集成测试源码位于 `tests/integration-tests/query_api.rs`，通过 package 的 `[[test]]` 显式注册。测试使用实际本地 TCP socket 和临时端口，不依赖外部业务系统。

- 配置：字段完整性、未知字段、版本、IPv4/IPv6、端口与退出时限边界、日志配置、环境覆盖及空值/非法值。
- HTTP：健康与版本 JSON、Content-Type、schema、UTC 启动时间稳定性、版本一致性、trace ID 唯一且不信任客户端覆盖。
- 错误：404、405、Allow、HEAD 空响应体及错误 trace 关联。
- 启动：CLI 帮助与版本无需有效配置；错误参数、缺失配置、非法环境值、端口占用的退出码及失败日志。
- 生命周期：在途请求完成、有界退出、超时取消 handler、退出后端口重新绑定。
- 实际进程：SIGTERM 和 SIGINT 都返回 0；JSON 启动、退出及请求日志可解析，日志与响应 trace ID 一致；4xx 使用 WARN；URL 参数、未知路径和 Authorization 内容不写入日志。

本环境 sandbox 禁止部分本地监听操作，集成测试在获得允许后执行；未将权限失败当作程序测试通过。

## Docker 交付与验证边界

已交付根 `compose.yaml`、`.dockerignore` 和 `gateway/query-api/Dockerfile`。镜像为多阶段构建，运行阶段非 root；程序通过 exec ENTRYPOINT 接收信号；配置只读挂载；Compose 使用只读根文件系统及 15 秒停止宽限时间，应用默认退出时限为 10 秒。

官方 Docker Hub 标签元数据已核对，构建与运行镜像按多架构 manifest 摘要固定，两者均提供 Linux amd64 与 arm64：

- Rust：`rust:1.96.1-bookworm@sha256:a339861ae23e9abb272cea45dfafde21760d2ce6577a70f8a926153677902663`
- Debian：`debian:bookworm-slim@sha256:3783cc01769c7b2b1b83a5c5ad96c815348e28ed7da68e2e3687004faa906251`

本机未安装 Docker CLI/运行时，未执行镜像构建、`docker compose config` 或容器启动。因此本版尚未通过完整 Docker 发布验收，也未验证 Linux amd64 或 arm64 容器运行。

## Docker 环境中的待验收步骤

从仓库根目录执行：

```sh
docker compose config
docker compose up --build -d
docker compose ps
docker compose logs query-api
curl -i http://127.0.0.1:8080/api/v1/health
curl -i http://127.0.0.1:8080/api/v1/version
curl -i http://127.0.0.1:8080/unknown
curl -i -X POST http://127.0.0.1:8080/api/v1/health
curl -I http://127.0.0.1:8080/api/v1/health
docker compose stop
docker inspect --format '{{.State.ExitCode}}' "$(docker compose ps -aq query-api)"
docker compose start query-api
docker compose down
```

正常停止退出码应为 0。进一步按开发文档验证配置修改重启、Compose 环境及端口修改后的容器重建、错误配置、缺失挂载、端口冲突、非 root 身份和只读文件系统。在 Linux amd64 环境执行并记录结果后才能完成基准发布验收；Apple Silicon 可额外验证原生 arm64。

README 未修改。运行说明保留在本版开发文档和本报告中。
