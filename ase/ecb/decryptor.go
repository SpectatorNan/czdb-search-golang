package ecb

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

type Decryptor ecb

func NewDecryptor(b cipher.Block) *Decryptor {
	return &Decryptor{b: b, blockSize: b.BlockSize()}
}
func (x *Decryptor) BlockSize() int { return x.blockSize }
func (x *Decryptor) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("ecb: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("ecb: output smaller than input")
	}

	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type DecryptedBlock struct {
	ClientId       int
	ExpirationDate int
	RandomSize     int
}

func Decrypt(encryptedBytes []byte, key string) ([]byte, error) {
	// Decode the base64 encoded key
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	// Check key length
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		return nil, errors.New("invalid key size")
	}

	// Initialize the Cipher with AES encryption
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	// Use ECB mode with PKCS5 padding
	decryptor := NewDecryptor(block)
	decryptedBytes := make([]byte, len(encryptedBytes))
	decryptor.CryptBlocks(decryptedBytes, encryptedBytes)
	return decryptedBytes, nil
}
