package config
import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient holds the MongoDB client instance
var MongoClient *mongo.Client

func ConnectMongoDB() {
	var err error
	mongoURI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)
	MongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)

	}
	// Check the connection
	err = MongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB")
}