package handlers

import (
	"movie-management-service/apperror"
	"movie-management-service/domain"
	"movie-management-service/model"
	"movie-management-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Deletes a particular movie.
func (mh MovieHandler) DeleteMovieHandler(c *gin.Context) {
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

	err := mh.domain.DeleteMovie(req.ID)

	if err != nil {
		if err == domain.ErrNotExists {
			c.JSON(http.StatusNoContent, nil)
			return
		}

		c.JSON(http.StatusInternalServerError, model.DefaultResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
