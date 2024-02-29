package model

type HeathResponse struct {
	Status   string `json:"status"` // "alive" or "dead"
	DBStatus bool   `json:"db"`
}
