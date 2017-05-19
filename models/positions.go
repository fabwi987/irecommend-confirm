package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Position struct {
	Idpositions uuid.UUID `json:"idpositions" bson:"idpositions"`
	Iduser      string    `json:"iduser" bson:"iduser"`
	Title       string    `json:"title" bson:"title"`
	Subtitle    string    `json:"subtitle" bson:"subtitle"`
	Text        string    `json:"text" bson:"text"`
	Enddate     time.Time `json:"enddate" bson:"enddate"`
	Reward      string    `json:"reward" bson:"reward"`
	Created     time.Time `json:"created" bson:"created"`
	Lastupdated time.Time `json:"lastupdated" bson:"lastupdated"`
	HREF        string    `json:"href" bson:"href"`
	Meta        string    `json:"meta" bson:"meta"`
}
