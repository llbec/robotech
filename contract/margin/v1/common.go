package v1

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func ABIEncode(str string) [32]byte {
	stringType := abi.Type{
		T: abi.StringTy,
	}
	arguments := abi.Arguments{
		{Type: stringType},
	}

	packedData, _ := arguments.Pack(str)

	var hash [32]byte
	copy(hash[:], crypto.Keccak256Hash(packedData).Bytes())

	return hash
}

func PositionEncode() []byte {
	stringType := abi.Type{
		T: abi.StringTy,
	}
	arguments := abi.Arguments{
		{Type: stringType},
	}

	packedData, _ := arguments.Pack("POSITION")

	return crypto.Keccak256Hash(packedData).Bytes()
}

func CalcPositionKey(account common.Address, id *big.Int) []byte {
	// 定义参数类型
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	addressType, _ := abi.NewType("address", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)

	// 创建 Arguments 切片
	arguments := abi.Arguments{
		{Type: bytes32Type},
		{Type: addressType},
		{Type: uint256Type},
	}

	var hash [32]byte
	copy(hash[:], PositionEncode())

	// 使用 Pack 方法进行编码
	packedData, _ := arguments.Pack(hash, account, id)

	return crypto.Keccak256Hash(packedData).Bytes()
}
