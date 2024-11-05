package main

import (
	"fmt"

	"github.com/aminerwx/repository/api"
	"github.com/aminerwx/repository/storage"
)

func main() {
	fmt.Println("Hello, repository")
	store := storage.NewMockStorage()
	srv := api.NewServer(store, ":3000")
	if err := srv.Start(); err != nil {
		panic(err)
	}
}
