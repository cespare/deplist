package c

/*
#include <stdio.h>
#include <stdlib.h>

void print(char *s) {
	printf("%s", s);
}
*/
import "C"

import "unsafe"

func C() {
	s := C.CString("Hello from C!\n")
	defer C.free(unsafe.Pointer(s))
	C.print(s)
}
