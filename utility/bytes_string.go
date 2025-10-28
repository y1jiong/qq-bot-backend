package utility

import "unsafe"

func BytesToString(b []byte) string {
	//return *(*string)(unsafe.Pointer(&b))
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func StringToBytes(s string) []byte {
	//sh := (*[2]uintptr)(unsafe.Pointer(&s))
	//bh := [3]uintptr{sh[0], sh[1], sh[1]}
	//return *(*[]byte)(unsafe.Pointer(&bh))
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
