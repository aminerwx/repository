package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aminerwx/repository/api"
	"github.com/aminerwx/repository/model"
	"github.com/aminerwx/repository/storage"
)

type createAccountTest struct {
	payload string
	want    int
}

var createAccountTests = []createAccountTest{
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
}

type getAccountTest struct {
	payload int
	want    string
}

func newHTTPMock(httpMethod, endpoint string, payload []byte, myHandler http.HandlerFunc) int {
	req := httptest.NewRequest(httpMethod, endpoint, bytes.NewBuffer(payload))
	w := httptest.NewRecorder()
	w.Header().Set("content-type", "application/json")
	handler := http.HandlerFunc(myHandler)
	handler.ServeHTTP(w, req)
	return w.Code
}

func runAccountHandlerTests(
	t *testing.T,
	httpMethod string,
	endpoint string,
	f func(w http.ResponseWriter, r *http.Request),
	funcName string,
	testCases []createAccountTest,
) {
	for i, test := range testCases {
		got := newHTTPMock(httpMethod, endpoint, []byte(test.payload), f)
		if got != test.want {
			t.Errorf("[%v] %s: got %v, want %v", i, funcName, got, test.want)
		}
	}
}

func TestCreateAccountHandler(t *testing.T) {
	store := storage.NewMockStorage()
	srv := api.NewServer(store, ":3333")
	runAccountHandlerTests(t, http.MethodPost, "/accounts", srv.CreateAccountHandler, "CreateAccountHandler", createAccountTests)
}

func TestGetAccountHandler(t *testing.T) {
	store := storage.NewMockStorage()
	srv := api.NewServer(store, ":3333")
	store.CreateAccount(model.Account{ID: 2, Username: "account_2", Password: "pwd2"})
	getAccountTests := []struct {
		id      int
		status  int
		payload string
	}{
		{0, 400, `"{\"message\": \"invalid request\", \"status\": 400}"`},
		{99, 404, `"{\"message\": \"not found\", \"status\": 404}"`},
		{2, 200, `{"id":2,"username":"account_2","password":"pwd2"}`},
	}
	for i, test := range getAccountTests {
		endpoint := fmt.Sprintf("%v", test.id)
		rr := httptest.NewRequest(http.MethodGet, "/accounts", bytes.NewBuffer([]byte(test.payload)))
		rr.SetPathValue("id", endpoint)
		w := httptest.NewRecorder()
		w.Header().Set("content-type", "application/json")
		handler := http.HandlerFunc(srv.GetAccountHandler)
		handler.ServeHTTP(w, rr)

		var acc model.Account
		body, _ := io.ReadAll(w.Body)
		json.Unmarshal(body, &acc)
		jsonResponse := string(body[:len(body)-1])
		if w.Code != test.status {
			t.Errorf("[%v] %v got: %v, want: %v", i, test, w.Code, test.status)
		}

		if test.payload != jsonResponse {
			t.Errorf("[%v] got: %v %v, want: %v %v", i, jsonResponse, len(jsonResponse), test.payload, len(test.payload))
		}
	}
}
