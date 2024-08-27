package routes

import (
	"github.com/dilippm92/bookingapplication/controllers"
	"github.com/dilippm92/bookingapplication/middlewares"
	"github.com/gin-gonic/gin"
)
func Ownerroutes(routerGroup *gin.RouterGroup){
	ownerGroup:= routerGroup.Group("/owner")
	{
		ownerGroup.POST("/add_theater",middlewares.JwtTokenVerify(),controllers.Addtheater)
	}
}