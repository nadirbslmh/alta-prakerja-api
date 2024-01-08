package models

import "encoding/json"

func UnmarshalFeedbackResponse(data []byte) (FeedbackResponse, error) {
	var r FeedbackResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *FeedbackResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type FeedbackResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    FeedbackData `json:"data"`
}

type FeedbackData struct {
	ID         int64  `json:"id"`
	RedeemCode string `json:"redeem_code"`
	Scope      string `json:"scope"`
	Sequence   int64  `json:"sequence"`
	Status     int64  `json:"status"`
	URLFile    string `json:"url_file"`
	Notes      string `json:"notes"`
}
