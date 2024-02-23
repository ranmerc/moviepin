package handlers

import (
	"moviepin/domain"
	"moviepin/model"
	"moviepin/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Responds with movie details along with its rating.
func (mh MovieHandler) GetMovieRatingHandler(c *gin.Context) {
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

	review, err := mh.domain.GetMovieRating(req.ID)

	if err == domain.ErrNotExists {
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, review)
}
