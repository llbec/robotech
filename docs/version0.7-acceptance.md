# v0.7 开发验证及验收记录

验证日期：2026-10-13。程序版本：0.7.0。

## 1. 本机验证

| 项目 | 结果 | 范围 |
| --- | --- | --- |
| 工作区 Rust 测试 | PASS：46 项 | 现货映射、精确金额、费用、别名去重、旧游标与原接口回归 |
| PostgreSQL 18 集成回归 | PASS：40 项 | 包含原 39 项数据库回归和 1 项现货完整链路测试 |
| 现货专项最终复验 | PASS：1 项 | 与上行现货测试重叠，不重复计数；增加了真实 HTTP 接收器和独立重新激活边界验证 |
| Python 验收辅助测试 | PASS：21 项 | 计数版本兼容、固定测试资源隔离及异常清理 |
| 编译与 Clippy | PASS | workspace 编译；all-targets，warnings 视为错误 |
| 实际服务进程 | PASS | 四服务启动；分类 stored 查询；非法类型 400；账户添加/暂停；补录 CLI 提交、幂等、状态；collector 重启保留暂停与补录任务 |
| 脚本、格式及文档 | PASS | Shell 语法、Python 语法、Rust 格式及 git diff --check |

现货固定样本及数据库测试确认：

1. 可读交易对名称和 @index 归为同一事实，同身份不同内容明确冲突。
2. BUY/SELL、基础资产手续费、计价资产返佣、非 USDC 计价和精确成交额符合固定案例。
3. 未确认费用币种、强制事件、快照和现货发布关闭时不生成候选。
4. 默认 live 展示合约；SPOT 展示现货；两类完整事实均入库。stored 在分类后分页，跨类型游标拒绝；旧游标默认 PERPETUAL。
5. 已完成旧任务的原始证据可补录现货；原任务回执、永续事实及原 HTTP 水位不改写。
6. 补录先扫描旧证据再执行指定时间范围的来源扫描；重复提交复用任务；暂停后的旧代次提交被拒绝，修复 BLOCKED 后可重新提交相同任务。
7. HTTP/补录不产生历史候选；新实时现货可以进入 outbox。首次投递补齐 published_at 后，HTTP 接收器收到的正文与持久化 wire_body 一致。
8. 现货重新启用递增自身激活代次，激活前成交被抑制，原永续激活边界保持。

固定数据位于 tests/fixtures/hyperliquid/spot，是可重复测试样本，不是对实际钱包的成交核验记录。

## 2. 服务器待验

本机未安装 Docker，未执行 Docker 镜像构建或 Compose 验收。真实账户现货、费用币种和 webhook 接收记录仍待服务器核对，不能将固定样本通过视为真实来源通过。

在项目根目录执行：

```sh
sh scripts/init-v0.7.sh
docker compose up --build -d
sh scripts/verify-v0.7.sh --fixture-tests --lifecycle --database-fault --wait-seconds 180
```

初始化保留原卷、账户和凭证。固定样本在独立测试库执行并清理该测试资源；不会向监控账户注入伪造成交。重启和停库参数会短暂中断服务。

脚本默认读取所有登记账户并查询库存，可加 --live 查询官方来源并保存新证据。没有实际现货或实时候选时对应项目 SKIP。使用本地验收接收器且存在已送达现货时，脚本核对 event_id 和正文摘要；外部接收方消费情况需另行核对。

需要验证真实发布时，先确保接收方支持 SPOT，然后明确设置 publishing.spot_enabled=true，同时保持 publishing.enabled=true。配置重启只开启后续合格实时候选，不追发旧现货。

## 3. 使用变更

- 成交查询新增 instrument_type=PERPETUAL/SPOT/ALL，默认 PERPETUAL；实时与库存功能仍由 /live、/stored 路径区分。
- 新响应 counts.counts_version=2，spot_records 表示成功标准化现货；历史回执不改写。完整持久化数量包含永续及现货。
- 现货费用按实际资产记录，不从原始数量或成交额中扣除，也不换算成 USDC。
- 默认不发布现货；能力边界和状态位于分账户发布状态的 instruments.SPOT。
- 补录进度独立保存，不回退原自动采集水位。COMPLETED 只说明任务扫描完成，source_coverage 仍不证明完整历史。

完整接口参数见 [api.md](api.md)，部署和手动补录见 [version0.7.md](version0.7.md)，服务器验收命令见 [verify.md](verify.md)。
