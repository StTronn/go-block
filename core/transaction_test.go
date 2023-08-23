package core

import (
	"ledger/common"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	account := &Account{
		Key:  "test-key",
		Name: "Test Account",
	}
	amount := big.NewInt(1000)

	entry := NewEntry(account, amount, common.Debit, common.Posted)
	assert.NotNil(t, entry)
	assert.Equal(t, "test-key", entry.Account.Key)
	assert.Equal(t, amount, entry.Amount)
	assert.Equal(t, common.Debit, entry.Direction)
}

func TestNewTransaction(t *testing.T) {
	entry1 := Entries{
		id:        "test-id-1",
		Account:   &Account{},
		Amount:    big.NewInt(1000),
		Direction: common.Debit,
	}

	entry2 := Entries{
		id:        "test-id-2",
		Account:   &Account{},
		Amount:    big.NewInt(500),
		Direction: common.Credit,
	}

	transaction := NewTransaction(entry1, entry2)
	assert.NotNil(t, transaction)
	assert.Equal(t, 2, len(transaction.entries))
}

// ... add more tests as needed.
