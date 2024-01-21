package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// the json tag is used to represent the JSON key
// pointers are used to represent fields that can be optional or nullable
// the bson id corresponds to the MongoDB client field id

type Table struct{
	ID                 primitive.ObjectID         `bson:"_id"`
	Number_of_guests   *int                    `json:"number_of_guests" validate:"required"`
	Table_number       *int                    `json:"table_number" validate:"required"`
	Created_at         time.Time               `json:"created_at"`
	Updated_at         time.Time               `json:"updated_at"`
	Table_id           string                  `json:"table_id"`
}  