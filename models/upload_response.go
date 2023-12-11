package models

import "encoding/json"

func UnmarshalUploadResponse(data []byte) (UploadResponse, error) {
	var r UploadResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UploadResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type UploadResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    UploadData `json:"data"`
}

type UploadData struct {
	ID         int64  `json:"id"`
	RedeemCode string `json:"redeem_code"`
	Scope      string `json:"scope"`
	Sequence   int64  `json:"sequence"`
	Status     int64  `json:"status"`
	URLFile    string `json:"url_file"`
}
