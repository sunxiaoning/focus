package ppayctl

import (
	"context"
	"encoding/json"
	ppayserv "focus/serv/ppay"
	"focus/types"
	ppaytype "focus/types/ppay"
	"net/http"
)

const pPayModule = types.ApiV1 + "/ppay"

var PPay = types.NewController(pPayModule+"/notify", http.MethodPost, pPayNotify)

func pPayNotify(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	var payNotifyReq ppaytype.PayResultNotifyReq
	if err := json.NewDecoder(req.Body).Decode(&payNotifyReq); err != nil {
		types.InvalidParamPanic("json invalid!")
	}
	ctx = context.WithValue(ctx, "payNotifyReq", &payNotifyReq)
	types.NewRestRestResponse(rw, ppayserv.ResultNotify(ctx))
}
