package handlers

import (
	"moviepin/domain"
	"moviepin/model"
	"moviepin/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Updates a particular movie.
func (mh MovieHandler) PatchMovieHandler(c *gin.Context) {
	var req model.MovieRequestUri

	if err := c.ShouldBindUri(&req); err != nil {
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

	// Represents partial movie details.
	var partialMovie map[string]interface{}

	if err := c.ShouldBindJSON(&partialMovie); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	existingMovie, err := mh.domain.GetMovie(req.ID)

	if err != nil {
		if err == domain.ErrNotExists {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Message: "movie not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
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
				c.JSON(http.StatusBadRequest, model.ErrorResponse{
					Message: "failed to assert type for field title",
				})
				return
			}
		case "releaseDate":
			if releaseDate, ok := value.(string); ok {
				if parsedReleaseDate, err := time.Parse(time.RFC3339, releaseDate); err == nil {
					existingMovie.ReleaseDate = parsedReleaseDate
				} else {
					c.JSON(http.StatusBadRequest, model.ErrorResponse{
						Message: "failed to assert type for field releaseDate",
					})
					return
				}
			} else {
				c.JSON(http.StatusBadRequest, model.ErrorResponse{
					Message: "failed to assert type for field releaseDate",
				})
				return
			}
		case "genre":
			if genre, ok := value.(string); ok {
				existingMovie.Genre = genre
			} else {
				c.JSON(http.StatusBadRequest, model.ErrorResponse{
					Message: "failed to assert type for field genre",
				})
				return
			}
		case "director":
			if director, ok := value.(string); ok {
				existingMovie.Director = director
			} else {
				c.JSON(http.StatusBadRequest, model.ErrorResponse{
					Message: "failed to assert type for field director",
				})
				return
			}
		case "description":
			if description, ok := value.(string); ok {
				existingMovie.Description = description
			} else {
				c.JSON(http.StatusBadRequest, model.ErrorResponse{
					Message: "failed to assert type for field description",
				})
				return
			}
		}
	}

	// Validate updated movie.
	if err := utils.Validate.Struct(existingMovie); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	err = mh.domain.UpdateMovie(req.ID, *existingMovie)

	if err != nil {
		if err == domain.ErrNotExists {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Message: "movie not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, existingMovie)
}
