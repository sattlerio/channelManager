package clients

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"fmt"
)

type StripeClient struct {
	Host   string
	ApiKey string
}

type StripeClientResponse struct {
	Status      string `json:"status; omitempty"`
	Message 	string `json:"message; omitempty"`
	ApiFeedback string `json:"api_feedback; omitempty"`
}

func ValidateStripeCredentials(client StripeClient) (bool, error, int) {
	var statusCode int

	url := client.Host

	values := map[string]string{"api_key": client.ApiKey}

	jsonValue, _ := json.Marshal(values)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		statusCode = 500
		return false, err, statusCode
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err, statusCode
	}
	if response.StatusCode != 200 {
		return false, nil, response.StatusCode
	}

	stripeData := StripeClientResponse{}
	jsonErr := json.Unmarshal(body, &stripeData)

	fmt.Println(jsonErr)

	if jsonErr != nil {
		statusCode = 500
		return false, jsonErr, statusCode
	}

	if stripeData.Status != "OK" {
		statusCode = 200
		return false, nil, statusCode
	}

	return true, err, statusCode
}
