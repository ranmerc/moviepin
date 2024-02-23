package routes

import (
	"moviepin/handlers"
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
			"Health",
			http.MethodGet,
			"/health",
			mh.HealthCheckHandler,
		},
		Route{
			"Get Movies",
			http.MethodGet,
			"/movies",
			mh.GetMoviesHandler,
		},
		Route{
			"Get Movie",
			http.MethodGet,
			"/movies/:movieID",
			mh.GetMovieHandler,
		},
		Route{
			"Post Movies",
			http.MethodPost,
			"/movies",
			mh.PostMoviesHandler,
		},
		Route{
			"Patch Movie",
			http.MethodPatch,
			"/movies/:movieID",
			mh.PatchMovieHandler,
		},
		Route{
			"Delete Movie",
			http.MethodDelete,
			"/movies/:movieID",
			mh.DeleteMovieHandler,
		},
		Route{
			"Put Movie",
			http.MethodPut,
			"/movies/:movieID",
			mh.PutMovieHandler,
		},
		Route{
			"Put Movies",
			http.MethodPut,
			"/movies",
			mh.PutMoviesHandler,
		},
		Route{
			"Get Movie Rating",
			http.MethodGet,
			"/movies/:movieID/rating",
			mh.GetMovieRatingHandler,
		},
	}
}

// AttachRoutes Attaches routes to the provided server
func AttachRoutes(server *gin.Engine, routes Routes) {
	for _, route := range routes {
		server.Handle(route.Method, route.Pattern, route.HandlerFunc)
	}
}
