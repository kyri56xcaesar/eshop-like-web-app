package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("./web/index.html"))

type PageVariables struct {
	LoginError    string
	RegisterError string
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{})
}

func handlerErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Something went wrong")
}

func HandleRootGet(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n------------------------------------\n")
	fmt.Printf("\nMethod type: %v\n", r.Method)

	// Serve the HTML form
	tmpl.Execute(w, PageVariables{LoginError: ""})
}

func HandleRootPost(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n------------------------------------\n")
	fmt.Printf("\nMethod type: %v\n", r.Method)

	//fmt.Print(r)

	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve and process the form data

	for key, value := range r.Form {
		fmt.Printf("%s = %s\n", key, value)
	}

	switch len(r.Form) {
	case 2:
		// Login form request
		username := r.FormValue("login-username")
		password := r.FormValue("login-password")

		// Filter characters

		// Send Login Request
		status, err := LoginRequest(username, password)
		fmt.Printf("Login request response: %d\n", status)

		if err != nil {
			fmt.Printf("%v", err)
		}

		if int(status/100) == 2 {
			// Successfull Login
			return

		} else {
			// Unsuccessful Login
			//loginError := "Invalid username or password"

			return
		}

		// If logged in, respond accordingly
	case 5:
		// Register form request
		username := r.FormValue("register-username")
		password := r.FormValue("register-password")
		password_r := r.FormValue("register-password-repeat")

		if password != password_r {
			fmt.Println("Passwords dont match.")

		}
		email := r.FormValue("register-email")
		role := r.FormValue("select-role")

		// Filter characters do checks, respond accordingly

		// Send Register request
		status, err := RegisterRequest(username, password, email, role)

		if err != nil {
			fmt.Printf("%v", err)
		}

		fmt.Printf("Register request response: %d\n", status)

		// If successful, respond accordingly

	default:
		http.Error(w, "Status Not allowed", http.StatusMethodNotAllowed)

	}

	// Validate user credintials and redirect accordingly.
	// Perhaps use html/template to form the ServedFile

	// Send request to other service

	// tmpl.Execute(w, PageVariables{})

}
