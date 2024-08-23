package ecb

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

type Encryptor ecb

func NewEncryptor(b cipher.Block) cipher.BlockMode {
	return &Encryptor{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

func (x *Encryptor) BlockSize() int { return x.blockSize }
func (x *Encryptor) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("ecb: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("ecb: output smaller than input")
	}

	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

func EncryptWithPkcs5Padding(data []byte, key string) ([]byte, error) {
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
	encryptor := NewEncryptor(block)
	paddedData := pkcs5Padding(data, block.BlockSize())
	encrypted := make([]byte, len(paddedData))
	encryptor.CryptBlocks(encrypted, paddedData)

	return encrypted, nil
}
