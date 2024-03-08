package handlers

import (
	"errors"
	"movie-management-service/apperror"
	"movie-management-service/domain"
	"movie-management-service/model"
	"movie-management-service/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// ErrAssertTitle represents an error when a field title failed be asserted.
	ErrAssertTitle = errors.New("failed to assert type for field title")

	// ErrAssertReleaseDate represents an error when a field releaseDate failed be asserted.
	ErrAssertReleaseDate = errors.New("failed to assert type for field releaseDate")

	// ErrAssertGenre represents an error when a field genre failed be asserted.
	ErrAssertGenre = errors.New("failed to assert type for field genre")

	// ErrAssertDirector represents an error when a field director failed be asserted.
	ErrAssertDirector = errors.New("failed to assert type for field director")

	// ErrAssertDescription represents an error when a field description failed be asserted.
	ErrAssertDescription = errors.New("failed to assert type for field description")
)

// Updates a particular movie.
func (mh MovieHandler) PatchMovieHandler(c *gin.Context) {
	var req model.MovieRequestUri

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.ErrorLogger.Print(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	// Represents partial movie details.
	var partialMovie map[string]interface{}

	if err := c.ShouldBindJSON(&partialMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apperror.CustomValidationError(err),
		})
		return
	}

	existingMovie, err := mh.domain.GetMovie(req.ID)

	if err != nil {
		if err == domain.ErrNotExists {
			c.JSON(http.StatusNotFound, model.DefaultResponse{
				Message: "movie not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.DefaultResponse{
			Message: err.Error(),
		})
		return
	}

	// Update fields of existing movie with partial request.
	for key, value := range partialMovie {
		switch key {
		case "title":
			if title, ok := value.(string); ok {
				existingMovie.Title = title
			} else {
				c.JSON(http.StatusBadRequest, model.DefaultResponse{
					Message: ErrAssertTitle.Error(),
				})
				return
			}
		case "releaseDate":
			if releaseDate, ok := value.(string); ok {
				if parsedReleaseDate, err := time.Parse(time.RFC3339, releaseDate); err == nil {
					existingMovie.ReleaseDate = parsedReleaseDate
				} else {
					c.JSON(http.StatusBadRequest, model.DefaultResponse{
						Message: ErrAssertReleaseDate.Error(),
					})
					return
				}
			} else {
				c.JSON(http.StatusBadRequest, model.DefaultResponse{
					Message: ErrAssertReleaseDate.Error(),
				})
				return
			}
		case "genre":
			if genre, ok := value.(string); ok {
				existingMovie.Genre = genre
			} else {
				c.JSON(http.StatusBadRequest, model.DefaultResponse{
					Message: ErrAssertGenre.Error(),
				})
				return
			}
		case "director":
			if director, ok := value.(string); ok {
				existingMovie.Director = director
			} else {
				c.JSON(http.StatusBadRequest, model.DefaultResponse{
					Message: ErrAssertDirector.Error(),
				})
				return
			}
		case "description":
			if description, ok := value.(string); ok {
				existingMovie.Description = description
			} else {
				c.JSON(http.StatusBadRequest, model.DefaultResponse{
					Message: ErrAssertDescription.Error(),
				})
				return
			}
		}
	}

	// Validate updated movie.
	if err := utils.Validate.Struct(existingMovie); err != nil {
		c.JSON(http.StatusBadRequest, model.DefaultResponse{
			Message: err.Error(),
		})
		return
	}

	err = mh.domain.UpdateMovie(req.ID, *existingMovie)

	if err != nil {
		if err == domain.ErrNotExists {
			c.JSON(http.StatusNotFound, model.DefaultResponse{
				Message: "movie not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.DefaultResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.MovieResponse{
		Movie: existingMovie,
	})
}
