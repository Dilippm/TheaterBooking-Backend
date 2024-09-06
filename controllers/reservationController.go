package controllers

import (
	"net/http"
	"time"

	"github.com/dilippm92/bookingapplication/models/queries"
	"github.com/dilippm92/bookingapplication/models/schemas"
	"github.com/gin-gonic/gin"
)

var reservationInput struct {
	Theater       string    `json:"theater"`       // Name of the movie
    SelectedSeats []string  `json:"selectedSeats"` // Selected seats for the reservation
    Time          time.Time `json:"time"`          // Show time
    Price         string    `json:"price"`         // Price of the reservation
    Date          string    `json:"date"`          // Date of the reservation
    User          string    `json:"user"`          // User ID making the reservation
	Movie		string		`json:"movie"`  
}

// function to add a new movie

func Addreservation(c *gin.Context) {
	if err := c.ShouldBindJSON(&reservationInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return

	}


	// Add a new theater
	reservation := schemas.Reservation{
		Theater:       reservationInput.Theater,
		SelectedSeats: reservationInput.SelectedSeats,
		Time:          reservationInput.Time,
Price: reservationInput.Price,
		Date: reservationInput.Date,
		User: reservationInput.User,
		Movie: reservationInput.Movie,
	}
	result, err := queries.AddReservation(reservation)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Return success message
	c.JSON(http.StatusCreated, gin.H{"message": "Reservation Added Successfully", "result": result.InsertedID})

}

// function to get  a single movie details
func GetReservation(c *gin.Context){
	time := c.Param("time")
	date:= c.Param("date")
	movie, err := queries.GetReservationsByTimeAndDate(date,time)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations"})
		return
	}

	c.JSON(http.StatusOK, movie)
}
