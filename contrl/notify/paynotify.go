package notifycontrl

import (
	"context"
	"focus/types"
	"net/http"
)

const NotifyModule = types.ApiV1 + "/notify"

var PPayNotify = types.NewController(NotifyModule+"/ppay", http.MethodGet, pPayNotify)

func pPayNotify(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	return nil
}
