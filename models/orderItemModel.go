package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	Quantity      *string            `json:"quantity" validate:"required,eq=S|eq=M|eq=L" bson:"quantity"`
	Unit_price    *float64           `json:"unit_price" validate:"required" bson:"unit_price"`
	Created_at    time.Time          `json:"created_at" bson:"created_at"`
	Updated_at    time.Time          `json:"updated_at" bson:"updated_at"`
	Food_id       *string            `json:"food_id" validate:"required" bson:"food_id"`
	Order_item_id string             `json:"order_item_id" bson:"order_item_id"`
	Order_id      string             `json:"order_id" validate:"required" bson:"order_id"`
}
