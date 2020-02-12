package ppayorderrepo

import (
	"context"
	"focus/cfg"
	orderstatusconst "focus/types/consts/orderstatus"
	ppaytype "focus/types/ppay"
	dbutil "focus/util/db"
	"github.com/jinzhu/gorm"
	"time"
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

func GetByOrderNo(ctx context.Context, payOrderNo string) *ppaytype.PPayOrderEntity {
	db := getDb(ctx)
	db = db.Where("pay_order_no = ?", payOrderNo)
	db = db.Where(normalQuery)
	var payOrderEntity ppaytype.PPayOrderEntity
	dbutil.NewDbExecutor(db.Find(&payOrderEntity))
	return &payOrderEntity
}

func GetWaittingByPayChannelAndAmount(ctx context.Context, payAmount string, payChannel string) *ppaytype.PPayOrderEntity {
	db := getDb(ctx)
	db = db.Where("pay_amount = ? and pay_channel = ?", payAmount, payChannel)
	db = db.Where("pay_status = 'P'")
	db = db.Where(normalQuery)
	var payOrderEntity ppaytype.PPayOrderEntity
	dbutil.NewDbExecutor(db.Find(&payOrderEntity))
	return &payOrderEntity
}

func UpdatePayResult(ctx context.Context, id int, payStatus string) int64 {
	db := getDb(ctx)
	db = db.Where("id = ?", id)
	db = db.Where("pay_status = 'P'")
	db = db.Where(normalQuery)
	return dbutil.NewDbExecutor(db.Update(ppaytype.PPayOrderEntity{PayStatus: orderstatusconst.S, FinishTime: time.Now()})).RowsAffected()
}

func Submit(ctx context.Context, payOrderNo string) int64 {
	db := getDb(ctx)
	db = db.Where("pay_order_no = ?", payOrderNo)
	db = db.Where("pay_status = 'I'")
	return dbutil.NewDbExecutor(db.Update("pay_status", "P")).RowsAffected()
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
