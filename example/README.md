# modest-mock example

This project shows a simple interface and how you can generate a mock with a `go generate`

```go
//go:generate modestmock -name=Bank -out=bank_mock.go
type Bank interface {
	CheckPin(cardNumber, pin int) (accountNumber int, success bool)
	Deposit(accountNumber string, amount int) (newBalance int, err error)
	Withdraw(accountNumber string, amount int) (newBalance int, err error)
}
```

See the tests (todo!) for examples of writing tests using the generated mock. 

## try it

`$ go generate`

`$ go test`