package monthlyPayment

import (
	"fmt"
	"sync"
)

func MonthlyPaymentAndUpdateBounds(
	mutex *sync.Mutex,
	wg *sync.WaitGroup,
	payment float32,
	lowerBound *float32,
	upperBound *float32,
	years int,
	loanAmount int,
	interestRate float32) {

	defer wg.Done()

	var amount float32 = float32(loanAmount)
	var months int = 0
	for months < years*12+1 {
		amount = float32(amount) + (amount * interestRate / 12) - float32(payment)
		months++
	}

	mutex.Lock()
	if amount < 0 && *upperBound == 0 { // For initializing upper bound
		*upperBound = payment
	} else if amount < 0 && payment < *upperBound {
		*upperBound = payment
	} else if amount > 0 && payment > *lowerBound {
		*lowerBound = payment
	}
	fmt.Printf("%f, %f, %f, %f \n", amount, *lowerBound, *upperBound, payment)
	mutex.Unlock()
}
