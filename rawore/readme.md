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
