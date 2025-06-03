package errs

import "errors"

var (
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrBadRequest     = errors.New("bad request")
	ErrNotFound       = errors.New("not found")
	ErrInternalServer = errors.New("internal server error")
	ErrInvalidRequest = errors.New("invalid request")
)
