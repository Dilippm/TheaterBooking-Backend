package queries

import (
	"context"
	"log"
	"time"

	"github.com/dilippm92/bookingapplication/config"
	"github.com/dilippm92/bookingapplication/models/schemas"
	"go.mongodb.org/mongo-driver/mongo"
)

// get user collection
func GetUserCollection()*mongo.Collection{
	if config.MongoClient == nil{
		log.Fatal("MongoDB client is not initialized")
	}
	return config.MongoClient.Database("movie").Collection("users")
}
func CreateUser(user schemas.User) (*mongo.InsertOneResult,error) {
	collection:=GetUserCollection()
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	// Perform the InsertOne operation with the context
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return nil, err
	}
	return result, nil
}