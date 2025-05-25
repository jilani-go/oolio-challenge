package constants

import "errors"

// Error constants
var (
	ErrClassNotFound     = errors.New("class not found")
	ErrDateOutOfRange    = errors.New("date is outside of class schedule")
	ErrInvalidDateFormat = errors.New("invalid date fromat")
)
