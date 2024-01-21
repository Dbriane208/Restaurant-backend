package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// pointers are used in fields that can be nullable or optional
// json tag is used to represent the JSON key
// the bson id is used to correspond to the MongoDB client field id

type Menu struct{
	ID             primitive.ObjectID          `bson:"_id"`
	Name           string                  `json:"name" validate:"required"`
	Category       string                  `json:"category" validate:"required"`
	Start_date    *time.Time               `json:"start_date"`
	End_date      *time.Time               `json:"end_date"`
	Created_at     time.Time               `json:"created_at"`
	Updated_at     time.Time               `json:"updated-at"`
	Menu_id        string                  `json:"menu_id"`
}