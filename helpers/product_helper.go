package helpers

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/Deatsilence/go-stocket/database"
	"github.com/Deatsilence/go-stocket/pkg/models"
	"github.com/Deatsilence/go-stocket/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var transactionCollection *mongo.Collection = database.OpenCollection(database.Client, "transaction")

func CreateTransactionForProduct(userID string, productID string, processtype types.ProcessTypes, amount uint) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

	return insertErr
}

func UpdateFilter(product models.Product) bson.M {
	update := bson.M{}
	if product.Name != nil {
		update["name"] = product.Name
		log.Println("Name: ", product.Name)
	}
	if product.Barcode != "" {
		update["barcode"] = product.Barcode
		log.Println("Barcode: ", product.Barcode)
	}
	if product.Description != nil {
		update["description"] = product.Description
		log.Println("Description: ", product.Description)
	}

	if product.Category >= 0 {
		update["category"] = product.Category
		log.Println("Category: ", product.Category)
	}

	update["stock"] = product.Stock
	log.Println("Stock: ", product.Stock)

	if product.Price >= 0.0 {
		update["price"] = product.Price
		log.Println("Price: ", product.Price)
	}
	update["updatedat"], _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	return update
}
