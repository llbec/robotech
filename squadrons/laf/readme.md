# LAF
## staking
### stake 用户质押X的USDT
<details>
<summary>
  1. 从用户地址转账X个USDT到staking合约
</summary>

  > 日志 usdt transer 事件: 用户->staking, x
</details>
<details>
<summary>
  2. 将一半的usdt通过swap换成LAF
</summary>

  > 日志 usdt transfer 事件: staking->Pair, x/2
  > 日志 LAF transfer 事件：Pair->staking, x/2=>laf
  > 日志 Pair sync 事件
  > 日志 Pair swap 事件：router->staking
</details>
<details>
<summary>
  3. 将收到的LAF与适量的USDT添加流动性
</summary>

  > 日志 usdt transfer 事件： staking->Pair, x/2=>laf=>u
  > 日志 LAF transfer 事件：staking->pair, x/2=>laf
  > 日志 pair transfer 事件：0->fee, lp fee
  > 日志 pair transfer 事件：0->0, lp
  > 日志 pair sync 事件
  > 日志 pair mint 事件：0->router, [x/2=>laf,x/2=>laf=>u]
</details>
<details>
<summary>
  4. 记录用户绑定关系
</summary>

  > 日志 referral bind 事件: 用户, up
</details>
<details>
<summary>
  5. 记录用户质押信息
</summary>

  > 日志 staking transfer 事件: 0->用户, x
  > 日志 staking staked 事件: 用户, x, ts, index, staketime
</details>

### unstake 用户取消质押X的USDT
