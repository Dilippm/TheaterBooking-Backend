package controllers

import (

	"net/http"
"fmt"
	"github.com/dilippm92/bookingapplication/helpers"
	"github.com/dilippm92/bookingapplication/models/queries"
	"github.com/dilippm92/bookingapplication/models/schemas"
	"github.com/gin-gonic/gin"
)

//test controller for auth routes
func TestSample(c *gin.Context)  {
	c.JSON(http.StatusOK,gin.H{"message":"Auth Test Route"})
}
var input struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	UserImage 		string 	`json:"UserImage"`
	Role			string	`json:"role"`
}
// funciton to register a user
func SignUp(c *gin.Context) {
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash the password
	hashedPassword, err := helpers.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create a new user struct
	user := schemas.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
		Role:input.Role,
	}

	// Try to find or create the user
	result, err := queries.FindOrCreateUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return success message
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "result": result.InsertedID})
}


//function to login a user and create jwt token
func UserLogin(c *gin.Context) {


	// Bind the JSON input to loginData struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(fmt.Errorf("failed to bind request body to user model: %v", err))
		c.AbortWithStatus(http.StatusBadRequest) // Bad Request is more appropriate for validation errors
		return
	}
	// Find the user by email
	user, err := queries.FindUserByEmail(input.Email)
	
	if err != nil {
		if err.Error() == fmt.Sprintf("user with email %s not found", input.Email) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}
	// Compare the password with the stored hash using the new function
	err = helpers.ComparePasswords(user.Password, input.Password)
	if err != nil {
		c.Error(fmt.Errorf("email or password wrong"))
		c.AbortWithStatus(http.StatusInternalServerError) // Internal Server Error for issues with server-side processing
		return
	}
	// Generate JWT token
	token, err := helpers.GenerateJWTToken(user.Id.Hex(), user.Email)
	if err != nil {
		c.Error(fmt.Errorf("token generation failed"))
		c.AbortWithStatus(http.StatusInternalServerError) // Internal Server Error for issues with server-side processing
		return
	}
	// Return the token to the client
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user":    user})
}

func UserUpdate(c *gin.Context){
	userId := c.Param("id")
// Bind the JSON input to loginData struct
if err := c.ShouldBindJSON(&input); err != nil {
	c.Error(fmt.Errorf("failed to bind request body to user model: %v", err))
	c.AbortWithStatus(http.StatusBadRequest) // Bad Request is more appropriate for validation errors
	return
}
// Find the user by email
user, err := queries.FindUserById(userId)

if err != nil {
	if err.Error() == fmt.Sprintf("user with id %s not found",userId) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not Found"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	}
	return
}
// Create a new user struct
updateUser := schemas.User{
	Username: input.Username,
	Email:    input.Email,
	Password: user.Password,
	UserImage: input.UserImage,
	Role:user.Role,
}

// Try to find or create the user
result, err := queries.UpdateUser(updateUser,userId)
if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	return
}

// Return success message
c.JSON(http.StatusCreated, gin.H{"message": "User Profile Updated Successfully", "result": result})
}


