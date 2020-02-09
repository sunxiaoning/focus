package servpricerepo

import (
	"context"
	"focus/cfg"
	servtype "focus/types/serv"
	dbutil "focus/util/db"
	"github.com/jinzhu/gorm"
)

const (
	tabName          = "service_price"
	normalPriceQuery = "status = 1"
)

func QueryByServiceId(ctx context.Context, serviceId int) []*servtype.PriceEntity {
	db := getDb()
	db = db.Where("service_id = ?", serviceId)
	db = db.Where(normalPriceQuery)
	var prices []*servtype.PriceEntity
	dbutil.NewDbExecutor(db.Find(&prices))
	return prices
}

func getDb() *gorm.DB {
	return cfg.FocusCtx.DB.Table(tabName)
}
