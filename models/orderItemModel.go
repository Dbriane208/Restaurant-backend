package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// the json tag is used to represent the JSON key
// pointers are used to represent fields that can be optional or nullable
// the bson id corresponds to the MongoDB client field id

type OrderItem struct{
	ID                 primitive.ObjectID     `bson:"_id"`
	Quantity          *string               `json:"quantity"`
	Unit_price         *float64              `json:"unit_price"`
	Created_at          time.Time            `json:"created_at"`
	Updated_at          time.Time            `json:"updated_at"`
	Food_id            *string               `json:"food_id" validated:"required"`
	Order_item_id       string               `json:"order_item_id"`
	Order_id            string               `json:"order_id" validate:"required"`
}