package routes

import "github.com/gin-gonic/gin"
func MainRoutes(router *gin.Engine){
	apiGroup:= router.Group("/api")
	Authroutes(apiGroup)
}