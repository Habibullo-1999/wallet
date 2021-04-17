package main

import (
	"github.com/Habibullo-1999/wallet/pkg/wallet"

)

func main() {
	svc := &wallet.Service{}

	svc.RegisterAccount("+992926421505")
	svc.RegisterAccount("+992926421502")
	svc.RegisterAccount("+992926421503")
	svc.RegisterAccount("+992926421504")
	svc.RegisterAccount("+992926421506")
	svc.RegisterAccount("+992926421507")
	svc.ExportToFile("../data/export.txt")
}