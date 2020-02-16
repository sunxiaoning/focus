package userserv

import (
	"context"
	"focus/cfg"
	memloginrepo "focus/repo/memlogin"
	"focus/types"
	userconsts "focus/types/consts/user"
	membertype "focus/types/member"
	aesutil "focus/util/aes"
	"strconv"
	"strings"
)

// 根据Ak校验
func CheckUserExistsByAk(ctx context.Context) context.Context {
	ak := ctx.Value(userconsts.AccessToken).(string)
	if ak == "" {
		panic(types.NewErr(types.NeedAuthError, "user need auth!"))
	}
	userinfo, err := aesutil.Decrypt(cfg.FocusCtx.Cfg.Server.SecretKey.AesKey, ak)
	if err != nil {
		panic(types.NewErr(types.NeedAuthError, "user need auth!"))
	}
	username := strings.Split(userinfo, ":")[1]
	if _, ok := cfg.FocusCtx.CurrentUser.Load(username); ok {
		return ctx
	}
	memberId, err := strconv.Atoi(strings.Split(userinfo, ":")[0])
	if err != nil {
		panic(types.NewErr(types.NeedAuthError, "user need auth!"))
	}
	currentUser := memloginrepo.GetByMemberId(ctx, memberId)
	if currentUser.ID == 0 {
		panic(types.NewErr(types.NeedAuthError, "user need auth!"))
	}
	cfg.FocusCtx.CurrentUser.Store(currentUser.UserName, currentUser)
	return ctx
}

// 根据用户名密码校验
func CheckUserExistsBypwd(ctx context.Context) *membertype.CurrentUserInfo {
	userlogin, ok := ctx.Value("userlogin").(*membertype.MemberLoginReq)
	if !ok {
		panic(types.NewErr(types.InvalidParamError, "userlogin param error!"))
	}

	// user already exists
	if user, ok := cfg.FocusCtx.CurrentUser.Load(userlogin.Username); ok {
		return user.(*membertype.CurrentUserInfo)
	}
	currentUser := memloginrepo.GetByUsernameAndPwd(ctx, userlogin.Username, userlogin.Passwd)
	if currentUser.ID == 0 {
		panic(types.NewErr(types.NotFound, "user not exists!"))
	}
	cfg.FocusCtx.CurrentUser.Store(currentUser.UserName, currentUser)
	return currentUser
}
