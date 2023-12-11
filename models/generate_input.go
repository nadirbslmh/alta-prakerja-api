package models

import "github.com/go-playground/validator/v10"

type GenerateInput struct {
	RedeemCode  string `json:"redeem_code" validate:"required"`
	Sequence    int    `json:"sequence" validate:"required,numeric"`
	RedirectURI string `json:"redirect_uri" validate:"required,uri"`
	Email       string `json:"email" validate:"required,email"`
}

func (r *GenerateInput) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	return err
}
