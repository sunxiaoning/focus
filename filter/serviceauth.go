package filter

import (
	"context"
	"fmt"
	"focus/cfg"
	resourceserv "focus/serv/resource"
	"focus/types"
	gtwtype "focus/types/gtw"
	resourcetype "focus/types/resource"
	dbutil "focus/util/db"
	timutil "focus/util/tim"
	"net/http"
	"time"
)

var ServiceAuth = &types.Filter{
	Order: 2,
	Paths: []string{
		types.ApiV1 + "/gtw",
	},
	Process: serviceAuth,
}

func serviceAuth(ctx context.Context, rw http.ResponseWriter, req *http.Request) context.Context {
	gtwReq := ctx.Value("gtwReq").(*gtwtype.GtWReq)
	defer func() {
		if r := recover(); r != nil {
			cfg.FocusCtx.MemberService.Store(gtwReq.MemberId, &memberservicecache{
				Timestamp:      timutil.DefFormat(time.Now().Add(time.Minute * 3)),
				MemberServices: nil,
			})
			panic(r)
		}
	}()
	memberServiceCache, ok := cfg.FocusCtx.MemberService.Load(gtwReq.MemberId)
	isCachedMemberResources := false
	if ok {
		timestamp, _ := timutil.DefParse(memberServiceCache.(*memberservicecache).Timestamp)
		if timestamp.After(time.Now()) {
			isCachedMemberResources = true
		}
	}
	if !isCachedMemberResources {
		var memberServices []*memberservicetype
		dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("member_service").Select("member_id,service_id,service_price_id").Where("member_id = ? and status = 1 and user_service_status = 'NORMAL'", gtwReq.MemberId).Find(&memberServices))
		if len(memberServices) == 0 {
			types.ErrPanic(types.NeedAuthError, "need auth!")
		}
		var priceIds []int
		for _, memberService := range memberServices {
			priceIds = append(priceIds, memberService.ServicePriceId)
		}
		var servicePrices []servicepricetype
		dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("service_price").Select("id,concurrency_number").Where("id in (?) and status = 1", priceIds).Find(&servicePrices))
		if len(servicePrices) == 0 {
			types.ErrPanic(types.DataDirty, fmt.Sprintf("priceIds:%v price not config!", priceIds))
		}
		servicePriceMap := make(map[int]int)
		for _, servicePrice := range servicePrices {
			servicePriceMap[servicePrice.ID] = servicePrice.ConcurrencyNumber
		}
		for _, memberService := range memberServices {
			memberService.ConcurrencyNumber = servicePriceMap[memberService.ServicePriceId]
		}
		memberServiceCache = &memberservicecache{
			Timestamp:      timutil.DefFormat(time.Now().Add(time.Minute * 3)),
			MemberServices: memberServices,
		}
		cfg.FocusCtx.MemberService.Store(gtwReq.MemberId, memberServiceCache)
	}
	memberServices := memberServiceCache.(*memberservicecache).MemberServices
	if memberServices == nil || len(memberServices) == 0 {
		types.ErrPanic(types.NeedAuthError, "need auth!")
	}
	hasAuth := false
	var concurrencyNumber int
	currentResource := resourceserv.FilterSingleResource(func(resource *resourcetype.Resource) bool {
		for _, memberService := range memberServiceCache.(*memberservicecache).MemberServices {
			if resource.ServiceId == memberService.ServiceId && resource.Path == gtwReq.ServUrl {
				hasAuth = true
				concurrencyNumber = memberService.ConcurrencyNumber
				return true
			}
		}
		return false
	})
	if !hasAuth {
		types.ErrPanic(types.NeedAuthError, "need auth!")
	}
	ctx = context.WithValue(ctx, "currentResource", &resourcetype.ResourceWithLimit{
		Resource:          currentResource,
		ConcurrencyNumber: concurrencyNumber,
	})
	return ctx
}

type memberservicetype struct {
	MemberId          int
	ServiceId         int
	ServicePriceId    int
	ConcurrencyNumber int
}

type memberservicecache struct {
	Timestamp      string
	MemberServices []*memberservicetype
}

type servicepricetype struct {
	ID                int
	ConcurrencyNumber int
}
