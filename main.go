package main

import (
	"sync"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/calculate", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func handler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)

	var cr CalculationRequest
	err := json.NewDecoder(r.Body).Decode(&cr)
	if (err != nil) {
		fmt.Println(err)
		panic(err)
	}

	fmt.Print(cr)

	w.Write([]byte(strconv.Itoa(int(monthlyPayment(cr)))))
}

type CalculationRequest struct {
	DurationInYears uint `json:",string"`
	LoanAmount uint `json:",string"`
	InterestRate float32 `json:",string"`
}

func monthlyPayment(calculationRequest CalculationRequest) int {

	lowerBound := 0
	upperBound := 0

	var mutex sync.Mutex

	// Establish an upper bound
	initialGuess := 0
	for upperBound == 0 {
		var wg sync.WaitGroup
		wg.Add(1)
		calculateForMonthlyPayment(
			&mutex,
			&wg,
			initialGuess,
			&lowerBound,
			&upperBound, 
			int(calculationRequest.DurationInYears), 
			int(calculationRequest.LoanAmount), 
			calculationRequest.InterestRate)
		initialGuess = initialGuess + 1000_00
		wg.Wait()
	}

	// Try different values between lower and upper bound
	for lowerBound != upperBound {
		fmt.Println("here we go")
		var wg sync.WaitGroup

		wg.Add(1)
		go calculateForMonthlyPayment(
			&mutex,
			&wg,
			(3*lowerBound + upperBound) / 4,
			&lowerBound,
			&upperBound, 
			int(calculationRequest.DurationInYears), 
			int(calculationRequest.LoanAmount), 
			calculationRequest.InterestRate)
		
		wg.Add(1)
		go calculateForMonthlyPayment(
			&mutex,
			&wg,
			(2*lowerBound + 2*upperBound) / 4,
			&lowerBound,
			&upperBound, 
			int(calculationRequest.DurationInYears), 
			int(calculationRequest.LoanAmount), 
			calculationRequest.InterestRate)

		wg.Add(1)
		go calculateForMonthlyPayment(
			&mutex,
			&wg,
			(lowerBound + 3*upperBound) / 4,
			&lowerBound,
			&upperBound, 
			int(calculationRequest.DurationInYears), 
			int(calculationRequest.LoanAmount), 
			calculationRequest.InterestRate)

		wg.Wait()
	}

	return lowerBound
}

func calculateForMonthlyPayment(
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
	} else if (amount > 10_00 && payment > *lowerBound) {
		*lowerBound = payment
	} else if (amount >= 0 && amount <= 10_00) {
		*lowerBound = payment
		*upperBound = payment
	}
	fmt.Printf("%f, %d, %d, %d \n", amount, *lowerBound, *upperBound, payment)
	mutex.Unlock()
}
