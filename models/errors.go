package models

import "errors"

var (
	// ErrNotFound is returned when a record is not found in the database
	ErrNotFound = errors.New("record not found")
	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")
	// ErrUnauthorized is returned when a user doesn't have permission
	ErrUnauthorized = errors.New("unauthorized access")
)
