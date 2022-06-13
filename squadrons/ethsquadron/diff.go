package ethsquadron

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/types"
)

func (eth *EthSquadron) BlockNumber() (uint64, error) {
	return eth.node.BlockNumber(context.TODO())
}

func (eth *EthSquadron) HeaderByNumber(height uint64) (*types.Header, error) {
	return eth.node.HeaderByNumber(context.TODO(), new(big.Int).SetUint64(height))
}

func (eth *EthSquadron) DiffPerMilliSecond(height, number uint64) (*big.Int, error) {
	diffs := []*big.Int{}
	x := new(big.Int)

	startH := height - number
	preHeader, err := eth.HeaderByNumber(startH)
	if err != nil {
		return nil, fmt.Errorf("HeaderByNumber error@(%v): %v", startH, err)
	}

	for i := 1; i <= int(number); i++ {
		curH := startH + uint64(i)
		header, err := eth.HeaderByNumber(curH)
		if err != nil {
			return nil, fmt.Errorf("HeaderByNumber error@(%v): %v", curH, err)
		}
		deltaS := new(big.Int).SetUint64(header.Time - preHeader.Time)
		dpms := new(big.Int).Div(preHeader.Difficulty, deltaS.Mul(deltaS, big.NewInt(1000)))
		diffs = append(diffs, dpms)
		preHeader = header
	}

	for _, v := range diffs {
		x.Add(x, v)
	}
	return x.Div(x, big.NewInt(int64(len(diffs)))), nil
}

func (eth *EthSquadron) MakeDiffHeader(time, height uint64, diff *big.Int) *types.Header {
	return &types.Header{
		UncleHash:  types.EmptyUncleHash,
		Time:       time,
		Number:     new(big.Int).SetUint64(height),
		Difficulty: diff,
	}
}

func (eth *EthSquadron) CalcDifficulty(time uint64, parent *types.Header) *big.Int {
	return ethash.CalcDifficulty(EthConfig, time, parent)
}
