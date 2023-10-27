// Package main is the runnable code /executable scripts.
// Every other package is a library.

// Go build -> creates an executable
// Go run -> creates an executable and runs it
// Go install -> creates an executable and installs it in the machine, making it global

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello world")

	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the enviroment")
	}

	fmt.Println("Port:", portString)
}

// Reverse reverses a string left to right
// Notice that we need to capitalize the first letter
// of the function

// If we don't then we won't be able to access this
// function outside of the mystrings package

func Reverse(s string) string {
	result := ""
	for _, v := range s {
		result = string(v) + result
	}
	return result
}

// Uppercase name functions are exported
