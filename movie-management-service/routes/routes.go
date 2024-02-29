package routes

import (
	"movie-management-service/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route Structure of new routes
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc gin.HandlerFunc
}

// Routes Array of all available routes
type Routes []Route

// NewRoutes returns the list of available routes
func NewRoutes(mh *handlers.MovieHandler) Routes {
	return Routes{
		Route{
			Name:        "Health",
			Method:      http.MethodGet,
			Pattern:     "/health",
			HandlerFunc: mh.HealthCheckHandler,
		},
		Route{
			Name:        "Get Movies",
			Method:      http.MethodGet,
			Pattern:     "/movies",
			HandlerFunc: mh.GetMoviesHandler,
		},
		Route{
			Name:        "Get Movie",
			Method:      http.MethodGet,
			Pattern:     "/movies/:movieID",
			HandlerFunc: mh.GetMovieHandler,
		},
		Route{
			Name:        "Post Movies",
			Method:      http.MethodPost,
			Pattern:     "/movies",
			HandlerFunc: mh.PostMoviesHandler,
		},
		Route{
			Name:        "Patch Movie",
			Method:      http.MethodPatch,
			Pattern:     "/movies/:movieID",
			HandlerFunc: mh.PatchMovieHandler,
		},
		Route{
			Name:        "Delete Movie",
			Method:      http.MethodDelete,
			Pattern:     "/movies/:movieID",
			HandlerFunc: mh.DeleteMovieHandler,
		},
		Route{
			Name:        "Put Movie",
			Method:      http.MethodPut,
			Pattern:     "/movies/:movieID",
			HandlerFunc: mh.PutMovieHandler,
		},
		Route{
			Name:        "Put Movies",
			Method:      http.MethodPut,
			Pattern:     "/movies",
			HandlerFunc: mh.PutMoviesHandler,
		},
		Route{
			Name:        "Get Movie Rating",
			Method:      http.MethodGet,
			Pattern:     "/movies/:movieID/rating",
			HandlerFunc: mh.GetMovieRatingHandler,
		},
	}
}

// AttachRoutes Attaches routes to the provided server
func AttachRoutes(server *gin.Engine, routes Routes) {
	for _, route := range routes {
		server.Handle(route.Method, route.Pattern, route.HandlerFunc)
	}
}
