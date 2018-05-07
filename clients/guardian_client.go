package clients

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

var err error

type GuardianClient struct {
	Host      string
	CompanyId string
	UserId    string
}

type GuardianResponse struct {
	Status string               `json:"status"`
	Data   GuardianResponseData `json:"data"`
}

type GuardianResponseData struct {
	UserPermission int `json:"user_permission"`
}

func CheckCompanyAndPermissionFromGuardian(client GuardianClient, permission int) (bool, error) {

	url := fmt.Sprintf("%s/%s/%s", client.Host, client.UserId, client.CompanyId)
	response, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}
	if response.StatusCode != 200 {
		return false, nil
	}

	guardianData := GuardianResponse{}
	jsonErr := json.Unmarshal(body, &guardianData)

	if jsonErr != nil {
		return false, jsonErr
	}

	if guardianData.Status != "OK" {
		return false, nil
	}


	if guardianData.Data.UserPermission < 0 || guardianData.Data.UserPermission > permission {
		return false, nil
	}

	return true, err
}
