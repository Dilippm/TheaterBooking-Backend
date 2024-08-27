package queries
import (
"context"
"time"
	"log"
	
	"github.com/dilippm92/bookingapplication/config"
		"github.com/dilippm92/bookingapplication/models/schemas"
	"go.mongodb.org/mongo-driver/mongo"
	
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
