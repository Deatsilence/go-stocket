package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PasswordReset struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Email     *string            `json:"email" validate:"required,email"`
	Code      string             `json:"code"`
	CreatedAt time.Time          `json:"createdat"`
	ExpiresAt time.Time          `json:"expiresat"`
}
