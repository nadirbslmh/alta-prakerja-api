package models

import (
	"gugcp/utils"

	"github.com/go-playground/validator/v10"
)

type UploadForm struct {
	UserID     int    `form:"userID" validate:"required,numeric"`
	Session    string `form:"session" validate:"required"`
	Batch      string `form:"batch" validate:"required"`
	RedeemCode string `form:"redeemCode"`
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
