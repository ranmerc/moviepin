package model

import "time"

type MovieReview struct {
	ID          string    `json:"ID" binding:"required,uuid"`
	Title       string    `json:"title" binding:"required"`
	ReleaseDate time.Time `json:"releaseDate" binding:"required"`
	Genre       string    `json:"genre" binding:"required,oneof=ACTION COMEDY DRAMA FANTASY HORROR SCI-FI THRILLER"`
	Director    string    `json:"director" binding:"required"`
	Description string    `json:"description" binding:"required,max=500"`
	Rating      float32   `json:"rating" binding:"required,gte=0,lte=5"`
}
