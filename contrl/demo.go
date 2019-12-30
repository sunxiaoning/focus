package contrl

import (
	"context"
	"focus/types"
	"net/http"
)

var Hi = types.NewController("/hi", http.MethodGet, hi)

func hi(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("hi,world!"))
}

var Hello = types.NewController("/hello", http.MethodGet, hello)

func hello(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("hello,world!"))
}

var Err = types.NewController("/err", http.MethodGet, err)

func err(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	// return errors.New("divide")
	types.NewErr(types.SystemError, "未知异常！")
}
