package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/llbec/robotech/armory/aave"
	"github.com/llbec/robotech/armory/flashloan"
	"github.com/llbec/robotech/logistics/db"
)

var (
	LendingPoolAddress   string
	DataProviderAddress  string
	FlashAdapterAddress  string
	OracleAddress        string
	LendingPool          *aave.LendingPool
	ProtocaldataProvider *aave.AaveProtocolDataProvider
	FlashLiquidationAdp  *flashloan.AaveLiquidationAdapter
	AaveOracle           *aave.AaveOracle
	Client               *ethclient.Client
	StartHeight          uint64
	LiquidatePeroid      int64
	BlockPeroid          int64
	sKey                 string
	rpcURL               string
)

var (
	curHeight       = uint64(0)
	workDir         string
	healthThreshold = big.NewInt(1e18)
	liquidateLevel1 = big.NewInt(2 * 1e8)
	Reserves        map[common.Address]*big.Int
)

const (
	actDeposit  = 1
	actWithdraw = 2
	actBorrow   = 3
	actRepay    = 4
	maxEnd      = 1000
	DmPort      = 1234
)

var (
	usrDb   db.DB
	logDb   db.DB
	debtors map[string]*big.Int
)
