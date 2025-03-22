# compile abigen
```
git clone https://github.com/ethereum/go-ethereum.git
cd go-ethereum/cmd/abigen
gobuild
mv abigen $GOPATH/bin/.
```

# generate abi file
find `abi=[...]` in `*.json`,copy `[...]`  
paste `[...]` to `*.abi`

# generate abi.go
```
abigen --abi rawore/aave/protocol/lendingpool/lendingpool.sol/lendingpool.abi --pkg aave --type LendingPool --out armory/aave/LendingPool.go

--abi: Mandatory path to the contract ABI to bind to
--pkg: Mandatory Go package name to place the Go code into
--type: Optional Go type name to assign to the binding struct
--out: Optional output path for the generated Go source file (not set = stdout)
```

## abigen
- 安装 
```bash
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```
- 添加环境变量
```bash
export PATH="$HOME/go/bin:$PATH"
```

- 生成abi文件
```bash
solc --abi Example.sol -o build
```

- 生成go文件
```bash
abigen --abi Example.abi --pkg main --type Example --out Example.go
```
--abi 目标abi文件
--pkg go的包名
--type go的结构名
--out 指定生成文件


## 参考文档
[智能合约的编译与ABI](https://goethereumbook.org/zh/smart-contract-compile/)
[go-ethereum tool abigen](https://geth.ethereum.org/docs/tools/abigen)