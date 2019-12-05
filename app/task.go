package app

import (
	"focus/cfg"
	"github.com/robfig/cron"
)

func InitTask() {
	task := cron.New()
	task.AddFunc(cfg.FocusCtx.Cfg.Database.CheckDBIntervalCron, ReConnectDB)
	task.Start()
	defer task.Stop()
	cfg.FocusCtx.Task = task
	select {}
}
