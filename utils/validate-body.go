package utils

import (
	"github.com/go-playground/validator/v10"
)

func CheckErrors(err error) []map[string]string {
	var errors []map[string]string

	for _, e := range err.(validator.ValidationErrors) {
		errEl := make(map[string]string, 0)
		errEl["field"] = e.Field()
		errEl["message"] = e.Tag()
		errors = append(errors, errEl)
	}

	return errors
}
