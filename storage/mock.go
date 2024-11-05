package storage

import (
	"errors"

	"github.com/aminerwx/repository/model"
)

type MockStorage struct {
	Accounts []model.Account
}

func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

func (m *MockStorage) GetAccount(id int) (model.Account, error) {
	for _, acc := range m.Accounts {
		if acc.ID == id {
			return acc, nil
		}
	}
	return model.Account{}, errors.New("account not found by given ID")
}

func (m *MockStorage) ListAccounts() ([]model.Account, error) {
	return m.Accounts, nil
}

func (m *MockStorage) CreateAccount(account model.Account) error {
	for _, acc := range m.Accounts {
		if acc.ID == account.ID {
			return errors.New("account id already exist")
		}
	}
	m.Accounts = append(m.Accounts, account)
	return nil
}

func (m *MockStorage) UpdateAccount(id int, account model.Account) error {
	for idx, acc := range m.Accounts {
		if acc.ID == id {
			m.Accounts[idx] = account
			return nil
		}
	}
	return errors.New("account not found")
}

func (m *MockStorage) DeleteAccount(id int) error {
	for idx, acc := range m.Accounts {
		if acc.ID == id {
			m.Accounts = append(m.Accounts[:idx], m.Accounts[idx+1:]...)
			return nil
		}
	}
	return errors.New("account not found")
}
