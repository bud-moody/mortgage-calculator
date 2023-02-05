package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"mortgageCalculator/monthlyPayment"
	"net/http"
	"strconv"
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
	var bounds monthlyPayment.Bounds
	chBounds := make(chan monthlyPayment.Bounds)

	// Establish an upper bound
	var initialGuess float32
	for bounds.Upper == 0 {
		go monthlyPayment.MonthlyPaymentAndUpdateBounds(
			chBounds,
			initialGuess,
			int(calculationRequest.DurationInYears),
			int(calculationRequest.LoanAmount),
			calculationRequest.InterestRate)

		chBounds <- monthlyPayment.Bounds{Lower: 0, Upper: 0}
		bounds = <-chBounds
		initialGuess = initialGuess + 1000_00
	}

	// Try different values between lower and upper bound
	for int(math.Round(float64(bounds.Lower))) != int(math.Round(float64(bounds.Upper))) {
		for i := 0; i < numThreads; i++ {
			go monthlyPayment.MonthlyPaymentAndUpdateBounds(
				chBounds,
				(float32(i+1)*bounds.Lower+float32(numThreads-i)*bounds.Upper)/(numThreads+1),
				int(calculationRequest.DurationInYears),
				int(calculationRequest.LoanAmount),
				calculationRequest.InterestRate)
		}
		chBounds <- bounds
		bounds = <-chBounds // I think the problem is that this can receive before the 2nd or 3rd goroutine
		fmt.Printf("Bounds: %d, %d \n", int(bounds.Lower), int(bounds.Upper))
	}

	return int(bounds.Lower)
}
