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

func createEntryFromLedgerEntryType(line LedgerEntryType, params map[string]string) ([]Entries, error) {
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

	// Convert string result to big.Int (for now assuming result is directly convertible)
	// You may want to improve this to handle complex expressions
	amount := new(big.Int)
	amount, _ = amount.SetString(builder.String(), 10)

	account, exists := AccountStore[line.AccountKey]
	if !exists {
		return nil, fmt.Errorf("Account with key %s not found", line.AccountKey)
	}

	// Create Entries
	// Direction and status can be inferred based on your system's logic
	// For simplicity, I'm setting direction as debit and status as posted
	// Adjust according to your needs
	entry := NewEntry(account, amount, common.Debit, common.Posted) // Adjust `common.Debit` as per your need

	return []Entries{*entry}, nil
}

func addTransactionEntry(ik string, ledgerIK string, entryType string, ledgerLines []LedgerEntryType, params map[string]string) *Transaction {
	//we are ignoring ledgerIk for now

	entriesList := []Entries{}

	for _, line := range ledgerLines {
		entries, err := createEntryFromLedgerEntryType(line, params)
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
