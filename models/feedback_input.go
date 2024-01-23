package models

import (
	"gugcp/utils"

	"github.com/go-playground/validator/v10"
)

type FeedbackInput struct {
	Notes string `json:"notes" validate:"required"`
	Score int    `json:"score" validate:"required,scoreValid"`
}

func (r *FeedbackInput) Validate() []*ValidationErrorResponse {
	var errors []*ValidationErrorResponse

	validate := validator.New()
	utils.RegisterScoreValidator(validate)
	err := validate.Struct(r)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationErrorResponse
			element.ErrorMessage = utils.GetValidationErrorMessage(err)
			element.Field = err.Field()
			errors = append(errors, &element)
		}
	}

	return errors
}
