package models

import "errors"

var (
	ErrNotFound = errors.New("Not Found")
	ErrConflict = errors.New("Conflict")
)
