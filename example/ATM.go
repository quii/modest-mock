package example

import "time"

type ATM struct {
	bank Bank
}

func (a *ATM) NewSession(cardNumber, pin int) (*Session, error) {
	account, _ := a.bank.CheckPin(cardNumber, pin)

	return &Session{
		bank:          a.bank,
		accountNumber: account,
		timeout:       2 * time.Minute,
	}, nil
}

type Session struct {
	bank          Bank
	accountNumber int
	timeout       time.Duration
}
