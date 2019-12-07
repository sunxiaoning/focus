package filter

import (
	"context"
	"focus/service/user"
	"focus/types"
	"focus/types/consts"
	"net/http"
)

var UserIdentityAuthor = &types.Filter{
	Order: 0,
	Paths: []string{
		"*",
	},
	ExculdePaths: pubPaths,
	Process:      userIdentityAuth,
}

func userIdentityAuth(ctx context.Context, rw http.ResponseWriter, req *http.Request) (context.Context, error) {
	accessTokenCookie, err := req.Cookie(consts.AccessToken)
	if err != nil || accessTokenCookie == nil {
		return checkAccessKey(ctx, rw, req)
	}
	ctx = context.WithValue(ctx, consts.AccessToken, accessTokenCookie.Value)
	return userservice.CheckUserExistsByAk(ctx)
}

func checkAccessKey(ctx context.Context, rw http.ResponseWriter, req *http.Request) (context.Context, error) {
	ctx = context.WithValue(ctx, consts.AccessToken, req.Header.Get(consts.AccessToken))
	return userservice.CheckUserExistsByAk(ctx)
}
