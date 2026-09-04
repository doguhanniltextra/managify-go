package domain

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrForbidden    = errors.New("forbidden")
	ErrBadRequest   = errors.New("bad request")
	ErrConflict     = errors.New("resource conflict")
	ErrUnauthorized = errors.New("unauthorized")
)
