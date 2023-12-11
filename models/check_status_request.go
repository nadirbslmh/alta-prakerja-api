package models

type CheckStatusRequest struct {
	RedeemCode string `json:"redeem_code"`
	Sequence   int    `json:"sequence"`
}
