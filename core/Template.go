package core

import (
	"encoding/json"
	"fmt"
	"html/template"
	"ledger/common"
	"math/big"
	"strings"
)

type LedgerEntryTemplate struct {
	Key        string           `json:"key"`
	AccountKey string           `json:"account"`
	Amount     string           `json:"amount"` // This is a string because it appears to be a templated form
	Direction  common.Direction `json:"direction"`
}

type LedgerTransactionTemplate struct {
	Type                  string                `json:"type"`
	LedgerEntriesTemplate []LedgerEntryTemplate `json:"lines"`
}

type LedgerEntries struct {
	Types []LedgerTransactionTemplate `json:"types"`
}

func UnmarshalLedgerEntryTemplate(ledgerEntryTemplateJson []byte) (*LedgerEntryTemplate, error) {
	ledgerEntryTemplate := &LedgerEntryTemplate{}
	err := json.Unmarshal([]byte(ledgerEntryTemplateJson), ledgerEntryTemplate)
	return ledgerEntryTemplate, err
}

func UnmarshalLedgerTransactionTemplate(ledgerTransactionTemplateJson []byte) (*LedgerTransactionTemplate, error) {
	ledgerTransactionTemplate := &LedgerTransactionTemplate{}
	err := json.Unmarshal([]byte(ledgerTransactionTemplateJson), ledgerTransactionTemplate)
	return ledgerTransactionTemplate, err
}

func (line LedgerEntryTemplate) createEntry(params map[string]string) ([]Entries, error) {
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

	accountKey, err := parseTemplateField(line.AccountKey, params)
	if err != nil {
		return nil, err
	}

	account, exists := AccountStore[accountKey]
	if !exists {
		return nil, fmt.Errorf("Account with key %s not found", line.AccountKey)
	}

	// Create Entries
	entry := NewEntry(account, total, line.Direction, common.Posted) // Adjust `common.Debit` as per your need

	return []Entries{*entry}, nil
}

func parseTemplateField(templateStr string, params map[string]string) (string, error) {
	tmpl, err := template.New("templateField").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	err = tmpl.Execute(&builder, params)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

func CreateTransaction(ik string, ledgerIK string, transactionType string, ledgerLines []LedgerEntryTemplate, params map[string]string) *Transaction {
	//we are ignoring ledgerIk for now

	entriesList := []Entries{}

	for _, line := range ledgerLines {
		entries, err := line.createEntry(params)
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

func (ledgertransaction LedgerTransactionTemplate) CreateTransaction(input TransactionInput) *Transaction {
	ik := input.Ledger.IK

	entriesList := []Entries{}

	for _, line := range ledgertransaction.LedgerEntriesTemplate {
		entries, err := line.createEntry(input.Parameters)
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
	Ledger     LedgerInfo        `json:"ledger"`
	Parameters map[string]string `json:"parameters"`
}

// TODO: maybeMoved to transaction or ledger.go in future
type LedgerInfo struct {
	IK      string `json:"ik"`
	Version string `json:"version,omitempty"` // `omitempty` will ignore the field if it's empty when encoding to JSON
}
