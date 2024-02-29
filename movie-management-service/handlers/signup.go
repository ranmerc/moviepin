package handlers

import (
	"movie-management-service/domain"
	"movie-management-service/model"
	"movie-management-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handles user sign-up
func (mh MovieHandler) SignUpHandler(c *gin.Context) {
	var req model.AuthRequest
	req.Username = c.PostForm("username")
	req.Password = c.PostForm("password")

	if err := utils.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if err := mh.domain.RegisterUser(req.Username, req.Password); err != nil {
		if err == domain.ErrUsernameExists {
			c.JSON(http.StatusConflict, model.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "sign-up successful"})
	}
}
