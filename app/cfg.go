package app

import (
	"focus/cfg"
	"focus/util"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

func InitCfg(runtimeConfig *cfg.RuntimeConfig) error {
	cfg.Cfg.Runtime = runtimeConfig
	defaultServer, err := cfg.DefaultServer.GetDefaultCfg(runtimeConfig.Env)
	serverConfig := defaultServer.(*cfg.ServerCfg)
	if err != nil {
		return err
	}
	if strings.TrimSpace(runtimeConfig.SecretKeyPath) != "" {
		serverConfig.SecretKey.FilePath = runtimeConfig.SecretKeyPath
	}
	rsaPriKeyFile := path.Join(serverConfig.SecretKey.FilePath, cfg.PriKeyFileName)
	serverConfig.SecretKey.RsAKey.PriKey, err = util.DefaultEncryptor.ParseKeyFromFile(rsaPriKeyFile)
	if err != nil {
		return err
	}
	rsaPubKeyFile := path.Join(serverConfig.SecretKey.FilePath, cfg.PubKeyFileName)
	serverConfig.SecretKey.RsAKey.PubKey, err = util.DefaultEncryptor.ParseKeyFromFile(rsaPubKeyFile)
	if err != nil {
		return err
	}
	aesKeyFile := path.Join(serverConfig.SecretKey.FilePath, cfg.AesKeyFileName)
	aesKeyBytes, err := ioutil.ReadFile(aesKeyFile)
	if err != nil {
		return err
	}
	aesKey, err := util.DefaultEncryptor.Decrypt(serverConfig.SecretKey.RsAKey.PriKey, string(aesKeyBytes))
	if err != nil {
		return err
	}
	serverConfig.SecretKey.AesKey = aesKey
	cfg.Cfg.Server = serverConfig
	defaultDatabase, err := cfg.DefaultDatabase.GetDefaultCfg(runtimeConfig.Env)
	if err != nil {
		return err
	}
	cfg.Cfg.Database = defaultDatabase.(*cfg.DatabaseConfig)
	if err := util.DefaultFileHelper.CreateDirectory(filepath.Dir(defaultServer.(*cfg.ServerCfg).LogFilePath)); err != nil {
		return err
	}
	return nil
}
