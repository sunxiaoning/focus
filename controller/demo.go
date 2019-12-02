package controller

import (
	"context"
	"focus/types"
	"net/http"
)

var Hi = NewController("/hi", http.MethodGet, hi)

func hi(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	rw.Write([]byte("hi,world!"))
	return nil
}

var Hello = NewController("/hello", http.MethodGet, hello)

func hello(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	rw.Write([]byte("hello,world!"))
	return nil
}

var Err = NewController("/err", http.MethodGet, err)

func err(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	// return errors.New("divide")
	return types.NewErr(types.SystemError, "未知异常！")
}
