package handlers

import (
	"errors"
	"movie-management-service/apperror"
	"movie-management-service/domain"
	"movie-management-service/model"
	"movie-management-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterHandler handles user registrations.
func (mh MovieHandler) RegisterHandler(c *gin.Context) {
	var req model.AuthRequest

	if err := c.ShouldBind(&req); err != nil {
		utils.ErrorLogger.Print(err)

		c.JSON(http.StatusBadRequest, model.ValidationErrorResponse{
			Message: apperror.CustomValidationError(err),
		})
		return
	}

	if err := mh.domain.RegisterUser(req.Username, req.Password); err != nil {
		if errors.Is(err, domain.ErrUsernameExists) {
			c.JSON(http.StatusConflict, model.DefaultResponse{
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.DefaultResponse{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, model.DefaultResponse{
			Message: "user registered successfully",
		})
	}
}
