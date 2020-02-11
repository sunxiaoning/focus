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
	db := getDb(ctx)
	db = db.Where("service_id = ?", serviceId)
	db = db.Where(normalPriceQuery)
	var prices []*servtype.PriceEntity
	dbutil.NewDbExecutor(db.Find(&prices))
	return prices
}

func GetById(ctx context.Context, id int) *servtype.PriceEntity {
	db := getDb(ctx)
	db = db.Where("id = ?", id)
	db = db.Where(normalPriceQuery)
	var priceEntity servtype.PriceEntity
	dbutil.NewDbExecutor(db.Find(&priceEntity))
	return &priceEntity
}

func getDb(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		tx = tx.Table(tabName)
		return tx
	}
	return cfg.FocusCtx.DB.Table(tabName)
}
