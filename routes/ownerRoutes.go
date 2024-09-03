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
		ownerGroup.GET("/get_theaters_owner/:id",middlewares.JwtTokenVerify(),controllers.GetAllOnwertheaters)
		ownerGroup.GET("/get_theater_details/:id",middlewares.JwtTokenVerify(),controllers.GetSpecificTheaterByid)
		ownerGroup.PUT("/update_theater/:id",middlewares.JwtTokenVerify(),controllers.UpdateTheater)
		ownerGroup.GET("/get_theater",middlewares.JwtTokenVerify(),controllers.GetTheaterByQuery)
	}
}