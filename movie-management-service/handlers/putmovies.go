package handlers

import (
	"movie-management-service/model"
	"movie-management-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Updates whole collection of movies.
func (mh MovieHandler) PutMoviesHandler(c *gin.Context) {
	var req model.MoviesRequestBody

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.DefaultResponse{
			Message: ErrInvalidBody.Error(),
		})
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.ErrorLogger.Print(err)

		c.JSON(http.StatusBadRequest, model.DefaultResponse{
			Message: ErrInvalidBody.Error(),
		})
		return
	}

	// If no movies are sent in request.
	if len(req.Movies) == 0 {
		c.JSON(http.StatusBadRequest, model.DefaultResponse{
			Message: ErrEmptyBody.Error(),
		})
		return
	}

	err := mh.domain.ReplaceMovies(req.Movies)

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DefaultResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
