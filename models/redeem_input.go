package models

import "github.com/go-playground/validator/v10"

type RedeemInput struct {
	UserID     int    `json:"user_id" validate:"required,numeric"`
	State      string `json:"state" validate:"required"`
	RedeemCode string `json:"redeem_code" validate:"required"`
	Sequence   int    `json:"sequence" validate:"required,numeric"`
}

func (r *RedeemInput) Validate() error {
	validate := validator.New()

	err := validate.Struct(r)

	return err
}
