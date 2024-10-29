package controllers

import (
	"net/http"
	"github.com/dilippm92/bookingapplication/models/queries"
	"github.com/gin-gonic/gin"
)
func GetUserBookingAnalytics(c *gin.Context){
user:=c.Param("id")
userBookingAnalytics,err:=queries.GetAllUserBookingsData(user)
if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
	return
}

c.JSON(http.StatusOK, userBookingAnalytics)
}



func GetAdminAnalytics(c *gin.Context){
	adminBookingAnalytics,err:=queries.GetAdminBookingData()
if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
	return
}

c.JSON(http.StatusOK, adminBookingAnalytics)
}


func GetOwnerAnalytics(c *gin.Context){
	user:=c.Param("id")
	theaters, err := queries.GetAllOnwertheater(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch theaters"})
		return
	}
	var theaterNames []string
	for _,theater:= range theaters{
		theaterNames = append(theaterNames, theater.TheaterName)
	}
	ownerBookingAnalytics,err:=queries.GetAllOwnerData(theaterNames)
if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
	return
}

c.JSON(http.StatusOK, ownerBookingAnalytics)
}