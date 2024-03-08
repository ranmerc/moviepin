package model

type DefaultResponse struct {
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Message []map[string]string `json:"message"`
}
