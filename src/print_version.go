package main

import (
	"fmt"
	"runtime"
)

// Function to print program version
func printV() {
	fmt.Printf("scrawler version %.2f %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
}
