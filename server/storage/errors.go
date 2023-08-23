package storage

import "errors"

var (
	// ErrObjectMiss is returned when object is not found in storage
	ErrObjectMiss = errors.New("object miss")
)
