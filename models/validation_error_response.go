package models

type ValidationErrorResponse struct {
	ErrorMessage string `json:"error_message"`
	Field        string `json:"field"`
}
