package testdata

type Bank interface {
	CheckPin(cardNumber int, pin int) (accountNumber int, success bool)
}
