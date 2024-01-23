package models

type TaskData struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Sequence  int    `json:"sequence"`
	Link      string `json:"link"`
	Scope     string `json:"scope"`
	Batch     string `json:"batch"`
	CourseTag string `json:"course_tag"`
	Feedback  string `json:"feedback"`
	Score     int    `json:"score"`
	Name      string `json:"name"`
}
