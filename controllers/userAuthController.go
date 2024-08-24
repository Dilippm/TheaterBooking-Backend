package controllers

import (
	"fmt"
	"net/http"

	"github.com/dilippm92/bookingapplication/helpers"
	"github.com/dilippm92/bookingapplication/models/queries"
	"github.com/dilippm92/bookingapplication/models/schemas"
	"github.com/gin-gonic/gin"
)

//test controller for auth routes
func TestSample(c *gin.Context)  {
	c.JSON(http.StatusOK,gin.H{"message":"Auth Test Route"})
}

// funciton to register a user
func SignUp(c *gin.Context){
var user schemas.User
if err:= c.ShouldBindJSON(&user);err!=nil{
	c.Error(fmt.Errorf("failed to bind request body to user model: %v", err))
		c.AbortWithStatus(http.StatusBadRequest) // Bad Request is more appropriate for validation errors
		return
}
// Compare Password and ConfirmPassword
if user.Password != user.ConfirmPassword {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Password and Confirm Password do not match",
	})
	return
}
// hash password 
hashedPassword,err:= helpers.HashPassword(user.Password)
if err != nil {
	c.Error(fmt.Errorf("failed to hash password: %v", err))
	c.AbortWithStatus(http.StatusInternalServerError) // Internal Server Error for issues with server-side processing
	return
}
user.Password = hashedPassword
// Call CreateUser to insert the user into the database
result, err := queries.CreateUser(user)
if err != nil {
	c.Error(fmt.Errorf("failed to create user in the database: %v", err))
	c.AbortWithStatus(http.StatusInternalServerError)
	return
}

// Return a success response with the result
c.JSON(http.StatusCreated, gin.H{
	"message": "User Created Successfully",
	"result":  result.InsertedID,
})
}