package mock

import (
	"errors"
	"movie-management-service/domain"
	"movie-management-service/model"
	"time"
)

var (
	// Movie is a mock movie.
	Movie = model.Movie{
		ID:          "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		Title:       "The Shawshank Redemption",
		ReleaseDate: time.Date(1994, time.September, 23, 0, 0, 0, 0, time.UTC),
		Genre:       "DRAMA",
		Director:    "Frank Darabont",
		Description: "Prisoners",
	}

	// MovieReview is a mock movie review.
	MovieReview = model.MovieReview{
		ID:          Movie.ID,
		Title:       Movie.Title,
		ReleaseDate: Movie.ReleaseDate,
		Genre:       Movie.Genre,
		Director:    Movie.Director,
		Description: Movie.Description,
		Rating:      3.5,
	}

	// MovieIDToFail is a mock movie id designed to fail for non-existence.
	MovieIDToFail = "6ba7b810-9dad-11d1-80b4-00c04fd430c9"
)

type ErrMock int

const (
	// DBStatusError is a mock error for DBStatus.
	DBStatusError ErrMock = iota

	// GetMoviesError is a mock error for GetMovies.
	GetMoviesError

	// GetMovieError is a mock error for GetMovie.
	GetMovieError
	// GetMovieNotExistsError is a mock error for GetMovie when movie id does not exist.
	GetMovieNotExistsError

	// AddMovieError is a mock error for AddMovie.
	AddMovieError

	// UpdateMovieError is a mock error for UpdateMovie.
	UpdateMovieError
	// UpdateMovieNotExistsError is a mock error for UpdateMovie when movie id does not exist.
	UpdateMovieNotExistsError

	// DeleteMovieError is a mock error for DeleteMovie.
	DeleteMovieError
	// DeleteMovieNotExistsError is a mock error for DeleteMovie when movie id does not exist.
	DeleteMovieNotExistsError

	// ReplaceMoviesError is a mock error for ReplaceMovies.
	ReplaceMoviesError

	// GetMovieRatingError is a mock error for GetMovieRating.
	GetMovieRatingError
	// GetMovieRatingNotExistsError is a mock error for GetMovieRating when movie id does not exist.
	GetMovieRatingNotExistsError

	// When there is no error. This is the default value.
	OK
)

// ServiceMock is a mock of the Service interface.
type ServiceMock struct {
	Err ErrMock
}

func NewServiceMock() ServiceMock {
	return ServiceMock{
		Err: OK,
	}
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
	if s.Err == GetMovieNotExistsError {
		return nil, domain.ErrNotExists
	}

	if s.Err == GetMovieError {
		return nil, domain.ErrFailedGetMovie
	}

	return &Movie, nil
}

// GetMovies returns a slice of all movies present.
func (s ServiceMock) GetMovies() ([]*model.Movie, error) {
	if s.Err == GetMoviesError {
		return nil, domain.ErrFailedGetMovies
	}

	return []*model.Movie{&Movie}, nil
}

// AddMovie adds a movie to the database.
func (s ServiceMock) AddMovie(movie model.Movie) error {
	if movie.ID == MovieIDToFail {
		return domain.ErrFailedAdd
	}

	if s.Err == AddMovieError {
		return domain.ErrFailedAdd
	}

	return nil
}

// UpdateMovie updates a movie in the database.
func (s ServiceMock) UpdateMovie(id string, movie model.Movie) error {
	if s.Err == UpdateMovieNotExistsError {
		return domain.ErrNotExists
	}

	if s.Err == UpdateMovieError {
		return domain.ErrFailedUpdate
	}

	return nil
}

// DeleteMovie deletes a movie from the database.
func (s ServiceMock) DeleteMovie(id string) error {
	if s.Err == DeleteMovieNotExistsError {
		return domain.ErrNotExists
	}

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
	if s.Err == GetMovieRatingNotExistsError {
		return nil, domain.ErrNotExists
	}

	if s.Err == GetMovieRatingError {
		return nil, domain.ErrFailedGetRating
	}

	return &MovieReview, nil
}

func (s ServiceMock) RegisterUser(username, password string) error {
	return nil
}

func (s ServiceMock) LoginUser(username, password string) error {
	return nil
}
