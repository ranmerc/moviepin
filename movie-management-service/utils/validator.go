package utils

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		tags := []string{"json", "uri", "form"}
		for _, key := range tags {
			tag := fld.Tag.Get(key)
			name := strings.SplitN(tag, ",", 2)[0]
			if name == "-" {
				return ""
			} else if len(name) != 0 {
				return name
			}
		}
		return ""
	})
}
