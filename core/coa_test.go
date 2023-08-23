package core

import (
	"encoding/json"
	"fmt"
	"ledger/common"
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
					"amount": "{{.sales_before_tax}} + {{.tax_payable}}",
					"direction" : "Debit"
				},
				{
					"key": "income-root",
					"account": "income-root",
					"amount": "{{.sales_before_tax}}",
					"direction" : "Credit"
				},
				{
					"key": "tax_payable",
					"account": "tax_payable",
					"amount": "{{.tax_payable}}",
					"direction" : "Credit"
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

func TestCreateAccountWithChildren(t *testing.T) {

	accountJSON := `{
	"accounts": [{
		"key": "parent",
		"name": "Parent Account",
		"children": ["child1", "child2"]
	}]
	}`

	// Unmarshal the JSON into the AccountInput struct
	chartOfAccounts := &ChartOfAccounts{}
	err := json.Unmarshal([]byte(accountJSON), &chartOfAccounts)
	assert.Nil(t, err)

	// Create the account with its children
	for _, account := range chartOfAccounts.Accounts {
		account := CreateAccount(account.Key, account.Name, account.Children...)

		// Validate that the parent account is correctly added to the store
		assert.NotNil(t, AccountStore["parent"])
		assert.Equal(t, "parent", AccountStore["parent"].Key)
		assert.Equal(t, "Parent Account", AccountStore["parent"].Name)

		// Validate that the children are correctly added to the store
		assert.NotNil(t, AccountStore["parent/child1"])
		assert.Equal(t, "parent/child1", AccountStore["parent/child1"].Key)

		assert.NotNil(t, AccountStore["parent/child2"])
		assert.Equal(t, "parent/child2", AccountStore["parent/child2"].Key)

		// Validate that the children are correctly linked to the parent
		assert.Equal(t, 2, len(account.Children))
		assert.Equal(t, "parent/child1", account.Children[0].Key)
		assert.Equal(t, "parent/child2", account.Children[1].Key)
	}
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
		assert.Equal(t, transaction.entries[0].Account.Name, "sales_to_bank")
		assert.Equal(t, transaction.entries[1].Account.Name, "income-root-bank")
		assert.Equal(t, transaction.entries[2].Account.Name, "tax_payable")

		assert.Equal(t, transaction.entries[0].Amount, big.NewInt(10500))
		assert.Equal(t, transaction.entries[1].Amount, big.NewInt(10000))
		assert.Equal(t, transaction.entries[2].Amount, big.NewInt(500))

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
	assert.Equal(t, transaction.entries[0].Account.Name, "sales_to_bank")
	assert.Equal(t, transaction.entries[1].Account.Name, "income-root-bank")
	assert.Equal(t, transaction.entries[2].Account.Name, "tax_payable")
	assert.Equal(t, transaction.entries[0].Amount, big.NewInt(10500))
	assert.Equal(t, transaction.entries[1].Amount, big.NewInt(10000))
	assert.Equal(t, transaction.entries[2].Amount, big.NewInt(500))
	assert.Equal(t, transaction.entries[0].Direction, common.Debit)
	assert.Equal(t, transaction.entries[1].Direction, common.Credit)
	assert.Equal(t, transaction.entries[2].Direction, common.Credit)
}
