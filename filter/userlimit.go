package filter

import (
	"context"
	"focus/cfg"
	"focus/types"
	userconsts "focus/types/consts/user"
	"golang.org/x/time/rate"
	"net/http"
)

var VisiterLimiter = &types.Filter{
	Order: 1,
	Paths: []string{
		types.ApiV1 + "/gtw",
	},
	Process: visitLimit,
}

func visitLimit(ctx context.Context, rw http.ResponseWriter, req *http.Request) context.Context {
	ak := ctx.Value(userconsts.AccessToken).(string)
	limiter, ok := cfg.FocusCtx.VisitorLimiter.Load(ak)
	if !ok {
		limiter = rate.NewLimiter(1, 1)
		cfg.FocusCtx.VisitorLimiter.Store(ak, limiter)
	}
	visitLimiter := limiter.(*rate.Limiter)
	if !visitLimiter.Allow() {
		types.ErrPanic(types.ExceedRateLimit, "exceed rate!")
	}
	return ctx
}
