package userserv

import (
	"context"
	"focus/cfg"
	"focus/types"
	userconsts "focus/types/consts/user"
	"focus/types/user"
	aesutil "focus/util/aes"
	"strings"
)

type CurrentUserInfo struct {
	ID       int64
	MemberId int64
	UserName string
}

func (currentUserInfo CurrentUserInfo) TableName() string {
	return "user_account"
}
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
	var user CurrentUserInfo
	cfg.FocusCtx.DB.Where("member_id = ? and status = 1", strings.Split(userinfo, ":")[0]).Find(&user)
	if user.ID == 0 {
		panic(types.NewErr(types.NeedAuthError, "user need auth!"))
	}
	cfg.FocusCtx.CurrentUser.Store(username, &user)
	return ctx
}

func CheckUserExistsBypwd(ctx context.Context) *CurrentUserInfo {
	userlogin, ok := ctx.Value("userlogin").(*usertype.UserLoginReq)
	if !ok {
		panic(types.NewErr(types.InvalidParamError, "userlogin param error!"))
	}

	// user already exists
	if user, ok := cfg.FocusCtx.CurrentUser.Load(userlogin.Username); ok {
		return user.(*CurrentUserInfo)
	}
	var user CurrentUserInfo
	cfg.FocusCtx.DB.Where("user_name = ? and passwd = ? and status = 1", userlogin.Username, userlogin.Passwd).First(&user)
	if user.ID == 0 {
		panic(types.NewErr(types.NotFound, "user not exists!"))
	}
	cfg.FocusCtx.CurrentUser.Store(user.UserName, &user)
	return &user
}
