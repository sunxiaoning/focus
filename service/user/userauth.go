package userservice

import (
	"context"
	"focus/cfg"
	"focus/types"
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
func CheckUserExistsByAk(ctx context.Context, ak string) error {
	if ak == "" {
		return types.NewErr(types.NeedAuthError, "user need auth!")
	}
	userinfo, err := util.AESUtil.Decrypt(cfg.FocusCtx.Cfg.Server.SecretKey.AesKey, ak)
	if err != nil {
		return types.NewErr(types.NeedAuthError, "user need auth!")
	}
	username := strings.Split(userinfo, ":")[1]
	if _, ok := cfg.FocusCtx.CurrentUser.Load(username); ok {
		return nil
	}
	var user CurrentUserInfo
	cfg.FocusCtx.DB.Where(map[string]interface{}{
		"id":     strings.Split(userinfo, ":")[0],
		"status": true,
	}).Find(&user)
	if user.ID == 0 {
		return types.NewErr(types.NeedAuthError, "user need auth!")
	}
	cfg.FocusCtx.CurrentUser.Store(username, user)
	return nil
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
	cfg.FocusCtx.DB.Where(map[string]interface{}{
		"user_name": userlogin.Username,
		"passwd":    userlogin.Passwd,
	}).Find(&user)
	if user.ID == 0 {
		return nil, types.NewErr(types.UserNotFound, "user not exists!")
	}
	cfg.FocusCtx.CurrentUser.Store(user.UserName, &user)
	return &user, nil
}
