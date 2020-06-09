package accounts

import (
	"time"
)

type STransaction struct {
	Amount    int       `json:"amount"`
	Operation string    `json:"operation"`
	Balance   int       `json:"balance"` // running balance
	CreatedAt time.Time `json:"createdat"`
}

type SAccount struct {
	ID        string         `json:"id"`
	UserID    string         `json:"userid"`
	CreatedAt time.Time      `json:"createdat"`
	Balance   int            `json:"balance"`
	Activity  []STransaction `json:"activity"`
}

type SBank struct {
	AccountSerial   int            `json:"accountserial"`
	AccountsCatalog map[string]int `json:"accountscatalog"`
	Accounts        []SAccount     `json:"accounts"`
}
