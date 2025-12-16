package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Text       string             `json:"text" bson:"text"`
	Title      string             `json:"title" bson:"title"`
	Created_at time.Time          `json:"created_at" bson:"created_at"`
	Updated_at time.Time          `json:"updated_at" bson:"updated_at"`
	Note_id    string             `json:"note_id" bson:"note_id"`
}
