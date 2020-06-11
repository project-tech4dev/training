package bank

import (
	"time"
)

type SBank struct {
	Version         int               `json:"version"`
	AccountSerial   int               `json:"accountserial"`
	AccountsCatalog map[string]string `json:"accountscatalog"`
	NextUserID      int               `json:"nextuserid"`
	UserIDCatalog   map[string]int    `json:"useridcatalog"`
	UserNameCatalog map[string]int    `json:"usernamecatalog"`
	Users           []SUser           `json:"users"`
}

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

type Role string

const (
	User        Role = "User"
	BankManager Role = "BankManager"
)

type SUser struct {
	ID        string     `json:"userid"`
	Name      string     `json:"name"`
	UserName  string     `json:"username"`
	Password  string     `json:"password"`
	Role      Role       `json:"role"`
	CreatedAt time.Time  `json:"createdat"`
	Accounts  []SAccount `json:"accounts"`
}

type Session struct {
	Token  string
	UserID string
	Role   Role
}

type AccountListDisplay struct {
	Name      string    `json:"name"`
	AccountID string    `json:"accountid"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"createdat"`
}
