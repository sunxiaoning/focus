package ppayctl

import (
	"context"
	"encoding/json"
	ppayserv "focus/serv/ppay"
	"focus/tx"
	"focus/types"
	ppaytype "focus/types/ppay"
	"net/http"
)

const pPayModule = types.ApiV1 + "/ppay"

var PPay = types.NewController(pPayModule+"/notify", http.MethodPost, pPayNotify)

// 客户端通知微信/支付宝支付结果
func pPayNotify(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	var payNotifyReq ppaytype.PayResultNotifyReq
	if err := json.NewDecoder(req.Body).Decode(&payNotifyReq); err != nil {
		types.InvalidParamPanic("json invalid!")
	}
	ctx = context.WithValue(ctx, "payNotifyReq", &payNotifyReq)
	types.NewRestRestResponse(rw, tx.NewTxManager().RunTx(ctx, ppayserv.ResultNotifyTx))
}
