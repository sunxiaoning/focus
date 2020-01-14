package gtwctl

import (
	"context"
	"fmt"
	"focus/types"
	userconsts "focus/types/consts/user"
	gtwtype "focus/types/gtw"
	httputil "focus/util/http"
	"net/http"
	"time"
)

var Gtw = types.NewController(types.ApiV1+"/gtw", http.MethodPost, gtw)

const ApiServer = "http://localhost:7001/api/v1"

var Timeout = time.Second * 1

func gtw(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	gtwReq := ctx.Value("gtwReq").(*gtwtype.GtWReq)
	servUrl := ApiServer + gtwReq.ServUrl
	reqHeaders := make(map[string]string)
	reqHeaders[userconsts.AccessToken] = req.Header.Get(userconsts.AccessToken)
	fmt.Println(servUrl)
	code, resp, err := httputil.PostJsonWithHeader(servUrl, reqHeaders, gtwReq.BizContent, time.Second*10)
	fmt.Println(code, string(resp), err)
	if code != http.StatusOK || err != nil {
		types.SystemPanic(err.Error())
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}
