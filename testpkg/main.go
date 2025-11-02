package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting application...")
	// Since this is already in the testpkg package, 
	// we don't need to import testpkg to call Run()
	// However, I notice that the Run function is not defined in this file
	// You would need to either define it here or import it from another package
	fmt.Println("Run function not implemented")
}