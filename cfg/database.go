package cfg

import (
	"focus/util"
	"time"
)

var alphaDatabase = &DatabaseConfig{
	Host:                "wuPq8Hw7CcWds+ou4mpb+FfO2m+Ga+7xzdNmKdLBu+A=",
	Port:                "++OyfMhlTMCqbmzDf3L/mA==",
	DBName:              "gRAe9ihlDR8E66uLy0+avA==",
	Username:            "4cRAe7/oBdX0EwSFtfypBA==",
	Password:            "xpDXApIMNdkMS5HjrO766g==",
	ConnMaxLifetime:     time.Second * 10,
	CheckDBIntervalCron: "0/10 * * * * ?",
	MaxIdleConns:        10,
	MaxOpenConns:        200,
}

var prodDatabase = &DatabaseConfig{
	Host:                "wuPq8Hw7CcWds+ou4mpb+FfO2m+Ga+7xzdNmKdLBu+A=",
	Port:                "++OyfMhlTMCqbmzDf3L/mA==",
	DBName:              "gRAe9ihlDR8E66uLy0+avA==",
	Username:            "4cRAe7/oBdX0EwSFtfypBA==",
	Password:            "xpDXApIMNdkMS5HjrO766g==",
	ConnMaxLifetime:     time.Second * 10,
	CheckDBIntervalCron: "0/10 * * * * ?",
	MaxIdleConns:        10,
	MaxOpenConns:        200,
}

var DefaultDatabase = newDefaultCfg(map[string]DefaultCfgVal{
	ENV_ALPHA: alphaDatabase,
	ENV_PROD:  prodDatabase,
}, decrpytCfg)

func decrpytCfg(cfgVal DefaultCfgVal) (DefaultCfgVal, error) {
	key := FocusCtx.Cfg.Server.SecretKey.AesKey
	databaseCfg := cfgVal.(*DatabaseConfig)
	host, err := util.AESUtil.Decrypt(key, databaseCfg.Host)
	if err != nil {
		return nil, err
	}
	databaseCfg.Host = host
	port, err := util.AESUtil.Decrypt(key, databaseCfg.Port)
	if err != nil {
		return nil, err
	}
	databaseCfg.Port = port
	dbname, err := util.AESUtil.Decrypt(key, databaseCfg.DBName)
	if err != nil {
		return nil, err
	}
	databaseCfg.DBName = dbname
	username, err := util.AESUtil.Decrypt(key, databaseCfg.Username)
	if err != nil {
		return nil, err
	}
	databaseCfg.Username = username
	dbpasswd, err := util.AESUtil.Decrypt(key, databaseCfg.Password)
	if err != nil {
		return nil, err
	}
	databaseCfg.Password = dbpasswd
	return databaseCfg, err
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
