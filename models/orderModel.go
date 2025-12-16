package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Order_Date time.Time          `json:"order_date" validate:"required" bson:"order_date"`
	Created_at time.Time          `json:"created_at" bson:"created_at"`
	Updated_at time.Time          `json:"updated_at" bson:"updated_at"`
	Order_id   string             `json:"order_id" bson:"order_id"`
	Table_id   *string            `json:"table_id" validate:"required" bson:"table_id"`
}
