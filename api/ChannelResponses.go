package api

import "channelManager.sattler.io/models"

type MultipleChannels struct {
	Meta		BasicResponse		`json:"info"`
	Channels 	[]models.Channels	`json:"body; omitempty"`
}

type SingleChannel struct {
	Meta 		BasicResponse 		`json:"info"`
	Channel 	models.Channels		`json:"body; omitempty"`
}