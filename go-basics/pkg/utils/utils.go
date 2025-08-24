package utils

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Person demonstrates a simple struct with methods.
type Person struct {
	Name string
	Age  int
}

// Greet returns a greeting for the person.
func (p Person) Greet() string {
	// Use unicode-aware title-casing instead of deprecated strings.Title
	titler := cases.Title(language.Und)
	return fmt.Sprintf("Hello, %s!", titler.String(p.Name))
}

// SetAge sets the person's age (pointer receiver example).
func (p *Person) SetAge(age int) {
	p.Age = age
}

// Summarize returns a short summary and demonstrates multiple return values.
func (p Person) Summarize() (string, bool) {
	if p.Age == 0 {
		return fmt.Sprintf("%s has unknown age", p.Name), false
	}
	return fmt.Sprintf("%s is %d years old", p.Name, p.Age), true
}

// SumInts sums a slice of ints (demonstrates slices and for-range).
func SumInts(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// MapKeys returns the keys of a string->int map as a slice (demonstrates maps).
func MapKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// ErrorExample demonstrates returning an error.
func ErrorExample(input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", fmt.Errorf("input must not be empty")
	}
	return strings.ToUpper(input), nil
}
