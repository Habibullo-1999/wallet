package wallet

import (
	"errors"

	"github.com/Habibullo-1999/wallet/pkg/types"
	"github.com/google/uuid"

)

var ErrPhoneRegistered = errors.New("Phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("Error Not Enough Balance")
var ErrPaymentNotFound = errors.New("Payment not found")


type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 101,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}

// func (s *Service) Deposit(accountID int64, amount types.Money) error {
// 	if amount <= 0 {
// 		return ErrAmountMustBePositive
// 	}

// 	var account *types.Account
// 	for _, acc := range s.accounts {
// 		if acc.ID == accountID {
// 			account = acc
// 			break
// 		}
// 	}

// 	if account == nil {
// 		return ErrAccountNotFound
// 	}

// 	account.Balance += amount
// 	return nil
// }

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}
	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:			paymentID,
		AccountID:	accountID,
		Amount:		amount,
		Category:	category,
		Status:		types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment,nil

}


func (s *Service)FindAccountByID(accountID int64) (*types.Account, error) {

	var account *types.Account
	for _, acc := range  s.accounts{
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil,ErrAccountNotFound
	}
	
	return account,nil
}

func (s *Service)FindPaymentByID(paymentID string) (*types.Payment, error) {
	var payment *types.Payment
	for _, acc := range  s.payments{
		if acc.ID == paymentID {
			payment = acc
			break
		}
	}
	if payment == nil {
		return nil,ErrPaymentNotFound
	}
	
	return payment,nil
}

func (s *Service)Reject(paymentID string) error {
	
	pay, err :=s.FindPaymentByID(paymentID)

	if err != nil{
		return ErrPaymentNotFound
	}

	acc,err :=s.FindAccountByID(pay.AccountID)

	if err!=nil {
		return err
	}

	pay.Status=types.PaymentStatusFail
	acc.Balance+=pay.Amount
	return nil
	
}
