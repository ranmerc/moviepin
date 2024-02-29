package routes

import (
	"net/http"
	"testing"

	"movie-management-service/handlers"
	"movie-management-service/mock"
)

func TestNewRoutes(t *testing.T) {
	mockService := &mock.ServiceMock{}

	movieHandler := handlers.NewMovieHandler(mockService)
	got := NewRoutes(movieHandler)

	expectedRoutes := []struct {
		Pattern string
		Method  string
	}{
		{"/health", http.MethodGet},
		{"/movies", http.MethodGet},
		{"/movies/:movieID", http.MethodGet},
		{"/movies", http.MethodPost},
		{"/movies/:movieID", http.MethodPatch},
		{"/movies/:movieID", http.MethodDelete},
		{"/movies/:movieID", http.MethodPut},
		{"/movies", http.MethodPut},
		{"/movies/:movieID/rating", http.MethodGet},
	}

	t.Run("all routes are present", func(t *testing.T) {
		for i := 0; i < min(len(got), len(expectedRoutes)); i++ {
			got := got[i]
			want := expectedRoutes[i]

			if got.Pattern != want.Pattern {
				t.Errorf("pattern expectations mismatched: \n want: %v \n got: %v", got.Pattern, want.Pattern)
			}

			if got.Method != want.Method {
				t.Errorf("method expectations mismatched: \n want: %v \n got: %v", got.Method, want.Method)
			}
		}
	})
}
