package domain

import "errors"

var (
	ErrValidation = errors.New("validation error")
	ErrNotFound   = errors.New("not found")
	ErrForbidden  = errors.New("forbidden")
)

// ValidationError wraps ErrValidation with a more specific message.
func ValidationError(msg string) error {
	return errors.Join(ErrValidation, errors.New(msg))
}

// ForbiddenError wraps ErrForbidden to indicate authorization failures.
func ForbiddenError(msg string) error {
	return errors.Join(ErrForbidden, errors.New(msg))
}
