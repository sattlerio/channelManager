package clients

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"fmt"
)

type BraintreeClient struct {
	Host      	string
	PublicKey 	string
	PrivateKey	string
	MerchantId 	string
	Sandbox 	bool
}

type BraintreeResponse struct {
	Status      string `json:"status; omitempty"`
	Message 	string `json:"message; omitempty"`
}

func ValidateBraintreeCredentials(client BraintreeClient) (bool, error, int) {

	var (
		statusCode 	int
		environment string
	)
	if client.Sandbox {
		environment = "sandbox"
	} else {
		environment = "production"
	}

	url := client.Host

	values := map[string]string{
		"merchant_id": client.MerchantId,
		"environment": environment,
		"public_key": client.PublicKey,
		"private_key": client.PrivateKey,
	}

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

	braintreeData := BraintreeResponse{}
	jsonErr := json.Unmarshal(body, &braintreeData)

	fmt.Println(jsonErr)

	if jsonErr != nil {
		statusCode = 500
		return false, jsonErr, statusCode
	}

	if braintreeData.Status != "INFO" {
		statusCode = 200
		return false, nil, statusCode
	}

	return true, err, statusCode
}
