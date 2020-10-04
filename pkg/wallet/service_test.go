package wallet

import (
	"fmt"
	"reflect"
	"testing"

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
		{amount: 10_000_00, category: "auto"},
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

func TestService_FindAccountById(t *testing.T) {
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
	if err== nil {
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
