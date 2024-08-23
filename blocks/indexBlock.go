package blocks

type IndexBlock struct {
	startIp []byte
	endIp   []byte
	dataPtr int64
	dataLen int64
	dbType  int
}

func NewIndexBlock(startIp []byte, endIp []byte, dataPtr int64, dataLen int64, dbType int) *IndexBlock {
	return &IndexBlock{startIp: startIp, endIp: endIp, dataPtr: dataPtr, dataLen: dataLen, dbType: dbType}
}

func GetIndexBlockLength(dbType int) int {
	if dbType == 4 {
		return 13
	} else {
		return 37
	}
}
