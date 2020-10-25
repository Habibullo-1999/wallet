package wallet

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Habibullo-1999/wallet/pkg/types"
	"github.com/google/uuid"

)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")
var ErrFileNotFound = errors.New("File not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
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

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	return nil, ErrAccountNotFound
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return ErrAccountNotFound
	}

	// зачисление средств пока не рассматриваем как платёж
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
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}

	return nil, ErrPaymentNotFound
}

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	return s.Pay(payment.AccountID, payment.Amount, payment.Category)
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favorite := &types.Favorite{
		ID:        uuid.New().String(),
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Name:      name,
		Category:  payment.Category,
	}

	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}

	return nil, ErrFavoriteNotFound
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}
	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			return
		}
	}()
	// if err != nil {
	// 	log.Print(err)
	// 	return ErrFileNotFound
	// }
	str := ""
	for _, data := range s.accounts {
		str += strconv.Itoa(int(data.ID)) + ";"
		str += string(data.Phone) + ";"
		str += strconv.Itoa(int(data.Balance)) + "|"
	}
	_, err = file.Write([]byte(str))
	if err != nil {
		return err
	}

	return nil
}
func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			return
		}
	}()

	content := make([]byte, 0)
	buf := make([]byte, 4)

	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		content = append(content, buf[:read]...)
	}
	accounts := strings.Split(string(content), "|")
	accounts = accounts[:len(accounts)-1]

	for _, acc := range accounts {
		splits := strings.Split(acc, ";")

		id, err := strconv.Atoi(splits[0])
		if err != nil {
			return err
		}

		balance, err := strconv.Atoi(splits[2])
		if err != nil {
			return err
		}

		account := &types.Account{
			ID:      int64(id),
			Phone:   types.Phone(splits[1]),
			Balance: types.Money(balance),
		}

		s.accounts = append(s.accounts, account)
	}
	return nil
}

// func (s *Service) Export(dir string) (err error) {

// 	if len(s.accounts) !=0 {

// 	file, err := os.OpenFile("../"+dir+"/accounts.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if cerr := file.Close(); cerr != nil {
// 			if err != nil {
// 				err = cerr
// 				log.Print(err)
// 			}
// 		}
// 	}()
// 	str := ""

// 	for _, data := range s.accounts {
// 		str += strconv.Itoa(int(data.ID)) + ";"
// 		str += string(data.Phone) + ";"
// 		str += strconv.Itoa(int(data.Balance)) + "\n"
// 	}
// 	err = ioutil.WriteFile(file.Name(), []byte(str), 0666)
// 	if err != nil {
// 		return err
// 	}
// 	}

// 	if len(s.payments) !=0 {
// 	file, err := os.OpenFile("../"+dir+"/payments.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if cerr := file.Close(); cerr != nil {
// 			if err != nil {
// 				err = cerr
// 				log.Print(err)
// 			}
// 		}
// 	}()
// 	str := ""

// 	for _, data := range s.payments {
// 		str += string(data.ID) + ";"
// 		str += strconv.Itoa(int(data.AccountID)) + ";"
// 		str += strconv.Itoa(int(data.Amount)) +";"
// 		str += string(data.Category)+";"
// 		str += string(data.Status)+"\n"
// 	}
// 	err = ioutil.WriteFile(file.Name(), []byte(str), 0666)
// 	if err != nil {
// 		return err
// 	}
// 	}

// 	if len(s.favorites) !=0 {

// 	file, err := os.OpenFile("../"+dir+"/favorites.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if cerr := file.Close(); cerr != nil {
// 			if err != nil {
// 				err = cerr
// 				log.Print(err)
// 			}
// 		}
// 	}()
// 	str := ""

// 	for _, data := range s.favorites {
// 		str += string(data.ID) + ";"
// 		str += strconv.Itoa(int(data.AccountID)) + ";"
// 		str += strconv.Itoa(int(data.Amount)) +";"
// 		str += string(data.Name)+";"
// 		str += string(data.Category)+"\n"

// 	}
// 	err = ioutil.WriteFile(file.Name(), []byte(str), 0666)
// 	if err != nil {
// 		return err
// 	}
// }
// 	return nil
// }

// func (s *Service) Import(dir string) (err error) {
// 	_, err = os.Stat("../../" + dir + "/accounts.dump")

// 	if err == nil {
// 	  content, err := ioutil.ReadFile("../../" + dir + "/accounts.dump")
// 	  if err != nil {
// 		return err
// 	  }

// 	  strArray := strings.Split(string(content), "\n")
// 	  if len(strArray) > 0 {
// 		strArray = strArray[:len(strArray)-1]
// 	  }
// 	  for _, v := range strArray {
// 		strArrAcount := strings.Split(v, ";")
// 		fmt.Println(strArrAcount)

// 		id, err := strconv.ParseInt(strArrAcount[0], 10, 64)
// 		if err != nil {
// 		  return err
// 		}
// 		balance, err := strconv.ParseInt(strArrAcount[2], 10, 64)
// 		if err != nil {
// 		  return err
// 		}
// 		flag := true
// 		for _, v := range s.accounts {
// 		  if v.ID == id {
// 			v.Phone = types.Phone(strArrAcount[1])
// 			v.Balance = types.Money(balance)
// 			flag = false
// 		  }
// 		}
// 		if flag {
// 		  account := &types.Account{
// 			ID:      id,
// 			Phone:   types.Phone(strArrAcount[1]),
// 			Balance: types.Money(balance),
// 		  }
// 		  s.accounts = append(s.accounts, account)
// 		}
// 	  }
// 	}

// 	_, err1 := os.Stat("../../" + dir + "/payments.dump")

// 	if err1 == nil {
// 	  content, err := ioutil.ReadFile("../../" + dir + "/payments.dump")
// 	  if err != nil {
// 		return err
// 	  }

// 	  strArray := strings.Split(string(content), "\n")
// 	  if len(strArray) > 0 {
// 		strArray = strArray[:len(strArray)-1]
// 	  }
// 	  for _, v := range strArray {
// 		strArrAcount := strings.Split(v, ";")
// 		fmt.Println(strArrAcount)

// 		id := strArrAcount[0]
// 		if err != nil {
// 		  return err
// 		}
// 		aid, err := strconv.ParseInt(strArrAcount[1], 10, 64)
// 		if err != nil {
// 		  return err
// 		}
// 		amount, err := strconv.ParseInt(strArrAcount[2], 10, 64)
// 		if err != nil {
// 		  return err
// 		}
// 		flag := true
// 		for _, v := range s.payments {
// 		  if v.ID == id {
// 			v.AccountID = aid
// 			v.Amount = types.Money(amount)
// 			v.Category = types.PaymentCategory(strArrAcount[3])
// 			v.Status = types.PaymentStatus(strArrAcount[4])
// 			flag = false
// 		  }
// 		}
// 		if flag {
// 		  data := &types.Payment{
// 			ID:        id,
// 			AccountID: aid,
// 			Amount:    types.Money(amount),
// 			Category:  types.PaymentCategory(strArrAcount[3]),
// 			Status:    types.PaymentStatus(strArrAcount[4]),
// 		  }
// 		  s.payments = append(s.payments, data)
// 		}
// 	  }
// 	}

// 	_, err2 := os.Stat("../../" + dir + "/favorites.dump")

// 	if err2 == nil {
// 	  content, err := ioutil.ReadFile("../../" + dir + "/favorites.dump")
// 	  if err != nil {
// 		return err
// 	  }

// 	  strArray := strings.Split(string(content), "\n")
// 	  if len(strArray) > 0 {
// 		strArray = strArray[:len(strArray)-1]
// 	  }
// 	  for _, v := range strArray {
// 		strArrAcount := strings.Split(v, ";")
// 		fmt.Println(strArrAcount)

// 		id := strArrAcount[0]
// 		if err != nil {
// 		  return err
// 		}
// 		aid, err := strconv.ParseInt(strArrAcount[1], 10, 64)
// 		if err != nil {
// 		  return err
// 		}
// 		amount, err := strconv.ParseInt(strArrAcount[2], 10, 64)
// 		if err != nil {
// 		  return err
// 		}
// 		flag := true
// 		for _, v := range s.favorites {
// 		  if v.ID == id {
// 			v.AccountID = aid
// 			v.Amount = types.Money(amount)
// 			v.Category = types.PaymentCategory(strArrAcount[3])
// 			flag = false
// 		  }
// 		}
// 		if flag {
// 		  data := &types.Favorite{
// 			ID:        id,
// 			AccountID: aid,
// 			Amount:    types.Money(amount),
// 			Category:  types.PaymentCategory(strArrAcount[3]),
// 		  }
// 		  s.favorites = append(s.favorites, data)
// 		}
// 	  }
// 	}

// 	return nil
// }

//Export method
func (s *Service) Export(dir string) error {
	if len(s.accounts) > 0 {
		file, err := os.OpenFile(dir+"/accounts.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

		defer func() {
			if cerr := file.Close(); cerr != nil {
				if err != nil {
					err = cerr
					log.Print(err)
				}
			}
		}()

		str := ""
		for _, v := range s.accounts {
			str += fmt.Sprint(v.ID) + ";" + string(v.Phone) + ";" + fmt.Sprint(v.Balance) + "\n"
		}
		file.WriteString(str)
	}
	if len(s.payments) > 0 {
		file, err := os.OpenFile(dir+"/payments.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

		defer func() {
			if cerr := file.Close(); cerr != nil {
				if err != nil {
					err = cerr
					log.Print(err)
				}
			}
		}()

		str := ""
		for _, v := range s.payments {
			str += fmt.Sprint(v.ID) + ";" + fmt.Sprint(v.AccountID) + ";" + fmt.Sprint(v.Amount) + ";" + fmt.Sprint(v.Category) + ";" + fmt.Sprint(v.Status) + "\n"
		}
		file.WriteString(str)
	}

	if len(s.favorites) > 0 {
		file, err := os.OpenFile(dir+"/favorites.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

		defer func() {
			if cerr := file.Close(); cerr != nil {
				if err != nil {
					err = cerr
					log.Print(err)
				}
			}
		}()

		str := ""
		for _, v := range s.favorites {
			str += fmt.Sprint(v.ID) + ";" + fmt.Sprint(v.AccountID) + ";" + fmt.Sprint(v.Amount) + ";" + fmt.Sprint(v.Category) + "\n"
		}
		file.WriteString(str)
	}

	return nil
}

func (s *Service) Import(dir string) error {

	_, err := os.Stat(dir + "/accounts.dump")

	if err == nil {
		content, err := ioutil.ReadFile(dir + "/accounts.dump")
		if err != nil {
			return err
		}

		strArray := strings.Split(string(content), "\n")
		if len(strArray) > 0 {
			strArray = strArray[:len(strArray)-1]
		}
		for _, v := range strArray {
			strArrAcount := strings.Split(v, ";")
			fmt.Println(strArrAcount)

			id, err := strconv.ParseInt(strArrAcount[0], 10, 64)
			if err != nil {
				return err
			}
			balance, err := strconv.ParseInt(strArrAcount[2], 10, 64)
			if err != nil {
				return err
			}
			flag := true
			for _, v := range s.accounts {
				if v.ID == id {
					v.Phone = types.Phone(strArrAcount[1])
					v.Balance = types.Money(balance)
					flag = false
				}
			}
			if flag {
				account := &types.Account{
					ID:      id,
					Phone:   types.Phone(strArrAcount[1]),
					Balance: types.Money(balance),
				}
				s.accounts = append(s.accounts, account)
			}
		}
	}

	_, err1 := os.Stat(dir + "/payments.dump")

	if err1 == nil {
		content, err := ioutil.ReadFile(dir + "/payments.dump")
		if err != nil {
			return err
		}

		strArray := strings.Split(string(content), "\n")
		if len(strArray) > 0 {
			strArray = strArray[:len(strArray)-1]
		}
		for _, v := range strArray {
			strArrAcount := strings.Split(v, ";")
			fmt.Println(strArrAcount)

			id := strArrAcount[0]
			if err != nil {
				return err
			}
			aid, err := strconv.ParseInt(strArrAcount[1], 10, 64)
			if err != nil {
				return err
			}
			amount, err := strconv.ParseInt(strArrAcount[2], 10, 64)
			if err != nil {
				return err
			}
			flag := true
			for _, v := range s.payments {
				if v.ID == id {
					v.AccountID = aid
					v.Amount = types.Money(amount)
					v.Category = types.PaymentCategory(strArrAcount[3])
					v.Status = types.PaymentStatus(strArrAcount[4])
					flag = false
				}
			}
			if flag {
				data := &types.Payment{
					ID:        id,
					AccountID: aid,
					Amount:    types.Money(amount),
					Category:  types.PaymentCategory(strArrAcount[3]),
					Status:    types.PaymentStatus(strArrAcount[4]),
				}
				s.payments = append(s.payments, data)
			}
		}
	}

	_, err2 := os.Stat(dir + "/favorites.dump")

	if err2 == nil {
		content, err := ioutil.ReadFile(dir + "/favorites.dump")
		if err != nil {
			return err
		}

		strArray := strings.Split(string(content), "\n")
		if len(strArray) > 0 {
			strArray = strArray[:len(strArray)-1]
		}
		for _, v := range strArray {
			strArrAcount := strings.Split(v, ";")
			fmt.Println(strArrAcount)

			id := strArrAcount[0]
			if err != nil {
				return err
			}
			aid, err := strconv.ParseInt(strArrAcount[1], 10, 64)
			if err != nil {
				return err
			}
			amount, err := strconv.ParseInt(strArrAcount[2], 10, 64)
			if err != nil {
				return err
			}
			flag := true
			for _, v := range s.favorites {
				if v.ID == id {
					v.AccountID = aid
					v.Amount = types.Money(amount)
					v.Category = types.PaymentCategory(strArrAcount[3])
					flag = false
				}
			}
			if flag {
				data := &types.Favorite{
					ID:        id,
					AccountID: aid,
					Amount:    types.Money(amount),
					Category:  types.PaymentCategory(strArrAcount[3]),
				}
				s.favorites = append(s.favorites, data)
			}
		}
	}

	return nil
}

func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {
	acc, err := s.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	var payments []types.Payment

	for _, pay := range s.payments {
		if pay.AccountID == acc.ID {
			data := types.Payment{
				ID:        pay.ID,
				AccountID: pay.AccountID,
				Amount:    pay.Amount,
				Category:  pay.Category,
				Status:    pay.Status,
			}
			payments = append(payments, data)
		}
	}
	return payments, nil
}

//HistoryToFiles ...
func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {

	str := ""

	if len(payments) > 0 && len(payments) <= records {
		file, _ := os.OpenFile(dir+"/payments.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		defer file.Close()

		for _, v := range payments {
			str += fmt.Sprint(v.ID) + ";" + fmt.Sprint(v.AccountID) + ";" + fmt.Sprint(v.Amount) + ";" + fmt.Sprint(v.Category) + ";" + fmt.Sprint(v.Status) + "\n"
		}
		file.WriteString(str)
	} else {
		k := 0 // limit on record
		t := 1 // count for files
		var file *os.File
		for _, v := range payments {
			if k == 0 {
				file, _ = os.OpenFile(dir+"/payments"+fmt.Sprint(t)+".dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			}
			k++
			str = fmt.Sprint(v.ID) + ";" + fmt.Sprint(v.AccountID) + ";" + fmt.Sprint(v.Amount) + ";" + fmt.Sprint(v.Category) + ";" + fmt.Sprint(v.Status) + "\n"
			_, _ = file.WriteString(str)
			if k == records { // если лимит был дастигнут, то обнулить "записи"
				str = ""
				t++
				k = 0
				file.Close()
			}
		}
	}
	return nil
}
func (s *Service) SumPayments(goroutines int) types.Money {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	sum := int64(0)
	kol := 0
	i := 0
	if goroutines == 0 {
		kol = len(s.payments)
	} else {
		kol = int(len(s.payments) / goroutines)
	}
	for i = 0; i < goroutines-1; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			val := int64(0)
			payments := s.payments[index*kol : (index+1)*kol]
			for _, payment := range payments {
				val += int64(payment.Amount)
			}
			mu.Lock()
			sum += val
			mu.Unlock()

		}(i)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		val := int64(0)
		payments := s.payments[i*kol:]
		for _, payment := range payments {
			val += int64(payment.Amount)
		}
		mu.Lock()
		sum += val
		mu.Unlock()

	}()
	wg.Wait()
	return types.Money(sum)
}

//FilterPayments ...
func (s *Service) FilterPayments(accountID int64, goroutines int) ([]types.Payment, error) {

	account, err := s.FindAccountByID(accountID)

	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	kol := 0
	i := 0
	var ps []types.Payment
	if goroutines == 0 {
		kol = len(s.payments)
	} else {
		kol = int(len(s.payments) / goroutines)
	}
	for i = 0; i < goroutines-1; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			var pays []types.Payment
			payments := s.payments[index*kol : (index+1)*kol]
			for _, v := range payments {
				if v.AccountID == account.ID {
					pays = append(pays, types.Payment{
						ID:        v.ID,
						AccountID: v.AccountID,
						Amount:    v.Amount,
						Category:  v.Category,
						Status:    v.Status,
					})
				}
			}
			mu.Lock()
			ps = append(ps, pays...)
			mu.Unlock()

		}(i)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		var pays []types.Payment
		payments := s.payments[i*kol:]
		for _, v := range payments {
			if v.AccountID == account.ID {
				pays = append(pays, types.Payment{
					ID:        v.ID,
					AccountID: v.AccountID,
					Amount:    v.Amount,
					Category:  v.Category,
					Status:    v.Status,
				})
			}
		}
		mu.Lock()
		ps = append(ps, pays...)
		mu.Unlock()

	}()
	wg.Wait()
	if len(ps) == 0 {
		return nil, nil
	}
	return ps, nil
}

//FilterPaymentsByFn ...
func (s *Service) FilterPaymentsByFn(filter func(payment types.Payment) bool, goroutines int) ([]types.Payment, error) {

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	kol := 0
	i := 0
	var ps []types.Payment
	if goroutines == 0 {
		kol = len(s.payments)
	} else {
		kol = int(len(s.payments) / goroutines)
	}
	for i = 0; i < goroutines-1; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			var pays []types.Payment
			payments := s.payments[index*kol : (index+1)*kol]
			for _, v := range payments {
				p := types.Payment{
					ID:        v.ID,
					AccountID: v.AccountID,
					Amount:    v.Amount,
					Category:  v.Category,
					Status:    v.Status,
				}

				if filter(p) {
					pays = append(pays, p)
				}
			}
			mu.Lock()
			ps = append(ps, pays...)
			mu.Unlock()

		}(i)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		var pays []types.Payment
		payments := s.payments[i*kol:]
		for _, v := range payments {

			p := types.Payment{
				ID:        v.ID,
				AccountID: v.AccountID,
				Amount:    v.Amount,
				Category:  v.Category,
				Status:    v.Status,
			}

			if filter(p) {
				pays = append(pays, p)
			}
		}
		mu.Lock()
		ps = append(ps, pays...)
		mu.Unlock()

	}()
	wg.Wait()
	if len(ps) == 0 {
		return nil, nil
	}
	return ps, nil
}

func (s *Service) SumPaymentsWithProgress() <-chan types.Progress {

	ch := make(chan types.Progress)

	size := 100_000
	parts := len(s.payments) / size
	wg := sync.WaitGroup{}
	i := 0
	if parts < 1 {
		parts = 1
	}
	for i := 0; i < parts; i++ {
		wg.Add(1)
		var payments []*types.Payment
		if len(s.payments) < size {
			payments = s.payments
		} else {
			payments = s.payments[i*size : (i+1)*size]
		}
		go func(ch chan types.Progress, data []*types.Payment) {
			defer wg.Done()
			val := types.Money(0)
			for _, v := range data {
				val += v.Amount
			}
			if len(s.payments) < size {
				ch <- types.Progress{
					Part:   len(data),
					Result: val,
				}
			}

		}(ch, payments)
	}
	if len(s.payments) > size {
		wg.Add(1)
		payments := s.payments[i*size:]
		go func(ch chan types.Progress, data []*types.Payment) {
			defer wg.Done()
			val := types.Money(0)
			for _, v := range data {
				val += v.Amount
			}
			ch <- types.Progress{
				Part:   len(data),
				Result: val,
			}

		}(ch, payments)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	return ch
}

func merge(channels []<-chan types.Progress) <-chan types.Progress {
	wg := sync.WaitGroup{}
	wg.Add(len(channels))
	merged := make(chan types.Progress)
	for _, ch := range channels {
		go func(ch <-chan types.Progress) {
			defer wg.Done()
			for val := range ch {
				merged <- val
			}
		}(ch)
	}
	go func() {
		defer close(merged)
		wg.Wait()
	}()
	return merged
}
