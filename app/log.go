package app

import (
	"github.com/scue/go-logrotate"
	"github.com/sirupsen/logrus"
)

const DefaultTimeFormat = "2006-01-02 15:04:05.000"

func InitLog() {
	serverConfig := FocusCtx.Cfg.Server
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: DefaultTimeFormat,
	})
	if serverConfig.Env == "prod" {
		logrus.SetLevel(logrus.InfoLevel)

		// 每天凌晨备份日志文件，日志最多保留30天
		writer := logrotate.New(serverConfig.LogFilePath, "0 0 0 1/1 * ?", 30)
		logrus.SetOutput(writer)
		go writer.CronTask()
	}
}
