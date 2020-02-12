package ppaynotifyrepo

import (
	"context"
	"focus/cfg"
	ppaytype "focus/types/ppay"
	dbutil "focus/util/db"
	"github.com/jinzhu/gorm"
)

const (
	tabName     = "personal_pay_notify"
	normalQuery = "status = 1"
)

func Create(ctx context.Context, payNotifyEntity *ppaytype.PPayNotifyEntity) {
	db := getDb(ctx)
	dbutil.NewDbExecutor(db.Create(payNotifyEntity))
}

func getDb(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		tx = tx.Table(tabName)
		return tx
	}
	return cfg.FocusCtx.DB.Table(tabName)
}
