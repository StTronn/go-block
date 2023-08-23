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
	Key      string   `json:"key"`
	Name     string   `json:"name,omitempty"`
	Children []string `json:"children,omitempty"`
}

type ChartOfAccounts struct {
	Accounts []*AccountType `json:"accounts"`
}

// Persistent Layer -------------------------------->
type AccountsStore map[string]*Account

var AccountStore = make(AccountsStore)

func CreateAccount(key, name string, childrens ...string) *Account {
	childrensAccount := make([]Account, len(childrens))

	for i, children := range childrens {
		fullKey := fmt.Sprintf("%s/%s", key, children)
		childrensAccount[i] = Account{Key: fullKey}
		AccountStore[fullKey] = &childrensAccount[i]
	}
	fmt.Printf("size of AccountStore %v", len(AccountStore))

	account := &Account{
		Key:      key,
		Name:     name,
		Children: childrensAccount,
	}
	AccountStore[key] = account
	return account

}
