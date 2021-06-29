package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/llbec/robotech/logistics/daemon"
	"github.com/llbec/robotech/squadrons/balasquadron/fallbackteam"
	"github.com/llbec/robotech/squadrons/chainlinksquadron"
)

var (
	fInit       bool
	fList       bool
	fRun        bool
	fTerminate  bool
	fGet        string
	fDir        string
	fSet        string
	gnode       *chainlinksquadron.ChainNode
	gclSquadron *chainlinksquadron.ChainLinkSquadron
	gfbTeam     *fallbackteam.FallBack
)

func init() {
	flag.BoolVar(&fInit, "i", false, "Generate a simple config file. default false")
	flag.BoolVar(&fList, "l", false, "list all assets price. default false")
	flag.BoolVar(&fRun, "r", false, "Running fallback process. default false")
	flag.BoolVar(&fTerminate, "t", false, "Terminate fallback process. default false")
	flag.StringVar(&fDir, "d", ".", "specify feeds file dir. defult is \".\"")
	flag.StringVar(&fGet, "g", "", "Read specify asset(uppercase letter) price. default is null")
	flag.StringVar(&fSet, "s", "", "Set specify asset(uppercase letter) price. default is null")
}

type AssetFeeds struct {
	name  string
	price *big.Int
}

func main() {
	flag.Parse()

	wd, e := os.Getwd()
	if e != nil {
		log.Printf("load config file failed: %v\n", e)
		return
	}
	if fInit {
		GeneratConfig(wd)
		return
	}
	LoadConfig(wd)

	//init global var
	gnode = chainlinksquadron.NewChainNode()
	for n, r := range cfg.rpcs {
		if n != "hsc" {
			gnode.AddNode(n, r)
		}
	}

	gclSquadron = chainlinksquadron.NewChainLinkSquadron(gnode, fDir)
	gfbTeam = fallbackteam.NewFallback(cfg.rpcs[cfg.fallbacknet], cfg.fallbackaddr)

	if fGet != "" {
		showPrice(fGet)
		return
	}

	if fList {
		//alist := gclSquadron.GetAssets(net, PRICE, testnet)
		for asset := range cfg.assetscontract {
			showPrice(asset)
		}
		return
	}

	if fTerminate {
		d := daemon.NewDaemon(1234, nil)
		d.Stop()
		return
	}

	sk := ""

	if fSet != "" {
		fmt.Println("Enter fallback manager's private key")
		fmt.Scanf("%s", &sk)
		gfbTeam.SetKey(sk)
		setPrice(fSet)
		return
	}

	if fRun {
		sk, e := getSecret()
		if e != nil {
			fmt.Println(e)
			return
		}
		gfbTeam.SetKey(sk)
		d := daemon.NewDaemon(1234, fallbackPrice)
		d.Run("fallback.log")
		return
	}

	flag.Usage()
}

func readPrice(asset string) (*big.Int, *big.Int, int, error) {
	pNet, d, e := gclSquadron.GetLastPrice(cfg.assetschanlinknet[asset], strings.ToUpper(asset), cfg.price, cfg.testnet)
	if e != nil {
		return big.NewInt(0), big.NewInt(0), 0, fmt.Errorf("read %s@%s price failed: %v", asset, cfg.assetschanlinknet[asset], e)
	}
	pFB, e := gfbTeam.GetPrice(cfg.assetscontract[asset])
	if e != nil {
		return big.NewInt(0), big.NewInt(0), 0, fmt.Errorf("read %s@Bala price failed: %v", asset, e)
	}
	return pNet, pFB, d, nil
}

func showPrice(asset string) {
	pNet, pFB, d, e := readPrice(asset)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Printf("%6s/%s : %7s(%14v) bala(%14v) %d\n", asset, cfg.price, cfg.assetschanlinknet[asset], pNet, pFB, d)
	}
}

func setPrice(astr string) {
	alist := gclSquadron.GetAssets(cfg.assetschanlinknet[astr], cfg.price, cfg.testnet)
	for _, asset := range alist {
		if asset == strings.ToUpper(astr) {
			pNet, pFB, d, e := readPrice(astr)
			if e != nil {
				fmt.Println(e)
				return
			}
			if pNet.Cmp(pFB) == 0 {
				fmt.Printf("price is same %s\n", pNet.String())
				return
			}
			e = gfbTeam.SetPrice(cfg.assetscontract[astr], pNet)
			if e != nil {
				fmt.Printf("set %s price(%v) failed: %v\n", astr, pNet.String(), e)
				return
			}
			fmt.Printf("set %s price from %s to %s on bala(%d)\n", astr, pFB.String(), pNet.String(), d)
			return
		}
	}
	if cfg.testnet {
		fmt.Printf("ChainLink %s testnet assets(%s) is %v\n", cfg.assetschanlinknet[astr], PRICE, alist)
	} else {
		fmt.Printf("ChainLink %s mainnet assets(%s) is %v\n", cfg.assetschanlinknet[astr], PRICE, alist)
	}
}

func fallbackPrice(chSig, chExit chan int) {
	tc := time.NewTicker(time.Duration(cfg.tickime) * time.Second)
	defer tc.Stop()

	for asset, addr := range cfg.assetscontract {
		p, e := gfbTeam.GetPrice(addr)
		if e != nil {
			log.Println(e)
			goto EXIT
		}
		assets = append(assets, AssetFeeds{name: asset, price: p})
	}
	log.Println("Start to watch price")
	for {
		select {
		case <-tc.C:
			for i, feeds := range assets {
				lastp, _, e := gclSquadron.GetLastPrice(cfg.assetschanlinknet[feeds.name], strings.ToUpper(feeds.name), cfg.price, cfg.testnet)
				if e != nil {
					log.Printf("read %s@%s price failed: %v\n", feeds.name, cfg.assetschanlinknet[feeds.name], e)
					continue
				}
				delta := time.Now().Unix() - lastupdate
				if isupdatePrice(lastp, feeds.price, delta) {
					e := gfbTeam.SetPrice(cfg.assetscontract[feeds.name], lastp)
					if e != nil {
						log.Printf("set %s: %v -> %v failed: %v\n", feeds.name, feeds.price, lastp.String(), e)
						continue
					}
					assets[i].price = lastp
					lastupdate = time.Now().Unix()
					log.Printf("update %s: %v -> %v, time %v\n", feeds.name, feeds.price, lastp, time.Unix(lastupdate, 0).Format("2006-01-02 15:04:05"))
				} else {
					log.Printf("skip %s: %v -> %v, delta time %v\n", feeds.name, feeds.price, lastp, delta)
				}
			}
		case <-chSig:
			goto EXIT
		}
	}
EXIT:
	chExit <- 1
}

//deltaPercent return |x-y|*100/x
func deltaPercent(x, y *big.Int) *big.Int {
	delta := big.NewInt(0).Sub(x, y)
	return delta.Abs(delta).Mul(delta, big.NewInt(100)).Div(delta, x)
}

func isupdatePrice(p, lp *big.Int, delta int64) bool {
	v := deltaPercent(p, lp)
	if v.Cmp(big.NewInt(int64(cfg.threhold))) >= 0 {
		return true
	}

	if p.Cmp(lp) != 0 && delta >= cfg.interval {
		return true
	}
	return false
}

func getSecret() (string, error) {
	f, err := ioutil.ReadFile("key")
	if err != nil {
		return "", err
	}
	return strings.Split(string(f), "\n")[0], nil
}
