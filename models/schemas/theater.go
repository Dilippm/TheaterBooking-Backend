package schemas

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Theater represents a theater document in MongoDB
type Theater struct {
	ID          primitive.ObjectID       `bson:"_id,omitempty"` // MongoDB's ObjectID
	TheaterName string                   `bson:"theaterName"`   // Name of the theater
	OwnerID     primitive.ObjectID       `bson:"ownerId"`       // Owner's ID, also an ObjectID
	Place       string                   `bson:"place"`         // Location of the theater
	State       string                   `bson:"state"`         // State where the theater is located
	Movie       string                   `bson:"movie"`         // Currently playing movie
	Rows        int                      `bson:"rows"`          // Number of rows in the theater
	Columns     int                      `bson:"columns"`       // Number of columns in the theater
	Seats       int                      `bson:"seats"`         // Total number of seats
	Price       map[string]float64       `bson:"price"`         // Price could vary, e.g., by type (e.g., regular, VIP)
	ShowTimings []time.Time              `bson:"showTimings"`   // Different show timings
}
