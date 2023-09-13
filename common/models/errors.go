package models

import "errors"

var (
	ErrDeleted     = errors.New("item was deleted")
	ErrUnknownType = errors.New("unknown type")
)
