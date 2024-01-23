package models

type Task struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	Session    string `json:"session"`
	Sequence   int    `json:"sequence"`
	Link       string `json:"link"`
	Batch      string `json:"batch"`
	CourseTag  string `json:"course_tag"`
	RedeemCode string `json:"redeem_code"`
	Scope      string `json:"scope"`
	Feedback   string `json:"feedback"`
	Score      int    `json:"score"`
}
