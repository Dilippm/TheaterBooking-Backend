package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)
//test controller for auth routes
func TestSample(c *gin.Context)  {
	c.JSON(http.StatusOK,gin.H{"message":"Auth Test Route"})
}