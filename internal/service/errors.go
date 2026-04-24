package service

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrLimitReached = errors.New("monitor limit reached (max 100)")
	ErrValidation   = errors.New("validation error")
)
