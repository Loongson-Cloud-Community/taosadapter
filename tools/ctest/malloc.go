package ctest

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"
import "unsafe"

func Malloc(size int) unsafe.Pointer {
	p := C.malloc(C.ulong(size))
	//for i := 0; i < size; i++ {
	//	*(*C.char)(unsafe.Pointer(uintptr(p)+uintptr(i))) = C.char(1)
	//}
	C.memset(p, 1, C.ulong(size))
	return p
}

func Free(p unsafe.Pointer) {
	C.free(p)
}
