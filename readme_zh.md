# ginrpc

用于 RPC 风格编码的 Gin 中间件

- 高性能 **使用泛型而不是反射，没有性能损失**
- 简单且可定制
- 自动绑定参数并验证
- 统一的响应/错误处理


## 安装

```
$ go get github.com/apicat/ginrpc
```

## 使用

### 基本方式

使用 `ginrpc.Handle` 转换你的rpc函数，仅此而已

```go
type In struct {
	ID int64      `uri:"id" binding:"required"`
}

type Out struct {
	Message string `json:"message"`
}

func rpcHandleDemo(c *gin.Context, in *In) (*Out, error) {
	msg := fmt.Sprintf(" myid = %d", in.ID)
	return &Out{Message: msg}, nil
}

func main() {
	e := gin.Default()
	e.POST("/example/:id", ginrpc.Handle(rpcHandleDemo))
}
```


### 自定义方式
通过中间件注册配置信息

- `ReponseRender` 自定义响应处理。 
- `AutomaticBinding` 你可能想自己处理绑定，因此可以禁用默认的自动绑定。默认是启用的
- `RequestBeforeHook` 添加一个前置钩子。 比如打印请求参数，或者自定义绑定可以在这里处理

[查看完整示例](examples/custom/main.go)

```go
e := gin.Default()
e.Use(ginrpc.AutomaticBinding(false), ginrpc.RequestBeforeHook(customBind))
```

