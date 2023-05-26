package ginrpc

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Empty like google.protobuf.Empty
type Empty struct{}

func IsEmpty(v any) bool {
	_, ok := v.(*Empty)
	return ok
}

var autoBindRequestParamTags = []string{
	"header",
	"query",
	"uri",
}

func findRequestParamTags(t reflect.Type) map[string]struct{} {
	m := make(map[string]struct{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			for k := range findRequestParamTags(field.Type) {
				m[k] = struct{}{}
			}
		} else {
			for _, v := range autoBindRequestParamTags {
				if _, ok := field.Tag.Lookup(v); ok {
					m[v] = struct{}{}
				}
			}
		}
	}
	return m
}

func verifyBindParams(ctx *gin.Context, obj any, tags map[string]struct{}) error {
	var err error
	// request param
	// in:[path(uri),query,header...]
	for bindTag := range tags {
		switch bindTag {
		case "query":
			// err = ctx.ShouldBindQuery(obj)
			// ShouldBindQuery 使用 form标签
			// 当不是get请求时 同时存在formdata和query的时候会有冲突
			// 在这里所有query都只使用query标签
			// example:
			// ```
			// {
			//   SomeVar string `query:"some_var"`
			// }
			// ```
			err = customShouldBindQuery(ctx, obj)
		case "header":
			err = ctx.ShouldBindHeader(obj)
		case "uri":
			err = ctx.ShouldBindUri(obj)
		}
	}

	if err != nil {
		return err
	}
	// request body content
	if ctx.Request.Method != http.MethodGet {
		return ctx.ShouldBind(obj)
	}
	return nil
}

func customShouldBindQuery(ctx *gin.Context, obj any) error {
	if err := binding.MapFormWithTag(
		obj, ctx.Request.URL.Query(), "query"); err != nil {
		return err
	}
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}
