package core

import (
	"ledger/common"
	"math/big"

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

// newTransaction creates a new Transaction with a unique id and a slice of entries.
func newTransaction(entries ...Entries) *Transaction {
	return &Transaction{
		id:      xid.New().String(),
		entries: entries,
	}
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
