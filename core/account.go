package core

const (
	AddressLength = 32
)

type Account struct {
	Key      string    `json:"key"`
	Name     string    `json:"name,omitempty"`
	Children []Account `json:"children,omitempty"`
}

// Persistent Layer -------------------------------->
type AccountsStore map[string]*Account

var AccountStore = make(AccountsStore)
