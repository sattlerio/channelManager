package models

import (
	"github.com/jinzhu/gorm"
)

type ChannelId 	 string


var Types = map[string][]string {
	"payment": []string{"Stripe"},
}

type Channels struct {
	gorm.Model

	Name 		string 		`json:"name" gorm:"size:250; not null"`
	CompanyId 	string 		`gorm:"size:250; not null" json:"company_id"`
	ChannelUuid string		`gorm:"size:250; unique; not null" json:"channel_uuid"`

	Type 		string		`gorm:"size:250; not null" json:"type"`

	Key 		string 		`gorm:"size:250; not null" json:"key"`
	PrivateKey 	ChannelId	`json:"private_key" gorm:"size:250; not null"`

	ChannelId 	string		`gorm:"size:250; not null" json:"channel_id"`

	Sandbox 	bool		`sql:"default:false" json:"sandbox" gorm:"not null"`
}
