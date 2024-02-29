package domain

import (
	"database/sql"
	"movie-management-service/model"
)

// Service is the interface that defines the movie service methods.
type Service interface {
	DBStatus() (bool, error)
	GetMovies() ([]*model.Movie, error)
	GetMovie(id string) (*model.Movie, error)
	AddMovie(movie model.Movie) error
	UpdateMovie(id string, movie model.Movie) error
	DeleteMovie(id string) error
	ReplaceMovies(movies []model.Movie) error
	GetMovieRating(id string) (*model.MovieReview, error)
}

// MovieService implements the Service interface.
type MovieService struct {
	db *sql.DB
}

func NewMovieService(db *sql.DB) *MovieService {
	return &MovieService{db}
}
