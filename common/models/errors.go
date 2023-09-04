package models

import "errors"

var (
	ErrDeleted  = errors.New("item was deleted")
	UnknownType = errors.New("unknown type")
)
