package handlers

import (
	"moviepin/model"
	"moviepin/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Adds list of movies sent in request.
func (mh MovieHandler) PostMoviesHandler(c *gin.Context) {
	var req model.MoviesRequestBody

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	type MovieStatus struct {
		Status string
		Movie  model.Movie
	}

	status := make(chan MovieStatus)

	for _, movie := range req.Movies {
		go func(movie model.Movie) {
			if err := mh.domain.AddMovie(movie); err != nil {
				status <- MovieStatus{
					Status: "failed",
					Movie:  movie,
				}
				return
			}

			status <- MovieStatus{
				Status: "success",
				Movie:  movie,
			}
		}(movie)
	}

	type postMoviesResponse struct {
		AddedMovies  []*model.Movie `json:"addedMovies,omitempty"`
		FailedMovies []*model.Movie `json:"failedMovies,omitempty"`
	}

	response := &postMoviesResponse{}

	for range req.Movies {
		result := <-status

		switch result.Status {
		case "failed":
			response.FailedMovies = append(response.FailedMovies, &result.Movie)
		case "success":
			response.AddedMovies = append(response.AddedMovies, &result.Movie)
		}
	}

	if len(response.FailedMovies) == len(req.Movies) {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: "failed to add movies",
		})
		return
	}

	if len(response.FailedMovies) > 0 {
		c.JSON(http.StatusMultiStatus, response)
		return
	} else {
		c.JSON(http.StatusCreated, response)
		return
	}
}
