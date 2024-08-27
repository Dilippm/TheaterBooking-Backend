package controllers

import (
	"fmt"
	"net/http"
	"time"

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