package routes

import (
	"github.com/dilippm92/bookingapplication/controllers"
	"github.com/dilippm92/bookingapplication/middlewares"
	"github.com/gin-gonic/gin"
)
func Reservationroutes(routerGroup *gin.RouterGroup){
	reservationGroup:= routerGroup.Group("/reservation")
	bookingGroup:=routerGroup.Group("/bookings")
	{
		reservationGroup.POST("/add_reservation",middlewares.JwtTokenVerify(),controllers.Addreservation)
		reservationGroup.GET("/get_reservation/:time/:date",middlewares.JwtTokenVerify(),controllers.GetReservation)
		
		
	}
	{
		bookingGroup.POST("/create-payment-intent",controllers.CreatePaymentIntent)
		bookingGroup.POST("/add_booking",controllers.AddBooking)
		bookingGroup.GET("/user_bookings/:id",controllers.GetUserBookings)
		bookingGroup.GET("/get_booking/:time/:date",middlewares.JwtTokenVerify(),controllers.GetBookings)
	}
}