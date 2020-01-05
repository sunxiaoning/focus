package notifycontrl

import (
	"context"
	"focus/types"
	"net/http"
)

const NotifyModule = types.ApiV1 + "/ppay"

var PPayNotify = types.NewController(NotifyModule+"/notify", http.MethodPost, pPayNotify)

func pPayNotify(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	types.NewRestRestResponse(rw, "SUCCESS")
}
