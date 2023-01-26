package monthlyPayment

import (
	"fmt"
	"sync"
)

const tolerance = 500_00 // TODO allow this to be lower by using decimals for working out monthly payment
func MonthlyPaymentAndUpdateBounds(
		mutex *sync.Mutex,
		wg *sync.WaitGroup,
		payment int, 
		lowerBound *int, 
		upperBound *int, 
		years int, 
		loanAmount int, 
		interestRate float32) {

	defer wg.Done()
			
	var amount float32 = float32(loanAmount)
	var months int = 0
	for months < years * 12 + 1 {
		amount = float32(amount) + (amount * interestRate / 12) - float32(payment)
		months++
	}

	mutex.Lock()
	if (*upperBound == 0 && amount < 0) {
		*upperBound = payment;
	} else if (amount < 0 && payment < *upperBound) {
		*upperBound = payment
	} else if (amount > tolerance && payment > *lowerBound) {
		*lowerBound = payment
	} else if (amount >= 0 && amount <= tolerance) {
		*lowerBound = payment
		*upperBound = payment
	}
	fmt.Printf("%f, %d, %d, %d \n", amount, *lowerBound, *upperBound, payment)
	mutex.Unlock()
}