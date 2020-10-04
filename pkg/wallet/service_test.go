package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Habibullo-1999/wallet/pkg/types"

)

// func TestService_Deposit_success(t *testing.T) {
// 	svc := &Service{}
// 	account, err := svc.RegisterAccount("+992926421505")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	result := svc.Deposit(account.ID, 10)
// 	if result != nil {
// 		switch result {
// 		case ErrAmountMustBePositive:
// 			fmt.Println("Сумма должна быт положительной")
// 		case ErrAccountNotFound:
// 			fmt.Println("Аккаунт пользователя не найден")
// 		}
// 		return
// 	}

// 	if !reflect.DeepEqual(nil, result) {
// 		t.Errorf("invalid result, expected: %v, actual: %v", nil, result)
// 	}
// }

// func TestService_Pay_succes(t *testing.T) {
// 	svc := &Service{}
// 	svc.RegisterAccount("+992926421505")

// 	result, err := svc.Pay(svc.nextAccountID, types.Money(100), "cat")
// 	if err != nil {
// 		switch err {
// 		case ErrAmountMustBePositive:
// 			fmt.Println("Сумма должна быт положительной")
// 		case ErrAccountNotFound:
// 			fmt.Println("Аккаунт пользователя не найден")
// 		case ErrNotEnoughBalance:
// 			fmt.Println("Ошибка! недостаточно баланса")
// 		}
// 		return 
// 	}
	

// 	if !reflect.DeepEqual(svc.payments , result) {
// 		t.Errorf("invalid result, expected: %v, actual: %v", svc.payments, result)
// 	}
	
// }
func TestService_RegisterAccount_success(t *testing.T) { 
	svc := Service{}
	svc.RegisterAccount("+9920000001")
  
	account, err := svc.FindAccountByID(svc.nextAccountID)
	if err != nil {
	  t.Errorf("\ngot > %v \nwant > nil", account)
	}
  }

func TestService_FindAccountById(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992926421505")

	result, err := svc.FindAccountByID(1)
	if err != nil{
		fmt.Println("Аккаунт не найдень")
	}
	myResult:= types.Account {
		ID:		1,
		Phone:	"+992926421505",
		Balance:101,
	}

	if !reflect.DeepEqual(&myResult , result) {
		t.Errorf("invalid result, expected: %v, actual: %v", &myResult, result)
	}
}
func TestService_FindPaymentByID_success(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992926421505")

	payment, err := svc.Pay(2, types.Money(110), "cat")

	if err != nil {
		switch err {
			case ErrAmountMustBePositive:
				fmt.Println("Сумма должна быт положительной")
			case ErrAccountNotFound:
				fmt.Println("Аккаунт пользователя не найден")
			case ErrNotEnoughBalance:
				fmt.Println("Ошибка! недостаточно баланса")
			}
			return 
	}

	result, err := svc.FindPaymentByID(payment.ID)
	if err != nil{
		fmt.Println("Платёж не найдень")
	}

	if !reflect.DeepEqual(payment, result) {
		t.Errorf("invalid result, expected: %v, actual: %v", payment, result)
	}
}

func TestService_Reject_success(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992926421505")

	payment, err := svc.Pay(1, types.Money(110), "cat")
	if err != nil {
		switch err {
			case ErrAmountMustBePositive:
				fmt.Println("Сумма должна быт положительной")
			case ErrAccountNotFound:
				fmt.Println("Аккаунт пользователя не найден")
			case ErrNotEnoughBalance:
				fmt.Println("Ошибка! недостаточно баланса")
			}
			return 
	}
	
	result:=svc.Reject(payment.ID)
	if !reflect.DeepEqual(nil, result) {
		t.Errorf("invalid result, expected: %v, actual: %v", nil, result)
	}



} 


