package bytex

import "encoding/binary"

func GetIntLong(b []byte, offset int) int64 {
	return getIntLong1(b, offset)
}

func GetInt1(b []byte, offset int) int {
	return int(b[offset] & 0x000000FF)
}

func getIntLong1(b []byte, offset int) int64 {
	return int64(b[offset]&0xFF) |
		(int64(b[offset+1]&0xFF) << 8) |
		(int64(b[offset+2]&0xFF) << 16) |
		(int64(b[offset+3]&0xFF) << 24)
}
func getIntLong2(b []byte) int64 {
	return int64(binary.LittleEndian.Uint32(b))
}
