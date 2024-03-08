package handlers

import (
	"movie-management-service/apperror"
	"movie-management-service/domain"
	"movie-management-service/model"
	"movie-management-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handles user login.
func (mh MovieHandler) LoginHandler(c *gin.Context) {
	var req model.AuthRequest
	req.Username = c.PostForm("username")
	req.Password = c.PostForm("password")

	if err := utils.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, model.ValidationErrorResponse{
			Message: apperror.CustomValidationError(err),
		})
		return
	}

	if err := mh.domain.LoginUser(req.Username, req.Password); err != nil {
		if err == domain.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, model.DefaultResponse{
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.DefaultResponse{
			Message: err.Error(),
		})
		return
	}

	token, err := mh.tokenClient.GenerateToken(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DefaultResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.LoginResponse{
		Token: token,
	})
}
