package routes

import (
	"github.com/dilippm92/bookingapplication/controllers"
	"github.com/dilippm92/bookingapplication/middlewares"
	"github.com/gin-gonic/gin"
)
func Reportroutes(routerGroup *gin.RouterGroup){
	reportGroup:= routerGroup.Group("/report")
	{
		
		reportGroup.GET("/get_owner_report/:id",middlewares.JwtTokenVerify(),controllers.GetOwnerReport)
		reportGroup.GET("/get_admin_report",middlewares.JwtTokenVerify(),controllers.GetAdminReport)
		
		
	}
}