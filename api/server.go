package api

import (
	"net/http"

	"github.com/aminerwx/repository/middleware"
	"github.com/aminerwx/repository/storage"
)

type Server struct {
	store storage.AccountRepository
	port  string
}

func NewServer(store storage.AccountRepository, port string) *Server {
	return &Server{store: store, port: port}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /accounts", s.CreateAccountHandler)
	mux.HandleFunc("GET /accounts", s.ListAccountsHandler)
	mux.HandleFunc("GET /accounts/{id}", s.GetAccountHandler)
	mux.HandleFunc("PUT /accounts/{id}", s.UpdateAccountHandler)
	mux.HandleFunc("DELETE /accounts/{id}", s.DeleteAccountHandler)
	wrappedMux := middleware.NewLogger(mux)
	return http.ListenAndServe(s.port, wrappedMux)
}
