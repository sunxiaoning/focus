package filter

import (
	"context"
	"focus/service/user"
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

const (
	AccessToken = "accessToken"
)

func process(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	accessTokenCookie, err := req.Cookie(AccessToken)
	if err != nil || accessTokenCookie == nil {
		return checkAccessKey(ctx, rw, req)
	}
	return userservice.CheckUserExistsByAk(ctx, accessTokenCookie.Value)
}

func checkAccessKey(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	ak := req.Header.Get(AccessToken)
	return userservice.CheckUserExistsByAk(ctx, ak)
}
