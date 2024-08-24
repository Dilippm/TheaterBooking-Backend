package queries

import (
	"context"
	"log"
	"time"
"fmt"
	"github.com/dilippm92/bookingapplication/config"
	"github.com/dilippm92/bookingapplication/models/schemas"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// get user collection
func GetUserCollection()*mongo.Collection{
	if config.MongoClient == nil{
		log.Fatal("MongoDB client is not initialized")
	}
	return config.MongoClient.Database("movie").Collection("users")
}
// find or create a user
func FindOrCreateUser(user schemas.User) (*mongo.InsertOneResult, error) {
	collection := GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user already exists
	var existingUser schemas.User
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err != nil && err != mongo.ErrNoDocuments {
		// If there's an error other than "no documents found", log it and return
		log.Printf("Failed to find user: %v", err)
		return nil, err
	}

	// If user already exists, return nil and an appropriate error or message
	if err == nil {
		log.Printf("User with email %s already exists", user.Email)
		return nil, fmt.Errorf("user with email %s already exists", user.Email)
	}

	// If user does not exist, create a new user
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return nil, err
	}

	return result, nil
}


// getuser by email
func FindUserByEmail(email string) (schemas.User, error) {
	collection := GetUserCollection()
	var user schemas.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return a custom error message if no document is found
			return schemas.User{}, fmt.Errorf("user with email %s not found", email)
		}
		// Log and return the error if there is a database error
		log.Printf("Failed to find user: %v", err)
		return schemas.User{}, err
	}

	return user, nil
}