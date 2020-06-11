package bank

import (
	"errors"
	"log"
	"time"
)

// flush is called after every modifying operation to write the bank structure to disk

// Create creates a new bank account.
// Inputs: userid string
//				 balance int (in INR * 100) E.g., a value of 300000 means a balance of INR 3000
// Returns: new Account ID
// A user can have multiple accounts
func CreateAccount(userid string, balance int) (SAccount, error) {
	accountid := NextAccountID()
	transaction := STransaction{balance, "OB", balance, time.Now().UTC()}
	transactions := make([]STransaction, 0)
	transactions = append(transactions, transaction)
	account := SAccount{accountid, userid, time.Now().UTC(), balance, transactions}
	return account, nil
}

// Credit credits the given account with the given amount.
// Inputs: accountid string
//				 amount int
// Returns: the new balance
// If the accountid doesn't exist, it returns an error.
func Credit(accountid string, amount int) (int, error) {
	account, err := FindAccount(accountid)
	if err != nil {
		// no such account
		err := errors.New("ENOSUCHACCOUNT")
		return 0, err
	}

	balance := account.Balance
	balance += amount
	account.Activity = append(account.Activity, STransaction{amount, "CR", balance, time.Now().UTC()})
	account.Balance = balance
	flush()
	// log.Printf("DEBIT: %d, %v\n%v\n", amount, account, Bank.Accounts[accountix])
	return balance, nil
}

// Debit debits the given account with the given amount.
// Inputs: accountid string
//				 amount int
// Returns: the new balance
// If the accountid doesn't exist, it returns an error.
// It returns an error if the account doesn't have sufficient funds
func Debit(accountid string, amount int) (int, error) {
	account, err := FindAccount(accountid)
	if err != nil {
		// no such account
		err := errors.New("ENOSUCHACCOUNT")
		return 0, err
	}

	balance := account.Balance
	if balance < amount {
		err := errors.New("ENOTENOUGHBALANCE")
		return 0, err
	}
	balance -= amount
	account.Activity = append(account.Activity, STransaction{amount, "DB", balance, time.Now().UTC()})
	account.Balance = balance
	flush()
	// log.Printf("DEBIT: %d, %v\n%v\n", amount, account, Bank.Accounts[accountix])
	return balance, nil
}

// Activity returns the account activity for the account
// Inputs: accountid string
// Returns: JSON Object
// 		{ accountid: accountid
//		}
func Activity(accountid string) ([]STransaction, error) {
	account, err := FindAccount(accountid)
	log.Printf("%v %s\n", account, accountid)
	if err != nil {
		// no such account
		err := errors.New("ENOSUCHACCOUNT")
		return nil, err
	}

	return account.Activity, nil
}

// List returns the list of accounts
func List() ([]AccountListDisplay, error) {

	accounts := make([]AccountListDisplay, 0)
	for _, u := range bank.Users {
		for _, a := range u.Accounts {
			accounts = append(accounts, AccountListDisplay{u.Name, a.ID, a.Balance, a.CreatedAt})
		}
	}
	return accounts, nil
	// select user.name, account.id, account.balance from user, account account left outer join user on account.userid = user.id
}
