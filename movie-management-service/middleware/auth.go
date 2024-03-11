package middleware

import (
	"errors"
	"movie-management-service/grpcclient"
	"movie-management-service/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	// ErrMissingAuthHeader is used when authorization header is missing
	ErrMissingAuthHeader = errors.New("authorization header is required")

	// ErrInvalidAuthHeader is used when authorization header is invalid
	ErrInvalidAuthHeader = errors.New("invalid authorization header")

	// ErrInvalidToken is used when token is invalid
	ErrInvalidToken = errors.New("invalid token")
)

type Auth interface {
	Middleware() gin.HandlerFunc
}

type AuthMiddleware struct {
	tokenClient grpcclient.TokenServiceGRPCClient
}

func NewAuthMiddleware(tokenClient grpcclient.TokenServiceGRPCClient) *AuthMiddleware {
	return &AuthMiddleware{
		tokenClient: tokenClient,
	}
}

func (am *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if len(authHeader) == 0 {
			c.JSON(http.StatusUnauthorized, model.DefaultResponse{
				Message: ErrMissingAuthHeader.Error(),
			})
			c.Abort()
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")

		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, model.DefaultResponse{
				Message: ErrInvalidAuthHeader.Error(),
			})
			c.Abort()
			return
		}

		token := authHeaderParts[1]

		if valid, err := am.tokenClient.VerifyToken(token); err != nil || !valid {
			c.JSON(http.StatusUnauthorized, model.DefaultResponse{
				Message: ErrInvalidToken.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
