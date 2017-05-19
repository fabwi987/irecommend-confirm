package models

import (
	"time"
)

type User struct {
	IdUser      string    `json:"iduser" bson:"iduser"`
	UserType    int       `json:"usertype" bson:"usertype"`
	Name        string    `json:"name" bson:"name"`
	Telephone   string    `json:"telephone" bson:"telephone"`
	Mail        string    `json:"mail" bson:"mail"`
	Picture     string    `json:"picture" bson:"picture"`
	Headline    string    `json:"headline" bson:"headline"`
	ProfileURL  string    `json:"profileURL" bson:"profileURL"`
	Created     time.Time `json:"created" bson:"created"`
	Lastupdated time.Time `json:"lastupdated" bson:"lastupdated"`
	HREF        string    `json:"href" bson:"href"`
	Meta        string    `json:"meta" bson:"meta"`
}
