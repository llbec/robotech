# 测试环境存储清理

在服务器项目根目录执行，按以下顺序逐步操作。

**本流程会永久删除本项目全部测试数据。代码、config 和 secrets 保留。**

## 1. 清理数据库、证据文件和本项目容器

```sh
docker compose --profile "*" down --volumes --remove-orphans
```

删除内容：

- PostgreSQL 卷 `trade-log-postgres`：成交、原始响应、采集水位、候选、发布队列及投递记录。
- 文件证据卷 `trade-log-evidence`。
- 验收接收器记录卷 `webhook-verification`。
- 本项目容器、容器输出日志和网络。

数据删除后不可恢复。配置目录和凭证目录是宿主机文件，不随上述数据卷删除。

## 2. 清理其他已停止的容器

```sh
docker container prune
```

根据提示确认。此命令影响整台服务器，删除所有已停止的容器，包括其输出日志和容器可写层；不删除运行中的容器或命名数据卷。

## 3. 清理构建缓存和旧镜像

先清理构建缓存：

```sh
docker builder prune -a
```

再清理悬空镜像：

```sh
docker image prune
```

如果空间仍不足，可进一步删除所有没有容器引用的镜像：

```sh
docker image prune -a
```

这些命令影响整台服务器的 Docker 缓存。清理后，后续部署可能需要重新下载镜像、依赖或重新编译。

## 4. 检查剩余空间及其他占用

```sh
df -h
df -i
docker system df
```

检查系统日志、下载缓存、备份及编译产物的大小：

```sh
du -xhd1 /var/log /var/cache /root 2>/dev/null
du -sh target 2>/dev/null
```

以上仅查看占用，不自动删除这些文件。不要直接删除 `/var/lib/docker/overlay2` 或 PostgreSQL 的 `pg_wal`。

## 5. 配置日志轮转并重新启动

重新创建容器前，在 `compose.yaml` 各服务中加入以下配置，与 image、volumes 等字段同级：

```yaml
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"
```

每个容器保留约 30MB 输出日志。日志配置需要重新创建容器才能生效，仅 restart 不够；下面的启动步骤会创建已删除的容器。

调整 `config/trade-collector.toml` 中的 `start_time`，设置需要重新采集的起点，避免再次回补大量旧数据。

```sh
sh scripts/init-v0.5.sh
docker compose up --build -d
```

需要验收接收器时，再执行：

```sh
docker compose --profile verify-webhook up -d webhook-receiver
```

检查数据库及空间：

```sh
docker compose ps
docker compose logs postgres --tail=20
df -h
```

日志轮转只限制容器输出日志。数据库、原始证据和投递记录仍会增长，后续需要单独设计保留期限和定期清理机制。

命令范围参考：[Compose down](https://docs.docker.com/reference/cli/docker/compose/down/)、[容器清理](https://docs.docker.com/reference/cli/docker/container/prune/)、[镜像清理](https://docs.docker.com/reference/cli/docker/image/prune/)、[日志轮转](https://docs.docker.com/engine/logging/drivers/json-file/)。
