package testdata

import "fmt"

type BankMock struct {
	Calls struct {
		CheckPin []BankMock_CheckPinArgs
	}

	Returns struct {
		CheckPin map[BankMock_CheckPinArgs]BankMock_CheckPinReturns
	}
}

func NewBankMock() *BankMock {
	newMock := new(BankMock)

	newMock.Returns.CheckPin = make(map[BankMock_CheckPinArgs]BankMock_CheckPinReturns)

	return newMock
}

func (b *BankMock) CheckPin(cardNumber int, pin int) (accountNumber int, success bool) {
	call := BankMock_CheckPinArgs{cardNumber, pin}
	b.Calls.CheckPin = append(b.Calls.CheckPin, call)

	if vals, exists := b.Returns.CheckPin[call]; exists {
		return vals.accountNumber, vals.success
	}

	panic(fmt.Sprintf("no return values found for args %+v, ive got %+v", call, b.Returns.CheckPin))
}

type BankMock_CheckPinArgs struct {
	cardNumber int
	pin        int
}

type BankMock_CheckPinReturns struct {
	accountNumber int
	success       bool
}
