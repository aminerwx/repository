package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/aminerwx/repository/model"
)

type Response struct {
	Message    string          `json:"message"`
	StatusCode int             `json:"status_code"`
	Data       []model.Account `json:"data"`
}

var (
	StatusBadRequestJSON = Response{
		Message:    "bad request",
		StatusCode: http.StatusBadRequest,
	}
	StatusConflictJSON = Response{
		Message:    "already exist",
		StatusCode: http.StatusConflict,
	}
	StatusNotFoundJSON = Response{
		Message:    "not found",
		StatusCode: http.StatusNotFound,
	}
	StatusOkJSON = Response{
		Message:    "success",
		StatusCode: http.StatusOK,
	}
	StatusCreatedJSON = Response{
		Message:    "success",
		StatusCode: http.StatusCreated,
	}
)

func hasEmptyField(account model.Account) bool {
	return account.ID == 0 || account.Username == "" || account.Password == ""
}

func (s *Server) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var account model.Account
	body, _ := io.ReadAll(r.Body)
	w.Header().Set("content-type", "application/json")

	err := json.Unmarshal(body, &account)
	if err != nil || account == (model.Account{}) || hasEmptyField(account) {
		w.WriteHeader(StatusBadRequestJSON.StatusCode)
		json.NewEncoder(w).Encode(StatusBadRequestJSON)
		return
	}

	err = s.store.CreateAccount(account)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(StatusConflictJSON)
		return
	}
	w.WriteHeader(StatusCreatedJSON.StatusCode)
	json.NewEncoder(w).Encode(StatusCreatedJSON)
}

func (s *Server) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if id == 0 {
		w.WriteHeader(StatusBadRequestJSON.StatusCode)
		json.NewEncoder(w).Encode(StatusBadRequestJSON)
		return
	}

	if err != nil {
		w.WriteHeader(StatusBadRequestJSON.StatusCode)
		json.NewEncoder(w).Encode(StatusBadRequestJSON)
		return
	}
	account, err := s.store.GetAccount(id)
	if err != nil {
		w.WriteHeader(StatusNotFoundJSON.StatusCode)
		json.NewEncoder(w).Encode(StatusNotFoundJSON)
		return
	}
	w.WriteHeader(StatusOkJSON.StatusCode)
	StatusOkJSON.Data = []model.Account{account}
	json.NewEncoder(w).Encode(StatusOkJSON)
}

func (s *Server) ListAccountsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	accounts, err := s.store.ListAccounts()
	if err != nil {
		w.WriteHeader(StatusNotFoundJSON.StatusCode)
		json.NewEncoder(w).Encode(StatusNotFoundJSON)
		return
	}
	w.WriteHeader(StatusOkJSON.StatusCode)
	StatusOkJSON.Data = accounts
	json.NewEncoder(w).Encode(StatusOkJSON)
}

func (s *Server) UpdateAccountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(StatusBadRequestJSON.StatusCode)
		json.NewEncoder(w).Encode(StatusBadRequestJSON)
		return
	}
	var newAccount model.Account
	json.NewDecoder(r.Body).Decode(&newAccount)
	err = s.store.UpdateAccount(id, newAccount)
	if err != nil {
		w.WriteHeader(StatusNotFoundJSON.StatusCode)
		json.NewEncoder(w).Encode(StatusNotFoundJSON)
		return
	}
	w.WriteHeader(StatusOkJSON.StatusCode)
	StatusOkJSON.Data = []model.Account{}
	json.NewEncoder(w).Encode(StatusOkJSON)
}

func (s *Server) DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(StatusBadRequestJSON.StatusCode)
		json.NewEncoder(w).Encode(StatusBadRequestJSON)
		return
	}
	err = s.store.DeleteAccount(id)
	if err != nil {
		w.WriteHeader(StatusNotFoundJSON.StatusCode)
		json.NewEncoder(w).Encode(StatusNotFoundJSON)
		return
	}
	w.WriteHeader(StatusOkJSON.StatusCode)
	StatusOkJSON.Data = []model.Account{}
	json.NewEncoder(w).Encode(StatusOkJSON)
}
