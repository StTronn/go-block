package core

import (
	"encoding/json"
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
					"account": "bank",
					"amount": "{{sales_before_tax}} + {{tax_payable}}"
				},
				{
					"key": "income_from_sales_before_tax",
					"account": "income-root/sales",
					"amount": "{{sales_before_tax}}"
				},
				{
					"key": "tax_payable",
					"account": "tax_payables",
					"amount": "{{tax_payable}}"
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
			"key": "income_from_sales_before_tax",
			"name": "income_from_sales_before_tax"
		},
		{
			"key": "tax_payable",
			"name": "tax_payable"
		}
	]

}
`

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
	root := &Root{}
	err := json.Unmarshal([]byte(ledgerTransactionsJson), root)
	assert.Nil(t, err)

	// for _, tt := range root.Transactions.Types {
	// 	params := map[string]string{
	// 		"sales_before_tax": "10000",
	// 		"tax_payable":      "500",
	// 	}
	// 	transaction := addTransactionEntry("entry-1", "ledger-1", tt.Type, tt.Lines, params)
	// 	assert.Equal(t, len(transaction.entries), 2)
	// 	assert.Equal(t, transaction.entries[0].account, "bank")
	// 	assert.Equal(t, transaction.entries[1].account, "income-root/sales")
	// 	assert.Equal(t, transaction.entries[2].account, "tax_paybles")

	// }

	// Further validation based on expected behavior of addTransactionEntry function
}
