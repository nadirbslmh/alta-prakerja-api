package models

type Task struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	Session    string `json:"session"`
	Link       string `json:"link"`
	Batch      string `json:"batch"`
	RedeemCode string `json:"redeem_code"`
	Scope      string `json:"scope"`
}
