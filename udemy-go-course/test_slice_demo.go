package main

import "fmt"

func modifySlice(s []float64) {
	if len(s) > 0 {
		s[0] = 999.99 // Modify the first element
	}
	fmt.Println("Inside function:", s)
}

func testSlice() {
	prices := []float64{10.99, 20.49, 5.99, 15.99, 25.99, 30.99}
	priceSlice := prices[1:2]

	fmt.Println("--- Slice Passing Demo ---")
	fmt.Println("Before function call:")
	fmt.Println("prices:", prices)
	fmt.Println("priceSlice:", priceSlice)

	modifySlice(priceSlice)

	fmt.Println("After function call:")
	fmt.Println("prices:", prices)
	fmt.Println("priceSlice:", priceSlice)
	fmt.Println("--- End Demo ---")
}
