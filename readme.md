# ginrpc

[中文](readme_zh.md)

A Gin middleware for RPC-Style coding

- High Performance **Use generics instead of reflection**
- Simple to use
- Automatic parameter binding
- Unified response/error handling


## Installation

```
$ go get github.com/apicat/ginrpc
```

## Usage

### base
```go
type In struct {
	ID int64      `uri:"id" binding:"required"`
}

type Out struct {
	Message string `json:"message"`
}

func rpcHandleDemo(c *gin.Context, in *In) (*Out, error) {
	return &Out{
		Message: fmt.Sprintf(" myid = %d", in.ID),
	}, nil
}

func main() {
	e := gin.Default()
	e.POST("/example/:id", ginrpc.Handle(rpcHandleDemo))
}

```


### custom
Inject custom configuration using middleware

- `ReponseRender` Customize the response
- `AutomaticBinding` Automatic binding default:enable
- `RequestBeforeHook` Add a front hook With Handle

[See the full example](examples/custom/main.go)

```go
e := gin.Default()
e.Use(ginrpc.AutomaticBinding(false), ginrpc.RequestBeforeHook(customBind))
```

