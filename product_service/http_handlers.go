package main

import (
	"fmt"
	"net/http"
)

func getRoot(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Arrived here~.")
	// fmt.Print(r)
	fmt.Print(r.Method)
	fmt.Print("\n\n")

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "post" {

	}

	if r.Method == "get" {

	}

	if r.Method == "put" {

	}

	if r.Method == "delete" {

	}

}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")

	fmt.Print(r.Method)
	fmt.Print("\n\n")

}
