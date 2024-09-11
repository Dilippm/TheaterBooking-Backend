package controllers

import (

	"net/http"
	"strconv"

	"github.com/dilippm92/bookingapplication/config"
	"github.com/dilippm92/bookingapplication/models/queries"
	"github.com/dilippm92/bookingapplication/models/schemas"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)


var bookingInput struct {
	Theater       string    `json:"theater"`       // Name of the movie
    SelectedSeats []string  `json:"selectedSeats"` // Selected seats for the reservation
	Time         string `json:"time"`          // Show time
    Price         string    `json:"price"`         // Price of the reservation
    Date          string    `json:"date"`          // Date of the reservation
    User          string    `json:"user"`          // User ID making the reservation
	 PaymentId          string    `json:"paymentId"`
	 Movie          string    `json:"movie"`
}
type RequestBody struct {
	Amount string `json:"amount"` // Amount should be in cents
}

func CreatePaymentIntent(c *gin.Context) {
	// Initialize Stripe
	config.InitStripe()

// Parse the request body
var requestBody RequestBody
if err := c.ShouldBindJSON(&requestBody); err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
	return
}

// Retrieve amount from request body and convert it to int64
amount, err := strconv.ParseInt(requestBody.Amount, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount format"})
		return
	}

	// Create a PaymentIntent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(string(stripe.CurrencyINR)),
		Description: stripe.String("Export transaction - testing"), // Add a description here,
		
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return client secret to the frontend
	c.JSON(http.StatusOK, gin.H{
		"clientSecret": pi.ClientSecret,
	})
}
func AddBooking(c *gin.Context) {
	

	if err := c.ShouldBindJSON(&bookingInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return

	}
	
    
	// Add a new theater
	booking := schemas.Booking{
		Theater:       bookingInput.Theater,
		SelectedSeats: bookingInput.SelectedSeats,
		Time:          bookingInput.Time,
Price: bookingInput.Price,
		Date: bookingInput.Date,
		User: bookingInput.User,
		PaymentId:bookingInput.PaymentId,
		Movie:bookingInput.Movie,
	}
	result, err := queries.AddBooking(booking)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorad": err.Error()})
		return
	}
	_,err =queries.UpdateWallet(bookingInput.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erroradmin": err.Error()})
		return
	}
	theater,err :=queries.FindTheaterByName(bookingInput.Theater)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ownerIDStr := theater.OwnerID.Hex()
	_,err = queries.UpdateWalletByUserId(ownerIDStr,bookingInput.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorowner": err.Error()})
		return
	}
	// Return success message
	c.JSON(http.StatusCreated, gin.H{"message": "Booking Done Successfully","result": result.InsertedID})
	
}

// get bookings by user

func GetUserBookings (c *gin.Context){
user:= c.Param("id")
bookings, err := queries.GetAllUserBookings(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}