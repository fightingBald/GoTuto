package domain

import "errors"

var (
	ErrValidation = errors.New("validation error")
	ErrNotFound   = errors.New("not found")
)

// ValidationError wraps ErrValidation with a more specific message.
func ValidationError(msg string) error {
	return errors.Join(ErrValidation, errors.New(msg))
}
