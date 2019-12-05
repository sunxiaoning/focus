package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"focus/types"
	"github.com/jmcvetta/randutil"
)

type aesEncrypter struct {
}

var AESUtil = &aesEncrypter{}

func (aesEncrypter *aesEncrypter) Encrypt(key string, origin string) (encrypt string, err error) {
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	blocksize := block.BlockSize()
	encrypter := cipher.NewCBCEncrypter(block, k[:blocksize])
	originbytes := pKCS7Padding([]byte(origin), blocksize)
	encryptbytes := make([]byte, len(originbytes))
	encrypter.CryptBlocks(encryptbytes, originbytes)
	return base64.StdEncoding.EncodeToString(encryptbytes), nil
}

func (aesEncrypter *aesEncrypter) Decrypt(key string, encrypt string) (origin string, err error) {
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
	originbytes = pKCS7UnPadding(originbytes)
	return string(originbytes), nil
}

func (aesEncrypter *aesEncrypter) GenerateKey(n int) (key string, err error) {
	switch n {
	default:
		return "", types.NewErr(types.InvalidParamError, "key length is invalid!")
	case 16, 24, 32:
		break
	}
	return randutil.AlphaString(n)
}

func pKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
