package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var client *http.Client

// Initialize a client in our service.
// Will be reused for all http requests sent.
func createClient() {
	client = &http.Client{}

}

// Send http request to keycloak in order to Register a user.
// Firstly must acquire the access token
func sendLoginUserRequest(username, password string) error {
	fmt.Println("Trying to log in as user...")

	// If there is no client initialized return..
	if client == nil {
		err := errors.New("no http client existing")
		return err
	}

	// Must aqcuire access token
	apiUrl := "http://localhost:8080/auth/realms/master/protocol/openid-connect/token"
	payload := strings.NewReader("client_secret=HkQiIZCJEC3hSvGoLeTMDGModjWwEfxY&grant_type=client_credentials&client_id=admin-cli")

	// create a new http request
	request, err := http.NewRequest("POST", apiUrl, payload)
	request.Header.Add("Accept", "*/*")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// send the request
	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

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

	var x map[string]interface{}

	json.Unmarshal([]byte(formattedData), &x)

	return nil

}

func RegisterRequest() {
	fmt.Println("Trying to register a user...")
}
