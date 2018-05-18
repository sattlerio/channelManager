package helpers

import (
	"channelManager.sattler.io/api"
	"os"
	"channelManager.sattler.io/clients"
	"strings"
)

func ValidateUserFromHeader(transactionId string, userId string) (bool, api.BasicResponse) {

	Init(os.Stderr, os.Stdout, os.Stdout, os.Stderr)

	var response = api.BasicResponse{}

	if len(transactionId) <= 0 || len(userId) <= 0 {
		Warning.Println("got no user Id and no transaction id in the header")
		response := api.BasicResponse{Status: "ERROR", StatusCode: 400, Message: "you have to be logged in to use this service", TransactionId: transactionId}

		return true, response
	}
	return false, response
}

func ValidateCompany(transactionId string, companyId string, userId string) (bool, api.BasicResponse, int) {
	Init(os.Stderr, os.Stdout, os.Stdout, os.Stderr)

	var response = api.BasicResponse{}

	if len(companyId) <= 0 {
		Info.Println(transactionId + ": no company id or shippping rule id in request")
		response := api.BasicResponse{Status: "ERROR", StatusCode: 400, Message: "you have to submit a valid company id as url param", TransactionId: transactionId}
		return true, response, 400
	}

	guardianClient := clients.GuardianClient{Host: os.Getenv("GUARDIAN_URL"), CompanyId: companyId, UserId: userId}

	guardianResponse, err := clients.CheckCompanyAndPermissionFromGuardian(guardianClient, 1)
	if err != nil {
		Info.Println(transactionId + ": guardian host responded with errror, abort transaction")
		Info.Println(err)
		response := api.BasicResponse{Status: "ERROR", StatusCode: 500, Message: "internal server error", TransactionId: transactionId}
		return true, response, 500
	}

	if !guardianResponse {
		Info.Println(transactionId + ": user is not allowed to access company / settings")
		response := api.BasicResponse{Status: "ERROR", StatusCode: 401, Message: "not allowed to access", TransactionId: transactionId}
		return true, response, 401
	}

	return false, response, 0
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.ToLower(b) == strings.ToLower(a) {
			return true
		}
	}
	return false
}
