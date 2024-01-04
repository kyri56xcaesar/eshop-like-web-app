package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var client *http.Client

// Initialize a client in our service.
// Will be reused for all http requests sent.
func createClient() {

	godotenv.Load(".env")

	client = &http.Client{}
}

// Send http request to keycloak in order to Register a user.
// Firstly must acquire the access token
func LoginRequest(username, password string) (int, error) {
	fmt.Println("Trying to log in as user...")

	client_secret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
	port := os.Getenv("KEYCLOAK_PORT")
	realm := os.Getenv("KEYCLOAK_REALM")

	if client_secret == "" || port == "" || realm == "" {
		log.Fatal("Environment variables not parsed correctly...")
	}

	// If there is no client initialized return..
	if client == nil {
		err := errors.New("no http client existing")
		return -1, err
	}

	url := fmt.Sprintf("http://localhost:%s/auth/realms/%s/protocol/openid-connect/token", port, realm)
	payload := strings.NewReader(fmt.Sprintf("username=%s&password=%s&client_secret=%s&grant_type=password&client_id=admin-cli", username, password, client_secret))

	// create a new http request
	request, err := http.NewRequest("POST", url, payload)
	//request.Header.Add("Accept", "*/*")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// send the request
	response, err := client.Do(request)

	if err != nil {
		return -1, err
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)

	if err != nil {
		return -1, err
	}

	formattedData, err := formatJSON(responseBody)

	if err != nil {
		return -1, err
	}

	status := response.StatusCode
	//fmt.Println("Status: ", status)
	//fmt.Println("Response body: ", formattedData)

	var data_json map[string]interface{}

	json.Unmarshal([]byte(formattedData), &data_json)

	if err != nil {
		return -1, err
	}

	return status, nil

}

func RegisterRequest(username, password, email, role string) (int, error) {
	fmt.Println("Trying to register a user...")

	client_secret := os.Getenv("KEYCLOAK_CLIENT_SECRET")
	port := os.Getenv("KEYCLOAK_PORT")
	realm := os.Getenv("KEYCLOAK_REALM")

	if client_secret == "" || port == "" || realm == "" {
		log.Fatal("Environment variables not parsed correctly...")
	}

	// If there is no client initialized return..
	if client == nil {
		err := errors.New("no http client existing")
		return -1, err
	}

	// Must aqcuire access token
	apiUrl := fmt.Sprintf("http://localhost:%s/auth/realms/master/protocol/openid-connect/token", port)
	payload := strings.NewReader(fmt.Sprintf("client_secret=%s&grant_type=client_credentials&client_id=admin-cli", client_secret))

	// create a new http request
	request, err := http.NewRequest("POST", apiUrl, payload)
	request.Header.Add("Accept", "*/*")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// send the request
	response, err := client.Do(request)

	if err != nil {
		return -1, err
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)

	if err != nil {
		return -1, err
	}

	formattedData, err := formatJSON(responseBody)

	if err != nil {
		return -1, err
	}

	//fmt.Println("Status: ", response.Status)
	//fmt.Println("Response body: ", formattedData)

	var data_json map[string]interface{}

	json.Unmarshal([]byte(formattedData), &data_json)

	if err != nil {
		return -1, err
	}

	access_token := data_json["access_token"]

	//fmt.Print(access_token)

	// Create a register user request.
	apiUrl = fmt.Sprintf("http://localhost:%s/auth/admin/realms/%s/users", port, realm)
	payload = strings.NewReader(fmt.Sprintf("{\n  \"email\":\"%s\",\n  \"enabled\":\"true\",\n  \"username\":\"%s\",\n  \"attributes\":{\n    \"client_id\":\"frontend_app\"\n  },\n  \"groups\":[\"%s\"],\n  \"credentials\":[{\"type\":\"password\",\"value\":\"%s\",\"temporary\":\"false\"}]\n}", email, username, role, password))

	reg_request, err := http.NewRequest("POST", apiUrl, payload)
	reg_request.Header.Add("Accept", "*/*")
	reg_request.Header.Add("Content-Type", "application/json")
	reg_request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", access_token))

	reg_response, err := client.Do(reg_request)

	if err != nil {
		return -1, err
	}

	defer reg_response.Body.Close()
	//body, _ := io.ReadAll(reg_response.Body)

	status := reg_response.StatusCode
	//fmt.Println("Status: ", status)

	//fmt.Println(reg_response)
	//fmt.Println(string(body))

	return status, nil
}

func LogoutRequest() {

}
