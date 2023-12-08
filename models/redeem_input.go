package models

type RedeemInput struct {
	UserID     int    `json:"user_id"`
	State      string `json:"state"`
	RedeemCode string `json:"redeem_code"`
	Sequence   int    `json:"sequence"`
	Status     byte   `json:"status"`
}
