package config

import "errors"

var (
	ErrEmptyAddress  = errors.New("address is required")
	ErrEmptyDatabase = errors.New("database path is required")
	ErrEmptyToken    = errors.New("admin token is required")
	ErrInvalidDate   = errors.New("fixed date must use YYYY-MM-DD")
)
