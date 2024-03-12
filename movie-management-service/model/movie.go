package model

import "time"

type Movie struct {
	ID          string    `json:"ID" binding:"required,uuid"`
	Title       string    `json:"title" binding:"required"`
	ReleaseDate time.Time `json:"releaseDate" binding:"required"`
	Genre       string    `json:"genre" binding:"required,oneof=ACTION COMEDY DRAMA FANTASY HORROR SCI-FI THRILLER"`
	Director    string    `json:"director" binding:"required"`
	Description string    `json:"description" binding:"required,max=500"`
}

type MoviesResponse struct {
	Movies []*Movie `json:"movies"`
}

type MovieResponse struct {
	Movie *Movie `json:"movie"`
}

type MovieRequestUri struct {
	ID string `uri:"movieID" binding:"required,uuid"`
}

type MoviesRequestBody struct {
	Movies []Movie `json:"movies" binding:"required,gt=0,dive,required"`
}

type MovieRequestBody struct {
	Movie Movie `json:"movie" binding:"required"`
}
