package accounts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"training/errors"

	"github.com/julienschmidt/httprouter"
)

func CreateAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

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

	accountid, err := Create(userid, int(balance))
	if err != nil {
		err = errors.Wrap(err, "Error Creating Account")
		e := errors.New()
		e.Error = err
		e.UserError = errors.UserError{http.StatusInternalServerError, "Internal Server Error"}
		panic(e)
	}

	var repl = struct {
		AccountId string `json:"accountid"`
	}{accountid}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repl)
}

func CreditAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
func DebitAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
func AccountList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	list, err := List()
	if err != nil {
		err = errors.Wrap(err, "Unmarshal Error")
		e := errors.New()
		e.Error = err
		e.UserError = errors.UserError{http.StatusInternalServerError, "Internal Server Error"}
		panic(e)
	}

	// fmt.Printf("List: Length is %d\n", len(list))
	var repl = struct {
		Accounts []SAccount `json:"accounts"`
	}{list}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repl)
}

func AccountActivity(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	accountid := params.ByName("id")
	// fmt.Printf("AccountActivity: %s", accountid)
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
