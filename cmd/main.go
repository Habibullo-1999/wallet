package main

import "github.com/Habibullo-1999/wallet/pkg/wallet"

func main() {
	s:=&wallet.Service{}
	s.RegisterAccount("+992930000001")
	s.RegisterAccount("+992930000002")
	s.RegisterAccount("+992000000003")
	s.RegisterAccount("+992000000004")
	s.RegisterAccount("+992000000005")
	s.ExportToFile("../data/message.txt")

}