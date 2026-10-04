# v0.0 验收报告

- 日期：2026-10-04（Asia/Shanghai）
- 状态：本地自动验证、用户服务器手动验证及公网接口复核通过
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

开发机没有 Docker，未在开发机执行容器构建。用户已在服务器完成手动验证并确认正常；提供的日志显示容器启动、优雅停止、重启和健康请求成功。服务器 CPU 架构尚未采集。

## 服务器验证记录

用户于 2026-10-04 确认服务器手动验证全部正常，包括 Compose 启动、正常和错误接口、配置修改与错误配置、正常停止重启及公网访问。此部分为用户提供的结果，助手没有远程 shell 访问。

助手于 2026-10-04 19:11（Asia/Shanghai）直接请求 http://120.77.207.116:8080：

| 检查 | 实测结果 |
| --- | --- |
| GET /api/v1/health | 200，status=ok，schema_version=1 |
| GET /api/v1/version | 200，version=0.0.0 |
| GET /not-found | 404，RESOURCE_NOT_FOUND |
| POST /api/v1/health | 405，METHOD_NOT_ALLOWED；Allow: GET, HEAD |
| HEAD /api/v1/health | 200，无响应体 |
| 重复健康请求 | started_at 固定为 2026-10-04T11:08:57.170Z，trace ID 各不相同 |
| 请求关联 | JSON 响应头与响应体的 trace ID 一致 |

公网复核未修改配置或操作容器。非 root 身份、只读文件系统、缺失挂载和 CPU 架构未独立检查，不将这些部署属性记为实测通过。

## 可重复执行的 Docker 验收步骤

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

正常停止退出码应为 0。进一步按开发文档验证配置修改重启、Compose 环境及端口修改后的容器重建、错误配置、缺失挂载、端口冲突、非 root 身份和只读文件系统。如需限定 Linux amd64 发布支持，补充记录服务器 CPU 架构；Apple Silicon 可额外验证原生 arm64。

README 未修改。运行说明保留在本版开发文档和本报告中。
