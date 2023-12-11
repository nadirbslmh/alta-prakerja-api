package utils

import "github.com/go-playground/validator/v10"

func GetValidationErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return "the email is invalid"
	case "numeric":
		return "the " + err.Field() + " must be numeric"
	case "uri":
		return "the uri is invalid"
	case "sequenceValid":
		return "the sequence value must be from 1-10 or 999"
	default:
		return "validation error in " + err.Field()
	}
}

func RegisterSequenceValidator(validate *validator.Validate) {
	validate.RegisterValidation("sequenceValid", sequenceValidator)
}

func sequenceValidator(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(int)

	isValidSequence := (value >= 1 && value <= 10) || value == 999

	return isValidSequence
}
