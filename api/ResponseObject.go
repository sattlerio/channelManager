package api

type BasicResponse struct {
	Status			string	`json:"status"`
	StatusCode 		int		`json:"status_code"`
	Message			string	`json:"message"`
	TransactionId 	string	`json:"transaction_id"`
}

type ChannelTypeResponse struct {
	Status 			string 							`json:"status"`
	StatusCode 		int								`json:"status_code"`
	Message 		string							`json:"message"`
	TransactionId 	string							`json:"transaction_id"`
	Data 			map[string][]string				`json:"data"`
}
