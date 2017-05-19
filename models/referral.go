package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Referral struct {
	Idreferrals    uuid.UUID `json:"idreferral" bson:"idreferral"`
	ReferralUserID string    `json:"referraluserid" bson:"referraluserid"`
	Name           string    `json:"name" bson:"name"`
	Telephone      string    `json:"telephone" bson:"telephone"`
	Mail           string    `json:"mail" bson:"mail"`
	Picture        string    `json:"picture" bson:"picture"`
	Headline       string    `json:"headline" bson:"headline"`
	ProfileURL     string    `json:"profileURL" bson:"profileURL"`
	Created        time.Time `json:"created" bson:"created"`
	Lastupdated    time.Time `json:"lastupdated" bson:"lastupdated"`
	HREF           string    `json:"href" bson:"href"`
	Meta           string    `json:"meta" bson:"meta"`
}
