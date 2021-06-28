package liquidation

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"path/filepath"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/llbec/robotech/armory/aave"
	"github.com/llbec/robotech/armory/flashloan"
	"github.com/llbec/robotech/logistics/db"
	"github.com/llbec/robotech/utils"

	"github.com/spf13/viper"
)

const (
	LENDINGPOOL          = "lendingpool"
	PROTOCALDATAPROVIDER = "protocaldataprovider"
	FLASHLIQUIDATIONADP  = "flashliquidationadp"
	AAVEORACLE           = "aaveOracle"
	HTTPRPC              = "http"
	WEBSOCKETRPC         = "ws"
)

type ContractObjs struct {
	lendingpoolAddr         common.Address
	LendingPool             *aave.LendingPool
	ptcdataprovideraddr     common.Address
	ProtocaldataProvider    *aave.AaveProtocolDataProvider
	flashliquidationadpaddr common.Address
	FlashLiquidationAdp     *flashloan.AaveLiquidationAdapter
	aaveoracleaddr          common.Address
	AaveOracle              *aave.AaveOracle
}

var (
	Reserves             map[string]common.Address
	Contracts            ContractObjs
	Clients              map[string]*ethclient.Client
	LiquidationThreshold *big.Int
	secretkey            string
)

type XCom struct {
	datas          map[common.Address]SAccount
	datamtx        sync.Mutex
	channelInt     chan int
	liqudationlist map[common.Address]accountdata
	dbs            db.DB
}

func NewXCom() (*XCom, error) {
	initEnv()
	xc := &XCom{}
	xc.datas = make(map[common.Address]SAccount)
	xc.channelInt = make(chan int)
	return xc, nil
}

func initEnv() {
	viper.AddConfigPath(xcPath())
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Fatal("InitEnv-ReadInConfig ", err)
	}

	// init rpc client
	Clients = make(map[string]*ethclient.Client)
	rpcs := viper.GetStringMapString("rpc")
	for net, url := range rpcs {
		Clients[net], err = ethclient.Dial(url)
		if err != nil {
			log.Fatalf("InitEnv-ethclient.Dial error %v, %v, %v", err, net, url)
		}
	}

	// init contract object
	Contracts = ContractObjs{}
	Contracts.lendingpoolAddr = common.HexToAddress(viper.GetString(LENDINGPOOL))
	Contracts.LendingPool, err = aave.NewLendingPool(Contracts.lendingpoolAddr, Clients[HTTPRPC])
	if err != nil {
		log.Fatal("InitEnv-NewLendingPool ", err)
	}
	Contracts.ptcdataprovideraddr = common.HexToAddress(viper.GetString(PROTOCALDATAPROVIDER))
	Contracts.ProtocaldataProvider, err = aave.NewAaveProtocolDataProvider(Contracts.ptcdataprovideraddr, Clients[HTTPRPC])
	if err != nil {
		log.Fatal("InitEnv-NewProtocalDataProvider ", err)
	}
	Contracts.aaveoracleaddr = common.HexToAddress(viper.GetString(AAVEORACLE))
	Contracts.AaveOracle, err = aave.NewAaveOracle(Contracts.aaveoracleaddr, Clients[HTTPRPC])
	if err != nil {
		log.Fatal("InitEnv-NewAaveOracle ", err)
	}
	Contracts.flashliquidationadpaddr = common.HexToAddress(viper.GetString(FLASHLIQUIDATIONADP))
	Contracts.FlashLiquidationAdp, err = flashloan.NewAaveLiquidationAdapter(Contracts.flashliquidationadpaddr, Clients[HTTPRPC])
	if err != nil {
		log.Fatal("InitEnv-NewFlashLiquidationAdp ", err)
	}

	// init reserves
	Reserves = make(map[string]common.Address)
	rsvs := viper.GetStringMapString("reserves")
	for rsv, addr := range rsvs {
		Reserves[rsv] = common.HexToAddress(addr)
	}

	//
	LiquidationThreshold = big.NewInt(viper.GetInt64("liquidationthreshold"))
	secretkey = viper.GetString("privatekey")

	//
	//fmt.Printf("Enter private key ...\n")
	//fmt.Scanln(&secretkey)
}

func GetAuther() *bind.TransactOpts {
	secret, err := crypto.HexToECDSA(secretkey)
	if err != nil {
		log.Fatal("InitEnv-HexToECDSA ", err)
	}
	publicKey := secret.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := Clients[HTTPRPC].PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal("InitEnv-PendingNonceAt ", err)
	}
	chainid, err := Clients[HTTPRPC].ChainID(context.Background())
	if err != nil {
		log.Fatal("InitEnv-ChainID ", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(secret, chainid)
	if err != nil {
		log.Fatal("InitEnv-NewKeyedTransactorWithChainID ", err)
	}
	gasPrice, err := Clients[HTTPRPC].SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal("InitEnv-SuggestGasPrice ", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(4400000) // in units
	auth.GasPrice = gasPrice
	return auth
}

func xcPath() string {
	xcpath, err := utils.Home()
	if err != nil {
		log.Fatal("xcPath-Home ", err)
	}
	return filepath.Join(xcpath, ".xcom")
}

func GeneratConfig() error {
	rpcs := make(map[string]string)
	rpcs[HTTPRPC] = "https://http-mainnet-node.huobichain.com"
	rpcs[WEBSOCKETRPC] = "wss://heco.apron.network/ws/v1/4927dcd4e4964c47a8f863aba70b61f6"
	viper.Set("rpc", rpcs)

	viper.Set(LENDINGPOOL, "0x74CA5081d0a0561d1f654737642f9F9a3DDB0492")
	viper.Set(PROTOCALDATAPROVIDER, "0x18C789412E7652d094505E249d985DFE14Bd772e")
	viper.Set(FLASHLIQUIDATIONADP, "0xAD7958B141fd151910A67e973A4e1173d38C5801")
	viper.Set(AAVEORACLE, "0x1be62C4A97a45B04628E2B4f38F2eC71cC709e24")

	rsvs := make(map[string]string)
	rsvs["hbtc"] = "0x66a79D23E58475D2738179Ca52cd0b41d73f0BEa"
	rsvs["hdot"] = "0xA2c49cEe16a5E5bDEFDe931107dc1fae9f7773E3"
	rsvs["heth"] = "0x64FF637fB478863B7468bc97D30a5bF3A428a1fD"
	rsvs["hfil"] = "0xae3a768f9aB104c69A7CD6041fE16fFa235d1810"
	rsvs["hltc"] = "0xecb56cf772B5c9A6907FB7d32387Da2fCbfB63b4"
	rsvs["husd"] = "0x0298c2b32eaE4da002a15f36fdf7615BEa3DA047"
	rsvs["mdx"] = "0x25D2e80cB6B86881Fd7e07dd263Fb79f4AbE033c"
	rsvs["usdt"] = "0xa71EdC38d189767582C38A3145b5873052c3e47a"
	rsvs["wht"] = "0x5545153CCFcA01fbd7Dd11C0b23ba694D9509A6F"
	viper.Set("reserves", rsvs)

	viper.Set("liquidationthreshold", 1e8)
	viper.Set("privatekey", "123456789abcdef")

	return viper.WriteConfigAs(filepath.Join(xcPath(), "config.yaml"))
}
