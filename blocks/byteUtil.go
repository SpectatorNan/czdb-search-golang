package blocks

func writeIntLong(b []byte, offset int, value int64) {
	b[offset] = byte((value >> 0) & 0xFF)
	b[offset+1] = byte((value >> 8) & 0xFF)
	b[offset+2] = byte((value >> 16) & 0xFF)
	b[offset+3] = byte((value >> 24) & 0xFF)
}
