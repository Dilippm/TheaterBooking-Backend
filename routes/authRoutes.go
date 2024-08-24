package routes

import (
	"github.com/dilippm92/bookingapplication/controllers"
	"github.com/gin-gonic/gin"
)

// Authroutes sets up authentication-related routes
func Authroutes(routerGroup *gin.RouterGroup) {
	authGroup := routerGroup.Group("/auth")
	{
		authGroup.GET("/test", controllers.TestSample) // Pass the handler function reference directly
	}
}
