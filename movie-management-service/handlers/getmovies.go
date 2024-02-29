package handlers

import (
	"movie-management-service/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Responds with details all the movies.
func (mh MovieHandler) GetMoviesHandler(c *gin.Context) {
	movies, err := mh.domain.GetMovies()

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DefaultResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.MoviesResponse{
		Movies: movies,
	})
}
