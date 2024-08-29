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
		"go.mongodb.org/mongo-driver/bson/primitive"
	
		
	)
// get movie collection
func GetMovieCollection()*mongo.Collection{
	if config.MongoClient == nil{
		log.Fatal("MongoDB client is not initialized")
	}
	return config.MongoClient.Database("movie").Collection("movies")
}


// Add a Movie
func AddMovie(movie schemas.Movie) (*mongo.InsertOneResult, error) {
	collection := GetMovieCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// add a new theater
	result, err := collection.InsertOne(ctx, movie)
	if err != nil {
		log.Printf("Failed to insert theater: %v", err)
		return nil, err
	}

	return result, nil
}
// get all movies added

func GetAllMovies()([]schemas.Movie, error){
	collection := GetMovieCollection()
	var movies []schemas.Movie
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	
	// Find all theaters with the given owner ID
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and decode each document into a Theater struct
	for cursor.Next(ctx) {
		var movie schemas.Movie
		if err := cursor.Decode(&movie); err != nil {
			log.Printf("Failed to decode theater: %v", err)
			return nil, err
		}
		movies = append(movies, movie)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	// If no theaters were found, return a custom error
	if len(movies) == 0 {
		return nil, fmt.Errorf("no theaters found ")
	}

	return movies, nil
}

//update  a movie 
func UpdateMovieById(movie schemas.Movie, id string) (*mongo.UpdateResult, error) {
	collection := GetMovieCollection()
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
			"movieName":   movie.MovieName,
			"description":         movie.Description,
			"language":         movie.Language,
			"releaseDate":         movie.ReleaseDate,
			"revenue":          movie.Revenue,
			"genre":       movie.Genre,
			"image":         movie.Image,
			"trailerId":         movie.TrailerId,
			
		},
	}

	// Perform the update operation
	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update movie: %v", err)
	}

	return updateResult, nil
}

// find a movie by id
func FindMovieById(id string) (schemas.Movie, error) {
	collection := GetMovieCollection()
	var movie schemas.Movie
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
// Convert string ID to ObjectID
objectID, err := primitive.ObjectIDFromHex(id)
if err != nil {
	return schemas.Movie{}, fmt.Errorf("invalid ID format: %v", err)
}
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&movie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return a custom error message if no document is found
			return schemas.Movie{}, fmt.Errorf("movie with Id %s not found", id)
		}
		// Log and return the error if there is a database error
		log.Printf("Failed to find theater: %v", err)
		return schemas.Movie{}, err
	}

	return movie, nil
}