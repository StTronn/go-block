package core

import (
	"fmt"
)

const (
	AddressLength = 32
)

type Account struct {
	Key      string    `json:"key"`
	Name     string    `json:"name,omitempty"`
	Children []Account `json:"children,omitempty"`
}

type AccountType struct {
	Key       string   `json:"key"`
	Name      string   `json:"name,omitempty"`
	Childrens []string `json:"children,omitempty"`
	template  []bool   `json:"template,omitempty"`
}

type ChartOfAccounts struct {
	Accounts []*AccountType `json:"accounts"`
}

// Persistent Layer -------------------------------->
type AccountsStore map[string]*Account

var AccountStore = make(AccountsStore)

func CreateAccount(accountType *AccountType) *Account {
	childrensAccount := make([]Account, len(accountType.Childrens))

	for i, children := range accountType.Childrens {
		fullKey := fmt.Sprintf("%s/%s", accountType.Key, children)
		childrensAccount[i] = Account{Key: fullKey}
		AccountStore[fullKey] = &childrensAccount[i]
	}
	fmt.Printf("size of AccountStore %v", len(AccountStore))

	account := &Account{
		Key:      accountType.Key,
		Name:     accountType.Name,
		Children: childrensAccount,
	}
	AccountStore[account.Key] = account
	return account

}
