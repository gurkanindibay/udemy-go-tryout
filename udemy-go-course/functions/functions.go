package functions

import "fmt"

type transformFn func(int) int

func transformNumbers(numbers []int, fn transformFn) []int {
	var result []int
	for _, n := range numbers {
		result = append(result, fn(n))
	}
	return result
}

func double(n int) int {
	return n * 2
}

func triple(n int) int {
	return n * 3
}

func TestFunctions() {
	numbers := []int{1, 2, 3, 4, 5}

	doubled := transformNumbers(numbers, double)
	tripled := transformNumbers(numbers, triple)

	fmt.Println("Doubled:", doubled)
	fmt.Println("Tripled:", tripled)

}
