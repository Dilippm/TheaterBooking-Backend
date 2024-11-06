package routes

import (
	"github.com/dilippm92/bookingapplication/controllers"
	"github.com/dilippm92/bookingapplication/middlewares"
	"github.com/gin-gonic/gin"
)
func Chatroutes(routerGroup *gin.RouterGroup){
	chatGroup:= routerGroup.Group("/chat")
	{
		chatGroup.GET("/",middlewares.JwtTokenVerify(),controllers.Getusers)
		chatGroup.POST("/send_message/:id",middlewares.JwtTokenVerify(),controllers.SendMessage)
		chatGroup.GET("/get_messages/:id",middlewares.JwtTokenVerify(),controllers.GetMessages)
	}
}