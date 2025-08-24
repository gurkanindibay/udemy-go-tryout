package collections

import "fmt"

func modifySlice(s []float64) {
	if len(s) > 0 {
		s[0] = 999.99 // Modify the first element
	}
	fmt.Println("Inside function:", s)
}

func TestLists() {
	prices := []float64{10.99, 20.49, 5.99, 15.99, 25.99, 30.99}

	priceSlice := prices[1:2]

	fmt.Println(prices)
	fmt.Println(priceSlice)

	fmt.Println(len(prices), cap(prices))         // 6 6
	fmt.Println(len(priceSlice), cap(priceSlice)) // 1 5

	// Demonstrate slice passing behavior
	fmt.Println("\n--- Slice Passing Demo ---")
	fmt.Println("Before function call:")
	fmt.Println("prices:", prices)
	fmt.Println("priceSlice:", priceSlice)

	modifySlice(priceSlice)

	fmt.Println("After function call:")
	fmt.Println("prices:", prices)
	fmt.Println("priceSlice:", priceSlice)
	fmt.Println("--- End Demo ---")

	discountPrices := []float64{1.99, 2.49, 0.99, 1.49, 2.99, 3.99}
	priceSlice = append(priceSlice, discountPrices...)

	fmt.Println(priceSlice)

	websites := map[string]string{
		"Google":   "https://www.google.com",
		"Facebook": "https://www.facebook.com",
		"Twitter":  "https://www.twitter.com",
	}

	fmt.Println(websites)

	for name, url := range websites {
		fmt.Printf("%s: %s\n", name, url)
	}

	if googleUrl, ok := websites["Google"]; ok {
		fmt.Println("Google URL:", googleUrl)
	} else {
		fmt.Println("Google URL not found")
	}

	delete(websites, "Twitter")
	fmt.Println(websites)

	websites["LinkedIn"] = "https://www.linkedin.com"
	fmt.Println(websites)
}
