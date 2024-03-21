package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Deatsilence/go-stocket/database"
	"github.com/Deatsilence/go-stocket/pkg/models"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email    string
	Name     string
	Surname  string
	UserType string
	UserId   string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var blacklistCollection *mongo.Collection = database.OpenCollection(database.Client, "blacklist")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, name string, surname string, userType string, userID string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:    email,
		Name:     name,
		Surname:  surname,
		UserType: userType,
		UserId:   userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(30)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(2)).Unix(),
		},
	}

	token, tokenErr := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, refreshTokenErr := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if tokenErr != nil {
		log.Panic(tokenErr)
		return
	}
	if refreshTokenErr != nil {
		log.Panic(refreshTokenErr)
		return
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (cliams *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = fmt.Sprintf("the token is invalid: %v", err.Error())
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token has expired"
		return
	}
	return claims, msg
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refreshtoken", Value: signedRefreshToken})

	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedat", Value: updatedAt})

	upsert := true
	filter := bson.M{"userid": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{{Key: "$set", Value: updateObj}},
		&opt,
	)

	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
}

func IsTokenBlacklisted(token string) bool {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	count, err := blacklistCollection.CountDocuments(ctx, bson.M{"token": token})
	if err != nil {
		log.Printf("error occured while checking if token is blacklisted: %v", err)
	}
	return count > 0
}

func BlacklistToken(token string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	blacklistedToken := &models.BlacklistedToken{
		Token:         token,
		BlacklistedAt: time.Now(),
	}

	_, err := blacklistCollection.InsertOne(ctx, blacklistedToken)

	if err != nil {
		log.Printf("error occured while blacklisting token: %v", err)
	}
}
