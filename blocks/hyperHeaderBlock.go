package blocks

import (
	"czdb-search-golang/bytex"
	"errors"
	"io"
	"strconv"
	"time"
)

type HyperHeaderBlock struct {
	Version            int64
	ClientId           int64
	EncryptedBlockSize int64
	EncryptedData      []byte
	DecryptedBlock     DecryptedBlock
}

func (b *HyperHeaderBlock) ToBytes() []byte {
	bs := make([]byte, 12)
	writeIntLong(bs, 0, b.Version)
	writeIntLong(bs, 4, b.ClientId)
	writeIntLong(bs, 8, b.EncryptedBlockSize)
	return bs
}

func (b *HyperHeaderBlock) GetHeaderSize() int64 {
	return 12 + b.EncryptedBlockSize + b.DecryptedBlock.RandomSize
}

func NewHyperHeaderBlockFromBytes(bs []byte) *HyperHeaderBlock {
	version := bytex.GetIntLong(bs, 0)
	clientId := bytex.GetIntLong(bs, 4)
	encryptedBlockSize := bytex.GetIntLong(bs, 8)

	return &HyperHeaderBlock{
		Version:            version,
		ClientId:           clientId,
		EncryptedBlockSize: encryptedBlockSize,
	}
}

func NewHyperHeaderBlockFromDecrypt(r io.Reader, key string) (*HyperHeaderBlock, error) {
	headerBytes := make([]byte, 12)
	if _, err := io.ReadFull(r, headerBytes); err != nil {
		return nil, err
	}

	version := bytex.GetIntLong(headerBytes, 0)
	clientId := bytex.GetIntLong(headerBytes, 4)
	encryptedBlockSize := bytex.GetIntLong(headerBytes, 8)

	encryptedBytes := make([]byte, encryptedBlockSize)
	if _, err := io.ReadFull(r, encryptedBytes); err != nil {
		return nil, err
	}

	decryptedBlock, err := NewDecryptBlock(encryptedBytes, key)
	if err != nil {
		return nil, err
	}
	if decryptedBlock.ClientId != clientId {
		return nil, errors.New("client id mismatch")
	}

	currentDate := time.Now().Format("060102")
	currentDateInt, err := strconv.Atoi(currentDate)
	if err != nil {
		return nil, err
	}
	expirationDate := decryptedBlock.ExpirationDate
	if expirationDate < int64(currentDateInt) {
		return nil, errors.New("DB is expired")
	}

	return &HyperHeaderBlock{
		Version:            version,
		ClientId:           clientId,
		EncryptedBlockSize: encryptedBlockSize,
		DecryptedBlock:     *decryptedBlock,
	}, nil

}
