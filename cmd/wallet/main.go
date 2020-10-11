package main

import "github.com/Habibullo-1999/wallet/pkg/wallet"

func main() {
	s:=*&wallet.Service{}
	s.RegisterAccount("926421505")
	s.RegisterAccount("926421506")
	s.RegisterAccount("926421507")
	s.RegisterAccount("926421508")
	s.RegisterAccount("926421509")
	s.ExportToFile("data/message.txt")

}