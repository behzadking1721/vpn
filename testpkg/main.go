package main

import (
	"fmt"
	"vpn-client/testpkg"
)

func main() {
	fmt.Println("Starting application...")
	testpkg.Run()
}
package testpkg

import (
	"fmt"
	"vpn-client/internal/managers"
)

func Run() {
	fmt.Println("Testing import...")
	cm := managers.NewConnectionManager()
	if cm != nil {
		fmt.Println("Import successful!")
	} else {
		fmt.Println("Import failed!")
	}
}