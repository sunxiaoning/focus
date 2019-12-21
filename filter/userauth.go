package filter

import (
	"context"
	"focus/service/user"
	"focus/types"
	userconsts "focus/types/consts/user"
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
	accessTokenCookie, err := req.Cookie(userconsts.AccessToken)
	if err != nil || accessTokenCookie == nil {
		return checkAccessKey(ctx, rw, req)
	}
	ctx = context.WithValue(ctx, userconsts.AccessToken, accessTokenCookie.Value)
	return userservice.CheckUserExistsByAk(ctx)
}

func checkAccessKey(ctx context.Context, rw http.ResponseWriter, req *http.Request) (context.Context, error) {
	ctx = context.WithValue(ctx, userconsts.AccessToken, req.Header.Get(userconsts.AccessToken))
	return userservice.CheckUserExistsByAk(ctx)
}
