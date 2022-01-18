package compliquidation_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"testing"
	"time"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/llbec/robotech/armory/compound"
	liquidation "github.com/llbec/robotech/squadrons/compoundsquadron/compliquidation"
)

func Test_CompoundApi(t *testing.T) {
	body := make(map[string]interface{})
	body["min_borrow_value_in_eth"] = map[string]string{"value": "1"}
	body["max_health"] = map[string]string{"value": "1.0"}

	bytesData, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	reader := bytes.NewReader(bytesData)

	u := url.URL{Scheme: "https", Host: "api.compound.finance", Path: "/api/v2/account"}
	request, err := http.NewRequest("POST", u.String(), reader)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("Accept", "application/json")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}

	resBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	str := (*string)(unsafe.Pointer(&resBytes))
	fmt.Println("rpc params:\n", *str)
}

func testEnv() *liquidation.CompLiquidation {
	rpcurl := "https://http-mainnet-node.huobichain.com"
	//wsurl := "wss://heco.apron.network/ws/v1/afdff8d504f14ee8a2643112ec9578aa"
	taddr := "0xb74633f2022452f377403B638167b0A135DB096d"
	faddr := "0xb74633f2022452f377403B638167b0A135DB096d"
	return liquidation.NewReadOnce(rpcurl, taddr, faddr, big.NewFloat(10))
}

func Test_BigMath(t *testing.T) {
	a, ok := new(big.Int).SetString("138726159538235265326902", 10)
	if !ok {
		t.Fatalf("138726159538235265326902 set int faild\n")
	}
	af := new(big.Float).SetInt(a)
	fmt.Println(af)
	fmt.Println(af.Quo(af, new(big.Float).SetFloat64(1e18)))
	f32, ac := af.Float32()
	f64, ac2 := af.Float64()
	fmt.Println(f32, ac, f64, ac2)
	if f64 > 1 {
		fmt.Println(f64)
	}
}

func Test_deleteArray(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i := 0; i < len(a); i++ {
		if a[i] == 5 || a[i] == 6 {
			a = append(a[:i], a[i+1:]...)
			i -= 1
		}
	}
	fmt.Println(a)
}

/*func Test_GetAccountAssetStatus(t *testing.T) {
	ld := testEnv()
	account, l, b := ld.GetAccountAssetStatus("0x744808B56E9B470CE901ABaa9A6261E92928cE82")
	fmt.Printf("%#v\n%v\n%v\n", account.String(), l, b)
}*/

func Test_GetAccountSnapshot(t *testing.T) {
	ld := testEnv()
	c, b, err := ld.GetAccountSnapshot(common.HexToAddress("0x744808B56E9B470CE901ABaa9A6261E92928cE82"), common.HexToAddress("0x0AD0bee939E00C54f57f21FBec0fBa3cDA7DEF58"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c, b)
}

func Test_GetAccount(t *testing.T) {
	ld := testEnv()
	addrs := []string{"0x7c0Fc56fE45c98cCF000DD548F457D2fCA9ECC6e",
		"0xf2A78433AB99230015Ae37A566349248e9514CfD"}
	for _, addr := range addrs {
		a, err := ld.GetAccount(common.HexToAddress(addr))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(a.String())
	}

}

func Test_LiquidationAccount(t *testing.T) {
	ld := testEnv()
	err := ld.LiquidationDebtor(common.HexToAddress("0xf2A78433AB99230015Ae37A566349248e9514CfD"))
	if err != nil {
		t.Fatal(err)
	}
}

func Test_GetZero(t *testing.T) {
	ld := testEnv()
	ctoken, err := compound.NewCToken(common.HexToAddress("0x5545153ccfca01fbd7dd11c0b23ba694d9509a6f"), ld.Chain.GetHTTPClient())
	if err != nil {
		t.Fatal(err)
	}

	b, err := ctoken.BalanceOf(nil, common.HexToAddress("0x0000000000000000000000000000000000000000"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("0x0000000000000000000000000000000000000000 blance is %v\n", b)
}

func Test_Prices(t *testing.T) {
	ld := testEnv()
	tokens := map[string]string{
		"FilDA": "0xe36ffd17b2661eb57144ceaef942d95295e637f0",
		"fHT":   "0x824151251B38056d54A15E56B73c54ba44811aF8",
		"fHFIL": "0x043aFB65e93500CE5BCbf5Bbb41FC1fDcE2B7518",
		"fHUSD": "0xB16Df14C53C4bcfF220F4314ebCe70183dD804c0",
		"fHBTC": "0xF2a308d3Aea9bD16799A5984E20FDBfEf6c3F595",
		"fETH":  "0x033F8C30bb17B47f6f1f46F3A42Cc9771CCbCAAE",
		"fHPT":  "0x749E0198f12559E7606987F8e7bD3AA1DE6d236E",
		"fELA":  "0x0AD0bee939E00C54f57f21FBec0fBa3cDA7DEF58",
		"fHDOT": "0xCca471B0d49c0d4835a5172Fd97ddDEA5C979100",
		"fHBCH": "0x09e3d97A7CFbB116B416Dae284f119c1eC3Bd5ea",
		"fHLTC": "0x4937A83Dc1Fa982e435aeB0dB33C90937d54E424",
		"fUSDT": "0xAab0C9561D5703e84867670Ac78f6b5b4b40A7c1",
		"fHBSV": "0x74f8d9b701bd4d8ee4ec812af82c71eb67b9ec75",
		"fHXTZ": "0xfea846a1284554036ac3191b5dfd786c0f4db611",
		"fAAVE": "0x73Fa2931e060F7d43eE554fd1De7F61115fE1751",
		"fUNI":  "0xAc9E3AE0C188eb583785246Fef37AEF9ea159fb7",
		"fSNX":  "0x88962975FDE8C7805fE0f38b7c91C18f4d55bb40",
		"fNEO":  "0x92701DA6A28Ca70aA5Dfca2B8Ae2b4B8a22a0C11",
	}
	for name, addr := range tokens {
		token := common.HexToAddress(addr)
		p, err := ld.Oracle.GetUnderlyingPrice(nil, token)
		if err != nil {
			t.Errorf("GetUnderlyingPrice(%v):%v\n", name, err)
		}
		fmt.Printf("%v: %v\n", name, p)
	}
	str := "fHT"
	ctoken, err := compound.NewCToken(common.HexToAddress(tokens[str]), ld.Chain.GetHTTPClient())
	if err != nil {
		t.Fatalf("NewCToken(%v):%v\n", str, err)
	}
	_, _, _, id, _ := ctoken.GetAccountSnapshot(nil, common.HexToAddress("0x744808B56E9B470CE901ABaa9A6261E92928cE82"))
	fmt.Printf("index: %v\n", id)
}

type Thread struct {
	a int
	b int
}

func (th *Thread) Adda(ch chan int) {
	for i := 0; i < 10; i++ {
		th.a += 1
		th.b = th.b - 1
		fmt.Printf("Adda %v:%v-%v\n", i, th.a, th.b)
		time.Sleep(time.Duration(1) * time.Second)
	}
	ch <- 1
}

func (th *Thread) Addb(ch chan int) {
	for i := 0; i < 10; i++ {
		th.b += 1
		th.a = th.a - 1
		fmt.Printf("Addb %v:%v-%v\n", i, th.a, th.b)
		time.Sleep(time.Duration(1) * time.Second)
	}
	ch <- 2
}

func Test_Thread(t *testing.T) {
	th := &Thread{a: 10, b: 0}
	ch := make(chan int, 2)
	go th.Adda(ch)
	//<-ch
	go th.Addb(ch)
	<-ch
	<-ch
	fmt.Println(th)
}

func Test_RpcHost(t *testing.T) {
	http, err := ethclient.Dial("https://heco.apron.network/v1/39407a13d0e54fdc99dde78ca1c41900")
	if err != nil {
		t.Error(err)
	} else {
		id, _ := http.ChainID(context.Background())
		fmt.Println(id)
	}

	ws, err := ethclient.Dial("wss://heco.apron.network/ws/v1/39407a13d0e54fdc99dde78ca1c41900")
	if err != nil {
		t.Error(err)
	} else {
		id, _ := ws.ChainID(context.Background())
		fmt.Println(id)
	}
}
