package server

import "errors"

var (
	// ErrInvalidLogin is returned when the login is invalid
	ErrInvalidLogin      = errors.New("invalid login")
	ErrInvalidSecret     = errors.New("invalid secret")
	ErrObjectMiss        = errors.New("object miss")
	ErrNoDataFound       = errors.New("no data found")
	ErrNothingToDelete   = errors.New("nothing to delete")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserNotAuthorized = errors.New("user not authorized")
	ErrInvalidRequest    = errors.New("invalid request")
)
