package queries
import (
	"context"
	"time"
		"log"
		"fmt"
	
	"github.com/dilippm92/bookingapplication/models/schemas"
		"github.com/dilippm92/bookingapplication/config"
		"go.mongodb.org/mongo-driver/mongo"
		"go.mongodb.org/mongo-driver/bson"
		
		
	)
// get movie collection
func GetBookingCollection()*mongo.Collection{
	if config.MongoClient == nil{
		log.Fatal("MongoDB client is not initialized")
	}
	return config.MongoClient.Database("movie").Collection("bookings")
}
// Add a reservation
func AddBooking(booking schemas.Booking) (*mongo.InsertOneResult, error) {
	collection := GetBookingCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
// Set creation time
booking.CreatedAt = time.Now()
	// add a new theater
	result, err := collection.InsertOne(ctx, booking)
	if err != nil {
		log.Printf("Failed to insert booking: %v", err)
		return nil, err
	}

	return result, nil
}

// get all bokkings of user by user id

func GetAllUserBookings(user string)([]schemas.Booking,error){
	collection:= GetBookingCollection()
	var bookings []schemas.Booking
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// find all bookings having user matching the id provided
	cursor,err:= collection.Find(ctx, bson.M{"user": user})
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	// Iterate over the cursor and decode each document into a booking struct
	for cursor.Next(ctx) {
		var booking schemas.Booking
		if err := cursor.Decode(&booking); err != nil {
			log.Printf("Failed to decode booking: %v", err)
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	// If no theaters were found, return a custom error
	if len(bookings) == 0 {
		return nil, fmt.Errorf("no theaters found for user ID %s", user)
	}

	return bookings, nil
}

func GetOwnerReport(theaterNames []string) ([]schemas.Booking,error){
// Get the booking collection
collection := GetBookingCollection()

// Define a slice to hold the results
var bookings []schemas.Booking

// Create a context with a timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Create a filter to match bookings for the given theater names
filter := bson.M{"theater": bson.M{"$in": theaterNames}}

// Perform the query
cursor, err := collection.Find(ctx, filter)
if err != nil {
	return nil, err
}
defer cursor.Close(ctx)

// Iterate through the cursor and decode each document
for cursor.Next(ctx) {
	var booking schemas.Booking
	if err := cursor.Decode(&booking); err != nil {
		return nil, err
	}
	bookings = append(bookings, booking)
}

// Check for errors from iterating over the cursor
if err := cursor.Err(); err != nil {
	return nil, err
}

// Return the results
return bookings, nil
}