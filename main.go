package main

import (
	"fmt"
	"os"
)

func main() {
	os.Exit(run())
}

func run() int {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: json-parser <file>")
		return 2
	}

	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return 2
	}

	isValid, _ := ValidateJSON(string(input))
	if isValid {
		fmt.Fprintln(os.Stdout, "Valid JSON")
		return 0
	}

	fmt.Fprintln(os.Stdout, "Invalid JSON")
	return 1
}
