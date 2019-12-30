package usercontrl

import (
	"context"
	"focus/cfg"
	userserv "focus/serv/user"
	"focus/types"
	userconsts "focus/types/consts/user"
	"focus/types/user"
	aesutil "focus/util/aes"
	"net/http"
	"strconv"
	"strings"
)

var Login = types.NewController(types.ApiV1+"/user/login", http.MethodGet, login)

func login(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	username := strings.TrimSpace(params.Get("username"))
	if username == "" {
		types.InvalidParamPanic("username can't be empty!")
	}
	passwd := strings.TrimSpace(params.Get("passwd"))
	if passwd == "" {
		types.InvalidParamPanic("passwd can't be empty!")
	}
	ctx = context.WithValue(ctx, "userlogin", &usertype.UserLoginReq{Username: username, Passwd: passwd})
	user := userserv.CheckUserExistsBypwd(ctx)
	accessToken, err := aesutil.Encrypt(cfg.FocusCtx.Cfg.Server.SecretKey.AesKey, strings.Join([]string{
		strconv.FormatInt(user.ID, 10),
		user.UserName,
	}, ":"))
	if err != nil {
		panic(err)
	}
	writeUserCookie(rw, accessToken)
	types.NewRestRestResponse(rw, &usertype.UserLoginRes{UserId: user.ID})
}

func writeUserCookie(rw http.ResponseWriter, accessToken string) {
	http.SetCookie(rw, &http.Cookie{
		Name:   userconsts.AccessToken,
		Value:  accessToken,
		Path:   types.ApiV1,
		MaxAge: 60 * 60 * 24,
	})
}
