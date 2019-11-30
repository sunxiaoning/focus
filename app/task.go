package app

import (
	"github.com/robfig/cron"
)

func InitTask() {
	task := cron.New()
	task.AddFunc(FocusCtx.Cfg.Database.CheckDBIntervalCron, ReConnectDB)
	task.Start()
	defer task.Stop()
	FocusCtx.Task = task
	select {}
}
