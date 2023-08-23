package common

const (
	HashLength    = 32
	AddressLength = 20
)

type Hash [HashLength]byte
type Direction int

const (
	Debit Direction = iota
	Credit
)

type Status int

const (
	Pending Status = iota
	Posted
)
