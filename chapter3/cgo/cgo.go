package main

/*
#include <stdio.h>

void print_hello(const char *name) {
	printf("Hello %s!\n", name);
}
*/
import "C"

func main() {
	cName := C.CString("World")
	C.print_hello(cName)
}
