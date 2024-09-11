package routes

import (
	"github.com/dilippm92/bookingapplication/controllers"
	"github.com/dilippm92/bookingapplication/middlewares"
	"github.com/gin-gonic/gin"
)
func Analyticsroutes(routerGroup *gin.RouterGroup){
	analyticsGroup:= routerGroup.Group("/analytics")
	{
		analyticsGroup.GET("/get_user_booking_data/:id",middlewares.JwtTokenVerify(),controllers.GetUserBookingAnalytics)
		analyticsGroup.GET("/get_admin_booking_data",middlewares.JwtTokenVerify(),controllers.GetAdminAnalytics)
		analyticsGroup.GET("/get_owner_booking_data/:id",middlewares.JwtTokenVerify(),controllers.GetOwnerAnalytics)
		
	}
}
