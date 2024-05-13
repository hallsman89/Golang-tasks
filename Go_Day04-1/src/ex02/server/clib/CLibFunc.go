package clib

/*
   #include <stdio.h>
   #include <stdlib.h>
   #include <string.h>
   #include "cow_says.h"
   #cgo CFLAGS: -Wall
*/
import "C"
import "unsafe"

func CowSayMoo() string {
	cowPtr := C.ask_cow(C.CString("Thank you!"))
	defer C.free(unsafe.Pointer(cowPtr))
	cow := C.GoString(cowPtr)
	return cow
}
