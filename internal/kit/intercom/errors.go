package intercom

import "errors"

var (
	ErrNoToken      = errors.New("token is required")
	ErrUserNotFound = errors.New("user not found")
	ErrTooManyUsers = errors.New("more than one user found")
)
