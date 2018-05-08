package api

import "channelManager/models"

type MultipleChannels struct {
	Meta		BasicResponse		`json:"info"`
	Channels 	[]models.Channels	`json:"body; omitempty"`
}
