package servrepo

import (
	"focus/cfg"
	pagetype "focus/types/page"
	servicetype "focus/types/service"
	dbutil "focus/util/db"
	"github.com/jinzhu/gorm"
)

const (
	tabName            = "service"
	normalServiceQuery = "service_status = 'FWZ' and status = 1"
)

// 查询最新服务信息
func QueryLatest(pageQuery pagetype.PageQuery, services *[]*servicetype.ServiceEntity, total *int) {

	db := getDb()

	// 参数
	db = db.Where(pageQuery.Params)

	// 服务状态过滤
	db = db.Where(normalServiceQuery)

	// 查询分页
	dbutil.NewDbExecutor(db).PageQuery(pageQuery.Page, total, services)
}

func getDb() *gorm.DB {
	return cfg.FocusCtx.DB.Table(tabName)
}
