package storage_test

import (
	"errors"
	"testing"

	"github.com/aminerwx/repository/model"
	"github.com/aminerwx/repository/storage"
)

type TestCase struct {
	value  model.Account
	expect model.Account
}

var MockTestData = []model.Account{
	{ID: 1, Username: "user1", Password: "pwd1"},
	{ID: 2, Username: "user2", Password: "pwd2"},
	{ID: 3, Username: "user3", Password: "pwd3"},
	{ID: 4, Username: "user4", Password: "pwd4"},
	{ID: 5, Username: "user5", Password: "pwd5"},
}

func TestMockGetAccount(t *testing.T) {
	store := storage.NewMockStorage()
	for _, data := range MockTestData {
		store.CreateAccount(data)
	}
	t.Run("should found account", func(t *testing.T) {
		user := MockTestData[2]
		got, _ := store.GetAccount(3)
		if user != got {
			t.Errorf("got %v, wants %v\n", got, user)
		}
	})
	t.Run("should not found account", func(t *testing.T) {
		got, _ := store.GetAccount(8)
		empty := model.Account{}
		if empty != got {
			t.Errorf("got %v, wants %v\n", got, empty)
		}
	})
}

func TestMockListAccounts(t *testing.T) {
	store := storage.NewMockStorage()
	t.Run("should return empty list of accounts", func(t *testing.T) {
		accounts, err := store.ListAccounts()
		if len(accounts) != 0 || err != nil {
			t.Errorf("got: %v, want: %v", len(accounts), 0)
		}
	})
}

func TestMockCreateAccount(t *testing.T) {
	store := storage.NewMockStorage()
	t.Run("should create account", func(t *testing.T) {
		user := MockTestData[2]
		got := store.CreateAccount(user)
		if got != nil {
			t.Errorf("got %v, wants %v\n", got, nil)
		}
	})
	t.Run("should not create account", func(t *testing.T) {
		user := MockTestData[0]
		got := store.CreateAccount(user)
		got = store.CreateAccount(user)
		ErrAccountExists := errors.New("account id already exist")
		if got.Error() != ErrAccountExists.Error() {
			t.Errorf("got %v, wants %v\n", got, ErrAccountExists)
		}
	})
}

func TestMockUpdateAccount(t *testing.T) {
	store := storage.NewMockStorage()
	t.Run("should update account", func(t *testing.T) {
		user := MockTestData[2]
		store.CreateAccount(user)
		got := store.UpdateAccount(3, MockTestData[0])
		if got != nil {
			t.Errorf("got %v, wants %v\n", got, nil)
		}
	})
	t.Run("should not update account", func(t *testing.T) {
		user := MockTestData[2]
		store.CreateAccount(user)
		got := store.UpdateAccount(10, MockTestData[0])
		errAccountNotFound := errors.New("account not found")
		if got.Error() != errAccountNotFound.Error() {
			t.Errorf("got %v, wants %v\n", got, errAccountNotFound)
		}
	})
}

func TestMockDeleteAccount(t *testing.T) {
	store := storage.NewMockStorage()
	t.Run("should delete account", func(t *testing.T) {
		user := MockTestData[2]
		store.CreateAccount(user)
		got := store.DeleteAccount(3)
		if got != nil {
			t.Errorf("got %v, wants %v\n", got, nil)
		}
	})
	t.Run("should not delete account", func(t *testing.T) {
		user := MockTestData[2]
		store.CreateAccount(user)
		got := store.DeleteAccount(9)
		errAccountNotFound := errors.New("account not found")
		if got.Error() != errAccountNotFound.Error() {
			t.Errorf("got %v, wants %v\n", got, errAccountNotFound)
		}
	})
}
