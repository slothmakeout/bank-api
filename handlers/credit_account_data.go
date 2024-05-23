package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"bank-api/data/models"
	services "bank-api/pkg/credit_account_data"

	"github.com/gorilla/mux"
)

type CreditAccountDataHandler struct {
	l          *log.Logger
	creditServ *services.CreditAccountDataService
}

func NewCreditAccountDataHandler(l *log.Logger, creditServ *services.CreditAccountDataService) *CreditAccountDataHandler {
	return &CreditAccountDataHandler{l, creditServ}
}

func (cah *CreditAccountDataHandler) GetAllCreditAccounts(rw http.ResponseWriter, req *http.Request) {
	cah.l.Println("Handled GET Credit Accounts")

	creditAccounts, err := cah.creditServ.GetAllCreditAccounts()
	if err != nil {
		http.Error(rw, "Unable to get credit accounts list", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(creditAccounts); err != nil {
		http.Error(rw, "Unable to encode JSON", http.StatusInternalServerError)
		return
	}
}

func (cah *CreditAccountDataHandler) GetCreditAccountByID(rw http.ResponseWriter, req *http.Request) {
	cah.l.Println("Handled GET Credit Account by ID")

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

	creditAccount, err := cah.creditServ.GetCreditAccountById(id)
	if err != nil {
		http.Error(rw, "Unable to get credit account by ID", http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(creditAccount); err != nil {
		http.Error(rw, "Unable to encode JSON", http.StatusInternalServerError)
		return
	}
}

func (cah *CreditAccountDataHandler) AddCreditAccount(rw http.ResponseWriter, req *http.Request) {
	cah.l.Println("Handled POST Credit Account")

	var creditAccountData models.CreditAccountData
	err := json.NewDecoder(req.Body).Decode(&creditAccountData)
	if err != nil {
		http.Error(rw, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	err = cah.creditServ.AddCreditAccount(&creditAccountData)
	if err != nil {
		http.Error(rw, "Unable to add credit account", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte("OK"))
}

func (cah *CreditAccountDataHandler) UpdateCreditAccount(rw http.ResponseWriter, req *http.Request) {
	cah.l.Println("Handled PUT Credit Account")

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

	var updatedCreditAccount models.CreditAccountData
	err = json.NewDecoder(req.Body).Decode(&updatedCreditAccount)
	if err != nil {
		http.Error(rw, "Unable to decode JSON", http.StatusBadRequest)
		return
	}

	err = cah.creditServ.UpdateCreditAccount(id, &updatedCreditAccount)
	if err != nil {
		http.Error(rw, "Unable to update credit account", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
}

func (cah *CreditAccountDataHandler) DeleteCreditAccount(rw http.ResponseWriter, req *http.Request) {
	cah.l.Println("Handled DELETE Credit Account")

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

	err = cah.creditServ.DeleteCreditAccount(id)
	if err != nil {
		http.Error(rw, "Unable to delete credit account", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (cah *CreditAccountDataHandler) GetAccountWithCreditData(rw http.ResponseWriter, req *http.Request) {
	cah.l.Println("Handled GET Account with Credit Data")

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

	accountWithCreditData, err := cah.creditServ.GetAccountWithCreditData(uint(id))
	if err != nil {
		http.Error(rw, "Unable to get account with credit data", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(accountWithCreditData); err != nil {
		http.Error(rw, "Unable to encode JSON", http.StatusInternalServerError)
		return
	}
}
