package errorx

import "errors"

var (
	ServerURLError = errors.New("server uri error")
	ClientKeyError = errors.New("not set client key")
)
