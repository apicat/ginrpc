package ginrpc

import (
	"errors"
	"fmt"
)

var errHandleType = errors.New("rpc handle must HandlerFunc type")

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
