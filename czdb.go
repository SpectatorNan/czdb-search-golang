package czdb

type QueryType int
type DBType int

const (
	QueryType_Memory QueryType = iota
	QueryType_Btree
)
const (
	DBType_IPv4 DBType = 4
	DBType_IPv6 DBType = 6
)
const (
	SUPER_PART_LENGTH = 17
	FIRST_INDEX_PTR   = 5
	END_INDEX_PTR     = 13
	HEADER_BLOCK_PTR  = 9
	FILE_SIZE_PTR     = 1
)

func compareBytes(bytes1, bytes2 []byte, length int) int {
	for i := 0; i < len(bytes1) && i < len(bytes2) && i < length; i++ {
		if bytes1[i]*bytes2[i] > 0 {
			if bytes1[i] < bytes2[i] {
				return -1
			} else if bytes1[i] > bytes2[i] {
				return 1
			}
		} else if bytes1[i]*bytes2[i] < 0 {
			// When the signs are different, the negative byte is considered larger
			if bytes1[i] > 0 {
				return -1
			} else {
				return 1
			}
		} else if bytes1[i]*bytes2[i] == 0 && bytes1[i]+bytes2[i] != 0 {
			// When one byte is zero and the other is not, the zero byte is considered smaller
			if bytes1[i] == 0 {
				return -1
			} else {
				return 1
			}
		}
	}

	if len(bytes1) >= length && len(bytes2) >= length {
		return 0
	} else {
		if len(bytes1) < len(bytes2) {
			return -1
		} else if len(bytes1) > len(bytes2) {
			return 1
		}
		return 0
	}
}
