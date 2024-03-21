package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id"`
	Barcode     string             `json:"barcode" validate:"required"`
	Name        *string            `json:"name" validate:"required,min=2,max=50"`
	Description *string            `json:"description" validate:"required,min=2,max=100"`
	Price       float64            `json:"price" validate:"required"`
	Stock       uint               `json:"stock" validate:"required"`
	CreatedAt   time.Time          `json:"createdat"`
	UpdatedAt   time.Time          `json:"updatedat"`
	ProductID   string             `json:"productid"`
}
