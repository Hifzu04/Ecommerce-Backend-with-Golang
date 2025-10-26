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

func Connectdb() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading env file: %s", err)
	}
	MONGO_URI := os.Getenv("MONGO_URI")

	//CONNECT TO DB
	clientOptions := options.Client().ApplyURI(MONGO_URI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection.
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatalf("Error pinging mongodb %v", err)
	} else {
		fmt.Println("Connected to mongodb!!")
	}

	return client
}

var Client *mongo.Client = Connectdb()

func UserData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Ecommerce").Collection(CollectionName)
	return collection
}

func ProductData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var productCollection *mongo.Collection = client.Database("Ecommerce").Collection(CollectionName)
	return productCollection
}
