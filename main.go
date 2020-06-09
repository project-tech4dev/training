package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"training/accounts"
	"training/errors"

	"github.com/julienschmidt/httprouter"
)

const (
	PORT = 9765
)

func main() {
	accounts.Init()
	router := httprouter.New()
	router.PanicHandler = ErrorHandler
	router.POST("/accounts", accounts.CreateAccount)
	router.GET("/accounts", accounts.AccountList)
	router.GET("/accounts/:id", accounts.AccountActivity)
	router.POST("/accounts/credit", accounts.CreditAccount)
	router.POST("/accounts/debit", accounts.DebitAccount)

	log.Printf("Listening on localhost:%d\n", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), router))
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	if _, ok := err.(*errors.Error); ok {
		e := err.(*errors.Error)
		fmt.Printf("%s\n", e.Error)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(e.UserError.ErrorCode)
		json.NewEncoder(w).Encode(e.UserError)
	}
	return
}
