package servorderrepo

import (
	"context"
	"focus/cfg"
	orderstatusconst "focus/types/consts/orderstatus"
	servtype "focus/types/serv"
	dbutil "focus/util/db"
	"github.com/jinzhu/gorm"
	"time"
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

func GetWaittingOrderByPayOrderNo(ctx context.Context, payOrderNo string) *servtype.OrderEntity {
	db := getDb(ctx)
	db = db.Where("order_no = ? and order_status = 'P'", payOrderNo)
	db = db.Where(normalQuery)
	var orderEntity servtype.OrderEntity
	dbutil.NewDbExecutor(db.Find(&orderEntity))
	return &orderEntity
}

func UpdateOrderStatusByPayOrderNoAndPayResult(ctx context.Context, payOrderNo string, toStatus string) int64 {
	db := getDb(ctx)
	db = db.Where("out_order_no = ?", payOrderNo)
	db = db.Where("order_status = 'P'")
	db = db.Where(normalQuery)
	return dbutil.NewDbExecutor(db.Updates(servtype.OrderEntity{OrderStatus: toStatus, FinishedTime: time.Now()})).RowsAffected()
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
