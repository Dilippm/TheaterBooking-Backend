package controllers

import (
	"fmt"
	"net/http"
	"time"
"log"
	"github.com/dilippm92/bookingapplication/models/queries"
	"github.com/dilippm92/bookingapplication/models/schemas"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
var theaterinput struct {
	TheaterName string                   `bson:"theaterName"`   // Name of the theater
	OwnerID     string       `bson:"ownerId"`       // Owner's ID, also an ObjectID
	Place       string                   `bson:"place"`         // Location of the theater
	State       string                   `bson:"state"`         // State where the theater is located
	Movie       string                   `bson:"movie"`         // Currently playing movie
	Rows        int                      `bson:"rows"`          // Number of rows in the theater
	Columns     int                      `bson:"columns"`       // Number of columns in the theater
	Seats       int                      `bson:"seats"`         // Total number of seats
	Price       map[string]float64       `bson:"price"`         // Price could vary, e.g., by type (e.g., regular, VIP)
	ShowTimings []time.Time              `bson:"showTimings"`   // Different show timings
}
// function to add a new theater

func Addtheater(c *gin.Context){
	if err := c.ShouldBindJSON(&theaterinput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return

	}
	// Convert the string to a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(theaterinput.OwnerID)
	if err != nil {
		fmt.Println("Invalid ObjectID:", err)
		return
	}
// Add a new theater
theater:= schemas.Theater{
	TheaterName: theaterinput.TheaterName,
	OwnerID: objectID,
	Place: theaterinput.Place,
	State: theaterinput.State,
	Movie: theaterinput.Movie,
Rows: theaterinput.Rows,
Columns: theaterinput.Columns,
Seats: theaterinput.Seats,
Price: theaterinput.Price,
ShowTimings: theaterinput.ShowTimings,

}
result,err:= queries.AddTheater(theater)
if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	return
}
// Return success message
c.JSON(http.StatusCreated, gin.H{"message": "Theater Added Successfully", "result": result.InsertedID})

}

// function to get all theaters for a owner id
func GetAllOnwertheaters(c *gin.Context){
	ownerID := c.Param("id")
	theaters, err := queries.GetAllOnwertheater(ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch theaters"})
		return
	}

	c.JSON(http.StatusOK, theaters)
}

// function to update a theater by id
func UpdateTheater(c *gin.Context) {
	theaterId := c.Param("id")

	// Bind the JSON input to input struct
	var input schemas.Theater
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(fmt.Errorf("failed to bind request body to theater model: %v", err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Find the theater by ID
	theater, err := queries.FindTheaterById(theaterId)
	if err != nil {
		if err.Error() == fmt.Sprintf("theater with id %s not found", theaterId) {
			c.JSON(http.StatusNotFound, gin.H{"error": "theater not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	// Create an update struct
	updateTheater := schemas.Theater{
		TheaterName:   input.TheaterName,
		OwnerID:       theater.OwnerID, // Assuming you don’t want to change the OwnerID
		Place:         input.Place,
		State:         input.State,
		Movie:         input.Movie,
		Rows:          input.Rows,
		Columns:       input.Columns,
		Seats:         input.Seats,
		Price:         input.Price,
		ShowTimings:   input.ShowTimings,
	}

	// Try to update the theater
	result, err := queries.UpdateTheaterDetail(updateTheater, theaterId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "Theater updated successfully", "result": result})
}

// function to get  a single theater details
func GetSpecificTheaterByid(c *gin.Context){
	id := c.Param("id")
	theaters, err := queries.FindTheaterById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch theater"})
		return
	}

	c.JSON(http.StatusOK, theaters)
}

// function to get theater details with given movie  id or theater name or place
func GetTheaterByQuery(c *gin.Context) {
    // Extract query parameters
    name := c.Query("name")
    place := c.Query("place")
    id := c.Query("id")

    // Call the service layer to get the theater
    theater, err := queries.GetTheatersByNamePlaceId(name, place, id)
    if err != nil {
        // Check if the error is a not found error or another type of error
        if err.Error() == "theater with name or place not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        } else {
            // Log and return internal server error
            log.Printf("Error fetching theater: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while fetching the theater"})
        }
        return
    }

    // Fetch the movie by its ID
    movie, err := queries.FindMovieById(id)
    if err != nil {
        log.Printf("Error fetching movie: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while fetching the movie"})
        return
    }

    // Return both theater and movie data as a combined response
    c.JSON(http.StatusOK, gin.H{
        "theater": theater,
        "movie":   movie,
    })
}
