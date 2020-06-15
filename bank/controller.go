package bank

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"training/errors"

	"github.com/julienschmidt/httprouter"
)

var tokens = make(map[string]Session)

const Bearer = "Bearer "

func getSession(r *http.Request) *Session {
	authHeader := r.Header.Get("Authorization")
	if (len(authHeader) == 0) || (len(authHeader) < len(Bearer)) {
		return nil
	}
	token := authHeader[len("Bearer "):]
	s, ok := tokens[token]
	if !ok {
		return nil
	}

	return &s
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	body, err := ioutil.ReadAll(r.Body)
	var u map[string]interface{}
	err = json.Unmarshal(body, &u)
	if err != nil {
		err = errors.Wrap(err, "Unmarshal Error")
		e := errors.New()
		e.Error = err
		e.UserError = errors.UserError{http.StatusInternalServerError, "Internal Server Error"}
		panic(e)
	}

	fullname := GetStringParam(u, "fullname", false)
	username := GetStringParam(u, "username", false)
	password := GetStringParam(u, "password", false)

	session := getSession(r)
	log.Printf("CUH: Session: %v\n", session)

	// if (session == nil) || (session.Role != BankManager) {
	// 	e := errors.New()
	// 	e.UserError = errors.UserError{http.StatusForbidden, "Permission Denied"}
	// 	panic(e)

	// }
	_, err = FindUserByName(username)
	if err == nil {
		e := errors.New()
		e.Error = err
		e.UserError = errors.UserError{http.StatusForbidden, "User Exists"}
		panic(e)
	}

	user := CreateUser(fullname, username, password)
	log.Printf("%s %v\n", user.ID, user)
	AddUser(user)

	var repl = struct {
		UserID string `json:"userid"`
	}{user.ID}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repl)
}

func LoginHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	body, err := ioutil.ReadAll(r.Body)
	var u map[string]interface{}
	err = json.Unmarshal(body, &u)
	if err != nil {
		err = errors.Wrap(err, "Unmarshal Error")
		e := errors.New()
		e.Error = err
		e.UserError = errors.UserError{http.StatusInternalServerError, "Internal Server Error"}
		panic(e)
	}

	username := GetStringParam(u, "username", false)
	password := GetStringParam(u, "password", false)

	user, err := FindUserByName(username)
	if err != nil {
		if err.Error() == "ENOSUCHUSER" {
			err = errors.Wrapf(err, "Error finding user %s", username)
			e := errors.New()
			e.Error = err
			e.UserError = errors.UserError{http.StatusForbidden, "Incorrect Username or Password"}
			panic(e)
		}
		err = errors.Wrapf(err, "Error finding user %s", username)
		e := errors.New()
		e.Error = err
		e.UserError = errors.UserError{http.StatusInternalServerError, "Internal Server Error"}
		panic(e)
	}

	if user.Password != password {
		e := errors.New()
		e.Error = err
		e.UserError = errors.UserError{http.StatusForbidden, "Incorrect Username or Password"}
		panic(e)
	}

	h := sha256.New()
	h.Write([]byte(user.Password))
	hash := fmt.Sprintf("%x", h.Sum(nil))

	tokens[hash] = Session{hash, user.ID, user.Role}
	log.Printf("%s %s\n", password, hash)

	var repl = struct {
		UserID string `json:"id"`

		Token string `json:"token"`
	}{
		user.ID,
		hash,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repl)
}

func CreateAccountHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	body, err := ioutil.ReadAll(r.Body)
	var u map[string]interface{}
	err = json.Unmarshal(body, &u)
	if err != nil {
		err = errors.Wrap(err, "Unmarshal Error")
		e := errors.New()
		e.Error = err
		e.UserError = errors.UserError{http.StatusInternalServerError, "Internal Server Error"}
		panic(e)
	}

	userid := GetStringParam(u, "userid", false)
	balance := GetFloatParam(u, "balance")

	user, err := FindUserByID(userid)
	if err != nil {
		if err.Error() == "ENOSUCHUSER" {
			e := errors.New()
			e.Error = err
			e.UserError = errors.UserError{http.StatusBadRequest, "No Such User"}
			panic(e)

		}
	}

	log.Printf("CAH: User: %v\n", user)
	if user.Role == BankManager {
		// Can't create account for the branch manager
		e := errors.New()
		e.UserError = errors.UserError{http.StatusForbidden, "Permission Denied"}
		panic(e)
	}

	session := getSession(r)
	log.Printf("Session: %v\n", session)
	if (session == nil) || (session.Role != BankManager) {
		// Only bank manager can create accounts
		e := errors.New()
		e.UserError = errors.UserError{http.StatusForbidden, "Permission Denied"}
		panic(e)
	}

	account, err := CreateAccount(userid, int(balance))
	if err != nil {
		err = errors.Wrap(err, "Error Creating Account")
		e := errors.New()
		e.Error = err
		e.UserError = errors.UserError{http.StatusInternalServerError, "Internal Server Error"}
		panic(e)
	}

	AddAccount(user, account)

	var repl = struct {
		AccountId string `json:"accountid"`
		Balance   int    `json:"balance"`
	}{account.ID, int(balance)}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repl)
}

func AccountListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	list, err := List()
	if err != nil {
		err = errors.Wrap(err, "Unmarshal Error")
		e := errors.New()
		e.Error = err
		e.UserError = errors.UserError{http.StatusInternalServerError, "Internal Server Error"}
		panic(e)
	}

	session := getSession(r)
	if (session == nil) || (session.Role != BankManager) {
		e := errors.New()
		e.UserError = errors.UserError{http.StatusForbidden, "Permission Denied"}
		panic(e)
	}

	// log.Printf("List: Length is %d\n", len(list))
	var repl = struct {
		Accounts []AccountListDisplay `json:"accounts"`
	}{list}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repl)
}

func UserListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type UL struct {
		UserID      string `json:"id"`
		FullName    string `json:"fullname"`
		UserName    string `json:"username"`
		CreatedAt   string `json:"createdat"`
		NumAccounts int    `json:"numaccounts"`
	}

	users := make([]UL, 0)

	session := getSession(r)
	if (session == nil) || (session.Role != BankManager) {
		e := errors.New()
		e.UserError = errors.UserError{http.StatusForbidden, "Permission Denied"}
		panic(e)
	}

	for _, u := range bank.Users {
		if u.Role == BankManager {
			continue
		}
		users = append(users, UL{u.ID, u.Name, u.UserName, u.CreatedAt.Format(time.RFC822Z), len(u.Accounts)})
	}
	// log.Printf("List: Length is %d\n", len(list))
	var repl = struct {
		Users []UL `json:"users"`
	}{users}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repl)
}

func UserHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userid := params.ByName("id")

	log.Printf("Userid: %s\n", userid)
	user, err := FindUserByID(userid)
	if err != nil {
		err = errors.Wrap(err, "Unmarshal Error")
		e := errors.New()
		e.Error = err
		e.UserError = errors.UserError{http.StatusInternalServerError, "No Such User"}
		panic(e)
	}

	session := getSession(r)
	log.Printf("Session: %v", session)
	if (session == nil) || ((session.Role != BankManager) && (session.UserID != userid)) {
		e := errors.New()
		e.UserError = errors.UserError{http.StatusForbidden, "Permission Denied"}
		panic(e)
	}

	// dereference the pointer and then change the field to send to the user so the changedoesn't persist in the store
	repluser := *user
	repluser.Password = "xxxxxxxxxxxxxxxxxxxxxxxx"

	var repl = struct {
		User SUser `json:"user"`
	}{
		repluser,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repl)
}

func CreditAccountHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	var u map[string]interface{}
	err = json.Unmarshal(body, &u)
	if err != nil {
		err = errors.Wrap(err, "Unmarshal Error")
		e := errors.New()
		e.Error = err
		e.UserError = errors.UserError{http.StatusInternalServerError, "Internal Server Error"}
		panic(e)
	}

	accountid := GetStringParam(u, "accountid", false)
	amount := GetFloatParam(u, "amount")

	account, err := FindAccount(accountid)
	if err != nil {
		e := errors.New()
		e.UserError = errors.UserError{http.StatusForbidden, "No Such Account"}
		panic(e)
	}
	session := getSession(r)
	if (session == nil) || (session.UserID != account.UserID) {
		e := errors.New()
		e.UserError = errors.UserError{http.StatusForbidden, "Permission Denied"}
		panic(e)
	}
	balance, err := Credit(accountid, int(amount))
	if err != nil {
		e := errors.New()
		if err.Error() == "ENOSUCHACCOUNT" {
			e.UserError = errors.UserError{http.StatusBadRequest, "No Such Account"}
		} else {
			e.UserError = errors.UserError{http.StatusInternalServerError, "Internal Server Error"}
		}
		err = errors.Wrap(err, "Credit Error")
		e.Error = err
		panic(e)
	}

	var repl = struct {
		Balance int `json:"balance"`
	}{balance}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repl)

	// get account id
	// get amount
	// balance, err := Credit(accountid, amount)
	// if err != nil {
	// 	return error to user
	// }
	// return balance
}
func DebitAccountHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	var u map[string]interface{}
	err = json.Unmarshal(body, &u)
	if err != nil {
		err = errors.Wrap(err, "Unmarshal Error")
		e := errors.New()
		e.Error = err
		e.UserError = errors.UserError{http.StatusInternalServerError, "Internal Server Error"}
		panic(e)
	}

	accountid := GetStringParam(u, "accountid", false)
	amount := GetFloatParam(u, "amount")

	account, err := FindAccount(accountid)
	if err != nil {
		e := errors.New()
		e.UserError = errors.UserError{http.StatusForbidden, "No Such Account"}
		panic(e)
	}
	session := getSession(r)
	if (session == nil) || (session.UserID != account.UserID) {
		e := errors.New()
		e.UserError = errors.UserError{http.StatusForbidden, "Permission Denied"}
		panic(e)
	}

	balance, err := Debit(accountid, int(amount))
	if err != nil {
		e := errors.New()

		if err.Error() == "ENOSUCHACCOUNT" {
			e.UserError = errors.UserError{http.StatusBadRequest, "No Such Account"}
		} else if err.Error() == "ENOTENOUGHBALANCE" {
			e.UserError = errors.UserError{http.StatusBadRequest, "Not Enough Balance"}
		} else {
			e.UserError = errors.UserError{http.StatusInternalServerError, "Internal Server Error"}
		}
		err = errors.Wrap(err, "Debit Error")
		e.Error = err
		panic(e)
	}

	var repl = struct {
		Balance int `json:"balance"`
	}{balance}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repl)
}

func AccountActivityHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	accountid := params.ByName("id")

	account, err := FindAccount(accountid)
	if err != nil {
		e := errors.New()
		e.UserError = errors.UserError{http.StatusForbidden, "No Such Account"}
		panic(e)
	}
	session := getSession(r)

	if (session == nil) || ((session.Role != BankManager) && (session.UserID != account.UserID)) {
		e := errors.New()
		e.UserError = errors.UserError{http.StatusForbidden, "Permission Denied"}
		panic(e)
	}

	// log.Printf("AccountActivity: %s", accountid)
	activity, err := Activity(accountid)
	if err != nil {
		e := errors.New()
		if err.Error() == "ENOSUCHACCOUNT" {
			e.UserError = errors.UserError{http.StatusBadRequest, "No Such Account"}
		} else {
			e.UserError = errors.UserError{http.StatusInternalServerError, "Internal Server Error"}
		}
		e.Error = err
		err = errors.Wrap(err, "Error Retrieving Activity")
		panic(e)
	}
	var repl = struct {
		AccountID string         `json:"accountid"`
		Activity  []STransaction `json:"activity"`
	}{accountid, activity}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repl)
}

func GetStringParam(u map[string]interface{}, name string, emptyOk bool) string {
	_, ok := u[name]
	if !ok {
		if !emptyOk {
			err := errors.NewError("")
			s := fmt.Sprintf("Received No %s field", name)
			err = errors.Wrap(err, s)
			e := errors.New()
			e.Error = err
			e.UserError = errors.UserError{http.StatusInternalServerError, s}
			panic(e)
		} else {
			return ""
		}
	}
	n, ok := u[name].(string)
	if (!ok) || (!emptyOk && n == "") {
		err := errors.NewError("")
		err = errors.Wrapf(err, "Received Empty or badly formed %s", name)
		e := errors.New()
		e.Error = err
		e.UserError = errors.UserError{http.StatusInternalServerError, "Internal Server Error"}
		panic(e)
	}
	return n
}

func GetFloatParam(u map[string]interface{}, name string) float64 {
	n, ok := u[name].(float64)
	if !ok {
		err := errors.NewError("")
		err = errors.Wrapf(err, "Error retriving %s", name)
		e := errors.New()
		e.Error = err
		e.UserError = errors.UserError{http.StatusInternalServerError, "Internal Server Error"}
		panic(e)
	}
	return n
}
