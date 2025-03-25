package liquidation_test

import (
	"math/big"
	"robotech/applications/marginv1/liquidation"
	"robotech/applications/marginv1/storage"
	v1 "robotech/contract/margin/v1"
	EventEmitter "robotech/contract/margin/v1/event"
	Router "robotech/contract/margin/v1/router"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func Test_UpdateCollector(t *testing.T) {
	// TODO: Implement test cases for liquidation
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
	mm0Positions := storage.GetPositions(mm0)
	mm1Positions := storage.GetPositions(mm1)
	//t.Logf("collector is %v", liquidation.PositionCollector)
	if len(mm0Positions) != 2 {
		t.Fatalf("Expected 2 positions for mm0, got %d", len(mm0Positions))
	}
	if len(mm1Positions) != 2 {
		t.Fatalf("Expected 2 positions for mm1, got %d", len(mm1Positions))
	}
	for _, position := range mm0Positions {
		if position.Account == event0.Account {
			if position.MemeToken != event0.MemeToken {
				t.Fatalf("Expected memeToken %v, got %v", event0.MemeToken, position.MemeToken)
			}
			if position.PositionId.Cmp(event0.PositionId) != 0 {
				t.Fatalf("Expected positionId %v, got %v", event0.PositionId, position.PositionId)
			}
		}
		if position.Account == event1.Account {
			if position.MemeToken != event1.MemeToken {
				t.Fatalf("Expected memeToken %v, got %v", event1.MemeToken, position.MemeToken)
			}
			if position.PositionId.Cmp(event1.PositionId) != 0 {
				t.Fatalf("Expected positionId %v, got %v", event1.PositionId, position.PositionId)
			}
		}
	}
	for _, position := range mm1Positions {
		if position.Account == event3.Account {
			if position.MemeToken != event3.MemeToken {
				t.Fatalf("Expected memeToken %v, got %v", event3.MemeToken, position.MemeToken)
			}
			if position.PositionId.Cmp(event3.PositionId) != 0 {
				t.Fatalf("Expected positionId %v, got %v", event3.PositionId, position.PositionId)
			}
		}
		if position.Account == event4.Account {
			if position.MemeToken != event4.MemeToken {
				t.Fatalf("Expected memeToken %v, got %v", event4.MemeToken, position.MemeToken)
			}
			if position.PositionId.Cmp(event4.PositionId) != 0 {
				t.Fatalf("Expected positionId %v, got %v", event4.PositionId, position.PositionId)
			}
		}
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
	mm0Positions = storage.GetPositions(mm0)
	//t.Logf("collector is %v", liquidation.PositionCollector)
	if len(mm0Positions) != 1 {
		t.Fatalf("Expected 1 positions for mm0, got %d", len(mm0Positions))
	}
}

func Test_PositionAbiEncode(t *testing.T) {
	t.Logf("test position abi encode")
	positionCode := v1.PositionEncode()
	if len(positionCode) != 32 {
		t.Fatalf("Expected 32 bytes, got %d", len(positionCode))
	}
	if !strings.EqualFold(common.Bytes2Hex(positionCode), "cdc7b43153b083746433bc289e68a932ed04e12dec99b96b53b46ac976c0acce") {
		t.Fatalf("position abi encode is %v", positionCode)
	}
	t.Logf("test position abi encode success: %v", common.Bytes2Hex(positionCode))
}

func Test_PositionKey(t *testing.T) {
	t.Logf("test position key")
	positionKey := v1.CalcPositionKey(common.HexToAddress("0xe87069a64C32e3bF7a40B2Ce2778995533359074"), big.NewInt(1))
	if len(positionKey) != 32 {
		t.Fatalf("Expected 32 bytes, got %d", len(positionKey))
	}
	if !strings.EqualFold(common.Bytes2Hex(positionKey), "2b3828b9814e62718c2da30af90c25e6b47ef1d2bcdb99e9232167adac83756d") {
		t.Fatalf("position key is %v", positionKey)
	}
	t.Logf("test position key success: %v", common.Bytes2Hex(positionKey))
}

func Test_MarginLevelThreshold(t *testing.T) {
	t.Logf("test margin level threshold")
	threshold := liquidation.MarginLevelThreshold()
	t.Logf("test margin level threshold success: %v", threshold)
}

func Test_CallLiquidation(t *testing.T) {
	usr := common.HexToAddress("0x90518217d485843B87E2aCa5af4af0eF928eeEBd")
	positions, err := liquidation.CallGetPositions(usr)
	if err != nil {
		t.Fatalf("CallGetPositions error: %v", err)
	}

	var liquidationParam []Router.LiquidationUtilsLiquidationParams
	for _, position := range positions {
		if position.MarginLevel.Cmp(liquidation.MarginLevelThreshold()) < 1 {
			liquidationParam = append(liquidationParam, Router.LiquidationUtilsLiquidationParams{
				Account:    position.Account,
				PositionId: position.Id,
			})
			t.Logf("add Liquid operation account(%v) with Id(%v)", position.Account, position.Id)
		}
	}

	err = liquidation.CallLiquidation(liquidationParam)
	if err != nil {
		t.Fatalf("CallLiquidation error: %v", err)
	}
}
