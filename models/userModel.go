package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// the json tag is used to represent the JSON key
// pointers are used to represent fields that can be optional or nullable
// the bson id corresponds to the MongoDB client field id

type User struct{
	ID                   primitive.ObjectID           `bson:"_id"`
	First_name           *string                 `json:"first_name" validate:"required,min=2,max=100"`
	Last_name            *string                 `json:"last_name"  validate:"required,min=2,max=100"`
	Password             *string                 `json:"password" validate:"required,min=6"`
	Email                *string                 `json:"email"   validate:"required"`
	Avatar               *string                 `json:"avatar"`
	Phone                *string                 `json:"phone"  validate:"required"`
	Token                *string                 `json:"token"`
	Refresh_token        *string                 `json:"refresh_token"`
	Created_at           time.Time               `json:"created_at"`
	Updated_at           time.Time               `json:"updated_at"`
	User_id              string                  `json:"user_id"`
}