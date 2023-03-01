package BloomFilter

import (
	"encoding/binary"
	"reflect"
	"unsafe"
)

func StrToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Uint16ToBytes(num uint16) []byte {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, num)
	return data
}

func Uint32ToBytes(num uint32) []byte {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, num)
	return data
}

func Uint64ToBytes(num uint64) []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, num)
	return data
}
