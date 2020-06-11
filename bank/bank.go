package bank

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"training/errors"
)

var bank SBank

const Bank_Json = "./data/bank.json"

var permissions = make(map[string][]string)

const Version = 1

// permissions["User"] = [""]

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
			log.Printf("Error creating %s\n", Bank_Json)
			return err
		}
	}
	b, err := ioutil.ReadFile(Bank_Json)
	if err != nil {
		log.Printf("Error reading %s\n", Bank_Json)
		return err
	}
	if len(b) > 0 {
		err = json.Unmarshal(b, &bank)
		if err != nil {
			log.Printf("Error unmarshaling %s\n", Bank_Json)
			return err
		}
	}

	if bank.Version != Version {
		bank = SBank{Version, 10001, make(map[string]string), 1000000, make(map[string]int), make(map[string]int), make([]SUser, 0)}
		user := CreateUser("Stockson Bonds", "bankmanager", "headhoncho")
		user.Role = BankManager
		AddUser(user)
		flush()
	}
	return nil
}

func NextUserID() string {
	bank.NextUserID += 1
	flush()
	return fmt.Sprintf("%d", bank.NextUserID)
}
func NextAccountID() string {
	bank.AccountSerial += 1
	flush()
	return fmt.Sprintf("%d", bank.AccountSerial)
}

func FindUserByName(name string) (*SUser, error) {
	u, ok := bank.UserNameCatalog[name]
	if !ok {
		err := errors.NewError("ENOSUCHUSER")
		return nil, err
	}
	return &bank.Users[u], nil
}

func FindUserByID(name string) (*SUser, error) {
	u, ok := bank.UserIDCatalog[name]
	if !ok {
		err := errors.NewError("ENOSUCHUSER")
		return nil, err
	}
	return &bank.Users[u], nil
}

func FindUserAccount(userid string, accountid string) (*SAccount, error) {
	user, _ := FindUserByID(userid)
	for ix, a := range user.Accounts {
		if a.ID == accountid {
			return &user.Accounts[ix], nil
		}
	}
	err := errors.NewError("ENOSUCHACCOUNT")
	return nil, err
}

func FindAccount(accountid string) (*SAccount, error) {
	u, ok := bank.AccountsCatalog[accountid]
	log.Printf("ACCT: %v, %s\n", u, ok)
	if !ok {
		err := errors.NewError("ENOSUCHACCOUNT")
		return nil, err
	}

	account, _ := FindUserAccount(u, accountid)

	return account, nil
}

func AddUser(user SUser) {
	bank.Users = append(bank.Users, user)
	bank.UserIDCatalog[user.ID] = len(bank.Users) - 1
	bank.UserNameCatalog[user.UserName] = len(bank.Users) - 1
	flush()
}

func AddAccount(user *SUser, account SAccount) {
	user.Accounts = append(user.Accounts, account)
	bank.AccountsCatalog[account.ID] = user.ID
	flush()
}

func flush() error {
	b, err := json.Marshal(bank)
	if err != nil {
		log.Printf("Error unmarshaling %s\n", Bank_Json)
		return err
	}
	ioutil.WriteFile(Bank_Json, b, 0666)
	return nil
}
