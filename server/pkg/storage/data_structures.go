package storage

import "errors"

var (
	// ErrObjectMiss is returned when the cache does not contain the requested key
	ErrObjectMiss = errors.New("object miss")
	// ErrInvalidLogin is returned when the login is invalid
	ErrInvalidLogin = errors.New("invalid login")
)
