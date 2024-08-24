package main

import (
	"github.com/dilippm92/bookingapplication/config"
	"github.com/dilippm92/bookingapplication/routes"
	"github.com/gin-gonic/gin"
)
func main()  {
	port:= config.PORT
	config.ConnectMongoDB()
	router := gin.New()
	routes.MainRoutes(router)
	router.Use(gin.Logger())
	router.Run(":" + port)
}