package app

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron"
	"net/http"
)

type Ctx struct {

	// 环境配置参数
	Cfg *Cfg

	// 数据库
	DB *gorm.DB

	// 路由
	Router *mux.Router

	// http服务
	HttpServer *http.Server

	// 定时任务
	Task *cron.Cron
}

// 应用上下文
var FocusCtx = &Ctx{}
