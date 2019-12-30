package gtwctl

import (
	"context"
	"focus/types"
	"net/http"
)

var Gtw = types.NewController(types.ApiV1+"/gtw", http.MethodGet, gtw)

func gtw(ctx context.Context, rw http.ResponseWriter, req *http.Request) {

	types.NewRestRestResponse(rw, "Hello!")
}
