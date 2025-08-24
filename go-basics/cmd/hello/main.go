package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/example/go-basics/pkg/utils"
)

func main() {
	name := flag.String("name", "World", "name to greet")
	repeat := flag.Int("repeat", 1, "number of times to print the greeting")
	flag.Parse()

	if *repeat < 1 {
		fmt.Fprintln(os.Stderr, "repeat must be >= 1")
		os.Exit(1)
	}

	p := utils.Person{Name: *name}
	for i := 0; i < *repeat; i++ {
		fmt.Println(p.Greet())
	}
}
