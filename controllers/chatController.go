package controllers

import (
	"fmt"
	"net/http"

	"github.com/dilippm92/bookingapplication/models/queries"
	"github.com/gin-gonic/gin"
)


func Getusers(c *gin.Context){
	// Get the logged-in user's ID from the context
	loggedInUserId := c.MustGet("userId").(string)

chatUsers,err := queries.GetChatUsers(loggedInUserId) 
if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
	return
}
c.JSON(http.StatusOK, chatUsers)

}

// SendMessage handles sending a message between two users
func SendMessage(c *gin.Context) {
	receiverIdStr := c.Param("id")
	senderId := c.MustGet("userId").(string)

	 // Parse the request body
	 var messagePayload struct {
        Message string `json:"message" binding:"required"`
    }

    if err := c.ShouldBindJSON(&messagePayload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }


	// Call MakeMessage to create the message and link it to the conversation
	result, err := queries.MakeMessage(receiverIdStr, senderId, messagePayload.Message)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "invalid receiver ID" || err.Error() == "invalid sender ID" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Respond with the created message
	c.JSON(http.StatusCreated, gin.H{"message": result})
	fmt.Println("Message sent:", result)
}

// GetMessages retrieves messages for a specific conversation
func GetMessages(c *gin.Context) {
	userToChatId := c.Param("id")                  // Get the recipient user ID from the URL parameters
	senderId := c.MustGet("userId").(string)       // Get the sender's ID from the context (assuming userId is set in middleware)

	

	// Find the conversation between the two users
	conversation, err := queries.FindConversation(senderId, userToChatId)
	if err != nil {
		
		fmt.Println("Error finding conversation: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Retrieve messages from the conversation
	messages, err := queries.PopulateMessages(conversation.Messages)
	if err != nil {
		fmt.Println("Error populating messages: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, messages) // Return the messages as JSON
}