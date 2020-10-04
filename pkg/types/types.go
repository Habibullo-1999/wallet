package types


// Money представляет собой в минимальных единицах (центы, копейки, дирамы и т.д.)
type Money int64

// PaymentCategory представляет собой категорию, в каторой был совершён платёж (авто, аптеки, рестораны и т.д.).
type PaymentCategory string

// PaymentStatus представляет собой статус платёжа
type PaymentStatus string


const (
  PaymentStatusOk PaymentStatus = "OK"
  PaymentStatusFail PaymentStatus = "FAIL"
  PaymentStatusInProgress PaymentStatus = "INPROGREES"
)

//Payment представляет информацию о платеже
type Payment struct {
  ID 			string
  AccountID 	int64
  Amount 		Money
  Category		PaymentCategory
  Status		PaymentStatus
}

type Phone string

type Account struct {
  ID     	int64
  Phone  	Phone
  Balance 	Money
}

type Favorite struct {
ID          string
AccountID   int64
Name        string
Amount      Money
Category    PaymentCategory
}

