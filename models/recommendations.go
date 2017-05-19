package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Recommendation struct {
	Idrecommendations uuid.UUID `json:"idrecommendations" bson:"idrecommendations"`
	User              *User     `json:"user" bson:"user"`
	Referral          *Referral `json:"referral" bson:"referral"`
	Position          *Position `json:"position" bson:"position"`
	Text              string    `json:"text" bson:"text"`
	Confirmed         bool      `json:"confirmed" bson:"confirmed"`
	Created           time.Time `json:"created" bson:"created"`
	Lastupdated       time.Time `json:"lastupdated" bson:"lastupdated"`
	HREF              string    `json:"href" bson:"href"`
	Meta              string    `json:"meta" bson:"meta"`
}
