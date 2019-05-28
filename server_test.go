package main

import (
	"fmt"
	"testing"
)

func TestGenerateContainer(t *testing.T) {
	repo := "https://github.com/react-cosmos/create-react-app-example.git"
	port := "3000"
	addr, err := generateContainer(repo, port)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(addr)
}
