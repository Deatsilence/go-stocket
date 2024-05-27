package helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"time"

	"github.com/Deatsilence/go-stocket/database"
	"github.com/Deatsilence/go-stocket/pkg/models" // replace with your actual package path

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var passwordResetCollection *mongo.Collection = database.OpenCollection(database.Client, "passwordreset")

var FROMMAIL string = os.Getenv("FROMMAIL")
var FROMMAILPASSWORD string = os.Getenv("FROMMAILPASSWORD")

// GenerateResetCode creates a reset code and stores it in the database
func GenerateResetCode(email string) error {
	source := rand.NewSource(time.Now().UnixNano())
	localRNG := rand.New(source)

	code := fmt.Sprintf("%06d", localRNG.Intn(1000000))
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	passwordReset := &models.PasswordReset{
		Email:     &email,
		Code:      code,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Minute * 1), // Code expires in 1 minute
	}

	_, err := passwordResetCollection.InsertOne(ctx, passwordReset)
	if err != nil {
		return err
	}

	// Send the code to the user's email
	SendEmail(email, code)

	return nil
}

// VerifyResetCode checks if the reset code is valid and not expired.
func ValidateResetCode(email string, code string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var passwordReset models.PasswordReset
	err := passwordResetCollection.FindOne(ctx, bson.M{"email": email, "code": code}).Decode(&passwordReset)
	if err != nil {
		return false, errors.New("invalid code")
	}
	return passwordReset.ExpiresAt.After(time.Now()), nil
}

func ResetUserPassword(email string, code string, newPassword string) error {
	valid, err := ValidateResetCode(email, code)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("expired code")
	}

	// Hash the new password before storing it
	hashedPassword := HashPassword(newPassword)

	// Implement UpdateUserPassword to update the user's password in the user collection
	err = UpdateUserPassword(email, hashedPassword)
	if err != nil {
		return err
	}

	// Optionally, delete the password reset document
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	_, err = passwordResetCollection.DeleteOne(ctx, bson.M{"code": code})
	return err
}

func UpdateUserPassword(email string, hashedPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"password": hashedPassword}}

	_, err := userCollection.UpdateOne(ctx, filter, update)
	return err
}

func SendEmail(toEmail string, code string) error {
	auth := smtp.PlainAuth("", FROMMAIL, FROMMAILPASSWORD, "smtp.gmail.com")

	msg := "Subject: Password Reset Code\n\nHere is your password reset code: " + code

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		FROMMAIL,
		[]string{toEmail},
		[]byte(msg),
	)

	if err != nil {
		log.Printf("Error while sending email: %v", err)
		return err
	}
	return nil
}
