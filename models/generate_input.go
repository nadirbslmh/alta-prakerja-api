package models

import (
	"gugcp/utils"

	"github.com/go-playground/validator/v10"
)

type GenerateInput struct {
	RedeemCode  string `json:"redeem_code" validate:"required"`
	Sequence    int    `json:"sequence" validate:"required,numeric"`
	RedirectURI string `json:"redirect_uri" validate:"required,uri"`
	Email       string `json:"email" validate:"required,email"`
}

func (r *GenerateInput) Validate() []*ValidationErrorResponse {
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
