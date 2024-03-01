package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file `%v`", err)
	}

	mongoDB := os.Getenv("MONGODB_URL")

	// serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// opts := options.Client().ApplyURI(mongoDB).SetServerAPIOptions(serverAPI)
	// client, err := mongo.Connect(context.TODO(), opts)

	// if err != nil {
	// 	log.Fatal(" Error connecting to database `%v`", err)
	// }
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		log.Fatal(" Error disconnecting from database `%v`", err)
	// 	}
	// }()
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDB))
	if err != nil {
		log.Fatalf(" Error connecting to database `%v`", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf(" Error connecting to database `%v`", err)
	}
	fmt.Println("Connected to MongoDB!")

	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("STOCKET").Collection(collectionName)
	return collection
}
