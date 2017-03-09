package example

import "testing"

func TestATM_NewSession(t *testing.T) {

	bank := &BankMock{}

	t.Run("it goes to the bank to check the users pin before starting a session", func(t *testing.T) {

		accountNumber := 12345678
		cardNumber := 12345677
		pin := 9999

		// set up the bank so the first time it is called then it returns our expected account number
		//todo: order based stuff sucks, should be based on arguments sent in
		bank.Returns.CheckPin = []struct{
			accountNumber int
			success bool
		}{
			{accountNumber, true},
		}

		atm := ATM{bank}


		session, err := atm.NewSession(cardNumber, pin)

		if err != nil {
			t.Fatal("Didnt expect it to fail")
		}

		if session==nil {
			t.Fatal("expected a session")
		}

		if session.accountNumber != accountNumber {
			t.Error("Account number was not set correctly")
		}

		if bank.Calls.CheckPin[0].cardNumber != cardNumber {
			t.Error("Bank was not called with correct card number")
		}

		if bank.Calls.CheckPin[0].pin != pin {
			t.Error("Bank was not called with correct pin")
		}


	})
}