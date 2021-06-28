package pricefeeds

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/llbec/robotech/utils"
)

func NewFeeds(n string, d int, c string) *Feeds {
	return &Feeds{
		token:    n,
		decimal:  d,
		contract: c,
	}
}

func (f *Feeds) String() string {
	return fmt.Sprintf("%6s: %2d, contract address is %s", f.token, f.decimal, f.contract)
}

func NewChainLinkFeeds(feedspath string) *ChainLinkFeeds {
	files, err := utils.FileList(feedspath, "feeds")
	if err != nil {
		panic(err)
	}
	pfs := &ChainLinkFeeds{}
	for _, fi := range files {
		pfeeds := &PriceFeeds{
			feeds: map[string][]*Feeds{},
		}
		strs := strings.Split(strings.Split(fi, ".")[0], "-")
		if len(strs) < 2 {
			log.Printf("%s> can't match format xxx-xxx.feeds\n", fi)
			continue
		}
		if strs[1] == "mainnet" {
			pfeeds.net = utils.NewChain(strs[0], false)
		} else {
			pfeeds.net = utils.NewChain(strs[0], true)
		}

		f, err := os.Open(path.Join(feedspath, fi))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		br := bufio.NewReader(f)
		for {
			line, _, err := br.ReadLine()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Printf("%s> read line failed: %v\n", fi, err)
				break
			}
			//log.Println(string(line))

			regToken := regexp.MustCompile("[0-9a-zA-Z]{2,10}")
			tks := regToken.FindAllString(string(line), 2)
			if len(tks) != 2 {
				log.Printf("%s> %v: read token name failed: %v\n", fi, string(line), tks)
				continue
			}
			regDecimal := regexp.MustCompile("[0-9]{1,2}")
			deci, err := strconv.Atoi(regDecimal.FindString(string(line)))
			if err != nil {
				log.Printf("%s> %v: read decimal(%d) failed: %v\n", fi, string(line), deci, err)
				continue
			}
			regAddr := regexp.MustCompile("0x[0-9a-zA-Z]{40}")
			addr := regAddr.FindString(string(line))
			nf := NewFeeds(tks[0], deci, addr)
			pfeeds.feeds[tks[1]] = append(pfeeds.feeds[tks[1]], nf)
			//log.Printf("Add feeds:%s-%s\n", nf.String(), tks[1])
		}
		pfs.nets = append(pfs.nets, pfeeds)
	}
	return pfs
}

func (pfs *ChainLinkFeeds) String() string {
	str := ""
	for _, pf := range pfs.nets {
		str += fmt.Sprintf("%s:\n", pf.net.String())
		for price, feeds := range pf.feeds {
			str += fmt.Sprintf("\t%s:\n", price)
			for i, f := range feeds {
				str += fmt.Sprintf("\t\t%3d: %s\n", i, f.String())
			}
		}
	}
	return str
}

func (pfs *ChainLinkFeeds) GetContract(chain string, testnet bool, asset, price string) (string, int, error) {
	for _, pf := range pfs.nets {
		if pf.net.Name() == chain && pf.net.TestNet() == testnet {
			for _, feeds := range pf.feeds[price] {
				if feeds.token == asset {
					return feeds.contract, feeds.decimal, nil
				}
			}
		}
	}
	return "", 0, fmt.Errorf("there's no %s / %s in %s", asset, price, chain)
}

func (pfs *ChainLinkFeeds) GetAssets(chain string, price string, testnet bool) (assets []string) {
	for _, pf := range pfs.nets {
		if pf.net.Name() == chain && pf.net.TestNet() == testnet {
			for _, feeds := range pf.feeds[price] {
				assets = append(assets, feeds.token)
			}
			return
		}
	}
	return
}
