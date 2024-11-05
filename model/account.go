package model

type Account struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func NewAccount(id int, username, password string) Account {
	return Account{ID: id, Username: username, Password: password}
}
