package blocks

import (
	"bytes"
	"errors"
	"github.com/vmihailenco/msgpack/v5"
)

type DataBlock struct {
	region  []byte
	dataPtr int64
}

func NewDataBlock(region []byte, dataPtr int64) *DataBlock {
	return &DataBlock{region: region, dataPtr: dataPtr}
}

func (b *DataBlock) GetRegion(geoMapData []byte, columnSelection int64) (string, error) {
	return b.unpack(geoMapData, columnSelection)
}

func (b *DataBlock) unpack(geoMapData []byte, columnSelection int64) (string, error) {
	var geoPosMixSize int64
	var otherData string

	regionUnpacker := msgpack.NewDecoder(bytes.NewReader(b.region))
	if err := regionUnpacker.Decode(&geoPosMixSize); err != nil {
		return "", err
	}
	if err := regionUnpacker.Decode(&otherData); err != nil {
		return "", err
	}

	if geoPosMixSize == 0 {
		return otherData, nil
	}

	dataLen := int((geoPosMixSize >> 24) & 0xFF)
	dataPtr := int(geoPosMixSize & 0x00FFFFFF)

	if dataPtr+dataLen > len(geoMapData) {
		return "", errors.New("data pointer out of bounds")
	}

	regionData := geoMapData[dataPtr : dataPtr+dataLen]
	geoColumnUnpacker := msgpack.NewDecoder(bytes.NewReader(regionData))

	columnNumber, err := geoColumnUnpacker.DecodeArrayLen()
	if err != nil {
		return "", err
	}

	var sb bytes.Buffer
	for i := 0; i < columnNumber; i++ {
		columnSelected := (columnSelection>>(i+1))&1 == 1
		var value string
		if err := geoColumnUnpacker.Decode(&value); err != nil {
			return "", err
		}
		if value == "" {
			value = "null"
		}
		if columnSelected {
			sb.WriteString(value)
			sb.WriteString("\t")
		}
	}

	return sb.String() + otherData, nil
}
