package servorderrepo

import (
	"context"
	"focus/cfg"
	orderstatusconst "focus/types/consts/orderstatus"
	servtype "focus/types/serv"
	dbutil "focus/util/db"
	"github.com/jinzhu/gorm"
)

const (
	tabName     = "service_order"
	normalQuery = "status = 1"
)

func GetByOrderNo(ctx context.Context, orderNo string) *servtype.OrderEntity {
	db := getDb(ctx)
	db = db.Where("order_no = ?", orderNo)
	db = db.Where(normalQuery)
	var orderEntity servtype.OrderEntity
	dbutil.NewDbExecutor(db.Find(&orderEntity))
	return &orderEntity
}

func Create(ctx context.Context, orderEntity *servtype.OrderEntity) {
	db := getDb(ctx)
	dbutil.NewDbExecutor(db.Create(orderEntity))
}

func Submit(ctx context.Context, id int, outOrderNo string) {
	db := getDb(ctx)
	db = db.Where("id = ? and order_status = 'I' ", id)
	db = db.Where(normalQuery)
	dbutil.NewDbExecutor(db.Update(map[string]interface{}{"order_status": orderstatusconst.P, "out_order_no": outOrderNo}))
}

func getDb(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		tx = tx.Table(tabName)
		return tx
	}
	return cfg.FocusCtx.DB.Table(tabName)
}
