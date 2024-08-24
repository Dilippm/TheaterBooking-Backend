package routes

import (
	"github.com/dilippm92/bookingapplication/controllers"
	"github.com/dilippm92/bookingapplication/middlewares"
	"github.com/gin-gonic/gin"
)

// Authroutes sets up authentication-related routes
func Authroutes(routerGroup *gin.RouterGroup) {
	authGroup := routerGroup.Group("/auth/user")
	{
		authGroup.GET("/test", controllers.TestSample) 
		authGroup.POST("/register", controllers.SignUp) 
		authGroup.POST("/login",controllers.UserLogin)
		authGroup.PUT("/update_profile",middlewares.JwtTokenVerify(),controllers.UserUpdate)
	}
	}
