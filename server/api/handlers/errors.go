package handlers

import "errors"

var (
	ErrInvalidType     = errors.New("invalid type")
	ErrNoDataFound     = errors.New("no data found")
	ErrNothingToDelete = errors.New("nothing to delete")
	ErrInvalidRequest  = errors.New("invalid request")
)
