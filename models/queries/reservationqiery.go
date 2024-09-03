package queries
import (
	"context"
	"time"
		"log"
		"fmt"
	"github.com/dilippm92/bookingapplication/models/schemas"
		"github.com/dilippm92/bookingapplication/config"
		"go.mongodb.org/mongo-driver/mongo"
		 "go.mongodb.org/mongo-driver/mongo/options"
			"go.mongodb.org/mongo-driver/bson"
		
	)
// get movie collection
func GetReservationCollection()*mongo.Collection{
	if config.MongoClient == nil{
		log.Fatal("MongoDB client is not initialized")
	}
	return config.MongoClient.Database("movie").Collection("reservations")
}
// Add a reservation
func AddReservation(reservation schemas.Reservation) (*mongo.InsertOneResult, error) {
	collection := GetReservationCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
// Set creation time
reservation.CreatedAt = time.Now()
	// add a new theater
	result, err := collection.InsertOne(ctx, reservation)
	if err != nil {
		log.Printf("Failed to insert reservation: %v", err)
		return nil, err
	}

	return result, nil
}


func GetReservationsByTimeAndDate(dateStr, timeStr string) ([]schemas.Reservation, error) {
	collection := GetReservationCollection()
	var reservations []schemas.Reservation

	// Parse the timeStr into a time.Time object
	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %v", err)
	}

	// Create a filter for the date and time
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a query to find documents matching both date and time
	filter := bson.M{
		"date": bson.M{
			"$eq": dateStr,
		},
		"time": bson.M{
			"$eq": parsedTime,
		},
	}

	// Use Find to get all matching documents
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to find reservations: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode the results into the reservations slice
	for cursor.Next(ctx) {
		var reservation schemas.Reservation
		if err := cursor.Decode(&reservation); err != nil {
			log.Printf("Failed to decode reservation: %v", err)
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	// Check if any error occurred during iteration
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor iteration error: %v", err)
		return nil, err
	}

	return reservations, nil
}

func CreateTTLIndex() error {
    collection := GetReservationCollection()

    // Define the TTL index options
    indexModel := mongo.IndexModel{
        Keys: bson.D{{Key: "createdAt", Value: 1}}, // Field to index
        Options: options.Index().SetExpireAfterSeconds(900), // 900 seconds = 15 minutes
    }

    // Create the index
    _, err := collection.Indexes().CreateOne(context.Background(), indexModel)
    return err
}