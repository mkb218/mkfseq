// Created by cgo - DO NOT EDIT

package mkfseq

import "unsafe"

import "os"

import _ "runtime/cgo"

type _ unsafe.Pointer

func _Cerrno(dst *os.Error, x int) { *dst = os.Errno(x) }
type _Ctypedef_real8 _Ctype_double
type _Ctypedef_complex8 _Ctype_struct___0
type _Ctype_double float64
type _Ctype_struct___0 struct {
	re	_Ctypedef_real8
	im	_Ctypedef_real8
}
type _Ctype_void [0]byte

func _Cfunc_fftc8_1024(*_Ctypedef_complex8)
