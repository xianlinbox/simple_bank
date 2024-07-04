package api

import "github.com/go-playground/validator/v10"

var validateAccountID validator.Func = func(fl validator.FieldLevel) bool {
	return fl.Field().Int() > 0
}