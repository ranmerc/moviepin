package handlers

import (
	"moviepin/domain"
	"moviepin/model"
	"moviepin/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Responds with details of particular movie.
func (mh MovieHandler) GetMovieHandler(c *gin.Context) {
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

	movie, err := mh.domain.GetMovie(req.ID)

	if err == domain.ErrNotExists {
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Message: "movie not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, movie)
}
