package filter

import (
	"context"
	"focus/types"
	"net/http"
	"sort"
	"strings"
)

var filters = []*types.Filter{
	UserIdentityAuthor, VisiterLimiter,
}

var (
	pubPaths = []string{
		"/login",
		"/hello",
	}
)

func Process(ctx context.Context, rw http.ResponseWriter, req *http.Request) context.Context {
	sort.Sort(types.FilterComparable(filters))

	// 过滤器执行
	for _, filter := range filters {
		if isMatched(filter, req) {
			ctx = filter.Process(ctx, rw, req)
		}
	}
	return ctx
}

func isMatched(filter *types.Filter, req *http.Request) (isMatched bool) {
	if filter.Paths == nil || len(filter.Paths) <= 0 {
		return false
	}
	if filter.ExculdePaths != nil && len(filter.ExculdePaths) >= 1 {
		for _, excludePath := range filter.ExculdePaths {
			if strings.Contains(req.URL.Path, excludePath) {
				return false
			}
		}
	}
	if filter.Paths[0] == "*" {
		return true
	}
	for _, path := range filter.Paths {
		if strings.Contains(req.URL.Path, path) {
			return true
		}
	}
	return false
}
