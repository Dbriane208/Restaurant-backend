package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// the json tag is used to represent the JSON key
// pointers are used to represent fields that can be optional or nullable
// the bson id corresponds to the MongoDB client field id

type Note struct{
	ID             primitive.ObjectID            `bson:"_id"`
	Text           string                  `json:"text"`
	Title          string                  `json:"title"`
	Created_at     time.Time               `json:"created_at"`
	Updated_at     time.Time               `json:"updated_at"`
	Note_id        string                  `json:"note_id"`
}