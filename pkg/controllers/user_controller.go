package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Deatsilence/go-stocket/database"
	helper "github.com/Deatsilence/go-stocket/helpers"
	"github.com/Deatsilence/go-stocket/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic("Error hashing password")
	}
	return string(hashedPassword)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Password is incorrect")
		check = false
	}
	return check, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking for the email"})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Email already exists"})
			return

		} else {
			password := HashPassword(*user.Password)
			user.Password = &password
		}

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.Name, *user.Surname, *user.UserType, *&user.UserID)
		user.Token = &token
		user.RefreshToken = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if error := c.BindJSON(&user); error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Password or Email is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.Name, *foundUser.Surname, *foundUser.UserType, *&foundUser.UserID)
		helper.UpdateAllTokens(token, refreshToken, foundUser.UserID)
		err = userCollection.FindOne(ctx, bson.M{"userid": foundUser.UserID}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, secondErr := strconv.Atoi(c.Query("page"))
		if secondErr != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}}}
		projectStage := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
		}}}
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while paginating users"})
		}

		var allUsers []bson.M

		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allUsers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userid")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"userid": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
