package memservrepo

import (
	"context"
	"focus/cfg"
	servtype "focus/types/serv"
	dbutil "focus/util/db"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	tabName     = "member_service"
	normalQuery = "status = 1"
)

func GetByMemIdAndServPriceId(ctx context.Context, memId int, servPriceId int) *servtype.MemberServiceEntity {
	db := getDb(ctx)
	db = db.Where("member_id = ? and service_price_id = ?", memId, servPriceId)
	db = db.Where(normalQuery)
	var memberServiceEntity servtype.MemberServiceEntity
	dbutil.NewDbExecutor(db.Find(&memberServiceEntity))
	return &memberServiceEntity
}

func Create(ctx context.Context, memserv *servtype.MemberServiceEntity) {
	db := getDb(ctx)
	dbutil.NewDbExecutor(db.Create(memserv))
}

func UpdateMemServiceStatus(ctx context.Context, id int, orderId int, deadlineTime time.Time) {
	db := getDb(ctx)
	db = db.Where("id = ?", id)
	db = db.Where(normalQuery)
	dbutil.NewDbExecutor(db.Update(map[string]interface{}{
		"order_id":      orderId,
		"deadline_time": deadlineTime,
	}))

}

func getDb(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		tx = tx.Table(tabName)
		return tx
	}
	return cfg.FocusCtx.DB.Table(tabName)
}
