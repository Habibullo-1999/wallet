package wallet

import (
	"fmt"
	"reflect"
	"testing"
	"log"

	"github.com/Habibullo-1999/wallet/pkg/types"
	"github.com/google/uuid"

)

type testService struct {
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

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func (s *testService) addAccount(date testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(date.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	err = s.Deposit(account.ID, date.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposit account, error = %v", err)
	}

	payments := make([]*types.Payment, len(date.payments))

	for i, payment := range date.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make account, error = %v", err)
		}
	}

	return account, payments, nil
}

func (s *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
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

func TestService_RegisterAccount_Fail(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9920000001")

	_, err := svc.RegisterAccount("+9920000001")
	if err != ErrPhoneRegistered {
		t.Errorf("%v", err)
	}
}

func TestService_RegisterAccount_success(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9920000001")

	account, err := svc.FindAccountByID(svc.nextAccountID)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", account)
	}
}

func TestService_Deposit_success(t *testing.T) {
	svc := Service{}
	newAcc, err := svc.RegisterAccount("+9920000001")
	if err != nil {
		t.Errorf("%v", err)
	}

	err = svc.Deposit(newAcc.ID, 1000)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
}

func TestService_Deposit_fail_amount(t *testing.T) {
	svc := Service{}
	newAcc, err := svc.RegisterAccount("+9920000001")
	if err != nil {
		t.Errorf("%v", err)
	}

	err = svc.Deposit(newAcc.ID, -1000)
	if err != ErrAmountMustBePositive {
		t.Errorf("%v", err)
		return
	}
}
func TestService_Deposit_fail_account(t *testing.T) {
	svc := Service{}
	_, err := svc.RegisterAccount("+9920000001")
	if err != nil {
		t.Errorf("%v", err)
	}

	err = svc.Deposit(10, 21000)
	if err != ErrAccountNotFound {
		t.Errorf("%v", err)
		return
	}
}

func TestService_Pay_success(t *testing.T) {
	svc := Service{}
	account, err := svc.RegisterAccount("+9920000001")
	if err != nil {
		t.Errorf("%v", err)
	}

	err = svc.Deposit(account.ID, 50_000)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	payment, err := svc.Pay(account.ID, 1000, "auto")
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(svc.payments[0], payment) {
		t.Errorf("invalid result, expected: %v, actual: %v", svc.payments[0], payment)
	}

}

func TestService_Pay_Fail_Amount(t *testing.T) {
	svc := Service{}
	account, err := svc.RegisterAccount("+9920000001")
	if err != nil {
		t.Errorf("%v", err)
	}

	err = svc.Deposit(account.ID, 50_000)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	_, err = svc.Pay(account.ID, -1000, "auto")
	if err != ErrAmountMustBePositive {
		t.Errorf("%v", err)
	}
}

func TestService_Pay_Fail_Account(t *testing.T) {
	svc := Service{}
	account, err := svc.RegisterAccount("+9920000001")
	if err != nil {
		t.Errorf("%v", err)
	}

	err = svc.Deposit(account.ID, 50_000)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	_, err = svc.Pay(12313, 1000, "auto")
	if err != ErrAccountNotFound {
		t.Errorf("%v", err)
	}
}

func TestService_Pay_Fail_Balance(t *testing.T) {
	svc := Service{}
	account, err := svc.RegisterAccount("+9920000001")
	if err != nil {
		t.Errorf("%v", err)
	}

	err = svc.Deposit(account.ID, 50_000)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	_, err = svc.Pay(account.ID, 51_000, "auto")
	if err != ErrNotEnoughBalance {
		t.Errorf("%v", err)
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
	svc.RegisterAccount("+992926421506")

	_, err := svc.FindAccountByID(3)
	if err != nil {
		fmt.Println("Аккаунт не найден")
	}

}

func TestService_FindPaymentByID_success(t *testing.T) {
	so := newTestService()
	_, payments, err := so.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("error = %v", err)
		return
	}

	payment := payments[0]

	got, err := so.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID:  error = %v", err)
		return
	}

	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID:  wrong payment returned = %v", err)
		return
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

func TestService_FindFavoriteByID_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	payment := payments[0]
	nFavorite, err := s.FavoritePayment(payment.ID, "score AlifAcademy")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	favorite, err := s.FindFavoriteByID(nFavorite.ID)
	if err != nil {
		t.Errorf("error = %v", err)
		return
	}

	if !reflect.DeepEqual(s.favorites[0], favorite) {
		t.Errorf("invalid result, expected: %v, actual: %v", s.favorites[0], favorite)
	}

}
func TestService_FindFavoriteByID_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("error = %v", err)
		return
	}

	_, err = s.FindFavoriteByID(uuid.New().String())
	if err != ErrFavoriteNotFound {
		t.Errorf("error = %v", err)
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

	payment, err := s.Pay(account.ID, 1000_00, "auto")
	if err != nil {
		t.Errorf("Reject(): can't create payment, error = %v", err)
		return
	}
	payment.AccountID = 132132131

	err = s.Reject(payment.ID)
	if err != ErrAccountNotFound {
		t.Errorf("%v", err)
		return
	}

}
func TestService_Reject_fail_account(t *testing.T) {
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
		t.Errorf("%v", err)
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

func TestService_Repeat_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("error = %v", err)
		return
	}
	_, err = s.Repeat(uuid.New().String())
	if err != ErrPaymentNotFound {
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

func TestService_FavoritePayment_fail(t *testing.T) {
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	payment := uuid.New().String()
	_, err = s.FavoritePayment(payment, "score AlifAcademy")
	if err == ErrFavoriteNotFound {
		t.Errorf("%v", err)
		return
	}

}
func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("error = %v", err)
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
		t.Error("PayFromFavorite(): must return error, returned nil")
		return
	}

	fmt.Println(payment.ID)
}

func TestService_PayFromFavorite_fail(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("error = %v", err)
		return
	}

	payment := payments[0]

	_, err = s.FavoritePayment(payment.ID, "score AlifAcademy")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	newFavorite := types.Favorite{
		ID: uuid.New().String(),
	}

	_, err = s.PayFromFavorite(newFavorite.ID)
	if err == nil {
		t.Error("PayFromFavorite(): must return error, returned nil")
		return
	}
	if err != ErrFavoriteNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}

}

func BenchmarkSumPayments(b *testing.B) {	
	s := newTestService()
	s.RegisterAccount("+992926421509")
	s.RegisterAccount("+992926421506")
	s.RegisterAccount("+992926421505")
	s.Deposit(1, 5_000_00)
	
	_, err := s.Pay(1, 5000, "cat")
	_, err = s.Pay(1, 5000, "auto")
	_, err = s.Pay(1, 5000, "auto")
	_, err = s.Pay(1, 5000, "shop")
	_, err = s.Pay(1, 5000, "food")
	_, err = s.Pay(1, 5000, "coffe")
	if err != nil {
		log.Print(err)
	}
	
	want:= types.Money(30000)	
	for i := 0; i < b.N; i++ {
		result := s.SumPayments(2)
		if result != want {
			b.Fatalf("Invalid result, dot %v, want %v", result, want)
		}
	}
}