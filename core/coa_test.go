package core

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

const ledgerTransactionsJson = `
{
	"transactions": {
		"types": [{
			"type": "sell_something",
			"lines": [{
					"key": "sales_to_bank",
					"account": "sales_to_bank",
					"amount": "{{.sales_before_tax}} + {{.tax_payable}}"
				},
				{
					"key": "income-root",
					"account": "income-root",
					"amount": "{{.sales_before_tax}}"
				},
				{
					"key": "tax_payable",
					"account": "tax_payable",
					"amount": "{{.tax_payable}}"
				}
			]
		}]
	}
}
`
const ChartOfAccountsJson = `
{
	"accounts": [{
			"key": "sales_to_bank",
			"name": "sales_to_bank"
		},
		{
			"key": "income-root",
			"name": "income-root-bank"
		},
		{
			"key": "tax_payable",
			"name": "tax_payable"
		}
	]

}
`

const transactionInput = `
{
  "type": "sell_something",
  "ledger": {
    "ik": "my-ledger-ik"
  },
  "parameters": {
    "sales_before_tax": "10000",
    "tax_payable": "500"
  }
}
`

func loadAccounts() {
	chartOfAccounts := &ChartOfAccounts{}
	err := json.Unmarshal([]byte(ChartOfAccountsJson), chartOfAccounts)

	if err != nil {
		fmt.Errorf("error parsing chart of accounts")
	}

	for _, account := range chartOfAccounts.Accounts {
		CreateAccount(account.Key, account.Name)
	}

}

func TestChartOfAccounts(t *testing.T) {
	chartOfAccounts := &ChartOfAccounts{}
	err := json.Unmarshal([]byte(ChartOfAccountsJson), chartOfAccounts)
	assert.Nil(t, err)

	for _, account := range chartOfAccounts.Accounts {
		CreateAccount(account.Key, account.Name)
		assert.Equal(t, AccountStore[account.Key].Key, account.Key)
		assert.Equal(t, AccountStore[account.Key].Name, account.Name)
	}
	assert.Equal(t, len(AccountStore), 3)

}

type Transactions struct {
	Types []LedgerTransactionType `json:"types"`
}

type Root struct {
	Transactions Transactions `json:"transactions"`
}

func TestAddTransactionEntry(t *testing.T) {
	loadAccounts()

	root := &Root{}
	err := json.Unmarshal([]byte(ledgerTransactionsJson), root)
	assert.Nil(t, err)

	for _, tt := range root.Transactions.Types {
		params := map[string]string{
			"sales_before_tax": "10000",
			"tax_payable":      "500",
		}
		assert.Equal(t, tt.Type, "sell_something")
		transaction := CreateTransaction("entry-1", "ledger-1", tt.Type, tt.Entries, params)
		assert.Equal(t, len(transaction.entries), 3)
		assert.Equal(t, transaction.entries[0].account.Name, "sales_to_bank")
		assert.Equal(t, transaction.entries[1].account.Name, "income-root-bank")
		assert.Equal(t, transaction.entries[2].account.Name, "tax_payable")

		assert.Equal(t, transaction.entries[0].amount, big.NewInt(10500))
		assert.Equal(t, transaction.entries[1].amount, big.NewInt(10000))
		assert.Equal(t, transaction.entries[2].amount, big.NewInt(500))

	}

	// Further validation based on expected behavior of addTransactionEntry function
}

func TestTransactionFromInput(t *testing.T) {
	loadAccounts() // Ensure accounts are loaded

	// Unmarshal transactionInput into TransactionInput struct
	var input TransactionInput
	err := json.Unmarshal([]byte(transactionInput), &input)
	assert.Nil(t, err)

	// Extracting transaction type and parameters
	transactionType := input.Type
	params := input.Parameters

	// Ensure we got the right transaction type from the input
	assert.Equal(t, transactionType, "sell_something")

	// Load transaction lines based on the transaction type
	root := &Root{}
	err = json.Unmarshal([]byte(ledgerTransactionsJson), root)
	assert.Nil(t, err)

	// Find the correct transaction type from the loaded ledger transactions
	var tt LedgerTransactionType
	for _, transaction := range root.Transactions.Types {
		if transaction.Type == transactionType {
			tt = transaction
			break
		}
	}

	// Creating and validating the transaction
	transaction := CreateTransaction("entry-from-input", "ledger-from-input", tt.Type, tt.Entries, params)
	assert.Equal(t, len(transaction.entries), 3)
	assert.Equal(t, transaction.entries[0].account.Name, "sales_to_bank")
	assert.Equal(t, transaction.entries[1].account.Name, "income-root-bank")
	assert.Equal(t, transaction.entries[2].account.Name, "tax_payable")
	assert.Equal(t, transaction.entries[0].amount, big.NewInt(10500))
	assert.Equal(t, transaction.entries[1].amount, big.NewInt(10000))
	assert.Equal(t, transaction.entries[2].amount, big.NewInt(500))
}
