package main

import (
	"fmt"

	"github.com/Habibullo-1999/wallet/pkg/wallet"

)

func main() {
	s := &wallet.Service{}
	err := s.ImportFromFile("../data/export.txt")
	if err == nil {
		fmt.Println("success")
	}
}
