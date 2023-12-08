package models

type Redeem struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	State      string `json:"state"`
	RedeemCode string `json:"redeem_code"`
	Sequence   int    `json:"sequence"`
	Status     uint8  `json:"status"`
}
