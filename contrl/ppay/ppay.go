package ppayctl

import (
	"context"
	"focus/types"
	"net/http"
)

const pPayModule = types.ApiV1 + "/ppay"

var PPay = types.NewController(pPayModule+"/notify", http.MethodPost, pPayNotify)

func pPayNotify(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	types.NewRestRestResponse(rw, "SUCCESS")
}
