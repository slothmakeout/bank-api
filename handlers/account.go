package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"bank-api/data/models"
	services "bank-api/pkg/account"

	"github.com/gorilla/mux"
)

type AccountHandler struct {
	l           *log.Logger
	accountServ *services.AccountService
}

func NewAccountHandler(l *log.Logger, accountServ *services.AccountService) *AccountHandler {
	return &AccountHandler{l, accountServ}
}

func (ah *AccountHandler) GetAllAccounts(rw http.ResponseWriter, req *http.Request) {
	ah.l.Println("Handled GET Accounts")

	accounts, err := ah.accountServ.GetAllAccounts()
	if err != nil {
		http.Error(rw, "Unable to get accounts list", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(accounts); err != nil {
		http.Error(rw, "Unable to encode JSON", http.StatusInternalServerError)
		return
	}
}

func (ah *AccountHandler) GetAccountById(rw http.ResponseWriter, req *http.Request) {
	ah.l.Println("Handled GET Account by ID")

	vars := mux.Vars(req)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(rw, "ID not provided in the request", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, "Invalid ID format", http.StatusBadRequest)
		return
	}

	account, err := ah.accountServ.GetAccountById(id)
	if err != nil {
		http.Error(rw, "Unable to get account by ID", http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(account); err != nil {
		http.Error(rw, "Unable to encode JSON", http.StatusInternalServerError)
		return
	}
}

func (ah *AccountHandler) AddAccount(rw http.ResponseWriter, req *http.Request) {
	ah.l.Println("Handled POST Account")
	var account models.Account
	err := json.NewDecoder(req.Body).Decode(&account)
	if err != nil {
		http.Error(rw, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	err = ah.accountServ.AddAccount(&account)
	if err != nil {
		http.Error(rw, "Unable to add account", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte("OK"))
}

func (ah *AccountHandler) UpdateAccount(rw http.ResponseWriter, req *http.Request) {
	ah.l.Println("Handled PUT Account")

	vars := mux.Vars(req)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(rw, "ID not provided in the request", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updatedAccount models.Account
	err = json.NewDecoder(req.Body).Decode(&updatedAccount)
	if err != nil {
		http.Error(rw, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	err = ah.accountServ.UpdateAccount(id, &updatedAccount)
	if err != nil {
		http.Error(rw, "Unable to update account", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
}

func (ah *AccountHandler) DeleteAccount(rw http.ResponseWriter, req *http.Request) {
	ah.l.Println("Handled DELETE Account")

	vars := mux.Vars(req)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(rw, "ID not provided in the request", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = ah.accountServ.DeleteAccount(id)
	if err != nil {
		http.Error(rw, "Unable to delete account", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
}
