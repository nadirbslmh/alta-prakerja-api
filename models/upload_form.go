package models

import (
	"gugcp/utils"

	"github.com/go-playground/validator/v10"
)

type UploadForm struct {
	RedeemCode string `form:"redeemCode" validate:"required"`
	Scope      string `form:"scope" validate:"required,oneof=tpm uk"`
	Sequence   int    `form:"sequence" validate:"required,numeric,gte=1,lte=999"`
}

func (r *UploadForm) Validate() []*ValidationErrorResponse {
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
