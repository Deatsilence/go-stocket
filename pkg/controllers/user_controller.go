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
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validateUser = validator.New()

func VerifyEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var isVerified bool = true

		var requestBody struct {
			Email *string `json:"email" validate:"required,email"`
			Code  string  `json:"code"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the code is valid and update the user's verified status
		valid, err := helper.ValidateResetCode(*requestBody.Email, requestBody.Code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			isVerified = false
		}
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "expired code"})
			isVerified = false
		}
		if !isVerified {
			var user models.User

			err = userCollection.FindOneAndDelete(ctx, bson.M{"email": *requestBody.Email}).Decode(&user)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			return
		}

		update := bson.M{
			"$set": bson.M{
				"isverified": true,
			},
		}

		_, updateErr := userCollection.UpdateOne(ctx, bson.M{"email": *requestBody.Email}, update)
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": updateErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Email successfully verified"})
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validateUser.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking for the email"})
		}

		if count > 0 {
			isVerified, isVerifyErr := helper.DeleteUnverified(*user.Email)
			log.Printf("email: %v", *user.Email)
			if isVerifyErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": isVerifyErr.Error()})
				return
			}
			if isVerified {
				c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
				return
			}
		}

		err = helper.GenerateResetCode(*user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verify code"})
			return
		}

		password := helper.HashPassword(*user.Password)
		user.Password = &password

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.Name, *user.Surname, *user.UserType, user.UserID)
		user.Token = &token
		user.RefreshToken = &refreshToken
		user.IsVerified = false

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := "User not created"
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
		defer cancel()

		var user models.User
		var foundUser models.User

		if error := c.BindJSON(&user); error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Password or Email is incorrect"})
			return
		}

		passwordIsValid, msg := helper.VerifyPassword(*user.Password, *foundUser.Password)

		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if !foundUser.IsVerified {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not verified"})
			return
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.Name, *foundUser.Surname, *foundUser.UserType, foundUser.UserID)
		helper.UpdateAllTokens(token, refreshToken, foundUser.UserID)
		err = userCollection.FindOne(ctx, bson.M{"userid": foundUser.UserID}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		recordPerPage, recordPageErr := strconv.Atoi(c.Query("recordPerPage"))

		if recordPageErr != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, pageErr := strconv.Atoi(c.Query("page"))

		if pageErr != nil || page < 1 {
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
			{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
		}}}
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})

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
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("userid")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"userid": userId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			return
		}

		if helper.IsTokenBlacklisted(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is already blacklisted"})
			return
		}

		helper.BlacklistToken(token)
		c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	}
}

func RequestPasswordReset() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var requestBody struct {
			Email *string `json:"email" validate:"required,email"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		count, countErr := userCollection.CountDocuments(ctx, bson.M{"email": *requestBody.Email})

		if countErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking for the email"})
			return
		}

		if count == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email does not exist"})
			return
		}

		// Generate and send a reset code
		err := helper.GenerateResetCode(*requestBody.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset code"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Reset code sent to email"})
	}
}

func ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody struct {
			Email       *string `json:"email" validate:"required,email"`
			Code        string  `json:"code" validate:"required,min=6,max=6"`
			NewPassword string  `json:"newPassword" validate:"required,min=6"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Reset the password
		err := helper.ResetUserPassword(*requestBody.Email, requestBody.Code, requestBody.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password successfully reset"})
	}
}
