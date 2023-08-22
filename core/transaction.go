package core

import (
	"fmt"
	"html/template"
	"ledger/common"
	"math/big"
	"strings"

	"github.com/rs/xid"
)

type Transaction struct {
	id      string
	entries []Entries
}

type Entries struct {
	id        string
	account   *Account
	amount    *big.Int
	direction common.Direction
}

// newEntries creates a new Entries with a unique id.
func NewEntry(Account *Account, amount *big.Int, direction common.Direction, status common.Status) *Entries {
	return &Entries{
		id:        xid.New().String(),
		account:   Account,
		amount:    amount,
		direction: direction,
	}
}

// NewTransaction creates a new Transaction with a unique id and a slice of entries.
func NewTransaction(entries ...Entries) *Transaction {
	return &Transaction{
		id:      xid.New().String(),
		entries: entries,
	}
}

type LedgerEntryType struct {
	Key        string `json:"key"`
	AccountKey string `json:"account"`
	Amount     string `json:"amount"` // This is a string because it appears to be a templated form
}

type LedgerTransactionType struct {
	Type    string            `json:"type"`
	Entries []LedgerEntryType `json:"lines"`
}

// PERSISTEN LT LAYER <------------------------------->
type LedgerEntries struct {
	Types []LedgerTransactionType `json:"types"`
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

type TransactionInput struct {
	Type       string            `json:"type"`
	Ledger     map[string]string `json:"ledger"`
	Parameters map[string]string `json:"parameters"`
}

//convert transactions into double ledger Transaction
//

// type Transaction struct {
// 	signer Account //TODO: change to hash
// 	to     Account //TODO: change to hash
// 	value  *big.Int
// 	status string
// }

// type Entry struct {
// 	account Account
// 	value   *big.Int
// }

// type TxData interface {
// 	txType() byte
// }
