package aesutil

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	key, err := GenerateKey(32)
	if err != nil {
		t.Errorf("generate key error! %v", err)
	}
	t.Logf("key:%s", key)
	origin := "123456"
	t.Logf("origin:%s", origin)
	encrypt, err := Encrypt(key, origin)
	if err != nil {
		t.Errorf("encrypt error! %v", err)
	}
	t.Logf("encrypt:%s", encrypt)
	decrypt, err := Decrypt(key, encrypt)
	if err != nil {
		t.Errorf("decrypt error!%v", err)
	}
	t.Logf("decrypt:%s", decrypt)
	if decrypt != origin {
		t.Error("decrypt don't equal to origin!")
	}
}
