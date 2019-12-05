package cfg

import (
	"github.com/mitchellh/go-homedir"
	"path"
)

const (
	SecretPath     = ".secret"
	PubKeyFileName = "rsa_pub_key.pem"
	PriKeyFileName = "rsa_pri_key.pem"
	AesKeyFileName = "aes.key"
)

type RsaKeyPair struct {
	PriKey string
	PubKey string
}

var alphaRsaKeyPair = &RsaKeyPair{}
var prodRsaKeyPair = &RsaKeyPair{}

type SecretKey struct {
	FilePath string
	AesKey   string
	RsAKey   *RsaKeyPair
}

var alphaSecretKey = &SecretKey{
	FilePath: func() string {
		home, err := homedir.Dir()
		if err != nil {
			return ""
		}
		return path.Join(home, SecretPath)
	}(),
	AesKey: "",
	RsAKey: alphaRsaKeyPair,
}

var prodSecretKey = &SecretKey{
	FilePath: func() string {
		home, err := homedir.Dir()
		if err != nil {
			return ""
		}
		return path.Join(home, ".secret")
	}(),
	AesKey: "",
	RsAKey: prodRsaKeyPair,
}

var DefaultSecretKey = newDefaultCfg(map[string]DefaultCfgVal{
	ENV_ALPHA: alphaSecretKey,
	ENV_PROD:  prodSecretKey,
}, nil)
