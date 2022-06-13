package ethsquadron

import (
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/llbec/robotech/utils"
)

func getEthNode() (*EthSquadron, error) {
	err := utils.LocalEnv()
	if err != nil {
		return nil, err
	}
	eth, err := NewEthSquadron(os.Getenv("RPC"))
	if err != nil {
		return nil, err
	}
	return eth, nil
}

func Test_ethNode(t *testing.T) {
	eth, err := getEthNode()
	if err != nil {
		t.Fatal(err)
	}
	ch, err := eth.BlockNumber()
	if err != nil {
		t.Fatal(err)
	}

	block, err := eth.HeaderByNumber(ch)
	if err != nil {
		t.Fatalf("error@(%v): %v\n", ch, err)
	}
	fmt.Println(block)
}

func Test_Diff(t *testing.T) {
	eth, err := getEthNode()
	if err != nil {
		t.Fatal(err)
	}
	ch, err := eth.BlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	preH := ch - 1
	preHeader, err := eth.HeaderByNumber(preH)
	if err != nil {
		t.Fatalf("error@(%v): %v\n", preH, err)
	}
	header, err := eth.HeaderByNumber(ch)
	if err != nil {
		t.Fatalf("error@(%v): %v\n", ch, err)
	}
	caDiff := eth.CalcDifficulty(header.Time, preHeader)
	if header.Difficulty.Cmp(caDiff) != 0 {
		t.Fatalf("\ntime from %v to %v\nexpect: %v\ncalc  : %v\n", preHeader.Time, header.Time, header.Difficulty, caDiff)
	}
}

func Test_MakeHeader(t *testing.T) {
	eth, err := getEthNode()
	if err != nil {
		t.Fatal(err)
	}
	ch, err := eth.BlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	preH := ch - 1
	preHeader, err := eth.HeaderByNumber(preH)
	if err != nil {
		t.Fatalf("error@(%v): %v\n", preH, err)
	}
	header, err := eth.HeaderByNumber(ch)
	if err != nil {
		t.Fatalf("error@(%v): %v\n", ch, err)
	}

	parent := eth.MakeDiffHeader(preHeader.Time, preHeader.Number.Uint64(), preHeader.Difficulty)

	caDiff := eth.CalcDifficulty(header.Time, parent)
	if header.Difficulty.Cmp(caDiff) != 0 {
		t.Fatalf("\ntime from %v to %v\nexpect: %v\ncalc  : %v\n", preHeader.Time, header.Time, header.Difficulty, caDiff)
	}
}

func Test_DiffPerMs(t *testing.T) {
	eth, err := getEthNode()
	if err != nil {
		t.Fatal(err)
	}
	h, err := eth.BlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	preH := h - 1
	preHeader, err := eth.HeaderByNumber(preH)
	if err != nil {
		t.Fatalf("error@(%v): %v\n", preH, err)
	}
	header, err := eth.HeaderByNumber(h)
	if err != nil {
		t.Fatalf("error@(%v): %v\n", h, err)
	}
	power, err := eth.DiffPerMilliSecond(h, 1)
	if err != nil {
		t.Fatal(err)
	}
	calcs := new(big.Int).Div(preHeader.Difficulty, power)
	calcs = new(big.Int).Div(calcs, big.NewInt(1000))
	fmt.Printf("block: %v\nexpect: %v\ncalc  : %v\n", h, header.Time-preHeader.Time, calcs)
	if calcs.Cmp(new(big.Int).SetUint64(header.Time-preHeader.Time)) != 0 {
		t.Fatalf("block: %v\nexpect: %v\ncalc  : %v\n", h, header.Time-preHeader.Time, calcs)
	}
}

func Test_DiffOracle(t *testing.T) {
	eth, err := getEthNode()
	if err != nil {
		t.Fatal(err)
	}
	h, err := eth.BlockNumber()
	if err != nil {
		t.Fatal(err)
	}
	power, err := eth.DiffPerMilliSecond(h, 1)
	if err != nil {
		t.Fatal(err)
	}
	header, err := eth.HeaderByNumber(h)
	if err != nil {
		t.Fatalf("error@(%v): %v\n", h, err)
	}
	pTime := header.Time
	pDiff := header.Difficulty
	fmt.Println("From", h)
	for i := 0; i < 10000; i++ {
		ms := new(big.Int).Div(pDiff, power)
		deltaTime := new(big.Int).Div(ms, big.NewInt(1000))
		nextTime := pTime + deltaTime.Uint64()
		fmt.Printf("next block will create at %v, after %v s\n", time.Unix(int64(nextTime), 0).Format("2006-01-02 15:04:05"), new(big.Float).Quo(new(big.Float).SetInt(ms), new(big.Float).SetInt64(1000)))
		// next
		parent := eth.MakeDiffHeader(pTime, h, pDiff)
		newDiff := eth.CalcDifficulty(nextTime, parent)
		fmt.Printf("delta difficulty: %v = %v - %v\n", new(big.Int).Sub(newDiff, pDiff), newDiff, pDiff)

		pDiff = newDiff
		pTime = nextTime
		h += 1
	}
}
