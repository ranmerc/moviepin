package handlers

import (
	"movie-management-service/domain"
	"movie-management-service/grpcclient"

	"github.com/gin-gonic/gin"
)

// Handler is the interface that defines the movie handler methods.
type Handler interface {
	HealthCheckHandler(c *gin.Context)
	GetMoviesHandler(c *gin.Context)
	GetMovieHandler(c *gin.Context)
	PostMoviesHandler(c *gin.Context)
	PatchMovieHandler(c *gin.Context)
	DeleteMovieHandler(c *gin.Context)
	PutMoviesHandler(c *gin.Context)
	PutMovieHandler(c *gin.Context)
	GetMovieRatingHandler(c *gin.Context)

	RegisterHandler(c *gin.Context)
	LoginHandler(c *gin.Context)
}

// MovieHandler implements the Handler interface.
type MovieHandler struct {
	domain      domain.Service
	tokenClient grpcclient.TokenServiceGRPCClient
}

func NewMovieHandler(domain domain.Service, tokenClient grpcclient.TokenServiceGRPCClient) *MovieHandler {
	return &MovieHandler{
		domain:      domain,
		tokenClient: tokenClient,
	}
}
