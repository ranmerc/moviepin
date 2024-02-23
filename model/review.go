package model

import "time"

type Review struct {
	ID           string    `json:"ID" validate:"required,uuid"`
	MovieID      string    `json:"movieID" validate:"required,uuid"`
	Rating       float32   `json:"rating" validate:"required,lte=5,gte=0"`
	ReviewText   string    `json:"reviewText" validate:"required,lte=500"`
	CreatedAtUTC time.Time `json:"createdAtUTC" validate:"required"`
	UpdatedAtUTC time.Time `json:"updatedAtUTC" validate:"required"`
}
