package ecb

import (
	"bytes"
	"crypto/cipher"
)

type ecb struct {
	b         cipher.Block
	blockSize int
}

func pkcs5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}
