package handlers

import (
	"errors"
	"movie-management-service/model"
	"movie-management-service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// ErrInvalidCredentials is the error message when the credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Handles user login.
func (mh MovieHandler) LoginHandler(c *gin.Context) {
	var req model.AuthRequest
	req.Username = c.PostForm("username")
	req.Password = c.PostForm("password")

	if err := utils.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, model.DefaultResponse{
			Message: err.Error(),
		})
		return
	}

	if err := mh.domain.LoginUser(req.Username, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, model.DefaultResponse{
			Message: ErrInvalidCredentials.Error(),
		})
	} else {
		c.JSON(http.StatusOK, model.DefaultResponse{
			Message: "login successful",
		})
	}
}
