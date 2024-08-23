package czdb

type Decryptor struct {
	key string
}

func NewDecryptor(key string) *Decryptor {
	return &Decryptor{key: key}
}
func (x *Decryptor) Decrypt(data []byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[i] = b ^ x.key[i%len(x.key)]
	}
	return result
}
