package main

import (
	// "fmt"
	"log"

	"github.com/Habibullo-1999/wallet/pkg/wallet"

)

func main() {
	 s := &wallet.Service{}
	s.RegisterAccount("+992926421508")
	s.RegisterAccount("+992926421506")
	s.RegisterAccount("+992926421505")
	s.Deposit(1,5_000_00)
	pay,err :=s.Pay(1,5000,"cat")
	if err != nil {
		log.Print(err)
	}
	s.FavoritePayment(pay.ID,"cat favorite")
	
	err = s.Export("data")
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(err)
	err = s.Import("data")
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(err)	

	
}
