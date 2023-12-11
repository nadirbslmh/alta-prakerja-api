package models

import "encoding/json"

func UnmarshalCheckStatusResponse(data []byte) (CheckStatusResponse, error) {
	var r CheckStatusResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CheckStatusResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CheckStatusResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    StatusData `json:"data"`
}

type StatusData struct {
	RedeemCode           string `json:"redeem_code"`
	CourseCode           string `json:"course_code"`
	InvoiceCode          string `json:"invoice_code"`
	ScheduleCode         string `json:"schedule_code"`
	Status               int64  `json:"status"`
	AttendanceStatus     int64  `json:"attendance_status"`
	RedeemAt             string `json:"redeem_at"`
	DPPlatform           string `json:"dp_platform"`
	ScheduleStartEnd     string `json:"schedule_start_end"`
	CourseType           int64  `json:"course_type"`
	CourseTypeLabel      string `json:"course_type_label"`
	CourseLearningMethod string `json:"course_learning_method"`
	IsOnlineAttendance   bool   `json:"is_online_attendance"`
	Sequence             int64  `json:"sequence"`
}
