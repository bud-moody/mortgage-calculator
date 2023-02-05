package monthlyPayment

import (
	"fmt"
)

type Bounds struct {
	Lower float32
	Upper float32
}

func MonthlyPaymentAndUpdateBounds( // ALWAYS receives AND sends to both cLower and cUpper
	chBounds chan Bounds,
	payment float32,
	years int,
	loanAmount int,
	interestRate float32) {

	var amount float32 = float32(loanAmount)
	var months int = 0
	for months < years*12+1 {
		amount = float32(amount) + (amount * interestRate / 12) - float32(payment)
		months++
	}

	bounds := <-chBounds
	if amount < 0 && bounds.Upper == 0 { // For initializing upper bound
		bounds.Upper = payment
	} else if amount < 0 && payment < bounds.Upper {
		bounds.Upper = payment
	} else if amount > 0 && payment > bounds.Lower {
		bounds.Lower = payment
	}

	fmt.Printf("%f, %f, %f, %f \n", amount, bounds.Lower, bounds.Upper, payment)
	chBounds <- bounds
}
