package common

import "fmt"

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

func (d *Direction) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case `"Debit"`:
		*d = Debit
	case `"Credit"`:
		*d = Credit
	default:
		return fmt.Errorf("invalid direction: %s", string(b))
	}
	return nil
}

type Status int

const (
	Pending Status = iota
	Posted
)

func (d *Status) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case `"Pending"`:
		*d = Pending
	case `"Posted"`:
		*d = Posted
	default:
		return fmt.Errorf("invalid direction: %s", string(b))
	}
	return nil
}
