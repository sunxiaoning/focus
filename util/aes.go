package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"focus/types"
	"github.com/jmcvetta/randutil"
)

type Encrypter interface {

	// 16, 24, 32
	GenerateKey(int) (key string, err error)
	Encrypt(key string, origin string)(encrypt string, err error)
	Decrypt(key string, encrypt string)(origin string, err error)
}

type AesEncrypter struct {
}

var AESUtil = &AesEncrypter{}

func (aesEncrypter *AesEncrypter) Encrypt(key string, origin string)(encrypt string, err error) {
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	blocksize := block.BlockSize()
	encrypter := cipher.NewCBCEncrypter(block, k[:blocksize])
	originbytes := PKCS7Padding([]byte(origin), blocksize)
	encryptbytes := make([]byte, len(originbytes))
	encrypter.CryptBlocks(encryptbytes, originbytes)
	return base64.StdEncoding.EncodeToString(encryptbytes), nil
}

func (aesEncrypter *AesEncrypter) Decrypt(key string, encrypt string)(origin string, err error) {
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	blocksize := block.BlockSize()
	decryptor := cipher.NewCBCDecrypter(block, k[:blocksize])
	encryptbytes, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", err
	}
	originbytes := make([]byte, len(encryptbytes))
	decryptor.CryptBlocks(originbytes, encryptbytes)
	originbytes = PKCS7UnPadding(originbytes)
	return string(originbytes), nil
}

func (aesEncrypter *AesEncrypter) GenerateKey(n int) (key string, err error)  {
	switch n {
	default:
		return "", types.NewErr(types.InvalidParamError, "key length is invalid!")
	case 16, 24, 32:
		break
	}
	return randutil.AlphaString(n)
}

func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
