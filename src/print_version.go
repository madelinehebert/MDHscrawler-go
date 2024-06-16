package main

import (
	"fmt"
	"runtime"
)

// Function to print program version
func printv() {
	fmt.Printf("scrawler version %.2f %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
}
