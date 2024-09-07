package controllers

import (
	"fmt"
	"time"
	"net/http"
"strconv"
	"github.com/dilippm92/bookingapplication/models/queries"
	"github.com/gin-gonic/gin"
)
type TheaterReport struct {
	Date  string `json:"date"`
	TotalPrice float64 `json:"totalPrice"`
}

func GetOwnerReport(c *gin.Context){
	ownerID := c.Param("id")

	// Fetch theaters for the given owner ID
	theaters, err := queries.GetAllOnwertheater(ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch theaters"})
		return
	}

	// Extract theater names
	var theaterNames []string
	for _, theater := range theaters {
		theaterNames = append(theaterNames, theater.TheaterName)
	}

	// Fetch bookings based on theater names
	bookings, err := queries.GetOwnerReport(theaterNames)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	// Group bookings by date and calculate total price
	priceByDate := make(map[string]float64)
	for _, booking := range bookings {
		// Parse date and calculate total price
		parsedDate, err := time.Parse("Mon Jan 02 2006", booking.Date)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}
		formattedDate := parsedDate.Format("2006-01-02") // Format date as YYYY-MM-DD

		price, err := strconv.ParseFloat(booking.Price, 64)
		if err != nil {
			fmt.Println("Error parsing price:", err)
			continue
		}

		// Aggregate total price for each date
		priceByDate[formattedDate] += price
	}

	// Transform aggregated data into response format
	var reportData []TheaterReport
	for date, totalPrice := range priceByDate {
		reportData = append(reportData, TheaterReport{
			Date:       date,
			TotalPrice: totalPrice,
		})
	}

	// Send response with grouped data
	c.JSON(http.StatusOK, gin.H{"theaterReports": reportData})
}

func GetAdminReport(c *gin.Context) {
		// Fetch owners
		owners, err := queries.GetOwnerIds()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch owners"})
			return
		}
	
		// Create a new slice to hold the filtered owners
		var filteredOwners []struct {
			ID         string  `json:"id"`
			Username   string  `json:"username"`
			Date       string  `json:"date"`
			TotalPrice float64 `json:"totalPrice"`
		}
	
		// Iterate over the owners
		for _, owner := range owners {
			// Fetch theaters for the given owner ID
			theaters, err := queries.GetAllOnwertheater(owner.Id.Hex())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch theaters"})
				return
			}
	
			// Extract theater names
			var theaterNames []string
			for _, theater := range theaters {
				theaterNames = append(theaterNames, theater.TheaterName)
			}
	
			// Fetch bookings based on theater names
			bookings, err := queries.GetOwnerReport(theaterNames)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
				return
			}
	
			// Group bookings by date and calculate total price
			priceByDate := make(map[string]float64)
			for _, booking := range bookings {
				// Parse date and calculate total price
				parsedDate, err := time.Parse("Mon Jan 02 2006", booking.Date)
				if err != nil {
					fmt.Println("Error parsing date:", err)
					continue
				}
				formattedDate := parsedDate.Format("2006-01-02") // Format date as YYYY-MM-DD
	
				price, err := strconv.ParseFloat(booking.Price, 64)
				if err != nil {
					fmt.Println("Error parsing price:", err)
					continue
				}
	
				// Aggregate total price for each date
				priceByDate[formattedDate] += price
			}
	
			// Create report data for this owner
			for date, totalPrice := range priceByDate {
				filteredOwners = append(filteredOwners, struct {
					ID         string  `json:"id"`
					Username   string  `json:"username"`
					Date       string  `json:"date"`
					TotalPrice float64 `json:"totalPrice"`
				}{
					ID:         owner.Id.Hex(),
					Username:   owner.Username,
					Date:       date,
					TotalPrice: totalPrice,
				})
			}
		}
	
		// Send response with filtered and grouped data
		c.JSON(http.StatusOK, gin.H{"ownerReport": filteredOwners})
}
