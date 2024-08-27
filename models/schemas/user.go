package schemas
import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the structure of a user document in MongoDB
type User struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username    string             `bson:"username" json:"username"`
	Email       string             `bson:"email" json:"email"`
	Password    string             `bson:"password" json:"password"`
	UserImage   string             `bson:"userimage" json:"userimage"`
	Reservation []string           `bson:"reservation" json:"reservation"`  
	Bookings    []string           `bson:"bookings" json:"bookings"`
	Role        string             `bson:"role" json:"role"`
}