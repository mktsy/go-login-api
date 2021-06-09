package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	Email     string        `json:"email" bson:"email"`
	Password  string        `json:"password" bson:"password"`
	Role      string        `json:"role" bson:"role"`
	CreatedAt time.Time     `json:"createAt" bson:"createAt"`
	UpdatedAt time.Time     `json:"updateAt" bson:"updateAt"`
}
