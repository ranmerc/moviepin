package model

import "time"

type Movie struct {
	ID          string    `json:"ID" validate:"required,uuid"`
	Title       string    `json:"title" validate:"required"`
	ReleaseDate time.Time `json:"releaseDate" validate:"required"`
	Genre       string    `json:"genre" validate:"required,oneof=ACTION COMEDY DRAMA FANTASY HORROR SCI-FI THRILLER"`
	Director    string    `json:"director" validate:"required"`
	Description string    `json:"description" validate:"required,max=500"`
}

type MoviesResponse struct {
	Movies []*Movie `json:"movies"`
}

type MovieResponse struct {
	Movie *Movie `json:"movie"`
}

type MovieRequestUri struct {
	ID string `uri:"movieID" validate:"required,uuid"`
}

type MoviesRequestBody struct {
	Movies []Movie `json:"movies" validate:"required,dive"`
}

type MovieRequestBody struct {
	Movies Movie `json:"movie" validate:"required"`
}
