package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var Bank SBank

const Bank_Json = "./data/bank.json"

// Init initializes the internal structure of the bank from a JSON file in the current directory.
// If the file doesn't exist, it creates the file and initializes an empty bank structure.
func Init() error {
	dir := filepath.Dir(Bank_Json)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(Bank_Json); os.IsNotExist(err) {
		_, err = os.Create(Bank_Json)
		if err != nil {
			fmt.Printf("Error creating %s\n", Bank_Json)
			return err
		}
	}
	b, err := ioutil.ReadFile(Bank_Json)
	if err != nil {
		fmt.Printf("Error reading %s\n", Bank_Json)
		return err
	}
	if len(b) > 0 {
		err = json.Unmarshal(b, &Bank)
		if err != nil {
			fmt.Printf("Error unmarshaling %s\n", Bank_Json)
			return err
		}
	} else {
		Bank = SBank{10001, make(map[string]int), make([]SAccount, 0)}
		flush()
	}
	return nil
}

// flush is called after every modifying operation to write the bank structure to disk
func flush() error {
	b, err := json.Marshal(Bank)
	if err != nil {
		fmt.Printf("Error unmarshaling %s\n", Bank_Json)
		return err
	}
	ioutil.WriteFile(Bank_Json, b, 0666)
	return nil
}

// Create creates a new bank account.
// Inputs: userid string
//				 balance int (in INR * 100) E.g., a value of 300000 means a balance of INR 3000
// Returns: new Account ID
// A user can have multiple accounts
func Create(userid string, balance int) (string, error) {
	defer flush()
	Bank.AccountSerial += 1

	transaction := STransaction{balance, "OB", balance, time.Now().UTC()}
	transactions := make([]STransaction, 0)
	transactions = append(transactions, transaction)
	accountid := fmt.Sprintf("%d", Bank.AccountSerial)
	account := SAccount{accountid, userid, time.Now().UTC(), balance, transactions}
	Bank.Accounts = append(Bank.Accounts, account)
	Bank.AccountsCatalog[accountid] = len(Bank.Accounts) - 1
	return accountid, nil
}

// Credit credits the given account with the given amount.
// Inputs: accountid string
//				 amount int
// Returns: the new balance
// If the accountid doesn't exist, it returns an error.
func Credit(accountid string, amount int) (int, error) {
	defer flush()
	accountix, ok := Bank.AccountsCatalog[accountid]
	if !ok {
		// no such account
		err := errors.New("ENOSUCHACCOUNT")
		return 0, err
	}
	account := &(Bank.Accounts[accountix])
	balance := account.Balance
	balance += amount
	account.Activity = append(account.Activity, STransaction{amount, "CR", balance, time.Now().UTC()})
	account.Balance = balance
	// fmt.Printf("CREDIT: %d, %v, %v", amount, account, Bank.Accounts[accountix])
	return balance, nil
}

// Debit debits the given account with the given amount.
// Inputs: accountid string
//				 amount int
// Returns: the new balance
// If the accountid doesn't exist, it returns an error.
// It returns an error if the account doesn't have sufficient funds
func Debit(accountid string, amount int) (int, error) {
	defer flush()
	// fmt.Printf("Debit: AccountId: %s, Amount: %d\n", accountid, amount)
	accountix, ok := Bank.AccountsCatalog[accountid]
	if !ok {
		// no such account
		err := errors.New("ENOSUCHACCOUNT")
		return 0, err
	}
	account := &(Bank.Accounts[accountix])
	balance := account.Balance
	if balance < amount {
		err := errors.New("ENOTENOUGHBALANCE")
		return 0, err
	}
	balance -= amount
	account.Activity = append(account.Activity, STransaction{amount, "DB", balance, time.Now().UTC()})
	account.Balance = balance
	// fmt.Printf("DEBIT: %d, %v\n%v\n", amount, account, Bank.Accounts[accountix])
	return balance, nil
}

// Activity returns the account activity for the account
// Inputs: accountid string
// Returns: JSON Object
// 		{ accountid: accountid
//		}
func Activity(accountid string) ([]STransaction, error) {
	accountix, ok := Bank.AccountsCatalog[accountid]
	if !ok {
		// no such account
		err := errors.New("ENOSUCHACCOUNT")
		return nil, err
	}
	account := Bank.Accounts[accountix]
	return account.Activity, nil
}

// List returns the list of accounts
func List() ([]SAccount, error) {
	return Bank.Accounts, nil
	// select user.name, account.id, account.balance from user, account account left outer join user on account.userid = user.id
}
