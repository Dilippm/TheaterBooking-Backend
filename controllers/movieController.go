package controllers

import (
	
	"net/http"
	
"fmt"
	"github.com/dilippm92/bookingapplication/models/queries"
	"github.com/dilippm92/bookingapplication/models/schemas"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
var movieInput struct {

	MovieName   string             `bson:"movieName"`      
	Description string             `bson:"description"`   
	Language    string             `bson:"language"`       
	ReleaseDate primitive.DateTime `bson:"releaseDate"`   
	Revenue     string             `bson:"revenue"`       
	Genre		string				`bson:"genre"`
	Image 		string 				`bson:"image"`
	TrailerId   string				`bson:"trailerId"`
}

// function to add a new movie

func Addmovie(c *gin.Context){
	if err := c.ShouldBindJSON(&movieInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return

	}
	
// Add a new theater
movie:= schemas.Movie{
	MovieName: movieInput.MovieName,
	
	Description: movieInput.Description,
	Language: movieInput.Language,
	ReleaseDate: movieInput.ReleaseDate,
	Revenue: movieInput.Revenue,
	Genre: movieInput.Genre,
	Image: movieInput.Image,
	TrailerId: movieInput.TrailerId,


}
result,err:= queries.AddMovie(movie)
if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	return
}
// Return success message
c.JSON(http.StatusCreated, gin.H{"message": "Movie Added Successfully", "result": result.InsertedID})

}


// function to get all movies 
func GetMovies(c *gin.Context){
	
	movies, err := queries.GetAllMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}

	c.JSON(http.StatusOK, movies)
}

// function to update a theater by id
func UpdateMovie(c *gin.Context) {
	movieId := c.Param("id")

	// Bind the JSON input to input struct
	var input schemas.Movie
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(fmt.Errorf("failed to bind request body to movie model: %v", err))
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Find the theater by ID
	_, err := queries.FindMovieById(movieId)
	if err != nil {
		if err.Error() == fmt.Sprintf("movie with id %s not found", movieId) {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	// Create an update struct
	updateMovie := schemas.Movie{
		MovieName: input.MovieName,
	
	Description: input.Description,
	Language: input.Language,
	ReleaseDate: input.ReleaseDate,
	Revenue: input.Revenue,
	Genre: input.Genre,
	Image: input.Image,
	TrailerId: input.TrailerId,
	}

	// Try to update the theater
	result, err := queries.UpdateMovieById(updateMovie, movieId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "Movie updated successfully", "result": result})
}
// function to get  a single movie details
func GetMovieById(c *gin.Context){
	id := c.Param("id")
	movie, err := queries.FindMovieById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movie"})
		return
	}

	c.JSON(http.StatusOK, movie)
}