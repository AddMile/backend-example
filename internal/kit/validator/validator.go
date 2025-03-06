package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	val *validator.Validate
}

func New() *Validator {
	val := validator.New(validator.WithRequiredStructEnabled())
	val.RegisterValidation("custom_time", time)
	val.RegisterValidation("custom_money", money)
	val.RegisterValidation("custom_slug", slug)

	return &Validator{val: val}
}

func (v *Validator) Validate(s any) error {
	return v.val.Struct(s)
}

// time checks if the field is in "HH:MM" format.
func time(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	match, err := regexp.MatchString(`^(?:[01]\d|2[0-3]):[0-5]\d$`, field)
	if err != nil {
		return false
	}

	return match
}

func slug(fl validator.FieldLevel) bool {
	match, err := regexp.MatchString(`^[a-z0-9]+(?:-[a-z0-9]+)*$`, fl.Field().String())
	if err != nil {
		return false
	}

	return match
}

func money(fl validator.FieldLevel) bool {
	match, err := regexp.MatchString(`^\d+(\.\d{2})?$`, fl.Field().String())
	if err != nil {
		return false
	}

	return match
}
