package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var client *http.Client

// Initialize a client in our service.
// Will be reused for all http requests sent.
func createClient() {
	client = &http.Client{}
}

// Send http request to keycloak in order to Register a user.
// Firstly must acquire the access token
func createUser(name, job string) error {
	fmt.Println("Creating user...")

	fmt.Printf("Client: %v\n", client)

	// If there is no client initialized return..
	if client == nil {
		err := errors.New("no http client existing")
		return err
	}

	// Must
	apiUrl := "http://localhost:8080/auth/realms/master/protocol/openid-connect/token"
	bodyContent := "client_secret=HkQiIZCJEC3hSvGoLeTMDGModjWwEfxY&grant_type=client_credentials&client_id=admin-cli"

	// create a new http request
	request, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer([]byte(bodyContent)))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	// send the request
	response, err := client.Do(request)

	if err != nil {
		return err
	}

	responseBody, err := io.ReadAll(response.Body)

	if err != nil {
		return err
	}

	formattedData, err := formatJSON(responseBody)

	if err != nil {
		return err
	}

	fmt.Println("Status: ", response.Status)
	fmt.Println("Response body: ", formattedData)

	// clean up memory after execution
	defer response.Body.Close()

	return nil

}
