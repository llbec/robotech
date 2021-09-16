package accounts

import "math/big"

type Account struct {
	Address          string
	CollateralAssets map[string]*big.Int
	DebtAssets       map[string]*big.Int
}

type AccountMana struct {
	Accounts map[string]*Account
}
