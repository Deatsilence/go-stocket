package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
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
		defer cancel()

		var product models.Product

		if err := c.BindJSON(&product); err != nil {
			log.Println("33 ERROR: ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validateProduct.Struct(product)
		if validationErr != nil {
			log.Println("40 ERROR: ", validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := productCollection.CountDocuments(ctx, bson.M{"barcode": product.Barcode})

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while checking for product"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Product already exists"})
			return
		}

		product.ID = primitive.NewObjectID()
		product.ProductID = product.ID.Hex()
		product.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		product.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		log.Println("Product: ", product)
		resultInsertionNumber, insertErr := productCollection.InsertOne(ctx, product)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while inserting product"})
			return
		}
		userID := c.GetString("userid")
		helper.CreateTransactionForProduct(userID, product.ProductID, types.Add, product.Stock)

		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func DeleteAProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		productID := c.Param("productid")

		var product models.Product

		err := productCollection.FindOneAndDelete(ctx, bson.M{"productid": productID}).Decode(&product)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while deleting product"})
			return
		}
		userID := c.GetString("userid")
		helper.CreateTransactionForProduct(userID, product.ProductID, types.Delete, product.Stock)

		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	}
}

func GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, recordPageErr := strconv.Atoi(c.Query("recordPerPage"))

		if recordPageErr != nil && recordPerPage < 1 {
			recordPerPage = 10
		}

		page, pageErr := strconv.Atoi(c.Query("page"))

		if pageErr != nil && page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}}}
		projectStage := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "product_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
		}}}
		result, err := productCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while paginating products"})
		}

		var allProducts []bson.M

		if err = result.All(ctx, &allProducts); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allProducts[0])
	}
}

func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		productID := c.Param("productid")

		var product models.Product

		err := productCollection.FindOne(ctx, bson.M{"productid": productID}).Decode(&product)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusOK, product)
	}
}

func UpdateAProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		productID := c.Param("productid")

		var product models.Product

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validateProduct.Struct(product)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}

		product.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		update := bson.M{
			"$set": bson.M{
				"name":        product.Name,
				"barcode":     product.Barcode,
				"description": product.Description,
				"category":    product.Category,
				"stock":       product.Stock,
				"price":       product.Price,
				"updatedat":   product.UpdatedAt,
			},
		}

		_, err := productCollection.UpdateOne(ctx, bson.M{"productid": productID}, update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while updating product"})
			return
		}
		userID := c.GetString("userid")
		helper.CreateTransactionForProduct(userID, product.ProductID, types.Update, product.Stock)

		c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
	}
}

func UpdateSomePropertiesOfProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		productID := c.Param("productid")

		var product models.Product

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		update := helper.UpdateFilter(product)

		updated := bson.M{"$set": update}

		_, err := productCollection.UpdateOne(ctx, bson.M{"productid": productID}, updated)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while updating product"})
			return
		}
		userID := c.GetString("userid")
		helper.CreateTransactionForProduct(userID, productID, types.Update, product.Stock)

		c.JSON(http.StatusOK, gin.H{"message": "Product updated partially successfully"})
	}
}
