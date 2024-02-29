package model

import "time"

type MovieReview struct {
	ID          string    `json:"ID" validate:"required,uuid"`
	Title       string    `json:"title" validate:"required"`
	ReleaseDate time.Time `json:"releaseDate" validate:"required"`
	Genre       string    `json:"genre" validate:"required,oneof=ACTION COMEDY DRAMA FANTASY HORROR SCI-FI THRILLER"`
	Director    string    `json:"director" validate:"required"`
	Description string    `json:"description" validate:"required,max=500"`
	Rating      float32   `json:"rating" validate:"required,gte=0,lte=5"`
}
