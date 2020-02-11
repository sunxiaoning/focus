package memberrepo

import (
	"context"
	"focus/cfg"
	membertype "focus/types/member"
	dbutil "focus/util/db"
	"github.com/jinzhu/gorm"
)

const (
	tabName     = "member"
	normalQuery = "status = 1"
)

func GetById(ctx context.Context, id int) *membertype.MemberEntity {
	db := getDb(ctx)
	db = db.Where("id = ?", id)
	db = db.Where(normalQuery)
	var memberEntity membertype.MemberEntity
	dbutil.NewDbExecutor(db.Find(&memberEntity))
	return &memberEntity
}

func getDb(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		tx = tx.Table(tabName)
		return tx
	}
	return cfg.FocusCtx.DB.Table(tabName)
}
