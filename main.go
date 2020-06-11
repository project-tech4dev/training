package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"training/bank"
	"training/errors"

	"github.com/julienschmidt/httprouter"
)

const (
	PORT = 9765
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.SetOutput(os.Stderr)
	bank.Init()
	router := httprouter.New()
	router.PanicHandler = ErrorHandler
	router.POST("/accounts", bank.CreateAccountHandler)
	router.GET("/accounts", bank.AccountListHandler)
	router.GET("/accounts/:id", bank.AccountActivityHandler)
	router.POST("/accounts/credit", bank.CreditAccountHandler)
	router.POST("/accounts/debit", bank.DebitAccountHandler)

	router.POST("/users", bank.CreateUserHandler)
	router.GET("/users", bank.UserListHandler)
	router.GET("/users/:id", bank.UserHandler)

	router.POST("/login", bank.LoginHandler)

	log.Printf("Listening on localhost:%d\n", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), router))
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	if _, ok := err.(*errors.Error); ok {
		e := err.(*errors.Error)
		log.Printf("Error: %s - %s\n", e.Error, e.UserError.ErrorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(e.UserError.ErrorCode)
		json.NewEncoder(w).Encode(e.UserError)
	}
	return
}
