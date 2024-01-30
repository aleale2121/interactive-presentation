package model

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrIDMismatch= errors.New("id mismached")
	ErrNoPollDisplayed = errors.New("not displayed")
	ErrRunOutOfIndex = errors.New("poll run out of index")
)
