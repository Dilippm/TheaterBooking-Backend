package main

import (
	"github.com/dilippm92/bookingapplication/config"
	"github.com/gin-gonic/gin"
)
func main()  {
	port:= config.PORT
	config.ConnectMongoDB()
	router := gin.New()
	router.Run(":" + port)
}