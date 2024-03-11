package apperror

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var (
	// errIsRequired is a custom error message for required fields.
	errIsRequired = errors.New("is required")

	// errIsUUID is a custom error message for UUID fields.
	errIsUUID = errors.New("should be an UUID")
)

var (
	// customErrors is a map of custom error messages for validation errors.
	customErrors = map[string]error{
		"movieID.required":     errIsRequired,
		"movieID.uuid":         errIsUUID,
		"username.required":    errIsRequired,
		"username.min":         errors.New("should be minimum 6 characters"),
		"password.required":    errIsRequired,
		"password.min":         errors.New("should be minimum 8 characters"),
		"ID.required":          errIsRequired,
		"ID.uuid":              errIsUUID,
		"title.required":       errIsRequired,
		"releaseDate.required": errIsRequired,
		"genre.required":       errIsRequired,
		"genre.oneof":          errors.New("should be one of ACTION, COMEDY, DRAMA, FANTASY, HORROR, SCI-FI, THRILLER"),
		"director.required":    errIsRequired,
		"description.required": errIsRequired,
		"description.max":      errors.New("should be maximum 500 characters"),
		"movies.required":      errIsRequired,
		"movies.gt":            errors.New("should have at least one movie"),
		"movie.required":       errIsRequired,
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
