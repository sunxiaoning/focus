package cfg

var alphaServer = &ServerCfg{
	ListenPort:  7001,
	Env:         ENV_ALPHA,
	LogFilePath: "/Users/william/logs/focus/app.log",
	SecretKey:   alphaSecretKey,
}

var prodServer = &ServerCfg{
	ListenPort:  7001,
	Env:         "prod",
	LogFilePath: "/Users/william/logs/focus/app.log",
	SecretKey:   prodSecretKey,
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
}
