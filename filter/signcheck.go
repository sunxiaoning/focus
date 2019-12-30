package filter

import (
	"context"
	"encoding/json"
	"focus/types"
	gtwtype "focus/types/gtw"
	"net/http"
	"sort"
)

var SignCheck = &types.Filter{
	Order: 1,
	Paths: []string{
		types.ApiV1 + "/gtw",
	},
	Process: signCheck,
}

func signCheck(ctx context.Context, rw http.ResponseWriter, req *http.Request) context.Context {
	var gtwReq gtwtype.GtWReq
	if err := json.NewDecoder(req.Body).Decode(&gtwReq); err != nil {
		types.InvalidParamPanic("invalid json format!")
	}
	params := []string{gtwReq.Timestamp, gtwReq.MemberId, gtwReq.ServUrl, gtwReq.BizContent}
	sort.Strings(params)
	/*	origin := strings.Join(params, ",")
		if rsautil.VerifySign(origin, gtwReq.Sign, "")*/
	return ctx
}
