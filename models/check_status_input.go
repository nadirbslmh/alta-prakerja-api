package models

import "github.com/go-playground/validator/v10"

type CheckStatusInput struct {
	RedeemCode string `json:"redeem_code" validate:"required"`
	State      string `json:"state" validate:"required"`
	Sequence   int    `json:"sequence" validate:"required,numeric"`
}

func (r *CheckStatusInput) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	return err
}
