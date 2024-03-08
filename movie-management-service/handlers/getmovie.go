package handlers

import (
	"movie-management-service/apperror"
	"movie-management-service/domain"
	"movie-management-service/model"
	"movie-management-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Responds with details of particular movie.
func (mh MovieHandler) GetMovieHandler(c *gin.Context) {
	var req model.MovieRequestUri

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ValidationErrorResponse{
			Message: apperror.CustomValidationError(err),
		})
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.ErrorLogger.Print(err)

		c.JSON(http.StatusBadRequest, model.ValidationErrorResponse{
			Message: apperror.CustomValidationError(err),
		})
		return
	}

	movie, err := mh.domain.GetMovie(req.ID)

	if err != nil {
		if err == domain.ErrNotExists {
			c.JSON(http.StatusNotFound, model.DefaultResponse{
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.DefaultResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.MovieResponse{
		Movie: movie,
	})
}
