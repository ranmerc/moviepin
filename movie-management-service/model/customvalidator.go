package model

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	// Registers a tag name function to enable use of tag name as field name in custom validation error message.
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
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
}
