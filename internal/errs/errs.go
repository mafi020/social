package errs

import "errors"

var (
	ErrNotFound       = errors.New("resource not found")
	ErrDuplicateEntry = errors.New("duplicate entry")
	ErrUnauthorized   = errors.New("unauthorized access")
)
