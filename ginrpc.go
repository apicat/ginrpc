package ginrpc

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// Handle for turning handlers into gin.handlefunc
func Handle[T any, M any](handler func(*gin.Context, *T) (*M, error)) gin.HandlerFunc {

	if !checkHandleParamTypeValid(reflect.TypeOf(handler)) {
		panic(errHandleType)
	}
	bindtags := findRequestParamTags(reflect.TypeOf(new(T)))
	return func(ctx *gin.Context) {
		var err error
		var in T
		paramPtr := &in
		opt := getOptFromContext(ctx)
		if !IsEmpty(paramPtr) && opt.autoBinding() {
			err = verifyBindParams(ctx, paramPtr, bindtags)
		}
		for _, hook := range opt.beforeHooks {
			if err = hook(ctx, paramPtr, err); err != nil {
				break
			}
		}
		if err != nil {
			var e *Error
			if !errors.As(err, &e) {
				e = NewError(http.StatusBadRequest, err)
			}
			responseRender(opt, ctx, nil, e)
			return
		}
		res, err := handler(ctx, paramPtr)
		responseRender(opt, ctx, res, err)
	}
}

func checkHandleParamTypeValid(tp reflect.Type) bool {
	return tp.In(1).Elem().Kind() == reflect.Struct && tp.Out(0).Elem().Kind() == reflect.Struct
}

func responseRender(opt *option, ctx *gin.Context, res any, err error) {
	r := opt.getRenderFunc()
	if err != nil {
		var e *Error
		if !errors.As(err, &e) {
			e = &Error{Err: err, Code: http.StatusOK}
		}
		if e.Code == 0 {
			e.Code = http.StatusOK
		}
		r(ctx, res, e)
		return
	}
	if !IsEmpty(res) {
		r(ctx, res, nil)
	}
}
