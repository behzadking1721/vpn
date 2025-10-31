package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("Go Environment Test")
	fmt.Println("==================")
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("Operating System: %s\n", runtime.GOOS)
	fmt.Printf("Architecture: %s\n", runtime.GOARCH)
	fmt.Println()
	fmt.Println("Environment is properly set up for VPN Client development!")
}