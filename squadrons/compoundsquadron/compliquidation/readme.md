# Comptroller
## Get Assets In
Get the list of markets an account is currently entered into. In order to supply collateral or borrow in a market, it must be entered first. Entered markets count towards account liquidity calculations.  

查询地址流动性的资产列表。该地址在这个资产中有抵押品或者贷款。  

```
function getAssetsIn(address account) view returns (address[] memory)
```
- `account`: The account whose list of entered markets shall be queried.
- `RETURN`: The address of each market which is currently entered into.

## Get Account Liquidity
Account Liquidity represents the USD value borrowable by a user, before it reaches liquidation. Users with a shortfall (negative liquidity) are subject to liquidation, and can’t withdraw or borrow assets until Account Liquidity is positive again.

For each market the user has entered into, their supplied balance is multiplied by the market’s collateral factor, and summed; borrow balances are then subtracted, to equal Account Liquidity. Borrowing an asset reduces Account Liquidity for each USD borrowed; withdrawing an asset reduces Account Liquidity by the asset’s collateral factor times each USD withdrawn.

Because the Compound Protocol exclusively uses unsigned integers, Account Liquidity returns either a surplus or shortfall.  

查询地址的流动性。流动性和差额至少一个为0.  

```
function getAccountLiquidity(address account) view returns (uint, uint, uint)
```
- `account`: The account whose liquidity shall be calculated.
- `RETURN`: Tuple of values (error, liquidity, shortfall). The error shall be 0 on success, otherwise an error code. A non-zero liquidity value indicates the account has available account liquidity. A non-zero shortfall value indicates the account is currently below his/her collateral requirement and is subject to liquidation. At most one of liquidity or shortfall shall be non-zero.  

## Close Factor
The percent, ranging from 0% to 100%, of a liquidatable account's borrow that can be repaid in a single liquidate transaction. If a user has multiple borrowed assets, the closeFactor applies to any single borrowed asset, not the aggregated value of a user’s outstanding borrowing. 

查询单次清算可以偿还债务的比例

```
function closeFactorMantissa() view returns (uint)
```
- `RETURN`: The closeFactor, scaled by 1e18, is multiplied by an outstanding borrow balance to determine how much could be closed.

## Liquidation Incentive
The additional collateral given to liquidators as an incentive to perform liquidation of underwater accounts. For example, if the liquidation incentive is 1.1, liquidators receive an extra 10% of the borrowers collateral for every unit they close.
```
function liquidationIncentiveMantissa() view returns (uint)
```
- `RETURN`: The liquidationIncentive, scaled by 1e18, is multiplied by the closed borrow amount from the liquidator to determine how much collateral can be seized.