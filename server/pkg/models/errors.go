package model

import "errors"

var (
	// ErrCacheMiss is returned when the cache does not contain the requested key
	ErrCacheMiss = errors.New("cache miss")
	// ErrInvalidLogin is returned when the login is invalid
	ErrInvalidLogin = errors.New("invalid login")
)
