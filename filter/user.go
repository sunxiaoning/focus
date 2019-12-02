package filter

import (
	"context"
	"focus/types"
	"net/http"
)

var User = &types.Filter{
	Order: 0,
	Paths: []string{
		"*",
	},
	ExculdePaths: pubPaths,
	Process:      process,
}

func process(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {

	return nil
}
