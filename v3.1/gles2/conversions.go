// Glow automatically generated OpenGL binding: http://github.com/go-gl/glow

package gles2

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// #include <stdlib.h>
import "C"

// Ptr takes a slice or pointer (to a singular scalar value or the first
// element of an array or slice) and returns its GL-compatible address.
func Ptr(data interface{}) unsafe.Pointer {
	if data == nil {
		return unsafe.Pointer(nil)
	}
	var addr unsafe.Pointer
	v := reflect.ValueOf(data)
	switch v.Type().Kind() {
	case reflect.Ptr:
		e := v.Elem()
		switch e.Kind() {
		case
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			addr = unsafe.Pointer(e.UnsafeAddr())
		}
	case reflect.Uintptr:
		addr = unsafe.Pointer(v.Pointer())
	case reflect.Slice:
		addr = unsafe.Pointer(v.Index(0).UnsafeAddr())
	default:
		panic(fmt.Sprintf("Unsupported type %s; must be a pointer, slice, or array", v.Type()))
	}
	return addr
}

// PtrOffset takes a pointer offset and returns a GL-compatible pointer.
// Useful for functions such as glVertexAttribPointer that take pointer
// parameters indicating an offset rather than an absolute memory address.
func PtrOffset(offset int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(offset))
}

// Str takes a null-terminated Go string and returns its GL-compatible address.
// This function reaches into Go string storage in an unsafe way so the caller
// must ensure the string is not garbage collected.
func Str(str string) *uint8 {
	if !strings.HasSuffix(str, "\x00") {
		panic("str argument missing null terminator: " + str)
	}
	header := (*reflect.StringHeader)(unsafe.Pointer(&str))
	return (*uint8)(unsafe.Pointer(header.Data))
}

// GoStr takes a null-terminated string returned by OpenGL and constructs a
// corresponding Go string.
func GoStr(cstr *uint8) string {
	return C.GoString((*C.char)(unsafe.Pointer(cstr)))
}

// Strs takes a list of Go strings (with or without null-termination) and
// returns their C counterpart.
//
// The returned free function must be called once you are done using the strings
// in order to free the memory.
//
// If no strings are provided as a parameter this function will panic.
func Strs(strs ...string) (cstrs **uint8, free func()) {
	if len(strs) == 0 {
		panic("Strs: expected at least 1 string")
	}

	// Allocate a contiguous array large enough to hold all the strings' contents.
	n := 0
	for i := range strs {
		n += len(strs[i])
	}
	data := C.malloc(C.size_t(n))

	// Copy all the strings into data.
	dataSlice := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(data),
		Len:  n,
		Cap:  n,
	}))
	css := make([]*uint8, len(strs)) // Populated with pointers to each string.
	offset := 0
	for i := range strs {
		copy(dataSlice[offset:offset+len(strs[i])], strs[i][:]) // Copy strs[i] into proper data location.
		css[i] = (*uint8)(unsafe.Pointer(&dataSlice[offset]))   // Set a pointer to it.
		offset += len(strs[i])
	}

	return (**uint8)(&css[0]), func() { C.free(data) }
}
