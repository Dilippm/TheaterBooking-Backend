package schemas

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Conversation represents the conversation model
type Conversation struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Participants []primitive.ObjectID `bson:"participants" json:"participants"`
	Messages     []primitive.ObjectID `bson:"messages" json:"messages"`
	CreatedAt    time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time            `bson:"updatedAt" json:"updatedAt"`
}
