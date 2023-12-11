package models

type UploadRequest struct {
	RedeemCode string `json:"redeem_code"`
	Scope      string `json:"scope"`
	Sequence   int    `json:"sequence"`
	FileURL    string `json:"url_file"`
}
