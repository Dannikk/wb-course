package custom_error

import "errors"

var (
	ErrNotFoundCache = errors.New("not found in cache")
	ErrNonOrderType  = errors.New("non order type")
)
