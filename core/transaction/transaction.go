package transaction

import (
	"math/big"
)

type transaction struct {
	signer string
	to     string
	value  *big.Int
}

type TxData interface {
	txType() byte
}
