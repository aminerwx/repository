package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/aminerwx/repository/model"
)

func hasEmptyField(account model.Account) bool {
	return account.ID == 0 || account.Username == "" || account.Password == ""
}

func (s *Server) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var account model.Account
	// json.NewDecoder(r.Body).Decode(&account)
	body, _ := io.ReadAll(r.Body)
	w.Header().Set("content-type", "application/json")

	err := json.Unmarshal(body, &account)
	if err != nil || account == (model.Account{}) || hasEmptyField(account) {
		response := map[string]string{"message": "bad request.", "status": "error"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err = s.store.CreateAccount(account)
	if err != nil {
		response := map[string]string{"message": "bad request.", "status": "error"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	response := map[string]string{"message": "account created.", "status": "ok"}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (s *Server) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if id == 0 {
		response := `{"message": "invalid request", "status": 400}`
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err != nil {
		response := `{"message": "invalid request", "status": 400}`
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	account, err := s.store.GetAccount(id)
	if err != nil {
		response := `{"message": "not found", "status": 404}`
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

func (s *Server) ListAccountsHandler(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.store.ListAccounts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accounts)
}

func (s *Server) UpdateAccountHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var newAccount model.Account
	json.NewDecoder(r.Body).Decode(&newAccount)
	err = s.store.UpdateAccount(id, newAccount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.store.DeleteAccount(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}
