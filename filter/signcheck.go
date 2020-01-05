package filter

import (
	"context"
	"encoding/json"
	"fmt"
	"focus/cfg"
	"focus/types"
	gtwtype "focus/types/gtw"
	membertype "focus/types/member"
	dbutil "focus/util/db"
	rsautil "focus/util/rsa"
	strutil "focus/util/strs"
	timutil "focus/util/tim"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var SignCheck = &types.Filter{
	Order: 1,
	Paths: []string{
		types.ApiV1 + "/gtw",
	},
	Process: signCheck,
}

func signCheck(ctx context.Context, rw http.ResponseWriter, req *http.Request) context.Context {
	var gtwReq gtwtype.GtWReq
	if err := json.NewDecoder(req.Body).Decode(&gtwReq); err != nil {
		types.InvalidParamPanic("invalid json format!")
	}
	validateParams(&gtwReq)
	_, pubKeyStr := getMemPriKey(gtwReq.MemberId)
	if pubKeyStr == "" {
		types.InvalidParamPanic(fmt.Sprintf("member: %s, priKey not exists!", gtwReq.MemberId))
	}
	origin := strings.Join([]string{gtwReq.Timestamp, gtwReq.MemberId, gtwReq.ServUrl, gtwReq.BizContent}, ",")
	pubKey, err := rsautil.ParseKeyFromString(pubKeyStr)
	if err != nil {
		types.ErrPanic(types.InvalidKeyFormat, "parse user pubKey error!")
	}
	if ok, err := rsautil.VerifySign(origin, gtwReq.Sign, pubKey); !ok || err != nil {
		types.ErrPanic(types.VerifySignError, fmt.Sprintf("verifySign failure, %v", err))
	}
	ctx = context.WithValue(ctx, "gtwReq", &gtwReq)
	return ctx
}

func getMemPriKey(memberId string) (string, string) {
	var priKey, pubKey string
	defer func() {
		if r := recover(); r != nil {
			cfg.FocusCtx.MemberSecretKey.Store(memberId, &memberprikeycache{
				Timestamp: timutil.DefFormat(time.Now().Add(time.Minute * 3)),
				PriKey:    "",
				PubKey:    "",
			})
		}
	}()
	mId, err := strconv.Atoi(memberId)
	if err != nil {
		types.InvalidParamPanic(fmt.Sprintf("invalid memberId:%s", memberId))
	}
	memberPriKeyCache, ok := cfg.FocusCtx.MemberSecretKey.Load(mId)
	isMemberPriKeyCached := false
	if ok {
		timestamp, _ := timutil.DefParse(memberPriKeyCache.(*memberprikeycache).Timestamp)
		if timestamp.After(time.Now()) {
			isMemberPriKeyCached = true
		}
	}
	if !isMemberPriKeyCached {
		var secretKey membertype.SecretKey
		dbutil.NewDbExecutor(cfg.FocusCtx.DB.Table("member_secret_key").Where("member_id = ? and status = 1", mId).Find(&secretKey))
		cfg.FocusCtx.MemberSecretKey.Store(mId, &memberprikeycache{
			Timestamp: timutil.DefFormat(time.Now().Add(time.Minute * 3)),
			PriKey:    secretKey.PriKey,
			PubKey:    secretKey.PubKey,
		})
		priKey, pubKey = secretKey.PriKey, secretKey.PubKey
	} else {
		memberPriKey := memberPriKeyCache.(*memberprikeycache)
		priKey, pubKey = memberPriKey.PriKey, memberPriKey.PubKey
	}
	return priKey, pubKey
}

func validateParams(g *gtwtype.GtWReq) {
	if strutil.IsBlank(g.Timestamp) {
		types.InvalidParamPanic("timestamp param can't be empty!")
	}
	if strutil.IsBlank(g.Sign) {
		types.InvalidParamPanic("sign param can't be empty!")
	}
	if strutil.IsBlank(g.MemberId) {
		types.InvalidParamPanic("memberId can't be empty!")
	}
	if strutil.IsBlank(g.ServUrl) {
		types.InvalidParamPanic("servUrl can't be empty!")
	}
	if strutil.IsBlank(g.BizContent) {
		types.InvalidParamPanic("bizContent param can't be empty!")
	}
}

type memberprikeycache struct {
	Timestamp string
	PriKey    string
	PubKey    string
}
