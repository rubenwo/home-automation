package pnglib

// #cgo CXXFLAGS: -g -Wall -O4 -march=native
// #cgo CFLAGS: -g -Wall -O4 -march=native
// #include <stdlib.h>
// #include "png.h"
import "C"

import (
	"unsafe"
)

func init() {
	C.fpng_init_wrapper()
}

func EncodeFromMat2d(data []uint8, w, h uint32, size int) []byte {
	var pngSize uintptr
	pngPtr := C.fpng_encode_with_data_wrapper8((*C.uint8_t)(&data[0]), C.uint32_t(w), C.uint32_t(h), C.int(size), (*C.size_t)(unsafe.Pointer(&pngSize)))
	b := C.GoBytes(unsafe.Pointer(pngPtr), C.int(pngSize))
	C.free(unsafe.Pointer(pngPtr))
	return b
}
