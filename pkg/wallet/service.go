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
var ErrFavoriteNotFound = errors.New("Favorite not found")


type Service struct {
	nextAccountID int64
	favorites	  []*types.Favorite
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
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	account.Balance += amount
	return nil
}

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
	for _, acc := range  s.accounts{
		if acc.ID == accountID {
			return acc,nil
		}
	}
	return nil,ErrAccountNotFound
	
}

func (s *Service)FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, acc := range  s.payments{
		if acc.ID == paymentID {
			return acc,nil
		}
	}
		return nil,ErrPaymentNotFound
	
}

func (s *Service)Reject(paymentID string) error {
	
	payment, err := s.FindPaymentByID(paymentID)

	if err != nil{
		return err
	}

	account,err := s.FindAccountByID(payment.AccountID)

	if err!=nil {
		return err
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return nil
	
}

 
func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)

	if err != nil{
		return nil,err
	}
	account,err := s.FindAccountByID(payment.AccountID)
	if err != nil{
		return nil,err
	}
	account.Balance -= payment.Amount
	paymentIDE := uuid.New().String()
	payments := &types.Payment{
		ID:			paymentIDE,
		AccountID:	payment.AccountID,
		Amount:		payment.Amount,
		Category:	payment.Category,
		Status:		payment.Status,
	}
	s.payments = append(s.payments, payments)
	return payments,nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error){
	payment, err := s.FindPaymentByID(paymentID)

	if err != nil{
		return nil,err
	}
	paymentIDE := uuid.New().String()
	favorite := &types.Favorite{
		ID:			paymentIDE,
		AccountID:	payment.AccountID,
		Name:		name,
		Amount:		payment.Amount,
		Category:	payment.Category,
	}	
	s.favorites = append(s.favorites, favorite)
	return favorite,nil
}
func (s *Service)FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, fav := range  s.favorites{
		if fav.ID == favoriteID {
			return fav,nil
		}
	}
		return nil,ErrPaymentNotFound
	
}

func (s *Service)PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil{
		return nil,ErrFavoriteNotFound
	}
	payment,err := s.Pay(favorite.AccountID,favorite.Amount,favorite.Category)
	if err != nil{
		return nil,err
	}
	return payment,nil
}