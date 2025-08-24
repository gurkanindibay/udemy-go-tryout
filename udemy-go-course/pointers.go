package main
import "fmt"

func pointer_test() {
	var age int = 42
	var agePointer *int = &age
	fmt.Println("Age:", *agePointer)
	fmt.Println("Age directly:", age)
	fmt.Println("Address of age:", agePointer)
	editAgeToAdultYears(agePointer)
	fmt.Println("Age after edit:", age)

}

func editAgeToAdultYears(agePtr *int) {
	*agePtr = *agePtr -18
	
}