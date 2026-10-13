# 验收命令

在服务器项目根目录，按顺序逐条执行；出现 FAIL 时先处理，再继续。v0.1、v0.2 的示例账户可替换；v0.3～v0.5 默认读取首个登记账户，可用 --account 指定；v0.6、v0.7 检查全部登记账户。

脚本在当前部署上使用新功能入口；验收旧版本部署时保留对应旧入口。执行命令不变。

以下命令包含重启和停库测试，会短暂中断服务。

```sh
# v0.0
sh scripts/verify-v0.0.sh --lifecycle

# v0.1
sh scripts/verify-v0.1.sh --account 0x010461c14e146ac35fe42271bdc1134ee31c703a --lifecycle

# v0.2
sh scripts/verify-v0.2.sh --account 0x010461c14e146ac35fe42271bdc1134ee31c703a --live --lifecycle --database-fault --wait-seconds 180

# v0.3
sh scripts/verify-v0.3.sh --live --lifecycle --database-fault --wait-seconds 180

# v0.4
sh scripts/verify-v0.4.sh --live --lifecycle --database-fault --wait-seconds 180

# v0.5
sh scripts/verify-v0.5.sh --fixture-tests --lifecycle --database-fault --wait-seconds 180

# v0.6
sh scripts/verify-v0.6.sh --fixture-tests --lifecycle --database-fault --wait-seconds 180

# v0.7
sh scripts/verify-v0.7.sh --fixture-tests --lifecycle --database-fault --wait-seconds 180
```

- `--fixture-tests`：v0.5～v0.7 独立数据库固定样本测试，首次会构建测试镜像。
- `--account`：v0.1、v0.2 必填，指定验收账户。
- `--live`：查询官方来源并保存结果；v0.1 默认执行，v0.7 可选。
- `--lifecycle`：停止并重新启动对应应用服务。
- `--database-fault`：停止数据库 15 秒，再启动并检查恢复。
- `--wait-seconds 180`：采集或恢复最多等待 180 秒。

不允许中断服务时，去掉 `--lifecycle` 和 `--database-fault`。

结果：PASS 通过，FAIL 失败，SKIP 未验证。
