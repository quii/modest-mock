package example

//go:generate modestmock -name=Bank -out=bank_mock.go
type Bank interface {
	CheckPin(cardNumber, pin int) (accountNumber int, success bool)
	Deposit(accountNumber string, amount int) (newBalance int, err error)
	Withdraw(accountNumber string, amount int) (newBalance int, err error)
}
