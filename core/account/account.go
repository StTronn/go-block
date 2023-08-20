package account

const (
	AddressLength = 32
)

type Account struct {
	Address
}

type Address [AddressLength]byte
