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
func GetConversationCollection()*mongo.Collection{
	if config.MongoClient == nil{
		log.Fatal("MongoDB client is not initialized")
	}
	return config.MongoClient.Database("movie").Collection("conversations")
}

func GetMessageCollection()*mongo.Collection{
	if config.MongoClient == nil{
		log.Fatal("MongoDB client is not initialized")
	}
	return config.MongoClient.Database("movie").Collection("messages")
}

// MakeMessage creates a new message between the sender and receiver and links it to the conversation
func MakeMessage(receiverId, senderId, message string) (schemas.Message, error) {
	messageCollection := GetMessageCollection()
	conversationCollection := GetConversationCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert string IDs to ObjectIDs
	receiverObjectID, err := primitive.ObjectIDFromHex(receiverId)
	if err != nil {
		return schemas.Message{}, fmt.Errorf("invalid receiver ID: %v", err)
	}
	senderObjectID, err := primitive.ObjectIDFromHex(senderId)
	if err != nil {
		return schemas.Message{}, fmt.Errorf("invalid sender ID: %v", err)
	}

	// Check if a conversation exists between the sender and receiver
	filter := bson.M{"participants": bson.M{"$all": []primitive.ObjectID{senderObjectID, receiverObjectID}}}
	var conversation schemas.Conversation
	err = conversationCollection.FindOne(ctx, filter).Decode(&conversation)

	if err == mongo.ErrNoDocuments {
		// Create a new conversation if none exists
		conversation = schemas.Conversation{
			ID:           primitive.NewObjectID(),
			Participants: []primitive.ObjectID{senderObjectID, receiverObjectID},
			Messages:     []primitive.ObjectID{}, // Empty messages initially
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		_, insertErr := conversationCollection.InsertOne(ctx, conversation)
		if insertErr != nil {
			return schemas.Message{}, fmt.Errorf("failed to create new conversation: %v", insertErr)
		}
		log.Println("New conversation created successfully")
	} else if err != nil {
		return schemas.Message{}, fmt.Errorf("error finding conversation: %v", err)
	}

	// Create the new message
	newMessage := schemas.Message{
		ID:         primitive.NewObjectID(),
		SenderID:   senderObjectID,
		ReceiverID: receiverObjectID,
		Message:    message,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Insert the new message into the messages collection
	_, err = messageCollection.InsertOne(ctx, newMessage)
	if err != nil {
		return schemas.Message{}, fmt.Errorf("failed to save message: %v", err)
	}

	// Update the conversation to add the new message ID
	update := bson.M{"$push": bson.M{"messages": newMessage.ID}, "$set": bson.M{"updatedAt": time.Now()}}
	_, err = conversationCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return schemas.Message{}, fmt.Errorf("failed to update conversation with new message ID: %v", err)
	}

	log.Println("Message created and linked to conversation successfully")
	return newMessage, nil
}

// FindConversation retrieves a conversation between two users
func FindConversation(senderID, recipientID string) (schemas.Conversation, error) {
    collection := GetConversationCollection()

    // Convert string IDs to ObjectIDs
    receiverObjectID, err := primitive.ObjectIDFromHex(recipientID)
    if err != nil {
        return schemas.Conversation{}, fmt.Errorf("invalid receiver ID: %v", err)
    }
    senderObjectID, err := primitive.ObjectIDFromHex(senderID)
    if err != nil {
        return schemas.Conversation{}, fmt.Errorf("invalid sender ID: %v", err)
    }

    // Filter to find the conversation
    filter := bson.M{"participants": bson.M{"$all": []primitive.ObjectID{senderObjectID, receiverObjectID}}}

    var conversation schemas.Conversation
    err = collection.FindOne(context.TODO(), filter).Decode(&conversation)

    // Check if an error occurred while finding the conversation
    if err != nil {
        if err == mongo.ErrNoDocuments {
            // Optionally return an empty conversation or handle it as needed
            return schemas.Conversation{}, nil // Returning an empty conversation if not found
        }
        return schemas.Conversation{}, fmt.Errorf("error retrieving conversation: %v", err)
    }

    return conversation, nil // Successfully found and returned the conversation
}

// PopulateMessages retrieves message documents based on the IDs in the conversation
func PopulateMessages(messageIDs []primitive.ObjectID) ([]schemas.Message, error) {
	collection := GetMessageCollection()
	cursor, err := collection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": messageIDs}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var messages []schemas.Message
	for cursor.Next(context.TODO()) {
		var message schemas.Message
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, cursor.Err()
}