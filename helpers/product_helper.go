package helpers

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/Deatsilence/go-stocket/database"
	"github.com/Deatsilence/go-stocket/pkg/models"
	"github.com/Deatsilence/go-stocket/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var transactionCollection *mongo.Collection = database.OpenCollection(database.Client, "transaction")

func CreateTransactionForProduct(userID string, productID string, processtype types.ProcessTypes, amount int) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	var transaction models.Transaction

	processTime, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Printf("Error while parsing time: %v", err)
	}

	transaction.ID = primitive.NewObjectID()
	transaction.TransactionID = transaction.ID.Hex()
	transaction.UserID = userID
	transaction.ProductID = productID
	transaction.ProcessType = strconv.Itoa(int(processtype))
	transaction.Amount = amount
	transaction.ProcessTime = processTime

	_, insertErr := transactionCollection.InsertOne(ctx, transaction)

	if insertErr != nil {
		log.Printf("Error while inserting transaction: %v", insertErr)
	}
	defer cancel()
	return insertErr
}
