package ginrpc

import (
	"fmt"
)

// Error is a response error with http status code and attributes
type Error struct {
	Code  int   // http statusCode
	Err   error // src error
	Attrs map[string]any
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %v", e.Code, e.Err)
}

func NewError(code int, err error) *Error {
	return &Error{code, err, nil}
}
