package mock

import (
	"errors"
	"moviepin/domain"
	"moviepin/model"
	"time"
)

var (
	Movie = model.Movie{
		ID:          "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		Title:       "The Shawshank Redemption",
		ReleaseDate: time.Date(1994, time.September, 23, 0, 0, 0, 0, time.UTC),
		Genre:       "Drama",
		Director:    "Frank Darabont",
		Description: "Prisoners",
	}

	MovieReview = model.MovieReview{
		ID:          Movie.ID,
		Title:       Movie.Title,
		ReleaseDate: Movie.ReleaseDate,
		Genre:       Movie.Genre,
		Director:    Movie.Director,
		Description: Movie.Description,
		Rating:      3.5,
	}
)

type ErrMock int

const (
	DBStatusError ErrMock = iota
	NotExistsError
	GetMoviesError
	GetMovieError
	AddMovieError
	UpdateMovieError
	DeleteMovieError
	ReplaceMoviesError
	GetMovieRatingError
	OK
)

// ServiceMock is a mock of the Service interface.
type ServiceMock struct {
	Err ErrMock
}

func NewServiceMock() ServiceMock {
	return ServiceMock{}
}

// DBStatus is a mock implementation of DBStatus
func (s *ServiceMock) DBStatus() (bool, error) {
	if s.Err == DBStatusError {
		return false, errors.New("mock DB connection error")
	}

	return true, nil
}

// GetMovie returns a movie by its id.
func (s ServiceMock) GetMovie(id string) (*model.Movie, error) {
	if s.Err == NotExistsError {
		return nil, domain.ErrNotExists
	}

	if s.Err == GetMovieError {
		return nil, domain.ErrFailedGetMovie
	}

	return &Movie, nil
}

// GetMovies returns a slice of all movies present.
func (s ServiceMock) GetMovies() ([]*model.Movie, error) {
	if s.Err == NotExistsError {
		return nil, domain.ErrNotExists
	}

	if s.Err == GetMoviesError {
		return nil, domain.ErrFailedGetMovies
	}

	return []*model.Movie{&Movie}, nil
}

// AddMovie adds a movie to the database.
func (s ServiceMock) AddMovie(movie model.Movie) error {
	if s.Err == AddMovieError {
		return domain.ErrFailedAdd
	}

	return nil
}

// UpdateMovie updates a movie in the database.
func (s ServiceMock) UpdateMovie(id string, movie model.Movie) error {
	if s.Err == UpdateMovieError {
		return domain.ErrFailedUpdate
	}

	return nil
}

// DeleteMovie deletes a movie from the database.
func (s ServiceMock) DeleteMovie(id string) error {
	if s.Err == DeleteMovieError {
		return domain.ErrFailedDelete
	}

	return nil
}

// ReplaceMovies replaces all movies in the database.
func (s ServiceMock) ReplaceMovies(movies []model.Movie) error {
	if s.Err == ReplaceMoviesError {
		return domain.ErrFailedReplace
	}

	return nil
}

// GetMovieRating returns a movie rating by its id.
func (s ServiceMock) GetMovieRating(id string) (*model.MovieReview, error) {
	if s.Err == NotExistsError {
		return nil, domain.ErrNotExists
	}

	if s.Err == GetMovieRatingError {
		return nil, domain.ErrFailedGetRating
	}

	return &MovieReview, nil
}
