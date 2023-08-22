package server

import "errors"

var (
	// ErrInvalidLogin is returned when the login is invalid
	ErrInvalidLogin  = errors.New("invalid login")
	ErrInvalidSecret = errors.New("invalid secret")
	ErrObjectMiss    = errors.New("object miss")
)
