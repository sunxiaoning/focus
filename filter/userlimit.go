package filter

import (
	"context"
	"focus/cfg"
	"focus/types"
	gtwtype "focus/types/gtw"
	resourcetype "focus/types/resource"
	"golang.org/x/time/rate"
	"net/http"
	"strings"
)

var VisiterLimiter = &types.Filter{
	Order: 3,
	Paths: []string{
		types.ApiV1 + "/gtw",
	},
	Process: visitLimit,
}

func visitLimit(ctx context.Context, rw http.ResponseWriter, req *http.Request) context.Context {
	gtwReq := ctx.Value("gtwReq").(*gtwtype.GtWReq)
	currentResource := ctx.Value("currentResource").(*resourcetype.ResourceWithLimit)
	limiter, ok := cfg.FocusCtx.VisitorLimiter.Load(strings.Join([]string{currentResource.Resource.Path, gtwReq.MemberId}, ":"))
	if !ok {
		limiter = rate.NewLimiter(rate.Limit(currentResource.ConcurrencyNumber), currentResource.ConcurrencyNumber)
		cfg.FocusCtx.VisitorLimiter.Store(strings.Join([]string{currentResource.Resource.Path, gtwReq.MemberId}, ":"), limiter)
	}
	visitLimiter := limiter.(*rate.Limiter)
	if !visitLimiter.Allow() {
		types.ErrPanic(types.ExceedRateLimit, "exceed rate!")
	}
	return ctx
}
