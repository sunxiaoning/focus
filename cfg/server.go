package cfg

import (
	"github.com/mitchellh/go-homedir"
	"path"
)

var alphaServer = &ServerCfg{
	ListenPort:   7001,
	Env:          ENV_ALPHA,
	LogFilePath:  "/Users/william/logs/focus/app.log",
	SecretKey:    alphaSecretKey,
	RootFilePath: "/Users/william/file",
}

var prodServer = &ServerCfg{
	ListenPort:   7001,
	Env:          "prod",
	LogFilePath:  "/Users/william/logs/focus/app.log",
	SecretKey:    prodSecretKey,
	RootFilePath: getRootFilePath(),
}

func getRootFilePath() string {
	home, err := homedir.Dir()
	if err != nil {
		return ""
	}
	return path.Join(home, ".files")
}

var DefaultServer = newDefaultCfg(map[string]DefaultCfgVal{
	ENV_ALPHA: alphaServer,
	ENV_PROD:  prodServer}, nil)

type ServerCfg struct {

	// 服务器监听端口
	ListenPort int

	// 服务器运行环境
	Env string

	// 日志文件路径
	LogFilePath string

	// 密钥
	SecretKey *SecretKey

	// 文件存储路径
	RootFilePath string
}
