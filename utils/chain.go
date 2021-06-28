package utils

import "fmt"

type Chain struct {
	name    string
	testnet bool
}

func NewChain(name string, test bool) *Chain {
	return &Chain{
		name:    name,
		testnet: test,
	}
}

func (c *Chain) Name() string {
	return c.name
}

func (c *Chain) TestNet() bool {
	return c.testnet
}

func (c *Chain) String() string {
	if c.testnet {
		return fmt.Sprintf("%s-testnet", c.name)
	}
	return fmt.Sprintf("%s-mainnet", c.name)
}
