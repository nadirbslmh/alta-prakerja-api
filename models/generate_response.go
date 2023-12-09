package models

import "encoding/json"

func UnmarshalGenerateURLResponse(data []byte) (GenerateURLResponse, error) {
	var r GenerateURLResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GenerateURLResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GenerateURLResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type Data struct {
	URL    string `json:"url"`
	State  string `json:"state"`
	Course Course `json:"course"`
}

type Course struct {
	Type           int64  `json:"type"`
	LearningMethod string `json:"learning_method"`
}
