package czdb

import (
	"czdb-search-golang/blocks"
	"czdb-search-golang/bytex"
	"errors"
	"fmt"
	"net"
	"os"
)

type DBSearcher struct {
	queryType            QueryType
	totalHeaderBlockSize int64
	dbType               DBType
	//raf                  *os.File
	czFile           *Cz88File
	HeaderSip        [][]byte
	HeaderPtr        []int
	headerLength     int
	columnSelection  int64
	firstIndexPtr    int64
	totalIndexBlocks int64
	dbBinStr         []byte
	geoMapData       []byte
}

func NewDBSearcher(dbFile string, key string, queryType QueryType) (*DBSearcher, error) {

	raf, err := os.Open(dbFile)
	if err != nil {
		return nil, err
	}

	headerBlock, err := blocks.NewHyperHeaderBlockFromDecrypt(raf, key)
	if err != nil {
		return nil, err
	}

	czFile := NewCz88File(raf, headerBlock.GetHeaderSize())

	czFile.Seek(0)
	superBytes := make([]byte, SUPER_PART_LENGTH)
	_, err = czFile.Read(superBytes)
	if err != nil {
		return nil, err
	}

	dbType := DBType_IPv4
	if superBytes[0]&1 != 0 {
		dbType = DBType_IPv6
	}

	searcher := &DBSearcher{
		queryType: queryType,
		dbType:    dbType,
		czFile:    czFile,
	}

	err = searcher.loadGetSetting(key)
	if err != nil {
		return nil, err
	}
	if queryType == QueryType_Memory {
		err = searcher.initializeMemorySearch()
		if err != nil {
			return nil, err
		}
	} else {
		err = searcher.initBtreeModeParam()
		if err != nil {
			return nil, err
		}
	}

	return searcher, nil
}

func (s *DBSearcher) loadGetSetting(key string) error {
	_, err := s.czFile.Seek(END_INDEX_PTR)
	if err != nil {
		return err
	}

	data := make([]byte, 4)
	_, err = s.czFile.Read(data)
	if err != nil {
		return err
	}

	endIndexPtr := bytex.GetIntLong(data, 0)

	columnSelectionPtr := endIndexPtr + int64(blocks.GetIndexBlockLength(int(s.dbType)))
	_, err = s.czFile.Seek(columnSelectionPtr)
	if err != nil {
		return err
	}

	_, err = s.czFile.Read(data)
	if err != nil {
		return err
	}

	s.columnSelection = bytex.GetIntLong(data, 0)
	if s.columnSelection == 0 {
		return nil
	}

	geoMapPtr := columnSelectionPtr + 4
	_, err = s.czFile.Seek(geoMapPtr)
	if err != nil {
		return err
	}

	_, err = s.czFile.Read(data)
	if err != nil {
		return err
	}

	geoMapSize := bytex.GetIntLong(data, 0)
	_, err = s.czFile.Seek(geoMapPtr + 4)
	if err != nil {
		return err
	}

	geoMapData := make([]byte, geoMapSize)
	_, err = s.czFile.Read(geoMapData)
	if err != nil {
		return err
	}

	d := NewDecryptor(key)
	s.geoMapData = d.Decrypt(geoMapData)

	return nil
}

func (s *DBSearcher) initializeMemorySearch() error {

	fileSize := s.czFile.Length()
	s.dbBinStr = make([]byte, fileSize)

	_, err := s.czFile.Seek(0)
	if err != nil {
		return err
	}

	_, err = s.czFile.Read(s.dbBinStr)
	if err != nil {
		return err
	}

	err = s.czFile.Close()
	if err != nil {
		return err
	}

	return s.initMemoryOrBinaryModeParam(s.dbBinStr, int(fileSize))
}

func (s *DBSearcher) initMemoryOrBinaryModeParam(bytes []byte, fileSize int) error {
	s.totalHeaderBlockSize = bytex.GetIntLong(bytes, HEADER_BLOCK_PTR)
	fileSizeInFile := bytex.GetIntLong(bytes, FILE_SIZE_PTR)
	if fileSizeInFile != int64(fileSize) {
		return fmt.Errorf("db file size error, expected [%d], real [%d]", fileSizeInFile, fileSize)
	}
	s.firstIndexPtr = bytex.GetIntLong(bytes, FIRST_INDEX_PTR)
	lastIndexPtr := bytex.GetIntLong(bytes, END_INDEX_PTR)
	s.totalIndexBlocks = (lastIndexPtr - s.firstIndexPtr) / int64(blocks.GetIndexBlockLength(int(s.dbType)))

	b := make([]byte, s.totalHeaderBlockSize)
	copy(b, bytes[SUPER_PART_LENGTH:SUPER_PART_LENGTH+s.totalHeaderBlockSize])
	s.initHeaderBlock(b)
	return nil
}

func (s *DBSearcher) initBtreeModeParam() error {
	_, err := s.czFile.Seek(0)
	if err != nil {
		return err
	}

	data := make([]byte, SUPER_PART_LENGTH)
	_, err = s.czFile.Read(data)
	if err != nil {
		return err
	}

	s.totalHeaderBlockSize = bytex.GetIntLong(data, HEADER_BLOCK_PTR)

	data = make([]byte, s.totalHeaderBlockSize)
	_, err = s.czFile.Read(data)
	if err != nil {
		return err
	}

	s.initHeaderBlock(data)
	return nil
}

func (s *DBSearcher) initHeaderBlock(b []byte) {
	indexLength := 20
	idx := 0
	for i := 0; i < len(b); i += indexLength {
		dataPtr := bytex.GetIntLong(b, i+16)
		if dataPtr == 0 {
			break
		}
		s.HeaderSip = append(s.HeaderSip, b[i:i+16])
		s.HeaderPtr = append(s.HeaderPtr, int(dataPtr))
		idx++
	}
	s.headerLength = idx
}
func (s *DBSearcher) Close() {
	if s.czFile != nil {
		_ = s.czFile.Close()
		s.czFile = nil
	}
	s.dbBinStr = nil
	s.HeaderSip = nil
	s.HeaderPtr = nil
	s.geoMapData = nil
}
func (s *DBSearcher) getIpBytes(ip string) ([]byte, error) {
	nip := net.ParseIP(ip)
	if s.dbType == 4 {
		ipv4 := nip.To4()
		if ipv4 == nil {
			return nil, errors.New(fmt.Sprintf("IP [%s] format error for %d", ip, s.dbType))
		}
		return ipv4, nil
	} else {
		ipv6 := nip.To16()
		if ipv6 == nil {
			return nil, errors.New(fmt.Sprintf("IP [%s] format error for %d", ip, s.dbType))
		}
		return ipv6, nil
	}
}
func (s *DBSearcher) Search(ip string) (string, error) {

	ipBytes, err := s.getIpBytes(ip)
	if err != nil {
		return "", err
	}

	var dataBlock *blocks.DataBlock
	if s.queryType == QueryType_Memory {
		dataBlock = s.memorySearch(ipBytes)
	} else {
		dataBlock = s.bTreeSearch(ipBytes)
	}
	if dataBlock == nil {
		return "", nil
	} else {
		return dataBlock.GetRegion(s.geoMapData, s.columnSelection)
	}
}
func (s *DBSearcher) memorySearch(ipBytes []byte) *blocks.DataBlock {

	blockLen := blocks.GetIndexBlockLength(int(s.dbType))

	sptr, eptr := s.searchInHeader(ipBytes)
	if sptr == 0 {
		return nil
	}

	l := 0
	h := (eptr - sptr) / blockLen
	sip := make([]byte, len(ipBytes))
	eip := make([]byte, len(ipBytes))
	var dataPtr, dataLen int

	for l <= h {
		m := (l + h) >> 1
		p := sptr + m*blockLen
		copy(sip, s.dbBinStr[p:p+len(ipBytes)])
		copy(eip, s.dbBinStr[p+len(ipBytes):p+2*len(ipBytes)])

		cmpStart := compareBytes(ipBytes, sip, len(ipBytes))
		cmpEnd := compareBytes(ipBytes, eip, len(ipBytes))

		if cmpStart >= 0 && cmpEnd <= 0 {
			dataPtr = int(bytex.GetIntLong(s.dbBinStr, p+2*len(ipBytes)))
			dataLen = int(bytex.GetInt1(s.dbBinStr, p+2*len(ipBytes)+4))
			break
		} else if cmpStart < 0 {
			h = m - 1
		} else {
			l = m + 1
		}
	}

	if dataPtr == 0 {
		return nil
	}

	region := make([]byte, dataLen)
	copy(region, s.dbBinStr[dataPtr:dataPtr+dataLen])
	return blocks.NewDataBlock(region, int64(dataPtr))
}

func (s *DBSearcher) searchInHeader(ipBytes []byte) (int, int) {
	low := 0
	high := s.headerLength - 1
	sptr := 0
	eptr := 0
	for low <= high {
		m := (low + high) >> 1
		cmp := compareBytes(ipBytes, s.HeaderSip[m], len(ipBytes))
		if cmp < 0 {
			high = m - 1
		} else if cmp > 0 {
			low = m + 1
		} else {
			sptr = s.HeaderPtr[m]
			eptr = s.HeaderPtr[m+1]
			break
		}
	}
	if low == 0 {
		return 0, 0
	}
	if low > high {
		if low < s.headerLength {
			sptr = s.HeaderPtr[low-1]
			eptr = s.HeaderPtr[low]
		} else if high >= 0 && high+1 < s.headerLength {
			sptr = s.HeaderPtr[high]
			eptr = s.HeaderPtr[high+1]
		} else {
			sptr = s.HeaderPtr[s.headerLength-1]
			eptr = sptr + blocks.GetIndexBlockLength(int(s.dbType))
		}
	}
	return sptr, eptr
}

func (s *DBSearcher) bTreeSearch(ipBytes []byte) *blocks.DataBlock {
	sptr, eptr := s.searchInHeader(ipBytes)

	if sptr == 0 {
		return nil
	}

	blockLen := eptr - sptr
	blen := blocks.GetIndexBlockLength(int(s.dbType))
	iBuffer := make([]byte, blockLen+blen)
	_, err := s.czFile.Seek(int64(sptr))
	if err != nil {
		return nil
	}

	_, err = s.czFile.Read(iBuffer)
	if err != nil {
		return nil
	}

	l := 0
	h := blockLen / blen
	sip := make([]byte, len(ipBytes))
	eip := make([]byte, len(ipBytes))
	var dataPtr, dataLen int

	for l <= h {
		m := (l + h) >> 1
		p := m * blen
		copy(sip, iBuffer[p:p+len(ipBytes)])
		copy(eip, iBuffer[p+len(ipBytes):p+2*len(ipBytes)])

		cmpStart := compareBytes(ipBytes, sip, len(ipBytes))
		cmpEnd := compareBytes(ipBytes, eip, len(ipBytes))

		if cmpStart >= 0 && cmpEnd <= 0 {
			dataPtr = int(bytex.GetIntLong(iBuffer, p+2*len(ipBytes)))
			dataLen = int(bytex.GetInt1(iBuffer, p+2*len(ipBytes)+4))
			break
		} else if cmpStart < 0 {
			h = m - 1
		} else {
			l = m + 1
		}
	}

	if dataPtr == 0 {
		return nil
	}

	_, err = s.czFile.Seek(int64(dataPtr))
	if err != nil {
		return nil
	}

	region := make([]byte, dataLen)
	_, err = s.czFile.Read(region)
	if err != nil {
		return nil
	}

	return blocks.NewDataBlock(region, int64(dataPtr))
}
