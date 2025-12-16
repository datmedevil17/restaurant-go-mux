package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
	ID               primitive.ObjectID `bson:"_id" json:"_id"`
	Number_of_guests *int               `json:"number_of_guests" validate:"required" bson:"number_of_guests"`
	Table_number     *int               `json:"table_number" validate:"required" bson:"table_number"`
	Created_at       time.Time          `json:"created_at" bson:"created_at"`
	Updated_at       time.Time          `json:"updated_at" bson:"updated_at"`
	Table_id         string             `json:"table_id" bson:"table_id"`
}
