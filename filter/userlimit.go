package filter

import (
	"context"
	"focus/cfg"
	"focus/service/resource"
	"focus/types"
	"focus/types/consts"
	"golang.org/x/time/rate"
	"net/http"
)

var VisiterLimiter = &types.Filter{
	Order: 1,
	Paths: []string{
		"*",
	},
	Process: visitLimit,
}

func visitLimit(ctx context.Context, rw http.ResponseWriter, req *http.Request) (context.Context, error) {
	ak := ctx.Value(consts.AccessToken).(string)
	limiter, ok := cfg.FocusCtx.VisitorLimiter.Load(ak)
	if !ok {
		limiter = rate.NewLimiter(1, 1)
		cfg.FocusCtx.VisitorLimiter.Store(ak, limiter)
	}
	visitLimiter := limiter.(*rate.Limiter)
	if !visitLimiter.Allow() {
		return ctx, types.NewErr(types.ExceedRateLimit, "exceed rate!")
	}
	resourceservice.QueryServiceResource(ctx)
	return ctx, nil
}
