package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         *string            `json:"name" validate:"required, min=2, max=30"`
	Surname      *string            `json:"surname" validate:"required, min=2, max=30"`
	Password     *string            `json:"password" validate:"required, min=6"`
	Email        *string            `json:"email" validate:"required, email"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refreshToken"`
	CreatedAt    time.Time          `json:"createdAt" `
	UpdatedAt    time.Time          `json:"updatedAt"`
	UserID       string             `json:"userId"`
}
