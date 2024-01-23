package utils

type ctxKey string

const (
	PageKey      ctxKey = "page"
	LimitKey     ctxKey = "limit"
	UsernameKey  ctxKey = "username"
	BatchKey     ctxKey = "batch"
	CourseTagKey ctxKey = "course_tag"
)

var COURSE_TAGS []string = []string{
	"MSO",
	"CV",
}
