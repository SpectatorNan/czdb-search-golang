package czdb

import (
	"io"
	"os"
)

type Cz88File struct {
	file   *os.File
	offset int64
}

func NewCz88File(file *os.File, offset int64) *Cz88File {
	return &Cz88File{file: file, offset: offset}
}

func (f *Cz88File) Seek(pos int64) (int64, error) {
	return f.file.Seek(pos+f.offset, io.SeekStart)
}

func (f *Cz88File) Read(data []byte) (int, error) {
	return f.file.Read(data)
}

func (f *Cz88File) Close() error {
	return f.file.Close()
}

func (f *Cz88File) ReadSignedBytes(data []int) (int, error) {
	bs := make([]byte, len(data))
	n, err := f.file.Read(bs)
	if err != nil {
		return n, err
	}
	for i, b := range bs {
		data[i] = f.signedByteToInt(b)
	}
	return n, nil
}

func (f *Cz88File) signedByteToInt(b byte) int {
	if b > 127 {
		return int(b) - 256 // Correct the interpretation for negative values
	}
	return int(b)
}

func (f *Cz88File) Length() int64 {
	info, _ := f.file.Stat()
	return info.Size() - f.offset
}
