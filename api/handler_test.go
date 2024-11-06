package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/aminerwx/repository/api"
	"github.com/aminerwx/repository/model"
	"github.com/aminerwx/repository/storage"
)

type createAccountTest struct {
	payload string
	want    int
}

var (
	StatusBadRequestJSON = api.Response{
		Message:    "bad request",
		StatusCode: http.StatusBadRequest,
	}
	StatusNotFoundJSON = api.Response{
		Message:    "not found",
		StatusCode: http.StatusNotFound,
	}
	StatusOkJSON = api.Response{
		Message:    "success",
		StatusCode: http.StatusOK,
	}
	StatusCreatedJSON = api.Response{
		Message:    "success",
		StatusCode: http.StatusCreated,
	}
)

func newHTTPMock(
	httpMethod,
	endpoint string,
	payload []byte,
	myHandler http.HandlerFunc,
) *httptest.ResponseRecorder {
	req := httptest.NewRequest(httpMethod, endpoint, bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	w.Header().Set("content-type", "application/json")
	handler := http.HandlerFunc(myHandler)
	handler.ServeHTTP(w, req)
	return w
}

func TestCreateAccountHandler(t *testing.T) {
	createAccountTests := []createAccountTest{
		{``, http.StatusBadRequest},
		{`{}`, http.StatusBadRequest},
		{`{"id":1}`, http.StatusBadRequest},
		{`{"password":"pwd1"}`, http.StatusBadRequest},
		{`{"username":"account_1"}`, http.StatusBadRequest},
		{`{"id":1,"password":"pwd"}`, http.StatusBadRequest},
		{`{"id":1,"username":"account_1"}`, http.StatusBadRequest},
		{`{"id":0,"username":"","password":""}`, http.StatusBadRequest},
		{`{"username":"account_1","password":"pwd1"}`, http.StatusBadRequest},
		{`{"id":2,"username":"account_2","password":"pwd2"}`, http.StatusCreated},
		{`{"id":2,"username":"account_2","password":"pwd2"}`, http.StatusConflict},
	}
	store := storage.NewMockStorage()
	srv := api.NewServer(store, ":3333")
	for i, test := range createAccountTests {
		req := httptest.NewRequest(
			http.MethodPost,
			"/accounts",
			bytes.NewBuffer([]byte(test.payload)),
		)
		w := httptest.NewRecorder()
		w.Header().Set("content-type", "application/json")
		handler := http.HandlerFunc(srv.CreateAccountHandler)
		handler.ServeHTTP(w, req)
		got := w.Code
		if got != test.want {
			t.Errorf("[%v] got %v, want %v", i, got, test.want)
		}
	}
}

func TestGetAccountHandler(t *testing.T) {
	store := storage.NewMockStorage()
	srv := api.NewServer(store, ":3333")
	store.CreateAccount(model.Account{ID: 2, Username: "account_2", Password: "pwd2"})
	getAccountTests := []struct {
		id   string
		want api.Response
	}{
		{"1/l", StatusBadRequestJSON},
		{"/l", StatusBadRequestJSON},
		{"0", StatusBadRequestJSON},
		{"99", StatusNotFoundJSON},
		{"2", api.Response{StatusCode: 200, Message: "success", Data: []model.Account{
			{ID: 2, Username: "account_2", Password: "pwd2"},
		}}},
	}
	for i, test := range getAccountTests {
		endpoint := fmt.Sprintf("%v", test.id)
		rr := httptest.NewRequest(http.MethodGet, "/accounts", nil)
		rr.SetPathValue("id", endpoint)
		w := httptest.NewRecorder()
		w.Header().Set("content-type", "application/json")
		handler := http.HandlerFunc(srv.GetAccountHandler)
		handler.ServeHTTP(w, rr)

		var res api.Response
		json.NewDecoder(w.Body).Decode(&res)
		if res.StatusCode != test.want.StatusCode {
			t.Errorf("[%v] %v got: %v, want: %v", i, test, res.StatusCode, test.want.StatusCode)
		}
		if len(res.Data) > 0 {
			if res.Data[0] != test.want.Data[0] {
				t.Errorf("[%v] got: %v, want: %v", i, res, test.want.Data[0])
			}
		}
	}
}

func TestUpdateAccountHandler(t *testing.T) {
	store := storage.NewMockStorage()
	srv := api.NewServer(store, ":3333")
	store.CreateAccount(model.Account{ID: 1, Username: "account_1", Password: "pwd1"})
	dummyAccount := model.Account{ID: 2, Username: "account_2", Password: "pwd2"}
	success := StatusOkJSON
	success.Data = []model.Account{dummyAccount}
	getAccountTests := []struct {
		id      string
		payload model.Account
		want    api.Response
	}{
		{"1/", dummyAccount, StatusBadRequestJSON},
		{"1/$%", dummyAccount, StatusBadRequestJSON},
		{"99", dummyAccount, StatusNotFoundJSON},
		{"1", dummyAccount, success},
	}
	for i, test := range getAccountTests {
		endpoint := fmt.Sprintf("%v", test.id)
		payload, _ := json.Marshal(dummyAccount)
		rr := httptest.NewRequest(http.MethodGet, "/accounts", bytes.NewBuffer([]byte(payload)))
		rr.SetPathValue("id", endpoint)
		w := httptest.NewRecorder()
		w.Header().Set("content-type", "application/json")
		handler := http.HandlerFunc(srv.UpdateAccountHandler)
		handler.ServeHTTP(w, rr)

		var res api.Response
		json.NewDecoder(w.Body).Decode(&res)
		if res.StatusCode != test.want.StatusCode || res.Message != test.want.Message {
			t.Errorf("[%v] %v got: %v, want: %v", i, test, res.StatusCode, test.want.StatusCode)
		}

	}
}

func TestDeleteAccountHandler(t *testing.T) {
	store := storage.NewMockStorage()
	srv := api.NewServer(store, ":3333")
	store.CreateAccount(model.Account{ID: 1, Username: "account_1", Password: "pwd1"})
	getAccountTests := []struct {
		id   string
		want api.Response
	}{
		{"1/", StatusBadRequestJSON},
		{"1/$%", StatusBadRequestJSON},
		{"99", StatusNotFoundJSON},
		{"1", StatusOkJSON},
	}
	for i, test := range getAccountTests {
		endpoint := fmt.Sprintf("%v", test.id)
		rr := httptest.NewRequest(http.MethodGet, "/accounts", nil)
		rr.SetPathValue("id", endpoint)
		w := httptest.NewRecorder()
		w.Header().Set("content-type", "application/json")
		handler := http.HandlerFunc(srv.DeleteAccountHandler)
		handler.ServeHTTP(w, rr)

		var res api.Response
		json.NewDecoder(w.Body).Decode(&res)
		if res.StatusCode != test.want.StatusCode || res.Message != test.want.Message {
			t.Errorf("[%v] got: %v, want: %v", i, w.Code, test.want.StatusCode)
		}
	}
}

func TestListAccountsHandler(t *testing.T) {
	store := storage.NewMockStorage()
	srv := api.NewServer(store, ":3333")
	store.CreateAccount(model.Account{ID: 1, Username: "account_1", Password: "pwd1"})
	store.CreateAccount(model.Account{ID: 2, Username: "account_2", Password: "pwd2"})
	want := StatusOkJSON
	want.Data = []model.Account{
		{ID: 1, Username: "account_1", Password: "pwd1"},
		{ID: 2, Username: "account_2", Password: "pwd2"},
	}
	rr := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	w := httptest.NewRecorder()
	w.Header().Set("content-type", "application/json")
	handler := http.HandlerFunc(srv.ListAccountsHandler)
	handler.ServeHTTP(w, rr)
	var res api.Response
	json.NewDecoder(w.Body).Decode(&res)
	if !reflect.DeepEqual(res, want) {
		t.Errorf("got: %v, want: %v", res, want)
	}
}
