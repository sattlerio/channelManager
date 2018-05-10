package models

import (
	"github.com/jinzhu/gorm"
)

var Types = map[string][]string {
	"payment": []string{"Stripe", "Braintree"},
}

type Channels struct {
	gorm.Model

	Name 		string 		`json:"name" gorm:"size:250; not null"`
	CompanyId 	string 		`gorm:"size:250; not null" json:"company_id"`
	ChannelUuid string		`gorm:"size:250; unique; not null" json:"channel_uuid"`

	Type 		string		`gorm:"size:250; not null" json:"type"`

	MerchantId  string 		`gorm:"size:250" json:"merchant_id"`
	Key 		string 		`gorm:"size:250; not null" json:"key"`
	PrivateKey 	string  	`json:"private_key" gorm:"size:250; not null"`

	ChannelId 	string		`gorm:"size:250; not null" json:"channel_id"`

	Sandbox 	bool		`sql:"default:false" json:"sandbox" gorm:"not null"`
}
