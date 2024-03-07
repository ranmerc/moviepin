package routes

import (
	"movie-management-service/handlers"
	"movie-management-service/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route Structure of new routes
type Route struct {
	Name        string
	Method      string
	Pattern     string
	Protected   bool
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
			Protected:   false,
			HandlerFunc: mh.HealthCheckHandler,
		},
		Route{
			Name:        "Get Movies",
			Method:      http.MethodGet,
			Pattern:     "/movies",
			Protected:   false,
			HandlerFunc: mh.GetMoviesHandler,
		},
		Route{
			Name:        "Get Movie",
			Method:      http.MethodGet,
			Pattern:     "/movies/:movieID",
			Protected:   false,
			HandlerFunc: mh.GetMovieHandler,
		},
		Route{
			Name:        "Post Movies",
			Method:      http.MethodPost,
			Pattern:     "/movies",
			Protected:   true,
			HandlerFunc: mh.PostMoviesHandler,
		},
		Route{
			Name:        "Patch Movie",
			Method:      http.MethodPatch,
			Pattern:     "/movies/:movieID",
			Protected:   true,
			HandlerFunc: mh.PatchMovieHandler,
		},
		Route{
			Name:        "Delete Movie",
			Method:      http.MethodDelete,
			Pattern:     "/movies/:movieID",
			Protected:   true,
			HandlerFunc: mh.DeleteMovieHandler,
		},
		Route{
			Name:        "Put Movie",
			Method:      http.MethodPut,
			Pattern:     "/movies/:movieID",
			Protected:   true,
			HandlerFunc: mh.PutMovieHandler,
		},
		Route{
			Name:        "Put Movies",
			Method:      http.MethodPut,
			Pattern:     "/movies",
			Protected:   true,
			HandlerFunc: mh.PutMoviesHandler,
		},
		Route{
			Name:        "Get Movie Rating",
			Method:      http.MethodGet,
			Pattern:     "/movies/:movieID/rating",
			Protected:   false,
			HandlerFunc: mh.GetMovieRatingHandler,
		},
		Route{
			Name:        "Login",
			Method:      http.MethodPost,
			Pattern:     "/login",
			Protected:   false,
			HandlerFunc: mh.LoginHandler,
		},
		Route{
			Name:        "Register",
			Method:      http.MethodPost,
			Pattern:     "/register",
			Protected:   false,
			HandlerFunc: mh.RegisterHandler,
		},
	}
}

// AttachRoutes Attaches routes to the provided server
func AttachRoutes(server *gin.Engine, auth middleware.Auth, routes Routes) {
	for _, route := range routes {
		if route.Protected {
			server.Handle(route.Method, route.Pattern, auth.Middleware(), route.HandlerFunc)
		} else {
			server.Handle(route.Method, route.Pattern, route.HandlerFunc)
		}
	}
}
