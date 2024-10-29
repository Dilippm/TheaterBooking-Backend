package queries
import (
"context"
"time"
	"log"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/dilippm92/bookingapplication/config"
		"github.com/dilippm92/bookingapplication/models/schemas"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	
)
// get Theater collection
func GetTheaterCollection()*mongo.Collection{
	if config.MongoClient == nil{
		log.Fatal("MongoDB client is not initialized")
	}
	return config.MongoClient.Database("movie").Collection("theaters")
}

// Add a Theater
func AddTheater(theater schemas.Theater) (*mongo.InsertOneResult, error) {
	collection := GetTheaterCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// add a new theater
	result, err := collection.InsertOne(ctx, theater)
	if err != nil {
		log.Printf("Failed to insert theater: %v", err)
		return nil, err
	}

	return result, nil
}

// get all theaters for owner
func GetAllOnwertheater(ownerId string)([]schemas.Theater, error){
	collection := GetTheaterCollection()
	var theaters []schemas.Theater
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(ownerId)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	// Find all theaters with the given owner ID
	cursor, err := collection.Find(ctx, bson.M{"ownerId": objectID})
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and decode each document into a Theater struct
	for cursor.Next(ctx) {
		var theater schemas.Theater
		if err := cursor.Decode(&theater); err != nil {
			log.Printf("Failed to decode theater: %v", err)
			return nil, err
		}
		theaters = append(theaters, theater)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	// If no theaters were found, return a custom error
	if len(theaters) == 0 {
		return nil, fmt.Errorf("no theaters found for owner ID %s", ownerId)
	}

	return theaters, nil
}

// update a theater detail

func UpdateTheaterDetail(theater schemas.Theater, id string) (*mongo.UpdateResult, error) {
	collection := GetTheaterCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	// Create a filter to find the document to update
	filter := bson.M{"_id": objectID}

	// Create an update document with the fields to be updated
	update := bson.M{
		"$set": bson.M{
			"theaterName":   theater.TheaterName,
			"place":         theater.Place,
			"state":         theater.State,
			"movie":         theater.Movie,
			"rows":          theater.Rows,
			"columns":       theater.Columns,
			"seats":         theater.Seats,
			"price":         theater.Price,
			"showTimings":   theater.ShowTimings,
		},
	}

	// Perform the update operation
	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update theater: %v", err)
	}

	return updateResult, nil
}

func FindTheaterById(id string) (schemas.Theater, error) {
	collection := GetTheaterCollection()
	var theater schemas.Theater
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
// Convert string ID to ObjectID
objectID, err := primitive.ObjectIDFromHex(id)
if err != nil {
	return schemas.Theater{}, fmt.Errorf("invalid ID format: %v", err)
}
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&theater)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return a custom error message if no document is found
			return schemas.Theater{}, fmt.Errorf("theater with Id %s not found", id)
		}
		// Log and return the error if there is a database error
		log.Printf("Failed to find theater: %v", err)
		return schemas.Theater{}, err
	}

	return theater, nil
}

func GetTheatersByNamePlaceId(name, place, id string) ([]schemas.Theater, error) {
	collection := GetTheaterCollection()
	var theaters []schemas.Theater
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Validate movie ID
	if id == "" {
		return nil, fmt.Errorf("movie ID cannot be empty")
	}

	// Construct the query based on the presence of name and place
	query := bson.M{"movie": bson.M{"$regex": id, "$options": "i"}}

	if name != "" {
		query["theaterName"] = bson.M{"$regex": name, "$options": "i"}
	}

	if place != "" {
		query["place"] = bson.M{"$regex": place, "$options": "i"}
	}

	cursor, err := collection.Find(ctx, query)
	if err != nil {
		// Log and return the error if there is a database error
		log.Printf("Failed to find theaters: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode the results into a slice of Theater structs
	if err = cursor.All(ctx, &theaters); err != nil {
		// Log and return the error if there is an issue with decoding
		log.Printf("Failed to decode theaters: %v", err)
		return nil, err
	}

	if len(theaters) == 0 {
		// Return a custom error message if no documents are found
		return theaters, nil
	}

	return theaters, nil
}


func FindTheaterByName(name string) (schemas.Theater, error) {
	collection := GetTheaterCollection()
	var theater schemas.Theater
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"theaterName": name}).Decode(&theater)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return a custom error message if no document is found
			return schemas.Theater{}, fmt.Errorf("theater with Id %s not found", name)
		}
		// Log and return the error if there is a database error
		log.Printf("Failed to find theater: %v", err)
		return schemas.Theater{}, err
	}

	return theater, nil
}


//find theater by movie id

func FindtheaterBymovieId (id string)(schemas.Theater,error){
	collection := GetTheaterCollection()
	var theater schemas.Theater
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"movie": id}).Decode(&theater)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return a custom error message if no document is found
			return schemas.Theater{},nil
		}
		// Log and return the error if there is a database error
		log.Printf("Failed to find theater: %v", err)
		return schemas.Theater{}, err
	}

	return theater, nil
}