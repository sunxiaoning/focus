package types

type Encrypter interface {
	Encrypt(key string, origin string) (encrypt string, err error)
	Decrypt(key string, encrypt string) (origin string, err error)
}

type KeyGenerator interface {

	// AES: 16, 24, 32
	GenerateKey(length int) (key string, err error)
}

type KeyPairGenerator interface {

	// RSA: 2048, 4096
	GenerateKeyPair(length int) (priKey string, pubKey string, err error)
}

type Signer interface {
	Sign(originData string, prvKey string) (string, error)
	VerifySign(originData string, sign string, pubKey string) (bool, error)
}
