package wallet

import (
	"fmt"
	"reflect"
	"testing"
	"log"

	"github.com/Habibullo-1999/wallet/pkg/types"
	"github.com/google/uuid"

)

type testServise struct {
	*Service
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{

	phone:   "+992926421505",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{amount: 10_000, category: "auto"},
	},
}

func newTestService() *testServise {
	return &testServise{Service: &Service{}}
}

func (s *testServise) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {

	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposit account, error = %v", err)
	}

	payments := make([]*types.Payment, len(data.payments))

	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}

	return account, payments, nil
}

func (s *testServise) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	account, err := s.RegisterAccount(phone)
	if err != nil {
		return nil, fmt.Errorf("can't register account, error = %v", err)
	}
	err = s.Deposit(account.ID, balance)
	if err != nil {
		return nil, fmt.Errorf("can't deposit account, error = %v", err)
	}
	return account, nil
}

func TestService_RegisterAccount_success(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9920000001")

	account, err := svc.FindAccountByID(svc.nextAccountID)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", account)
	}
}

func TestService_FindAccountById_success(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992926421505")

	result, err := svc.FindAccountByID(1)
	if err != nil {
		fmt.Println("Аккаунт не найдень")
	}
	myResult := types.Account{
		ID:      1,
		Phone:   "+992926421505",
		Balance: 0,
	}

	if !reflect.DeepEqual(&myResult, result) {
		t.Errorf("invalid result, expected: %v, actual: %v", &myResult, result)
	}
}
func TestService_FindAccountById_fail(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992926421505")

	_, err := svc.FindAccountByID(5)
	if err != ErrAccountNotFound {
		t.Errorf("Deposit(): must return ErrAccountNotFound, returned = %v", err)
		return
	}

}
func TestService_FindPaymentByID_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("error = %v", err)
		return
	}
	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID:  error = %v", err)
		return
	}

	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID:  wrong payment returned = %v", err)
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("error = %v", err)
		return
	}
	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Error("FindPaymentByID(): must return error, returned nil")
		return
	}
	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}

}

func TestService_Reject_success(t *testing.T) {
	s := newTestService()
	account, err := s.addAccountWithBalance("+992926421505", 10_000_00)
	if err != nil {
		t.Errorf("error = %v", err)
		return
	}

	payment, err := s.Pay(account.ID, 1000_00, "auto")
	if err != nil {
		t.Errorf("Reject(): can't create payment, error = %v", err)
		return
	}

	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject():  error = %v", err)
		return
	}

}
func TestService_Reject_fail(t *testing.T) {
	s := newTestService()
	account, err := s.addAccountWithBalance("+992926421505", 10_000_00)
	if err != nil {
		t.Errorf("error = %v", err)
		return
	}

	_, err = s.Pay(account.ID, 1000_00, "auto")
	if err != nil {
		t.Errorf("Reject(): can't create payment, error = %v", err)
		return
	}

	err = s.Reject(uuid.New().String())
	if err != ErrPaymentNotFound {
		t.Errorf("Reject(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}
	if err != ErrPaymentNotFound {
		t.Errorf("Reject(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}

}

func TestService_Repeat_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("error = %v", err)
		return
	}
	payment := payments[0]
	_, err = s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): error = %v", err)
		return
	}

}

func TestService_FavoritePayment_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	payment := payments[0]
	_, err = s.FavoritePayment(payment.ID, "score AlifAcademy")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

}

func TestService_FindFavoriteByID_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	payment := payments[0]
	favorite, err := s.FavoritePayment(payment.ID, "score AlifAcademy")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	_, err = s.FindFavoriteByID(favorite.ID)
	
	if err != nil {
		t.Errorf("%v",err)
		return
	}


}
func TestService_FindFavoriteByID_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("%v", err)
		return
	}


	_, err = s.FindFavoriteByID(uuid.New().String())
	
	if err != ErrFavoriteNotFound {
		t.Errorf("FindPaymentByID(): must return ErrFavoriteNotFound, returned = %v", err)
		return
	}

}
func TestService_Deposit_AmountMustBePositive(t *testing.T) {
	s := newTestService()
	account, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	err = s.Deposit(account.ID,-1_000_00)
	if err != ErrAmountMustBePositive {
		t.Errorf("Deposit(): must return ErrAmountMustBePositive, returned = %v", err)
		return
	}	
}

func TestService_Deposit_ErrAccountNotFound(t *testing.T) {
	s := newTestService()
	account, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	err = s.Deposit(account.ID+1,1_000_00)
	if err != ErrAccountNotFound {
		t.Errorf("Deposit(): must return ErrAccountNotFound, returned = %v", err)
		return
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	payment := payments[0]
	favorite, err := s.FavoritePayment(payment.ID, "score AlifAcademy")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	_, err = s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
}

func TestService_Export_success_user(t *testing.T) {
	var svc Service

	svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")
	svc.RegisterAccount("+992000000003")

	err := svc.ExportToFile("export.txt")
	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}

}

func TestService_Import_success_user(t *testing.T) {
	var svc Service

	err := svc.ImportFromFile("export.txt")

	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}

}

func TestService_Export_success(t *testing.T) {
	svc := Service{}

	svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")
	svc.RegisterAccount("+992000000003")
	svc.RegisterAccount("+992000000004")

	err := svc.Export("data")
	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}

	err = svc.Import("data")
	if err != nil {
		t.Errorf("method ExportToFile returned not nil error, err => %v", err)
	}
}

func TestService_ExportHistory_success_user(t *testing.T) {
	svc := Service{}

	acc, err := svc.RegisterAccount("+992000000001")

	if err != nil {
		t.Errorf("method RegisterAccount returned not nil error, account => %v", acc)
	}

	err = svc.Deposit(acc.ID, 100_00)
	if err != nil {
		t.Errorf("method Deposit returned not nil error, error => %v", err)
	}

	_, err = svc.Pay(acc.ID, 1, "Cafe")
	_, err = svc.Pay(acc.ID, 2, "Auto")
	_, err = svc.Pay(acc.ID, 3, "MarketShop")
	if err != nil {
		t.Errorf("method Pay returned not nil error, err => %v", err)
	}

	payments, err := svc.ExportAccountHistory(acc.ID)
	if err != nil {
		t.Errorf("method ExportAccountHistory returned not nil error, err => %v", err)
	}

	err = svc.HistoryToFiles(payments, "../../data", 2)
	if err != nil {
		t.Errorf("method HistoryToFiles returned not nil error, err => %v", err)
	}
}

func BenchmarkSumPayment_user(b *testing.B) {
	var svc Service

	account, err := svc.RegisterAccount("+992000000001")

	if err != nil {
		b.Errorf("method RegisterAccount returned not nil error, account => %v", account)
	}

	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
		b.Errorf("method Deposit returned not nil error, error => %v", err)
	}

	_, err = svc.Pay(account.ID, 1, "Cafe")
	_, err = svc.Pay(account.ID, 2, "Cafe")
	_, err = svc.Pay(account.ID, 3, "Cafe")
	_, err = svc.Pay(account.ID, 4, "Cafe")
	_, err = svc.Pay(account.ID, 5, "Cafe")
	_, err = svc.Pay(account.ID, 6, "Cafe")
	_, err = svc.Pay(account.ID, 7, "Cafe")
	_, err = svc.Pay(account.ID, 8, "Cafe")
	_, err = svc.Pay(account.ID, 9, "Cafe")
	_, err = svc.Pay(account.ID, 10, "Cafe")
	_, err = svc.Pay(account.ID, 11, "Cafe")
	if err != nil {
		b.Errorf("method Pay returned not nil error, err => %v", err)
	}

	want := types.Money(66)

	got := svc.SumPayments(2)
	if want != got {
		b.Errorf(" error, want => %v got => %v", want, got)
	}

}
func BenchmarkSumPaymentsWithProgress_user(b *testing.B) {
	var svc Service

	account, err := svc.RegisterAccount("+992000000001")

	if err != nil {
		b.Errorf("method RegisterAccount returned not nil error, account => %v", account)
	}

	err = svc.Deposit(account.ID, 10000000_0000000)
	if err != nil {
		b.Errorf("method Deposit returned not nil error, error => %v", err)
	}

	for i := 0; i < 1000; i++ {
		svc.Pay(account.ID, types.Money(i), "Cafe")
	} 

	ch := svc.SumPaymentsWithProgress()

	 s, ok := <-ch

	if !ok {
		b.Errorf(" method SumPaymentsWithProgress ok not closed => %v", ok)
	} 

	log.Println("=======>>>>>",s) 

}