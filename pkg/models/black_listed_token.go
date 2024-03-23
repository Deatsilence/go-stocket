package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BlacklistedToken represents the structure of the token document in the database.
type BlacklistedToken struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Token         string             `json:"token"`
	BlacklistedAt time.Time          `json:"blacklistedAt"`
}
