package main

import (
	"crypto/aes"
	"crypto/cipher"
)

type Crypto struct {
	stream cipher.Stream
	key    []byte
	iv     []byte
}

func NewCrypto(key, iv []byte) *Crypto {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	return &Crypto{
		stream: cipher.NewCTR(block, iv),
		key:    key,
		iv:     iv,
	}
}

func (c *Crypto) Do(data []byte) []byte {
	ciphertext := make([]byte, len(data))
	c.stream.XORKeyStream(ciphertext, data)
	return ciphertext
}
