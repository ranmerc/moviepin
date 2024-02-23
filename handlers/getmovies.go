package handlers

import (
	"moviepin/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// ErrFailedToGetMovie is returned when failed to get movie.
	ErrFailedToGetMovie = "failed to get movie"
)

// Responds with all the movies.
func (mh MovieHandler) GetMoviesHandler(c *gin.Context) {
	movies, err := mh.domain.GetMovies()

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.MoviesGetResponse{
		Movies: movies,
	})
}
