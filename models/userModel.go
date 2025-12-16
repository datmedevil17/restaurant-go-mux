package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	First_name    *string            `json:"first_name" validate:"required,min=2,max=100" bson:"first_name"`
	Last_name     *string            `json:"last_name" validate:"required,min=2,max=100" bson:"last_name"`
	Password      *string            `json:"password" validate:"required,min=6" bson:"password"`
	Email         *string            `json:"email" validate:"email,required" bson:"email"`
	Avatar        *string            `json:"avatar" bson:"avatar"`
	Phone         *string            `json:"phone" validate:"required" bson:"phone"`
	Token         *string            `json:"token" bson:"token"`
	Refresh_Token *string            `json:"refresh_token" bson:"refresh_token"`
	Created_at    time.Time          `json:"created_at" bson:"created_at"`
	Updated_at    time.Time          `json:"updated_at" bson:"updated_at"`
	User_id       string             `json:"user_id" bson:"user_id"`
}
