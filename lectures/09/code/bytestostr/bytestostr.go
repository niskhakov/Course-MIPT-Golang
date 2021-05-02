package bytestostr

import (
	"reflect"
	"unsafe"
)

func BytesToStr(b []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	strHeader := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}
	return *(*string)(unsafe.Pointer(&strHeader))
}
