package controllers

import (
	"net/http"
	"os"
	"encoding/json"
	"channelManager/helpers"
	"github.com/jinzhu/gorm"
	"github.com/gorilla/mux"
	"channelManager/api"
	"channelManager/models"
	"io/ioutil"
	"fmt"
	"channelManager/clients"
	"strconv"
	"github.com/rs/xid"
)

var DbConn *gorm.DB

func FetchAllChannels(w http.ResponseWriter, r *http.Request) {

	helpers.Init(os.Stderr, os.Stdout, os.Stdout, os.Stderr)

	helpers.Info.Println("new request for fetching all channels")

	transactionId := r.Header.Get("x-transactionid")
	userId := r.Header.Get("x-user-uuid")

	if status, response := helpers.ValidateUserFromHeader(transactionId, userId); status {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(response)
		return
	}

	params := mux.Vars(r)
	company_id := params["company_id"]

	if status, response, statusCode := helpers.ValidateCompany(transactionId, company_id, userId); status {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}
	helpers.Info.Println(transactionId + " user is allowed to access, continue with request")

	helpers.Info.Println(transactionId + ": got new request for ping route")

	response := api.ChannelTypeResponse{StatusCode: 200, Status: "OK", Message: "fetched successfully data",
		TransactionId: transactionId, Data: models.Types}
	json.NewEncoder(w).Encode(response)
	return
}

func CreateNewChannel(w http.ResponseWriter, r *http.Request) {

	helpers.Init(os.Stderr, os.Stdout, os.Stdout, os.Stderr)

	helpers.Info.Println("new request to create a channel")

	transactionId := r.Header.Get("x-transactionid")
	userId := r.Header.Get("x-user-uuid")

	if status, response := helpers.ValidateUserFromHeader(transactionId, userId); status {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(response)
		return
	}

	params := mux.Vars(r)
	company_id := params["company_id"]
	channelType := params["type"]
	channelId := params["channel"]

	if status, response, statusCode := helpers.ValidateCompany(transactionId, company_id, userId); status {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	if _, ok := models.Types[channelType]; len(channelType) <= 0 || !ok {
		helpers.Info.Println("abort transaction because of wrong channel type: " + channelType)
		w.WriteHeader(400)
		response := api.BasicResponse{Status: "ERROR", StatusCode: 400, Message: "you have to submit the channel type as url param",
			TransactionId: transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	if !helpers.StringInSlice(channelId, models.Types[channelType]) {
		helpers.Info.Println("abort transaction because channel id is not valid")
		w.WriteHeader(400)
		response := api.BasicResponse{Status: "OK", StatusCode: 400, Message: "channel id does not exist", TransactionId: transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	helpers.Info.Println(transactionId + ": channel is valid and can be created... continue")

	var channel models.Channels
	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &channel)

	if err != nil {
		helpers.Info.Println(transactionId + ": not possible to parse channel in post body")
		helpers.Info.Println(err)
		w.WriteHeader(400)
		response := api.BasicResponse{Status: "ERROR", StatusCode: 400,
			Message: "please submit a valid channel object", TransactionId: transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	helpers.Info.Println("request to create new Channel for " + channelType + "and " + channelId)

	channel.CompanyId = company_id
	channel.Type = channelType
	channel.ChannelId = channelId
	channel.ChannelUuid = xid.New().String()

	fmt.Println(channel)

	helpers.Info.Println(transactionId + ": successfully received channel going to validate to the stripe service")

	stripeClient := clients.StripeClient{Host: "http://localhost:8080/channels/payments/stripe/validate_credentials",
		ApiKey: channel.Key}
	success, err, statusCode := clients.ValidateStripeCredentials(stripeClient)

	if !success || err != nil {
		helpers.Info.Println("error with strip client communication because of status " + strconv.Itoa(statusCode))
		w.WriteHeader(statusCode)
		response := api.BasicResponse{StatusCode: statusCode, Status: "ERROR",
			Message: "invalid response from stripe handler", TransactionId: transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	DbConn.NewRecord(channel)

	if DbConn.Save(&channel).Error != nil {
		helpers.Warning.Println(transactionId + ": error with database connection")
		w.WriteHeader(500)
		response := api.BasicResponse{StatusCode: statusCode, Status: "ERROR",
			Message: "not possible to write into db", TransactionId: transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := api.BasicResponse{Status: "OK", StatusCode: 200,
		Message: "successfully created channel", TransactionId: transactionId}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
	return
}

func FetchAllChannelsFromCompany(w http.ResponseWriter, r *http.Request) {

	helpers.Init(os.Stderr, os.Stdout, os.Stdout, os.Stderr)

	helpers.Info.Println("new request to fetch all channels")

	transactionId := r.Header.Get("x-transactionid")
	userId := r.Header.Get("x-user-uuid")

	if status, response := helpers.ValidateUserFromHeader(transactionId, userId); status {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(response)
		return
	}

	params := mux.Vars(r)
	company_id := params["company_id"]
	channelType := params["type"]

	if status, response, statusCode := helpers.ValidateCompany(transactionId, company_id, userId); status {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	if _, ok := models.Types[channelType]; len(channelType) <= 0 || !ok {
		helpers.Info.Println("abort transaction because of wrong channel type: " + channelType)
		w.WriteHeader(400)
		response := api.BasicResponse{Status: "ERROR", StatusCode: 400, Message: "you have to submit the channel type as url param",
			TransactionId: transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	channels := []models.Channels{}

	DbConn.Where("company_id = ? AND type = ?", company_id, channelType).Find(&channels)

	metaResponse := api.BasicResponse{StatusCode: 200, Status: "OK",
		Message: "", TransactionId: transactionId}
	response := api.MultipleChannels{Meta: metaResponse, Channels: channels}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
	return
}

func DeleteChannelById(w http.ResponseWriter, r *http.Request) {

	helpers.Init(os.Stderr, os.Stdout, os.Stdout, os.Stderr)

	helpers.Info.Println("new request to fetch all channels")

	transactionId := r.Header.Get("x-transactionid")
	userId := r.Header.Get("x-user-uuid")

	if status, response := helpers.ValidateUserFromHeader(transactionId, userId); status {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(response)
		return
	}

	params := mux.Vars(r)
	company_id := params["company_id"]
	channelType := params["type"]
	channelUuid := params["channel_uuid"]

	if status, response, statusCode := helpers.ValidateCompany(transactionId, company_id, userId); status {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	if _, ok := models.Types[channelType]; len(channelType) <= 0 || !ok {
		helpers.Info.Println("abort transaction because of wrong channel type: " + channelType)
		w.WriteHeader(400)
		response := api.BasicResponse{Status: "ERROR", StatusCode: 400, Message: "you have to submit the channel type as url param",
			TransactionId: transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	channel := models.Channels{}

	if err := DbConn.Where("company_id = ? AND type = ? AND channel_uuid = ?", company_id, channelType, channelUuid).Find(&channel).Error; err != nil {
		helpers.Info.Println(transactionId + ": not possible to query for channel")
		helpers.Info.Println(err)
		w.WriteHeader(404)
		response := api.BasicResponse{Status: "ERROR", StatusCode: 404, Message: "channel does not exist", TransactionId: transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := DbConn.Delete(&channel).Error; err != nil {
		helpers.Info.Println(transactionId + ": internal error, not able to delete channel")
		helpers.Info.Println(err)
		w.WriteHeader(500)
		response := api.BasicResponse{StatusCode: 500, Status: "ERROR", Message: "internal server error", TransactionId: transactionId}
		json.NewEncoder(w).Encode(response)
		return
	}

	helpers.Info.Println(transactionId + ": successfully deleted channel")
	response := api.BasicResponse{Status: "OK", StatusCode: 200,
		Message: "successfully deleted channel", TransactionId: transactionId}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
	return
}
