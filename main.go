package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mortgageCalculator/monthlyPayment"
	"net/http"
	"strconv"
	"sync"
)

const numThreads = 4

func main() {
	http.HandleFunc("/calculate", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)

	var cr CalculationRequest
	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Print(cr)

	w.Write([]byte(strconv.Itoa(int(findMonthlyPaymentIteratively(cr)))))
}

type CalculationRequest struct {
	DurationInYears uint    `json:",string"`
	LoanAmount      uint    `json:",string"`
	InterestRate    float32 `json:",string"`
}

func findMonthlyPaymentIteratively(calculationRequest CalculationRequest) int {
	var lowerBound float32
	var upperBound float32
	var mutex sync.Mutex

	// Establish an upper bound
	var initialGuess float32
	for upperBound == 0 {
		var wg sync.WaitGroup
		wg.Add(1)
		monthlyPayment.MonthlyPaymentAndUpdateBounds(
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
	for int(lowerBound) != int(upperBound) {
		var wg sync.WaitGroup
		for i := 0; i < numThreads; i++ {
			wg.Add(1)
			go monthlyPayment.MonthlyPaymentAndUpdateBounds(
				&mutex,
				&wg,
				(float32(i+1)*lowerBound+float32(numThreads-i)*upperBound)/(numThreads+1),
				&lowerBound,
				&upperBound,
				int(calculationRequest.DurationInYears),
				int(calculationRequest.LoanAmount),
				calculationRequest.InterestRate)
		}

		wg.Wait()
	}

	return int(lowerBound)
}
