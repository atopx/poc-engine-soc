package utils

import "unsafe"

// StrToBytes string to bytes
func StrToBytes(s string) []byte {
	t := (*[2]uintptr)(unsafe.Pointer(&s))
	b := [3]uintptr{t[0], t[1], t[1]}
	return *(*[]byte)(unsafe.Pointer(&b))
}

// BytesToStr bytes to string
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
