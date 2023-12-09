package models

type GenerateInput struct {
	RedeemCode  string `json:"redeem_code"`
	Sequence    int    `json:"sequence"`
	RedirectURI string `json:"redirect_uri"`
	Email       string `json:"email"`
}
