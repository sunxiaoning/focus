package cfg

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron"
	"net/http"
	"sync"
)

const (
	ENV_ALPHA = "alpha"
	ENV_PROD  = "prod"
)

type RuntimeConfig struct {
	Env           string
	SecretKeyPath string
}

type cfg struct {
	Runtime  *RuntimeConfig
	Server   *ServerCfg
	Database *DatabaseConfig
}

var Cfg = &cfg{}

type ctx struct {

	// 环境配置参数
	Cfg *cfg

	// 数据库
	DB *gorm.DB

	// 路由
	Router *mux.Router

	// http服务
	HttpServer *http.Server

	// 定时任务
	Task *cron.Cron

	// 用户信息缓存
	CurrentUser *sync.Map

	// 访问限流器
	VisitorLimiter *sync.Map
}

var UserMap = &sync.Map{}

var VisitorLimiter = &sync.Map{}

// 应用上下文
var FocusCtx = &ctx{Cfg, nil, nil, nil, nil, UserMap, VisitorLimiter}

// 默认配置
type DefaultCfgVal interface{}

// 默认配置文件解密
type DefaultCfgValDecryptor func(cfgVal DefaultCfgVal) (DefaultCfgVal, error)

// 默认配置
type DefaultCfg struct {
	cfgValMap      map[string]DefaultCfgVal
	decryptHandler DefaultCfgValDecryptor
}

func (defaultCfg *DefaultCfg) GetDefaultCfg(env string) (cfgVal DefaultCfgVal, err error) {
	cfgVal = defaultCfg.cfgValMap[ENV_ALPHA]
	if env != "" {
		cfgVal = defaultCfg.cfgValMap[env]
	}
	if defaultCfg.decryptHandler != nil {
		return defaultCfg.decryptHandler(cfgVal)
	}
	return cfgVal, nil
}

func newDefaultCfg(cfgValMap map[string]DefaultCfgVal, decryptHandler DefaultCfgValDecryptor) *DefaultCfg {
	return &DefaultCfg{cfgValMap, decryptHandler}
}
