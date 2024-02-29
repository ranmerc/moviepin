package handlers

import (
	"errors"
	"movie-management-service/model"
	"movie-management-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// ErrFailedAdd is returned when all movies failed to be added.
	ErrFailedAdd = errors.New("failed to add movies")
)

// Adds list of movies sent in request.
func (mh MovieHandler) PostMoviesHandler(c *gin.Context) {
	var req model.MoviesRequestBody

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Message: ErrInvalidBody.Error(),
		})
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.ErrorLogger.Print(err)

		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Message: ErrInvalidBody.Error(),
		})
		return
	}

	// If no movies are sent in request.
	if len(req.Movies) == 0 {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Message: ErrEmptyBody.Error(),
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

	// When all movies failed to be added.
	if len(response.FailedMovies) == len(req.Movies) {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: ErrFailedAdd.Error(),
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
