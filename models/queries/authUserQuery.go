package queries

import (
	"context"
	"log"
	"time"
"fmt"
"strconv"
	"github.com/dilippm92/bookingapplication/config"
	"github.com/dilippm92/bookingapplication/models/schemas"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserResponse struct {
	ID       string `json:"id" bson:"_id"`       // Adjust the field tag as needed
	Username string `json:"username" bson:"username"`
	Role     string `json:"role" bson:"role"`
	UserImage string `json:"userimage" bson:"userimage"`
}
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

func UpdateUser(user schemas.User, id string) (*mongo.UpdateResult, error) {
	collection := GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	// Create a filter to find the document to update (by userID)
	filter := bson.M{"_id": objectID}

	// Create an update document with the fields to be updated
	update := bson.M{
		"$set": bson.M{
			"username":   user.Username,
			"email":      user.Email,
			"password":   user.Password,
			"userimage":  user.UserImage,
			"role":user.Role,
		},
	}

	// Perform the update operation
	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return updateResult, nil
}

// getuser by email
func FindUserById(id string) (schemas.User, error) {
	collection := GetUserCollection()
	var user schemas.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
// Convert string ID to ObjectID
objectID, err := primitive.ObjectIDFromHex(id)
if err != nil {
	return schemas.User{}, fmt.Errorf("invalid ID format: %v", err)
}
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return a custom error message if no document is found
			return schemas.User{}, fmt.Errorf("user with Id %s not found", id)
		}
		// Log and return the error if there is a database error
		log.Printf("Failed to find user: %v", err)
		return schemas.User{}, err
	}

	return user, nil
}

// get owner ids by role
func GetOwnerIds()([]schemas.User,error){
	collection := GetUserCollection()
	var users []schemas.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define the query to filter users by role "owner"
	filter := bson.M{"role": "owner"}

	// Define the projection to include only the _id and username fields
	projection := bson.M{"_id": 1, "username": 1}

	// Perform the query
	cursor, err := collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each document into the user slice
	for cursor.Next(ctx) {
		var user schemas.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}


func UpdateWallet(price string)(*mongo.UpdateResult, error){
	collection := GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert price to float64
	priceFloat, err := strconv.ParseFloat(price, 64)
	
	if err != nil {
		return nil, fmt.Errorf("invalid price format: %v", err)
	}

	// Find user with role "admin"
	filter := bson.M{"role": "admin"}
	
	// Find the current wallet value
	var user struct {
		Wallet float64 `bson:"wallet"`
	}
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("failed to find user: %v", err)
	}
	
	// Calculate new wallet value
	newWalletValue := user.Wallet + (priceFloat * 0.2)

	// Update the wallet field with the new value
	update := bson.M{
		"$set": bson.M{"wallet": newWalletValue},
	}

	// Perform the update
	result, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return nil, fmt.Errorf("update failed: %v", err)
	}

	return result, nil
}


func UpdateWalletByUserId(id string, price string) (*mongo.UpdateResult, error) {
	collection := GetUserCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert price to float64
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid price format: %v", err)
	}


	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID format: %v", err)
	}

	// Find the current wallet value
	var user struct {
		Wallet float64 `bson:"wallet"`
	}
	filter := bson.M{"_id": objectID}
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("failed to find user: %v", err)
	}
	
	// Calculate new wallet value
	newWalletValue := user.Wallet + (priceFloat * 0.8)
	
	// Update the wallet field with the new value
	update := bson.M{
		"$set": bson.M{"wallet": newWalletValue},
	}

	// Perform the update
	result, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return nil, fmt.Errorf("update failed: %v", err)
	}

	return result, nil
}

// GetChatUsers retrieves users who are either "owner" or "admin", excluding the user with the specified ID
func GetChatUsers(id string) ([]schemas.User, error) {
	
	collection := GetUserCollection()
	var users []schemas.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	// Define the query to filter users by role "owner" or "admin", excluding the specified user ID
	filter := bson.M{
		"$or": []bson.M{
			{"role": "owner"},
			{"role": "admin"},
		},
		"_id": bson.M{"$ne": objectID}, // Exclude the user with the specified ObjectID
	}

	// Define the projection to include only the _id, username, role, and userimage fields
	projection := bson.M{"_id": 1, "username": 1, "role": 1, "userimage": 1}

	// Perform the query
	cursor, err := collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err // Return the error if the query fails
	}
	defer cursor.Close(ctx)

	// Decode the cursor results into the users slice
	for cursor.Next(ctx) {
		var user schemas.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err // Return the error if decoding fails
		}
		users = append(users, user)
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		return nil, err // Return the error if there was an issue with the cursor
	}

	return users, nil // Return the list of users
}
