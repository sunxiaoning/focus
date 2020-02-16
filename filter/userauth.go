package filter

import (
	"context"
	"focus/serv/user"
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

// 用户身份校验
func userIdentityAuth(ctx context.Context, rw http.ResponseWriter, req *http.Request) context.Context {
	accessTokenCookie, err := req.Cookie(userconsts.AccessToken)
	if err != nil || accessTokenCookie == nil {
		return checkAccessKey(ctx, rw, req)
	}
	ctx = context.WithValue(ctx, userconsts.AccessToken, accessTokenCookie.Value)
	return userserv.CheckUserExistsByAk(ctx)
}

func checkAccessKey(ctx context.Context, rw http.ResponseWriter, req *http.Request) context.Context {
	ctx = context.WithValue(ctx, userconsts.AccessToken, req.Header.Get(userconsts.AccessToken))
	return userserv.CheckUserExistsByAk(ctx)
}
