package main

import (
	"fmt"
	"math"
)

func investmentCalculator() {
	const inflationRate = 6.5
	var investmentAmount, years float64 = 1000.0, 10.0
	expectedReturnRate := 5.5

	fmt.Print("Investment Amount:")
	fmt.Scan(&investmentAmount)

	fmt.Print("Years:")
	fmt.Scan(&years)

	fmt.Print("Expected Return Rate:")
	fmt.Scan(&expectedReturnRate)

	futureValue, inflationAdjustedValue := calculateFutureValues(investmentAmount, expectedReturnRate, years, inflationRate)
	fmt.Printf("Future value of the investment: %.2f\n", futureValue)


	var stringToPrint = fmt.Sprintf("Future value adjusted for inflation: %.2f\n", inflationAdjustedValue)
	fmt.Print(stringToPrint)

	fmt.Println(`Multi Line is also enabled
	 You can enter more than one line here`)
}

func calculateFutureValues(investmentAmount, expectedReturnRate, years, inflationRate float64) (float64, float64) {
	futureValue := investmentAmount * math.Pow(1+expectedReturnRate/100, years)
	inflationAdjustedValue := futureValue / math.Pow(1+inflationRate/100, years)
	return futureValue, inflationAdjustedValue
}