package utils

type ctxKey string

const (
	PageKey     ctxKey = "page"
	LimitKey    ctxKey = "limit"
	UsernameKey ctxKey = "username"
)
