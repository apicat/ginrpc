package ginrpc

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// Handle for turning handlers into gin.handlefunc
func Handle[T any, M any](handler func(*gin.Context, *T) (M, error)) gin.HandlerFunc {
	paramType := reflect.TypeOf((*T)(nil)).Elem()
	if paramType.Kind() != reflect.Struct {
		panic("handler must be of type `func(*gin.Context,*T)(M,error)`` and the request parameter must be a struct pointer")
	}
	bindtags := findRequestParamTags(paramType)
	return func(ctx *gin.Context) {
		var err error
		var in T
		paramPtr := &in
		opt := getOptFromContext(ctx)
		if !IsEmpty(paramPtr) && opt.autoBinding() {
			err = verifyBindParams(ctx, paramPtr, bindtags)
			if err != nil {
				opt.getRenderFunc()(ctx, nil, warpError(err, http.StatusBadRequest))
				return
			}
		}
		for _, hook := range opt.beforeHooks {
			if err = hook(ctx, paramPtr, err); err != nil {
				break
			}
		}
		if err != nil {
			opt.getRenderFunc()(ctx, nil, warpError(err, http.StatusBadRequest))
			return
		}
		res, err := handler(ctx, paramPtr)
		opt.getRenderFunc()(ctx, res, warpError(err, 0))
	}
}

func warpError(err error, code int) *Error {
	if err == nil {
		return nil
	}
	var e *Error
	if !errors.As(err, &e) {
		e = NewError(code, err)
	}
	if e.Code == 0 && code > 0 {
		e.Code = code
	}
	return e
}
