package main

import (
	"github.com/dilippm92/bookingapplication/config"
	"github.com/dilippm92/bookingapplication/models/queries"
	"github.com/dilippm92/bookingapplication/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust origin check for security if needed
	},
}

var userSocketMap = make(map[string]*websocket.Conn) // {userId: *websocket.Conn}

// WebSocket handler
func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	userId := c.Query("userId")
	userSocketMap[userId] = conn
	defer delete(userSocketMap, userId)

	// Notify other users about the online status
	broadcastOnlineUsers()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("Received message from %s: %s", userId, msg)

		// Broadcast the message to all connected clients
		broadcastMessage(userId, msg)
	}
}

// Function to broadcast online users
func broadcastOnlineUsers() {
	onlineUsers := make([]string, 0, len(userSocketMap))
	for userId := range userSocketMap {
		onlineUsers = append(onlineUsers, userId)
	}
	for _, conn := range userSocketMap {
		if err := conn.WriteJSON(onlineUsers); err != nil {
			log.Println("Error sending online users:", err)
			conn.Close()
		}
	}
}

// Function to broadcast messages
func broadcastMessage(senderId string, msg []byte) {
	for userId, conn := range userSocketMap {
		if userId != senderId {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("Error sending message to %s: %v", userId, err)
				conn.Close()
				delete(userSocketMap, userId)
			}
		}
	}
}

func main()  {
	port:= config.PORT
	config.ConnectMongoDB()
	queries.CreateTTLIndex()
	router := gin.New()
	 // CORS configuration
	 router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"}, // Allow requests from this origin
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allow these methods
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allow these headers
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))
	router.Use(gin.Logger())
	routes.MainRoutes(router)
	// WebSocket route
	router.GET("/ws", handleWebSocket)
	router.Run(":" + port)
}