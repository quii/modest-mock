package example

import "testing"

const success = true
const fail = false

func TestATM_NewSession_Checks_Pin(t *testing.T) {

	// inject the mocked bank into the ATM that we're testing
	bank := NewBankMock()
	atm := ATM{bank}

	accountNumber := 12345678
	cardNumber := 12345677
	pin := 9999

	// set up mock. the key of the map are your args to CheckPin, the values are.. your values!
	bank.Returns.CheckPin = BankMock_CheckPinReturnsMap{
		{cardNumber, pin}: {accountNumber, success},
		{123123, 5678}:    {0, fail},
	}

	// call the method you want to test
	session, err := atm.NewSession(cardNumber, pin)

	// make assertions
	if err != nil {
		t.Fatal("Didnt expect it to fail")
	}

	if session == nil {
		t.Fatal("expected a session")
	}

	if session.accountNumber != accountNumber {
		t.Error("Account number was not set correctly")
	}

	// demo-ing the "spy" functionality here
	if bank.Calls.CheckPin[0].cardNumber != cardNumber {
		t.Error("Bank was not called with correct card number")
	}

	if bank.Calls.CheckPin[0].pin != pin {
		t.Error("Bank was not called with correct pin")
	}
}
