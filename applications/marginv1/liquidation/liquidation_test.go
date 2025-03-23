package liquidation_test

import (
	"math/big"
	"robotech/applications/marginv1/liquidation"
	EventEmitter "robotech/contract/margin/v1/event"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func Test_UpdateCollector(t *testing.T) {
	// TODO: Implement test cases for liquidation
	count := len(liquidation.PositionCollector)
	if count != 0 {
		t.Fatalf("Expected empty position collector, got %d", count)
	}
	mm0 := common.HexToAddress("0x0000000000000000000000000000000000000000")
	mm1 := common.HexToAddress("0x0000000000000000000000000000000000000001")
	accout0 := common.HexToAddress("0x0000000000000000000000000000000000000000")
	accout1 := common.HexToAddress("0x0000000000000000000000000000000000000001")
	accout3 := common.HexToAddress("0x000000000000000000000000000000000000003")
	accout4 := common.HexToAddress("0x000000000000000000000000000000000000004")
	event0 := &EventEmitter.EventEmitterPosition{
		Account:        accout0,
		ActionType:     common.Big0,
		BaseToken:      common.HexToAddress("0x0000000000000000000000000000000000000000"),
		MemeToken:      mm0,
		PositionId:     common.Big0,
		BaseCollateral: big.NewInt(10),
		BaseDebtScaled: big.NewInt(11),
		MemeCollateral: big.NewInt(12),
		MemeDebtScaled: big.NewInt(13),
	}
	liquidation.UpdateCollector(event0)
	event1 := &EventEmitter.EventEmitterPosition{
		Account:        accout1,
		ActionType:     common.Big0,
		BaseToken:      common.HexToAddress("0x0000000000000000000000000000000000000000"),
		MemeToken:      mm0,
		PositionId:     big.NewInt(1),
		BaseCollateral: big.NewInt(20),
		BaseDebtScaled: big.NewInt(21),
		MemeCollateral: big.NewInt(22),
		MemeDebtScaled: big.NewInt(23),
	}
	liquidation.UpdateCollector(event1)
	event3 := &EventEmitter.EventEmitterPosition{
		Account:        accout3,
		ActionType:     common.Big0,
		BaseToken:      common.HexToAddress("0x0000000000000000000000000000000000000000"),
		MemeToken:      mm1,
		PositionId:     big.NewInt(0),
		BaseCollateral: big.NewInt(30),
		BaseDebtScaled: big.NewInt(31),
		MemeCollateral: big.NewInt(32),
		MemeDebtScaled: big.NewInt(33),
	}
	liquidation.UpdateCollector(event3)
	event4 := &EventEmitter.EventEmitterPosition{
		Account:        accout4,
		ActionType:     common.Big0,
		BaseToken:      common.HexToAddress("0x0000000000000000000000000000000000000000"),
		MemeToken:      mm1,
		PositionId:     big.NewInt(1),
		BaseCollateral: big.NewInt(40),
		BaseDebtScaled: big.NewInt(41),
		MemeCollateral: big.NewInt(42),
		MemeDebtScaled: big.NewInt(43),
	}
	liquidation.UpdateCollector(event4)
	//t.Logf("collector is %v", liquidation.PositionCollector)
	if len(liquidation.PositionCollector) != 2 {
		t.Fatalf("Expected 2 positions, got %d", len(liquidation.PositionCollector))
	}
	if len(liquidation.PositionCollector[mm0]) != 2 {
		t.Fatalf("Expected 2 positions for mm0, got %d", len(liquidation.PositionCollector[mm0]))
	}
	if len(liquidation.PositionCollector[mm1]) != 2 {
		t.Fatalf("Expected 2 positions for mm1, got %d", len(liquidation.PositionCollector[mm1]))
	}
	if liquidation.PositionCollector[mm0][accout0].PositionId.Cmp(common.Big0) != 0 {
		t.Fatalf("Expected positionId 0, got %d", liquidation.PositionCollector[mm0][accout0].PositionId)
	}
	if liquidation.PositionCollector[mm1][accout4].PositionId.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("Expected positionId 1, got %d", liquidation.PositionCollector[mm1][accout4].PositionId)
	}
	event1 = &EventEmitter.EventEmitterPosition{
		Account:        accout1,
		ActionType:     common.Big0,
		BaseToken:      common.HexToAddress("0x0000000000000000000000000000000000000000"),
		MemeToken:      mm0,
		PositionId:     big.NewInt(1),
		BaseCollateral: big.NewInt(0),
		BaseDebtScaled: big.NewInt(0),
		MemeCollateral: big.NewInt(0),
		MemeDebtScaled: big.NewInt(0),
	}
	liquidation.UpdateCollector(event1)
	//t.Logf("collector is %v", liquidation.PositionCollector)
	if len(liquidation.PositionCollector[mm0]) != 1 {
		t.Fatalf("Expected 1 positions for mm0, got %d", len(liquidation.PositionCollector[mm0]))
	}
}
