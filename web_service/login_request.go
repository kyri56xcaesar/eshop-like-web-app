package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func send() {

	url := "http://localhost:8080/auth/realms/eshop_project/protocol/openid-connect/token"

	payload := strings.NewReader("username=test&password=test&client_id=client-front&client_secret=FovwjKxMhL4YSjOftIXavlDzkcAqX9Dr&grant_type=password")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "Thunder Client (https://www.thunderclient.com)")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}
