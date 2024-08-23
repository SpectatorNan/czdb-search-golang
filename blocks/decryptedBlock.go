package blocks

import (
	"czdb-search-golang/ase/ecb"
	"czdb-search-golang/bytex"
)

type DecryptedBlock struct {
	ClientId       int64
	ExpirationDate int64
	RandomSize     int64
}

func (b *DecryptedBlock) ToBytes() []byte {
	bytes := make([]byte, 16)
	writeIntLong(bytes, 0, (b.ClientId<<20)|b.ExpirationDate)
	writeIntLong(bytes, 8, b.RandomSize)
	return bytes
}

func (b *DecryptedBlock) Encrypt(data []byte, key string) ([]byte, error) {
	return ecb.EncryptWithPkcs5Padding(data, key)
}
func NewDecryptBlock(data []byte, key string) (*DecryptedBlock, error) {
	dbytes, err := ecb.Decrypt(data, key)
	if err != nil {
		return nil, err
	}

	// Parse the decrypted bytes
	decryptedBlock := &DecryptedBlock{}
	decryptedBlock.ClientId = bytex.GetIntLong(dbytes, 0) >> 20
	decryptedBlock.ExpirationDate = bytex.GetIntLong(dbytes, 0) & 0xFFFFF
	decryptedBlock.RandomSize = bytex.GetIntLong(dbytes, 4)

	return decryptedBlock, nil
}
