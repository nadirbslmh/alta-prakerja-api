package models

type CheckStatusInput struct {
	RedeemCode string `json:"redeem_code"`
	State      string `json:"state"`
	Sequence   int    `json:"sequence"`
}
