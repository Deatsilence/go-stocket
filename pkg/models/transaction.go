package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID            primitive.ObjectID `bson:"_id"`
	UserID        string             `json:"userid"`        /// The user who made the transaction
	ProductID     string             `json:"productid"`     /// The product that the transaction is made
	ProcessType   string             `json:"processtype"`   /// The type of the transaction (add, remove, update, delete)
	Amount        uint               `json:"amount"`        /// The amount of the product
	ProcessTime   time.Time          `json:"processtime"`   /// The time of the transaction
	TransactionID string             `json:"transactionid"` /// The id of the transaction
}
