package storage

import "github.com/aminerwx/repository/model"

type AccountRepository interface {
	GetAccount(int) (model.Account, error)
	ListAccounts() ([]model.Account, error)
	CreateAccount(model.Account) error
	UpdateAccount(int, model.Account) error
	DeleteAccount(int) error
}
