package helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Deatsilence/go-stocket/pkg/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("usertype")
	err = nil

	if userType != role {
		err = errors.New("unauthorized to access this route")
		return err
	}
	return err
}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("usertype")
	uid := c.GetString("userid")

	err = nil

	if userType == "USER" && uid != userId {
		err = errors.New("unauthorized to access this route")
		return err
	}
	err = CheckUserType(c, userType)
	return err
}

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
		msg = fmt.Sprintf("passwords do not match: %v", err)
		check = false
	}
	return check, msg
}

func DeleteUnverified(email string) (bool, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User
	userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if !user.IsVerified {
		log.Printf("Deleting user with email: %v", email)
		_, err := userCollection.DeleteOne(ctx, bson.M{"email": email})
		if err != nil {
			return user.IsVerified, err
		}
	}

	return user.IsVerified, nil
}
