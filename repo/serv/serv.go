package servrepo

import (
	"context"
	"focus/cfg"
	pagetype "focus/types/page"
	servtype "focus/types/serv"
	dbutil "focus/util/db"
	"github.com/jinzhu/gorm"
)

const (
	tabName            = "service"
	normalServiceQuery = "service_status = 'FWZ' and status = 1"
)

// 查询最新服务信息
func QueryLatest(ctx context.Context, pageQuery pagetype.PageQuery) (services []*servtype.ServiceEntity, total int) {

	db := getDb(ctx)

	// 参数
	db = db.Where(pageQuery.Params)

	// 服务状态过滤
	db = db.Where(normalServiceQuery)

	// 查询分页
	dbutil.NewDbExecutor(db).PageQuery(ctx, pageQuery.Page, &total, &services)
	return services, total
}

// 获取详情
func GetById(ctx context.Context, serviceId int) *servtype.ServiceEntity {
	db := getDb(ctx)
	db = db.Where("id = ?", serviceId)
	db = db.Where(normalServiceQuery)
	service := &servtype.ServiceEntity{}
	dbutil.NewDbExecutor(db.Find(service))
	return service
}

func getDb(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		tx = tx.Table(tabName)
		return tx
	}
	return cfg.FocusCtx.DB.Table(tabName)
}
