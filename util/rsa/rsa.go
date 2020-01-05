package rsautil

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"focus/types"
	"io/ioutil"
	"strings"
)

const (
	PKCS1 = "PKCS1"
	PKCS8 = "PKCS8"
)

type rsaEncryptor struct {
	keyFormat string
	block     bool
}

var rsaEncryptors = map[string]*rsaEncryptor{
	PKCS1: {PKCS1, false},
	PKCS8: {PKCS8, false},
}

var blockRsaEncryptors = map[string]*rsaEncryptor{
	PKCS1: {PKCS1, true},
	PKCS8: {PKCS8, true},
}

var DefaultEncryptor = &rsaEncryptor{PKCS8, false}

func NewRsaEncryptor(keyFormat string, block bool) (*rsaEncryptor, error) {
	var encryptor *rsaEncryptor
	if strings.TrimSpace(keyFormat) == "" {
		encryptor = rsaEncryptors[PKCS8]
		if block {
			encryptor = blockRsaEncryptors[PKCS8]
		}
	}
	if err := checkKeyFormat(keyFormat); err != nil {
		return nil, err
	}
	encryptor = rsaEncryptors[keyFormat]
	if block {
		encryptor = blockRsaEncryptors[keyFormat]
	}
	return encryptor, nil
}

func (rsaEncryptor *rsaEncryptor) GenerateKeyPair(length int) (priKey string, pubKey string, err error) {
	var keyPairGen types.KeyPairGenerator
	keyPairGen, err = newEncrypter(rsaEncryptor.keyFormat)
	if err != nil {
		return "", "", err
	}
	return keyPairGen.GenerateKeyPair(2048)
}

func (rsaEncryptor *rsaEncryptor) Encrypt(pubKey string, origin string) (encrypt string, err error) {
	var encrypter types.Encrypter
	if rsaEncryptor.block {
		encrypter, err = newBlockEncrypter(rsaEncryptor.keyFormat)
	} else {
		encrypter, err = newEncrypter(rsaEncryptor.keyFormat)
	}
	if err != nil {
		return "", err
	}
	return encrypter.Encrypt(pubKey, origin)
}

func (rsaEncryptor *rsaEncryptor) Decrypt(priKey string, encrypt string) (origin string, err error) {
	var encrypter types.Encrypter
	if rsaEncryptor.block {
		encrypter, err = newBlockEncrypter(rsaEncryptor.keyFormat)
	} else {
		encrypter, err = newEncrypter(rsaEncryptor.keyFormat)
	}
	if err != nil {
		return "", err
	}
	return encrypter.Decrypt(priKey, encrypt)
}

func Encrypt(key string, origin string) (encrypt string, err error) {
	return DefaultEncryptor.Encrypt(key, origin)
}
func Decrypt(key string, encrypt string) (origin string, err error) {
	return DefaultEncryptor.Decrypt(key, encrypt)
}

func (rsaEncryptor *rsaEncryptor) Sign(originData string, prvKey string) (string, error) {
	priKeyBytes, err := base64.StdEncoding.DecodeString(prvKey)
	if err != nil {
		return "", err
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(priKeyBytes)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	h.Write([]byte([]byte(originData)))
	hashed := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}
func (rsaEncryptor *rsaEncryptor) VerifySign(originData string, sign string, pubKey string) (bool, error) {
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}
	public, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return false, err
	}
	pub, err := x509.ParsePKIXPublicKey(public)
	if err != nil {
		return false, err
	}
	hashed := sha256.Sum256([]byte(originData))
	err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, hashed[:], signBytes)
	if err != nil {
		return false, err
	}
	return true, nil
}

func Sign(originData string, prvKey string) (string, error) {
	return DefaultEncryptor.Sign(originData, prvKey)
}
func VerifySign(originData string, sign string, pubKey string) (bool, error) {
	return DefaultEncryptor.VerifySign(originData, sign, pubKey)
}

func (rsaEncryptor *rsaEncryptor) ParseKeyFromFile(fileName string) (string, error) {
	priKeyBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	block, rest := pem.Decode(priKeyBytes)
	if block == nil && len(rest) == len(priKeyBytes) {
		return "", types.NewErr(types.InvalidParamError, "priKey format error!")
	}
	return base64.StdEncoding.EncodeToString(block.Bytes), nil
}

func (rsaEncryptor *rsaEncryptor) ParseKeyFromString(str string) (string, error) {
	block, rest := pem.Decode([]byte(str))
	if block == nil && len(rest) == len([]byte(str)) {
		return "", types.NewErr(types.InvalidParamError, "priKey format error!")
	}
	return base64.StdEncoding.EncodeToString(block.Bytes), nil
}

func ParseKeyFromFile(fileName string) (string, error) {
	return DefaultEncryptor.ParseKeyFromFile(fileName)
}

func ParseKeyFromString(str string) (string, error) {
	return DefaultEncryptor.ParseKeyFromString(str)
}

func checkKeyFormat(keyFormat string) error {
	if strings.TrimSpace(keyFormat) != PKCS1 && strings.TrimSpace(keyFormat) != PKCS8 {
		return types.NewErr(types.InvalidParamError, "keyFormat is invalid!")
	}
	return nil
}

type encrypter struct {
	keyFormat string
}

var encryptors = map[string]*encrypter{
	PKCS1: {PKCS1},
	PKCS8: {PKCS8},
}

func newEncrypter(keyFormat string) (*encrypter, error) {
	if err := checkKeyFormat(keyFormat); err != nil {
		return nil, err
	}
	return encryptors[keyFormat], nil
}

func (encrypter *encrypter) GenerateKeyPair(length int) (priKey string, pubKey string, err error) {
	key, err := rsa.GenerateKey(rand.Reader, length)
	if err != nil {
		return "", "", err
	}
	var priKeyBytes, pubKeyBytes []byte
	if encrypter.keyFormat == PKCS8 {
		priKeyBytes, err = x509.MarshalPKCS8PrivateKey(key)
		if err != nil {
			return "", "", err
		}
		pubKeyBytes, err = x509.MarshalPKIXPublicKey(&key.PublicKey)
		if err != nil {
			return "", "", err
		}
	} else {
		priKeyBytes = x509.MarshalPKCS1PrivateKey(key)
		pubKeyBytes = x509.MarshalPKCS1PublicKey(&key.PublicKey)
	}
	priKey = base64.StdEncoding.EncodeToString(priKeyBytes)
	pubKey = base64.StdEncoding.EncodeToString(pubKeyBytes)
	return priKey, pubKey, nil
}

func (encrypter *encrypter) Encrypt(pubKey string, origin string) (encrypt string, err error) {
	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return "", err
	}
	var key interface{}
	key, err = x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	encryptBytes, err := rsa.EncryptPKCS1v15(rand.Reader, key.(*rsa.PublicKey), []byte(origin))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptBytes), nil
}

func (encrypter *encrypter) Decrypt(priKey string, encrypt string) (origin string, err error) {
	priKeyBytes, err := base64.StdEncoding.DecodeString(priKey)
	if err != nil {
		return "", err
	}
	encryptBytes, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", err
	}
	var key interface{}
	if encrypter.keyFormat == PKCS8 {
		key, err = x509.ParsePKCS8PrivateKey(priKeyBytes)
	} else {
		key, err = x509.ParsePKCS1PrivateKey(priKeyBytes)
	}
	if err != nil {
		return "", err
	}
	decryptBytes, err := rsa.DecryptPKCS1v15(rand.Reader, key.(*rsa.PrivateKey), encryptBytes)
	if err != nil {
		return "", err
	}
	return string(decryptBytes), nil
}

type blockEncrypter encrypter

var blockEncryptors = map[string]*blockEncrypter{
	PKCS1: {PKCS1},
	PKCS8: {PKCS8},
}

func newBlockEncrypter(keyFormat string) (*blockEncrypter, error) {
	if err := checkKeyFormat(keyFormat); err != nil {
		return nil, err
	}
	return blockEncryptors[keyFormat], nil
}

func (blockEncrypter *blockEncrypter) Encrypt(pubKey string, origin string) (encrypt string, err error) {
	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return "", err
	}
	originBytes := []byte(origin)
	if err != nil {
		return "", err
	}
	var key interface{}
	key, err = x509.ParsePKIXPublicKey(pubKeyBytes)
	keySize, srcSize := key.(*rsa.PublicKey).Size(), len(originBytes)
	offSet, once := 0, keySize-11
	buffer := bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + once
		if endIndex > srcSize {
			endIndex = srcSize
		}
		bytesOnce, err := rsa.EncryptPKCS1v15(rand.Reader, key.(*rsa.PublicKey), originBytes[offSet:endIndex])
		if err != nil {
			return "", err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func (blockEncrypter *blockEncrypter) Decrypt(priKey string, encrypt string) (origin string, err error) {
	priKeyBytes, err := base64.StdEncoding.DecodeString(priKey)
	if err != nil {
		return "", err
	}
	encryptBytes, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", err
	}
	var key interface{}
	if blockEncrypter.keyFormat == PKCS8 {
		key, err = x509.ParsePKCS8PrivateKey(priKeyBytes)
	} else {
		key, err = x509.ParsePKCS1PrivateKey(priKeyBytes)
	}
	if err != nil {
		return "", err
	}
	keySize := key.(*rsa.PrivateKey).Size()
	srcSize := len(encryptBytes)
	var offSet = 0
	var buffer = bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + keySize
		if endIndex > srcSize {
			endIndex = srcSize
		}
		bytesOnce, err := rsa.DecryptPKCS1v15(rand.Reader, key.(*rsa.PrivateKey), encryptBytes[offSet:endIndex])
		if err != nil {
			return "", err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	return string(buffer.Bytes()), nil
}
