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
func LoginRequest(username, password string) error {
	fmt.Println("Trying to log in as user...")

	client_secret := "PcOdI78QbDPUlJExmVO3hieF6gRJFAsz"

	// If there is no client initialized return..
	if client == nil {
		err := errors.New("no http client existing")
		return err
	}

	// Must aqcuire access token
	url := "http://localhost:8080/auth/admin/realms/master/protocol/openid-connect/token"
	payload := strings.NewReader(fmt.Sprintf("username=%s&password=%s&client_secret=%s&grant_type=client_credentials&client_id=admin-cli", username, password, client_secret))

	// create a new http request
	request, err := http.NewRequest("POST", url, payload)
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

	var data_json map[string]interface{}

	json.Unmarshal([]byte(formattedData), &data_json)

	if err != nil {
		return err
	}

	return nil

}

func RegisterRequest(username, password, email, role string) error {
	fmt.Println("Trying to register a user...")

	client_secret := "PcOdI78QbDPUlJExmVO3hieF6gRJFAsz"

	// If there is no client initialized return..
	if client == nil {
		err := errors.New("no http client existing")
		return err
	}

	// Must aqcuire access token
	apiUrl := "http://localhost:8080/auth/realms/master/protocol/openid-connect/token"
	payload := strings.NewReader(fmt.Sprintf("client_secret=%s&grant_type=client_credentials&client_id=admin-cli", client_secret))

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

	var data_json map[string]interface{}

	json.Unmarshal([]byte(formattedData), &data_json)

	if err != nil {
		return err
	}

	access_token := data_json["access_token"]

	fmt.Print(access_token)

	return nil
}
