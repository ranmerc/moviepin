package apperror

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var (
	customErrors = map[string]error{
		"movieID.required":     errors.New("is required"),
		"movieID.uuid":         errors.New("should be an UUID"),
		"username.required":    errors.New("is required"),
		"username.min":         errors.New("should be minimum 6 characters"),
		"password.required":    errors.New("is required"),
		"password.min":         errors.New("should be minimum 8 characters"),
		"ID.required":          errors.New("is required"),
		"ID.uuid":              errors.New("should be an UUID"),
		"title.required":       errors.New("is required"),
		"releaseDate.required": errors.New("is required"),
		"genre.required":       errors.New("is required"),
		"genre.oneof":          errors.New("should be one of ACTION, COMEDY, DRAMA, FANTASY, HORROR, SCI-FI, THRILLER"),
		"director.required":    errors.New("is required"),
		"description.required": errors.New("is required"),
		"description.max":      errors.New("should be maximum 500 characters"),
		"movies.required":      errors.New("is required"),
		"movies.gt":            errors.New("should have at least one movie"),
		"movie.required":       errors.New("is required"),
	}
)

// CustomValidationError converts validation and json marshal error into custom error type.
func CustomValidationError(err error) map[string]string {
	errs := make(map[string]string, 0)
	switch errTypes := err.(type) {
	case validator.ValidationErrors:
		for _, e := range errTypes {
			key := e.Field() + "." + e.Tag()

			if v, ok := customErrors[key]; ok {
				errs[e.Field()] = v.Error()
			} else {
				errs[e.Field()] = fmt.Sprintf("custom message is not available: %v", err)
			}
		}
		return errs
	case *json.UnmarshalTypeError:
		errs[errTypes.Field] = fmt.Sprintf("%v cannot be a %v", errTypes.Field, errTypes.Value)
		return errs
	}
	errs["unknown"] = fmt.Sprintf("unsupported custom error for: %v", err)
	return errs
}
