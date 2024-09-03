package routes

import (
	"github.com/dilippm92/bookingapplication/controllers"
	"github.com/dilippm92/bookingapplication/middlewares"
	"github.com/gin-gonic/gin"
)
func Reservationroutes(routerGroup *gin.RouterGroup){
	reservationGroup:= routerGroup.Group("/reservation")
	{
		reservationGroup.POST("/add_reservation",middlewares.JwtTokenVerify(),controllers.Addreservation)
		reservationGroup.GET("/get_reservation/:time/:date",middlewares.JwtTokenVerify(),controllers.GetReservation)
		
		
	}
}