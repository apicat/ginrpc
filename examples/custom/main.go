package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/apicat/ginrpc"
	"github.com/gin-gonic/gin"
)

type In struct {
	ID int64 `uri:"id" binding:"required"`
}

type Out struct {
	Message string `json:"message"`
}

func rpcHandle(ctx *gin.Context, in *In) (*Out, error) {
	return &Out{Message: fmt.Sprintf(" myid = %d", in.ID)}, nil
}

func main() {
	e := gin.Default()
	e.Use(
		ginrpc.ReponseRender(MyRender),
		ginrpc.RequestBeforeHook(globalRequestDataLog),
	)
	e.GET("/ginrpc/:id", ginrpc.RequestBeforeHook(doubleID), ginrpc.Handle(rpcHandle))

	e.GET("/", ginrpc.Handle(func(ctx *gin.Context, t *ginrpc.Empty) (*ginrpc.Empty, error) {
		ctx.String(http.StatusOK, "hello,world")
		// response *ginrpc.Empty 不会触发ResponseRender
		return &ginrpc.Empty{}, nil
	}))
	_ = e.Run(":8080")
}

func globalRequestDataLog(ctx *gin.Context, in any, lastBindErr error) error {
	if ginrpc.IsEmpty(in) {
		return nil
	}
	log.Printf(
		"method=%s path=%s requestData=%+v",
		ctx.Request.Method,
		ctx.Request.URL.Path,
		in,
	)
	return nil
}

func doubleID(ctx *gin.Context, in any, lastBindErr error) error {
	if lastBindErr != nil {
		return lastBindErr
	}
	paramPtr, ok := in.(*In)
	if ok {
		paramPtr.ID *= 2
	}
	return nil
}

func MyRender(ctx *gin.Context, respdata any, err *ginrpc.Error) {
	if err != nil {
		ctx.JSON(err.Code, gin.H{
			"errcode": 1,
			"errmsg":  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"errcode": 0,
		"data":    respdata,
	})
}
