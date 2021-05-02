package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var f = 1.0
	fmt.Printf("%#016x\n", *(*uint64)(unsafe.Pointer(&f)))

	var x struct {
		a bool
		b int16
		c []int
	}
	// equivalent to pb := &x.b
	pb := (*int16)(unsafe.Pointer(
		uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
	*pb = 42
	fmt.Println(x.b)
}
