package schemas
import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct{
	ID          primitive.ObjectID `bson:"_id,omitempty"`  // MongoDB's ObjectID
	SenderID    primitive.ObjectID `bson:"senderId" json:"senderId" binding:"required"`
	ReceiverID  primitive.ObjectID `bson:"receiverId" json:"receiverId" binding:"required"`
	Message     string             `bson:"message" json:"message" binding:"required"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`

}