package core

import (
	"fmt"
	"html/template"
	"ledger/common"
	"math/big"
	"strings"
)

const (
	AddressLength = 32
)

type Account struct {
	Key      string     `json:"key"`
	Name     string     `json:"name,omitempty"`
	Children []*Account `json:"children,omitempty"`
}

type ChartOfAccounts struct {
	Accounts []*Account `json:"accounts"`
}

type LedgerEntryType struct {
	Key        string `json:"key"`
	AccountKey string `json:"account"`
	Amount     string `json:"amount"` // This is a string because it appears to be a templated form
}

type LedgerTransactionType struct {
	Type  string            `json:"type"`
	Lines []LedgerEntryType `json:"lines"`
}

// PERSISTEN LT LAYER <------------------------------->
type LedgerEntries struct {
	Types []LedgerTransactionType `json:"types"`
}

func getLedgerEntries() {

}

func createEntry(line LedgerEntryType, params map[string]string) ([]Entries, error) {
	// Use text/template to evaluate the amount string
	tmpl, err := template.New("amountCalc").Parse(line.Amount)
	if err != nil {
		return nil, err
	}

	var builder strings.Builder
	err = tmpl.Execute(&builder, params)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(builder.String(), "+") // Splitting by '+'
	total := big.NewInt(0)
	for _, part := range parts {
		amount := new(big.Int)
		trimmedPart := strings.TrimSpace(part)
		amount, _ = amount.SetString(trimmedPart, 10)
		total = total.Add(total, amount)
	}

	account, exists := AccountStore[line.AccountKey]
	if !exists {
		return nil, fmt.Errorf("Account with key %s not found", line.AccountKey)
	}

	// Create Entries
	entry := NewEntry(account, total, common.Debit, common.Posted) // Adjust `common.Debit` as per your need

	return []Entries{*entry}, nil
}

func CreateTransaction(ik string, ledgerIK string, transactionType string, ledgerLines []LedgerEntryType, params map[string]string) *Transaction {
	//we are ignoring ledgerIk for now

	entriesList := []Entries{}

	for _, line := range ledgerLines {
		entries, err := createEntry(line, params)
		if err != nil {
			fmt.Println("Error creating entries:", err)
			continue
		}
		entriesList = append(entriesList, entries...)
	}

	// Create a transaction with the entries
	transaction := Transaction{
		id:      ik,
		entries: entriesList,
	}

	// Add the transaction to your ledger (omitting that part for brevity)
	// ...

	// Printing for demo purposes
	fmt.Println(transaction)
	return &transaction
}

// Persistent Layer -------------------------------->
type AccountsStore map[string]*Account

var AccountStore = make(AccountsStore)

func CreateAccount(key, name string, childrens ...string) *Account {
	childrensAccount := []*Account{}

	for i, children := range childrens {
		childrensAccount[i] = &Account{Key: children}
		AccountStore[children] = childrensAccount[i]
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
