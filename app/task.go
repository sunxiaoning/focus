package app

import (
	"focus/cfg"
	ppayserv "focus/serv/ppay"
	"github.com/robfig/cron"
)

func InitTask() {
	task := cron.New()
	task.AddFunc(cfg.FocusCtx.Cfg.Database.CheckDBIntervalCron, ReConnectDB)

	// 10s 执行一次，通知Biz 支付结果
	task.AddFunc("0/10 * * * * ?", ppayserv.NotifyBiz)
	task.Start()
	defer task.Stop()
	cfg.FocusCtx.Task = task
	select {}
}

func InitGtwTask() {
	task := cron.New()
	task.AddFunc(cfg.FocusCtx.Cfg.Database.CheckDBIntervalCron, ReConnectDB)
	task.Start()
	defer task.Stop()
	cfg.FocusCtx.Task = task
	select {}
}
