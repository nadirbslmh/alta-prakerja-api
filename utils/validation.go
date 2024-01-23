package utils

import (
	"slices"

	"github.com/go-playground/validator/v10"
)

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
	case "courseTagValid":
		return "the course tag is invalid"
	case "scoreValid":
		return "the score value is invalid"
	case "oneof":
		return "the value of " + err.Field() + " must be tpm or uk"
	case "gte":
		return "the value of " + err.Field() + " must be greater than or equal 1"
	case "lte":
		return "the value of " + err.Field() + " must be less than or equal 999"
	default:
		return "validation error in " + err.Field()
	}
}

func RegisterSequenceValidator(validate *validator.Validate) {
	validate.RegisterValidation("sequenceValid", sequenceValidator)
}

func RegisterCourseTagValidator(validate *validator.Validate) {
	validate.RegisterValidation("courseTagValid", courseTagValidator)
}

func RegisterScoreValidator(validate *validator.Validate) {
	validate.RegisterValidation("scoreValid", taskScoreValidator)
}

func sequenceValidator(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(int)

	isValidSequence := (value >= 1 && value <= 10) || value == 999

	return isValidSequence
}

func courseTagValidator(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(string)

	isValidCourseTag := slices.Contains[[]string](COURSE_TAGS, value)

	return isValidCourseTag
}

func taskScoreValidator(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(int)

	scores := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	isValidScore := slices.Contains[[]int](scores, value)

	return isValidScore
}
