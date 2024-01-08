package models

import "encoding/json"

func UnmarshalFeedbackRequest(data []byte) (FeedbackRequest, error) {
	var r FeedbackRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *FeedbackRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type FeedbackRequest struct {
	RedeemCode string `json:"redeem_code"`
	Scope      string `json:"scope"`
	Sequence   int64  `json:"sequence"`
	Notes      string `json:"notes"`
	URLFile    string `json:"url_file"`
}
