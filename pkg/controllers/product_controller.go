package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Deatsilence/go-stocket/database"
	helper "github.com/Deatsilence/go-stocket/helpers"
	"github.com/Deatsilence/go-stocket/pkg/models"
	"github.com/Deatsilence/go-stocket/types"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var productCollection *mongo.Collection = database.OpenCollection(database.Client, "product")
var validateProduct = validator.New()

func AddAProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var product models.Product

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validateProduct.Struct(product)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}

		count, err := productCollection.CountDocuments(ctx, bson.M{"barcode": product.Barcode})
		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while checking for product"})
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Product already exists"})
			return
		}

		product.ID = primitive.NewObjectID()
		product.ProductID = product.ID.Hex()
		product.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		product.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		resultInsertionNumber, insertErr := productCollection.InsertOne(ctx, product)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while inserting product"})
			return
		}
		userID := c.GetString("userid")
		helper.CreateTransactionForProduct(*&userID, *&product.ProductID, types.Add, product.Stock)

		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}
