package types

import (
	"context"
	"net/http"
)

type Handle func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error

type Controller struct {
	Path   string
	Method string
	Handle Handle
}

func NewController(path string, method string, handle Handle) *Controller {
	return &Controller{
		Path:   path,
		Method: method,
		Handle: handle,
	}
}
