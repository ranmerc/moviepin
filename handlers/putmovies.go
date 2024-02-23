package handlers

import (
	"moviepin/model"
	"moviepin/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Updates whole collection of movies.
func (mh MovieHandler) PutMoviesHandler(c *gin.Context) {
	var req model.MoviesRequestBody

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := mh.domain.ReplaceMovies(req.Movies)

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
