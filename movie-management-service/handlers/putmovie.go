package handlers

import (
	"movie-management-service/domain"
	"movie-management-service/model"
	"movie-management-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Updates a particular movie.
func (mh MovieHandler) PutMovieHandler(c *gin.Context) {
	var req model.MovieRequestUri

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.DefaultResponse{
			Message: ErrInvalidID.Error(),
		})
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.ErrorLogger.Print(err)

		c.JSON(http.StatusBadRequest, model.DefaultResponse{
			Message: ErrInvalidID.Error(),
		})
		return
	}

	var movie model.Movie

	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, model.DefaultResponse{
			Message: ErrInvalidBody.Error(),
		})
		return
	}

	if err := utils.Validate.Struct(movie); err != nil {
		utils.ErrorLogger.Print(err)

		c.JSON(http.StatusBadRequest, model.DefaultResponse{
			Message: ErrInvalidBody.Error(),
		})
		return
	}

	err := mh.domain.UpdateMovie(req.ID, movie)

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
		Movie: &movie,
	})
}
