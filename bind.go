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
	if v == nil {
		return true
	}
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
			err = binding.MapFormWithTag(obj, ctx.Request.URL.Query(), "query")
		case "header":
			err = binding.MapFormWithTag(obj, ctx.Request.Header, "header")
		case "uri":
			if len(ctx.Params) > 0 {
				m := make(map[string][]string)
				for _, v := range ctx.Params {
					m[v.Key] = []string{v.Value}
				}
				err = binding.MapFormWithTag(obj, m, "uri")
			}
		}
	}

	if err != nil {
		return err
	}
	// request body content
	if ctx.Request.Method != http.MethodGet {
		b := binding.Default(ctx.Request.Method, ctx.ContentType())
		if err = ctx.ShouldBindWith(obj, b); err != nil {
			return err
		}
	}
	return binding.Validator.ValidateStruct(obj)
}
