package main

import (
	"testing"
)

func TestFindMonthlyPaymentIteratively(t *testing.T) {
	cr := CalculationRequest{}
	cr.DurationInYears = 20
	cr.InterestRate = 0.01
	cr.LoanAmount = 500_000_00

	actual := findMonthlyPaymentIteratively(cr)
	if actual != 2_289_60 {
		t.Fatalf("Expected: %d but got %d", 2_289_60, actual)
	}
}
