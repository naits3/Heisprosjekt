// Created by cgo - DO NOT EDIT

package main

import "unsafe"

import _ "runtime/cgo"

import "syscall"

var _ syscall.Errno
func _Cgo_ptr(ptr unsafe.Pointer) unsafe.Pointer { return ptr }

type _Ctype_int int32

type _Ctype_void [0]byte

var _cgo_runtime_cgocall_errno func(unsafe.Pointer, uintptr) int32
var _cgo_runtime_cmalloc func(uintptr) unsafe.Pointer


var _cgo_56271688278c_Cfunc_test_func unsafe.Pointer
func _Cfunc_test_func() (r1 _Ctype_int) {
	_cgo_runtime_cgocall_errno(_cgo_56271688278c_Cfunc_test_func, uintptr(unsafe.Pointer(&r1)))
	return
}
