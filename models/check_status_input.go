package models

import (
	"gugcp/utils"

	"github.com/go-playground/validator/v10"
)

type CheckStatusInput struct {
	RedeemCode string `json:"redeem_code" validate:"required"`
	State      string `json:"state" validate:"required"`
	Sequence   int    `json:"sequence" validate:"required,numeric,sequenceValid"`
}

func (r *CheckStatusInput) Validate() []*ValidationErrorResponse {
	var errors []*ValidationErrorResponse

	validate := validator.New()
	utils.RegisterSequenceValidator(validate)
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
