package main

import (
	"github.com/dilippm92/bookingapplication/config"
	"github.com/dilippm92/bookingapplication/routes"
	"github.com/gin-gonic/gin"
	 "github.com/gin-contrib/cors"
)
func main()  {
	port:= config.PORT
	config.ConnectMongoDB()
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
	
	router.Run(":" + port)
}