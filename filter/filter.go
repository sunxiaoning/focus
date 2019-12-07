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

func Process(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	sort.Sort(types.FilterComparable(filters))

	// 过滤器执行
	var err error
	var process bool
	for _, filter := range filters {
		if process, err = isMatched(filter, req); err != nil {
			return err
		}
		if !process {
			continue
		}
		if ctx, err = filter.Process(ctx, rw, req); err != nil {
			return err
		}
	}
	return nil
}

func isMatched(filter *types.Filter, req *http.Request) (isMatched bool, err error) {
	if filter.Paths == nil || len(filter.Paths) <= 0 {
		return false, err
	}
	if filter.ExculdePaths != nil && len(filter.ExculdePaths) >= 1 {
		for _, excludePath := range filter.ExculdePaths {
			if strings.Contains(req.URL.Path, excludePath) {
				return false, nil
			}
		}
	}
	if filter.Paths[0] == "*" {
		return true, nil
	}
	for _, path := range filter.Paths {
		if strings.Contains(req.URL.Path, path) {
			return true, nil
		}
	}
	return false, nil
}
