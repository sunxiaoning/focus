package ppayorderrepo

import (
	"context"
	"focus/cfg"
	ppaytype "focus/types/ppay"
	dbutil "focus/util/db"
	"github.com/jinzhu/gorm"
)

const (
	tabName     = "personal_pay_order"
	normalQuery = "status = 1"
)

func GetByOutTradeNo(ctx context.Context, outTradeNo string) *ppaytype.PPayOrderEntity {
	db := getDb(ctx)
	db = db.Where("out_trade_no = ?", outTradeNo)
	db = db.Where(normalQuery)
	var payOrderEntity ppaytype.PPayOrderEntity
	dbutil.NewDbExecutor(db.Find(&payOrderEntity))
	return &payOrderEntity
}

func Create(ctx context.Context, payOrder *ppaytype.PPayOrderEntity) {
	db := getDb(ctx)
	dbutil.NewDbExecutor(db.Create(payOrder))
}

func getDb(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		tx = tx.Table(tabName)
		return tx
	}
	return cfg.FocusCtx.DB.Table(tabName)
}
