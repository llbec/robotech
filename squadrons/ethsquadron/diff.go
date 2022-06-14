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

var (
	expDiffPeriod = big.NewInt(100000)
	big1          = big.NewInt(1)
	big2          = big.NewInt(2)
	big9          = big.NewInt(9)
	bigMinus99    = big.NewInt(-99)
)

func (eth *EthSquadron) BombDiff(height *big.Int) *big.Int {
	bombDiff := new(big.Int)
	bombDelayFromParent := new(big.Int).Sub(big.NewInt(10_700_000), big1)

	fakeBlockNumber := new(big.Int)
	if height.Cmp(bombDelayFromParent) >= 0 {
		fakeBlockNumber = fakeBlockNumber.Sub(height, bombDelayFromParent)
	}
	// for the exponential factor
	periodCount := fakeBlockNumber
	periodCount.Div(periodCount, expDiffPeriod)

	// the exponential factor, commonly referred to as "the bomb"
	// diff = diff + 2^(periodCount - 2)
	if periodCount.Cmp(big1) > 0 {
		bombDiff.Sub(periodCount, big2)
		bombDiff.Exp(big2, bombDiff, nil)
	}

	return bombDiff
}

func (eth *EthSquadron) OffsetDiff(time uint64, parent *types.Header) *big.Int {
	bigTime := new(big.Int).SetUint64(time)
	bigParentTime := new(big.Int).SetUint64(parent.Time)
	x := new(big.Int)

	// (2 if len(parent_uncles) else 1) - (block_timestamp - parent_timestamp) // 9
	x.Sub(bigTime, bigParentTime)
	x.Div(x, big9)
	if parent.UncleHash == types.EmptyUncleHash {
		x.Sub(big1, x)
	} else {
		x.Sub(big2, x)
	}

	// max((2 if len(parent_uncles) else 1) - (block_timestamp - parent_timestamp) // 9, -99)
	if x.Cmp(bigMinus99) < 0 {
		x.Set(bigMinus99)
	}
	return x
}
