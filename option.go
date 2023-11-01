package ginrpc

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type option struct {
	disAutobinding bool
	render         func(*gin.Context, any, *Error)
	beforeHooks    []func(*gin.Context, any, error) error
}

func (o *option) getRenderFunc() func(*gin.Context, any, *Error) {
	if o.render != nil {
		return o.render
	}
	return defaultResponseRender
}

func (o *option) autoBinding() bool {
	return o.disAutobinding
}

// ReponseRender custom response output
func ReponseRender(r func(*gin.Context, any, *Error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		setOptToContext(ctx, &option{render: r})
	}
}

// // RequestBeforeHook The input parameter error represents a binding verification error.
// // If you want to prevent this error, you need to return the error as nil, otherwise, please return the error of the input parameter
func RequestBeforeHook(hook ...func(*gin.Context, any, error) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(hook) > 0 {
			setOptToContext(ctx, &option{beforeHooks: hook})
		}
	}
}

// AutomaticBinding enable/disable automatic binding and parameter validation
func AutomaticBinding(v bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		setOptToContext(ctx, &option{disAutobinding: v})
	}
}

const ctxOptKey = "_ginrpc_opt"

var defaultOption = &option{}

func getOptFromContext(ctx *gin.Context) *option {
	if o, ok := ctx.Get(ctxOptKey); ok {
		return o.(*option)
	}
	return defaultOption
}

func setOptToContext(ctx *gin.Context, opt *option) {
	if o, ok := ctx.Get(ctxOptKey); ok {
		preOpt := o.(*option)
		if opt.render != nil {
			preOpt.render = opt.render
		}
		if opt.disAutobinding {
			preOpt.disAutobinding = opt.disAutobinding
		}
		if len(opt.beforeHooks) > 0 {
			preOpt.beforeHooks = append(preOpt.beforeHooks, opt.beforeHooks...)
		}
		opt = nil
	} else {
		ctx.Set(ctxOptKey, opt)
	}
}

func defaultResponseRender(ctx *gin.Context, data any, err *Error) {
	statusCode := http.StatusOK
	if err != nil {
		resbody := make(map[string]any)
		if err.Attrs != nil {
			for k := range err.Attrs {
				resbody[k] = err.Attrs[k]
			}
		}
		resbody["message"] = err.Error()
		if err.Code > 0 {
			statusCode = err.Code
		}
		ctx.JSON(statusCode, resbody)
	} else {
		ctx.JSON(http.StatusOK, data)
	}
}
