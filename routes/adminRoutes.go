package routes

import (
	"github.com/dilippm92/bookingapplication/controllers"
	"github.com/dilippm92/bookingapplication/middlewares"
	"github.com/gin-gonic/gin"
)
func Adminroutes(routerGroup *gin.RouterGroup){
	adminGroup:= routerGroup.Group("/admin")
	{
		adminGroup.POST("/add_movie",middlewares.JwtTokenVerify(),controllers.Addmovie)
		adminGroup.GET("/get_movies",middlewares.JwtTokenVerify(),controllers.GetMovies)
		adminGroup.GET("/get_movie_details/:id",middlewares.JwtTokenVerify(),controllers.GetMovieById)
		adminGroup.PUT("/update_movie/:id",middlewares.JwtTokenVerify(),controllers.UpdateMovie)
		adminGroup.GET("/get_latest_movies",controllers.GetLatestMovies)
		adminGroup.DELETE("/delete_movie/:id",middlewares.JwtTokenVerify(),controllers.DeleteMOvieById)
	}
}