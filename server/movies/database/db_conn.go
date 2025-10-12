package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client = DBInstance()

func DBInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: Trouble finding .env file")
	}

	MongoDb := os.Getenv("MONGODB_URI")
	if MongoDb == "" {
		log.Fatal("Error: MONGODB_URI not found")
	}

	fmt.Println("MongoDB URI: ", MongoDb)
	clientOptions := options.Client().ApplyURI(MongoDb)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil
	}

	return client
}

func OpenCollection(collectionName string) *mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: Trouble finding .env file")
	}

	dbName := os.Getenv("DATABASE_NAME")
	fmt.Println("Database name: ", dbName)
	collection := Client.Database(dbName).Collection(collectionName)
	if collection == nil {
		return nil
	}

	return collection
}
