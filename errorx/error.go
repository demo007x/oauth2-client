package errorx

import "errors"

var (
	ServerURLError        = errors.New("server uri error")
	ClientKeyError        = errors.New("not set client key")
	SecretKeyError        = errors.New("not set secret")
	CodeEmptyError        = errors.New("code is empty")
	RequestServerURLError = errors.New("request server url error")
)
