package accounts

import (
	"log"
	"math/big"
)

func NewAccount() Account {
	collaterals := make(map[string]*big.Int)
	debts := make(map[string]*big.Int)
	return Account{
		Address:          "",
		CollateralAssets: collaterals,
		DebtAssets:       debts,
	}
}

func NewAccountMana() *AccountMana {
	accounts := make(map[string]*Account)
	return &AccountMana{Accounts: accounts}
}

func (aMana *AccountMana) Deposit(address string, asset string, amount *big.Int) {
	if _, exist := aMana.Accounts[address]; exist {
		if _, assetExist := aMana.Accounts[address].CollateralAssets[asset]; assetExist {
			aMana.Accounts[address].CollateralAssets[asset].Add(
				aMana.Accounts[address].CollateralAssets[asset],
				amount)
		} else {
			aMana.Accounts[address].CollateralAssets[asset] = amount
		}
	} else {
		account := NewAccount()
		account.Address = address
		account.CollateralAssets[asset] = amount
		aMana.Accounts[address] = &account
	}
}

func (aMana *AccountMana) Withdraw(address string, asset string, amount *big.Int) {
	if _, exist := aMana.Accounts[address]; exist {
		if _, assetExist := aMana.Accounts[address].CollateralAssets[asset]; assetExist {
			aMana.Accounts[address].CollateralAssets[asset].Sub(
				aMana.Accounts[address].CollateralAssets[asset],
				amount)
		} else {
			aMana.Accounts[address].CollateralAssets[asset] = big.NewInt(0).Sub(big.NewInt(0), amount)
			//log should not happen
			log.Printf("No exit asset(%v): address(%v) withdraw %v", asset, address, amount)
		}
	} else {
		account := NewAccount()
		account.Address = address
		account.CollateralAssets[asset] = big.NewInt(0).Sub(big.NewInt(0), amount)
		aMana.Accounts[address] = &account
		//log should not happen
		log.Printf("No exit address(%v): withdraw @ %v - %v", address, asset, amount)
	}
}

func (aMana *AccountMana) Borrow(address string, asset string, amount *big.Int) {
	if _, exist := aMana.Accounts[address]; exist {
		if _, assetExist := aMana.Accounts[address].DebtAssets[asset]; assetExist {
			aMana.Accounts[address].DebtAssets[asset].Add(
				aMana.Accounts[address].DebtAssets[asset],
				amount)
		} else {
			aMana.Accounts[address].DebtAssets[asset] = amount
		}
	} else {
		account := NewAccount()
		account.Address = address
		account.DebtAssets[asset] = amount
		aMana.Accounts[address] = &account
		//log should not happen
		log.Printf("No exit address(%v): borrow @ %v - %v", address, asset, amount)
	}
}

func (aMana *AccountMana) Repay(address string, asset string, amount *big.Int) {
	if _, exist := aMana.Accounts[address]; exist {
		if _, assetExist := aMana.Accounts[address].DebtAssets[asset]; assetExist {
			aMana.Accounts[address].DebtAssets[asset].Sub(
				aMana.Accounts[address].DebtAssets[asset],
				amount)
		} else {
			aMana.Accounts[address].DebtAssets[asset] = big.NewInt(0).Sub(big.NewInt(0), amount)
			//log should not happen
			log.Printf("No exit asset(%v): address(%v) repay %v", asset, address, amount)
		}
	} else {
		account := NewAccount()
		account.Address = address
		account.DebtAssets[asset] = big.NewInt(0).Sub(big.NewInt(0), amount)
		aMana.Accounts[address] = &account
		//log should not happen
		log.Printf("No exit address(%v): repay @ %v - %v", address, asset, amount)
	}
}

func (aMana *AccountMana) GetDebtors() []string {
	var debtors []string

	for _, account := range aMana.Accounts {
		if account.Isdebtor() {
			debtors = append(debtors, account.Address)
		}
	}

	return debtors
}

func (account *Account) Isdebtor() bool {
	for _, amount := range account.DebtAssets {
		if amount.Cmp(big.NewInt(0)) != 0 {
			return true
		}
	}
	return false
}
