package gtwctl

import (
	"context"
	"focus/types"
	"net/http"
	"time"
)

var Gtw = types.NewController(types.ApiV1+"/gtw", http.MethodPost, gtw)

const ApiServer = "http://127.0.0.1:7001/"

var Timeout = time.Second * 1

func gtw(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	/*gtwReq := ctx.Value("gtwReq").(*gtwtype.GtWReq)
	servUrl := path.Join(ApiServer, types.ApiV1, gtwReq.ServUrl)*/
}
