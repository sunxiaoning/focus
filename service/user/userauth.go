package userservice

import (
	"context"
	"focus/cfg"
	"focus/types"
	userconsts "focus/types/consts/user"
	"focus/types/user"
	"focus/util"
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
func CheckUserExistsByAk(ctx context.Context) (context.Context, error) {
	ak := ctx.Value(userconsts.AccessToken).(string)
	if ak == "" {
		return ctx, types.NewErr(types.NeedAuthError, "user need auth!")
	}
	userinfo, err := util.AESUtil.Decrypt(cfg.FocusCtx.Cfg.Server.SecretKey.AesKey, ak)
	if err != nil {
		return ctx, types.NewErr(types.NeedAuthError, "user need auth!")
	}
	username := strings.Split(userinfo, ":")[1]
	if _, ok := cfg.FocusCtx.CurrentUser.Load(username); ok {
		return ctx, nil
	}
	var user CurrentUserInfo
	cfg.FocusCtx.DB.Where("member_id = ? and status = 1", strings.Split(userinfo, ":")[0]).Find(&user)
	if user.ID == 0 {
		return ctx, types.NewErr(types.NeedAuthError, "user need auth!")
	}
	cfg.FocusCtx.CurrentUser.Store(username, user)
	return ctx, nil
}

func CheckUserExistsBypwd(ctx context.Context) (*CurrentUserInfo, error) {
	userlogin, ok := ctx.Value("userlogin").(*usertype.UserLoginReq)
	if !ok {
		return nil, types.NewErr(types.InvalidParamError, "userlogin param error!")
	}

	// user already exists
	if user, ok := cfg.FocusCtx.CurrentUser.Load(userlogin.Username); ok {
		return user.(*CurrentUserInfo), nil
	}
	var user CurrentUserInfo
	cfg.FocusCtx.DB.Where("user_name = ? and passwd = ? and status = 1", userlogin.Username, userlogin.Passwd).First(&user)
	if user.ID == 0 {
		return nil, types.NewErr(types.NotFound, "user not exists!")
	}
	cfg.FocusCtx.CurrentUser.Store(user.UserName, &user)
	return &user, nil
}
