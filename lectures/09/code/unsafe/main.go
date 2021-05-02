package main

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println(unsafe.Sizeof(int(0)))
	fmt.Println(unsafe.Sizeof(int8(0)))
	fmt.Println(unsafe.Sizeof(Foo{}))
	fmt.Println(unsafe.Sizeof([4]int{}))
	fmt.Println(unsafe.Sizeof("hello, world! it's me!"))
	fmt.Println(unsafe.Sizeof(make([]int, 100, 200)))

	fmt.Println("===")

	fmt.Println(unsafe.Sizeof(struct {
		bool
		float64
		int16
	}{}))
	fmt.Println(unsafe.Sizeof(struct {
		float64
		int16
		bool
	}{}))
	fmt.Println(unsafe.Sizeof(struct {
		bool
		int16
		float64
	}{}))

	fmt.Println("===")
	fmt.Println(unsafe.Alignof(int8(0)))
	fmt.Println(unsafe.Alignof(int16(0)))
	fmt.Println(unsafe.Alignof(int32(0)))

	fmt.Println("===")
	fmt.Println(unsafe.Offsetof(Foo{}.Int))
	fmt.Println(unsafe.Offsetof(Foo{}.Float))
}

type Foo struct {
	Int   int
	Float float64
}
