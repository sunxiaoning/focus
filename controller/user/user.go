package usercontroller

import (
	"context"
	"focus/cfg"
	"focus/service/user"
	"focus/types"
	userconsts "focus/types/consts/user"
	"focus/types/user"
	"focus/util"
	"net/http"
	"strconv"
	"strings"
)

var Login = types.NewController(types.ApiV1+"/user/login", http.MethodGet, login)

func login(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
	params := req.URL.Query()
	username := strings.TrimSpace(params.Get("username"))
	if username == "" {
		return types.NewErr(types.InvalidParamError, "username can't be empty!")
	}
	passwd := strings.TrimSpace(params.Get("passwd"))
	if passwd == "" {
		return types.NewErr(types.InvalidParamError, "passwd can't be empty!")
	}
	ctx = context.WithValue(ctx, "userlogin", &usertype.UserLoginReq{username, passwd})
	user, err := userservice.CheckUserExistsBypwd(ctx)
	if err != nil {
		return err
	}
	accessToken, err := util.AESUtil.Encrypt(cfg.FocusCtx.Cfg.Server.SecretKey.AesKey, strings.Join([]string{
		strconv.FormatInt(user.ID, 10),
		user.UserName,
	}, ":"))
	if err != nil {
		return err
	}
	writeUserCookie(rw, accessToken)
	return types.NewRestRestResponse(rw, &usertype.UserLoginRes{UserId: user.ID})
}

func writeUserCookie(rw http.ResponseWriter, accessToken string) {
	http.SetCookie(rw, &http.Cookie{
		Name:   userconsts.AccessToken,
		Value:  accessToken,
		MaxAge: 60 * 60 * 24,
	})
}
