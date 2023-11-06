package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Print("Welcome to the Product Service.\n")

	// Load environment variables using godotenv
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	fmt.Printf("Listening on PORT =: %v\n", portString)

	log.Fatal(http.ListenAndServe(portString, nil))

}
