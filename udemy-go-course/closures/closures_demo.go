package closures

import "fmt"

// Example 1: Basic closure that captures a variable
func createCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// Example 2: Closure with parameters
func createMultiplier(factor int) func(int) int {
	return func(value int) int {
		return value * factor
	}
}

// Example 3: Closure that modifies captured variables
func createAccumulator() func(int) int {
	sum := 0
	return func(value int) int {
		sum += value
		return sum
	}
}

// Example 4: Multiple closures sharing the same variable
func createSharedCounter() (func() int, func() int, func()) {
	count := 0

	increment := func() int {
		count++
		return count
	}

	decrement := func() int {
		count--
		return count
	}

	reset := func() {
		count = 0
	}

	return increment, decrement, reset
}

// Example 5: Closure in a loop (common gotcha)
func createFunctionSlice() []func() int {
	functions := make([]func() int, 3)

	// Wrong way - all closures will capture the same variable
	for i := 0; i < 3; i++ {
		functions[i] = func() int {
			return i // This captures the loop variable i
		}
	}

	return functions
}

// Correct way to handle closures in loops
func createFunctionSliceCorrect() []func() int {
	functions := make([]func() int, 3)

	for i := 0; i < 3; i++ {
		// Create a new variable in each iteration
		val := i
		functions[i] = func() int {
			return val // This captures the local val variable
		}
	}

	return functions
}

func TestClosures() {
	fmt.Println("=== Go Closures Demo ===\n")

	// Example 1: Counter
	fmt.Println("1. Basic Counter Closure:")
	counter := createCounter()
	fmt.Println("First call:", counter())  // 1
	fmt.Println("Second call:", counter()) // 2
	fmt.Println("Third call:", counter())  // 3

	// Create another counter - it has its own state
	counter2 := createCounter()
	fmt.Println("New counter:", counter2())     // 1
	fmt.Println("Original counter:", counter()) // 4
	fmt.Println()

	// Example 2: Multiplier
	fmt.Println("2. Multiplier Closure:")
	double := createMultiplier(2)
	triple := createMultiplier(3)
	fmt.Println("Double 5:", double(5)) // 10
	fmt.Println("Triple 5:", triple(5)) // 15
	fmt.Println()

	// Example 3: Accumulator
	fmt.Println("3. Accumulator Closure:")
	acc := createAccumulator()
	fmt.Println("Add 10:", acc(10)) // 10
	fmt.Println("Add 20:", acc(20)) // 30
	fmt.Println("Add 5:", acc(5))   // 35
	fmt.Println()

	// Example 4: Shared counter
	fmt.Println("4. Shared Counter Closures:")
	inc, dec, reset := createSharedCounter()
	fmt.Println("Increment:", inc()) // 1
	fmt.Println("Increment:", inc()) // 2
	fmt.Println("Decrement:", dec()) // 1
	reset()
	fmt.Println("After reset, increment:", inc()) // 1
	fmt.Println()

	// Example 5: Loop closure gotcha
	fmt.Println("5. Loop Closure Gotcha:")
	wrongFuncs := createFunctionSlice()
	fmt.Println("Wrong way results:")
	for i, fn := range wrongFuncs {
		fmt.Printf("Function %d returns: %d\n", i, fn())
	}

	correctFuncs := createFunctionSliceCorrect()
	fmt.Println("Correct way results:")
	for i, fn := range correctFuncs {
		fmt.Printf("Function %d returns: %d\n", i, fn())
	}
	fmt.Println()

	// Example 6: Immediate invocation
	fmt.Println("6. Immediately Invoked Closure:")
	result := func(x, y int) int {
		return x + y
	}(10, 20)
	fmt.Println("10 + 20 =", result)

	// Example 7: Closure as callback
	fmt.Println("\n7. Closure as Callback:")
	numbers := []int{1, 2, 3, 4, 5}
	processNumbers(numbers, func(n int) int {
		return n * n // Square each number
	})
}

func processNumbers(numbers []int, processor func(int) int) {
	fmt.Print("Processed numbers: ")
	for _, num := range numbers {
		fmt.Print(processor(num), " ")
	}
	fmt.Println()
}
