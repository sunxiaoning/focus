package memloginrepo

import (
	"context"
	"focus/cfg"
	membertype "focus/types/member"
	dbutil "focus/util/db"
	"github.com/jinzhu/gorm"
)

const (
	tabName     = "member_login"
	normalQuery = "status = 1"
)

func GetByMemberId(ctx context.Context, memberId int) *membertype.CurrentUserInfo {
	db := getDb(ctx)
	db = db.Where("member_id = ?", memberId)
	db = db.Where(normalQuery)
	var currentUser membertype.CurrentUserInfo
	dbutil.NewDbExecutor(db.Find(&currentUser))
	return &currentUser
}

func GetByUsernameAndPwd(ctx context.Context, username string, pwd string) *membertype.CurrentUserInfo {
	db := getDb(ctx)
	db = db.Where("user_name = ? and passwd = ? ", username, pwd)
	var currentUser membertype.CurrentUserInfo
	dbutil.NewDbExecutor(db.Find(&currentUser))
	return &currentUser
}

func getDb(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		tx = tx.Table(tabName)
		return tx
	}
	return cfg.FocusCtx.DB.Table(tabName)
}
