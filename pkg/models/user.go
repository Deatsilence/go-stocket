package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         *string            `json:"name" validate:"required,min=2,max=30"`
	Surname      *string            `json:"surname" validate:"required,min=2,max=30"`
	Password     *string            `json:"password" validate:"required,min=8"`
	Email        *string            `json:"email" validate:"required,email"`
	Token        *string            `json:"token"`
	UserType     *string            `json:"usertype" validate:"required,eq=ADMIN|eq=USER"`
	RefreshToken *string            `json:"refreshtoken"`
	IsVerified   bool               `json:"isverified"`
	CreatedAt    time.Time          `json:"createdat" `
	UpdatedAt    time.Time          `json:"updatedat"`
	UserID       string             `json:"userid"`
}
