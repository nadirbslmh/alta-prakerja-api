package models

import (
	"gugcp/utils"

	"github.com/go-playground/validator/v10"
)

type RedeemInput struct {
	UserID     int    `json:"user_id" validate:"required,numeric"`
	State      string `json:"state" validate:"required"`
	RedeemCode string `json:"redeem_code" validate:"required"`
	Sequence   int    `json:"sequence" validate:"required,numeric"`
}

func (r *RedeemInput) Validate() []*ValidationErrorResponse {
	var errors []*ValidationErrorResponse

	validate := validator.New()
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
