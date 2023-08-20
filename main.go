package main

import (
	"fmt"
	"os"
)

func main() {
	// os.Args provides access to raw command-line arguments.
	// os.Args[0] is the name of the program itself, os.Args[1] is the first argument, and so on.
	if len(os.Args) < 2 {
		fmt.Println("You must provide an argument!")
		return
	}

	// Get the argument and print it.
	argument := os.Args[1]
	fmt.Println(argument)
}
