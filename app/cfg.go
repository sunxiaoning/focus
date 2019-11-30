package app

import (
	"focus/types"
	"focus/util"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	LogFilePath = "/Users/william/logs/focus/app.log"
)

type RuntimeConfig struct {
	Env string
	AesKey string
}

type ServerConfig struct {

	// 服务器监听端口
	ListenPort int

	// 服务器运行环境
	Env string

	// 日志文件路径
	LogFilePath string
}

var defaultServer = &ServerConfig{
	ListenPort: 7001,
	Env:        "alpha",
}

type DatabaseConfig struct {
	Host                string
	Port                string
	DBName              string
	Username            string
	Password            string
	ConnMaxLifetime     time.Duration
	CheckDBIntervalCron string
	MaxIdleConns        int
	MaxOpenConns        int
}

var defaultDatabase = &DatabaseConfig{
	Host: "wuPq8Hw7CcWds+ou4mpb+FfO2m+Ga+7xzdNmKdLBu+A=",
	Port: "++OyfMhlTMCqbmzDf3L/mA==",
	DBName:              "gRAe9ihlDR8E66uLy0+avA==",
	Username:            "4cRAe7/oBdX0EwSFtfypBA==",
	Password:            "xpDXApIMNdkMS5HjrO766g==",
	ConnMaxLifetime:     time.Second * 10,
	CheckDBIntervalCron: "0/10 * * * * ?",
	MaxIdleConns:        10,
	MaxOpenConns:        200,
}

type Cfg struct {
	Server   *ServerConfig
	Database *DatabaseConfig
}

func InitCfg(runtimeConfig *RuntimeConfig) error {
	if strings.TrimSpace(runtimeConfig.Env) != "" {
		defaultServer.Env = runtimeConfig.Env
	}
	if defaultServer.Env == "prod" {
		if err := os.MkdirAll(filepath.Dir(LogFilePath), 0755); err != nil {
			return err
		}
		defaultServer.LogFilePath = LogFilePath
	}
	if len(strings.TrimSpace(runtimeConfig.AesKey)) <= 0 {
		return types.NewErr(types.SystemError, "aeskey can't be empty!")
	}
	if err := decrpytCfg(runtimeConfig) ; err != nil {
		return err
	}
	FocusCtx.Cfg = &Cfg{
		Server:   defaultServer,
		Database: defaultDatabase,
	}
	return nil
}

func decrpytCfg(runtimeConfig *RuntimeConfig) error {
	host, err := util.AESUtil.Decrypt(runtimeConfig.AesKey, defaultDatabase.Host)
	if err != nil {
		return err
	}
	defaultDatabase.Host = host
	port, err := util.AESUtil.Decrypt(runtimeConfig.AesKey, defaultDatabase.Port)
	if err != nil {
		return err
	}
	defaultDatabase.Port = port
	dbname, err := util.AESUtil.Decrypt(runtimeConfig.AesKey, defaultDatabase.DBName)
	if err != nil {
		return err
	}
	defaultDatabase.DBName = dbname
	username, err := util.AESUtil.Decrypt(runtimeConfig.AesKey, defaultDatabase.Username)
	if err != nil {
		return err
	}
	defaultDatabase.Username = username
	dbpasswd, err := util.AESUtil.Decrypt(runtimeConfig.AesKey, defaultDatabase.Password)
	if err != nil {
		return err
	}
	defaultDatabase.Password = dbpasswd
	return nil
}
