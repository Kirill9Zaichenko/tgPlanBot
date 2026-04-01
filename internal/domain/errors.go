package domain

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrTaskNotFound        = errors.New("task not found")
	ErrTaskRequestNotFound = errors.New("task request not found")
)
