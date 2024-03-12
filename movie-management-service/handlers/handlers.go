package handlers

import (
	"movie-management-service/domain"
	"movie-management-service/grpcclient"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

func init() {
	// Registers a tag name function to enable use of tag name as field name in custom validation error message.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		tags := []string{"json", "uri", "form"}
		for _, key := range tags {
			tag := fld.Tag.Get(key)
			name := strings.SplitN(tag, ",", 2)[0]
			if name == "-" {
				return ""
			} else if len(name) != 0 {
				return name
			}
		}
		return ""
	})
}

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
