package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Account represents a bank account
type Account struct {
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
}

var accounts = map[string]Account{
	"12345": {AccountNumber: "12345", Balance: 1000.0},
	"67890": {AccountNumber: "67890", Balance: 500.0},
}

func main() {
	r := mux.NewRouter()

	// Create a new account
	r.HandleFunc("/account", createAccount).Methods("POST")

	// Check account balance
	r.HandleFunc("/account/{accountNumber}", checkBalance).Methods("GET")

	// Close an account
	r.HandleFunc("/account/{accountNumber}", closeAccount).Methods("DELETE")

	http.Handle("/", r)

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	var newAccount Account
	err := json.NewDecoder(r.Body).Decode(&newAccount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accounts[newAccount.AccountNumber] = newAccount
	w.WriteHeader(http.StatusCreated)
}

func checkBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountNumber := vars["accountNumber"]

	account, found := accounts[accountNumber]
	if !found {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	response, _ := json.Marshal(account)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func closeAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountNumber := vars["accountNumber"]

	_, found := accounts[accountNumber]
	if !found {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	delete(accounts, accountNumber)
	w.WriteHeader(http.StatusNoContent)
}
