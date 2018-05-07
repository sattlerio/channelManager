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
		TransactionId: transactionId, Data:models.Types}
	json.NewEncoder(w).Encode(response)
	return
}
