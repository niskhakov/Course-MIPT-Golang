package cgo

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -L/usr/lib -lbz2
#include <bzlib.h>
#include <stdlib.h>
int bz2compress(bz_stream *s, int action,
		char *in, unsigned *inlen, char *out, unsigned *outlen);
*/
import "C"

import (
	"io"
	"unsafe"
)

type writer struct {
	w      io.Writer // underlying output stream
	stream *C.bz_stream
	outBuf [64 * 1024]byte
}

func NewWriter(out io.Writer) io.WriteCloser {
	const (
		blockSize  = 9
		verbosity  = 0
		workFactor = 30
	)
	w := &writer{
		w:      out,
		stream: (*C.bz_stream)(C.malloc(C.sizeof_bz_stream)),
	}
	C.BZ2_bzCompressInit(w.stream, blockSize, verbosity, workFactor)
	return w
}

func (w *writer) Write(data []byte) (int, error) {
	if w.stream == nil {
		panic("closed")
	}
	var total int // uncompressed bytes written

	for len(data) > 0 {
		inlen, outlen := C.uint(len(data)), C.uint(cap(w.outBuf))
		C.bz2compress(w.stream, C.BZ_RUN,
			(*C.char)(unsafe.Pointer(&data[0])), &inlen,
			(*C.char)(unsafe.Pointer(&w.outBuf)), &outlen)
		total += int(inlen)
		data = data[inlen:]
		if _, err := w.w.Write(w.outBuf[:outlen]); err != nil {
			return total, err
		}
	}

	return total, nil
}

func (w *writer) Close() error {
	if w.stream == nil {
		panic("closed")
	}
	defer func() {
		C.BZ2_bzCompressEnd(w.stream)
		C.free(unsafe.Pointer(w.stream))
		w.stream = nil
	}()
	for {
		inlen, outlen := C.uint(0), C.uint(cap(w.outBuf))
		r := C.bz2compress(w.stream, C.BZ_FINISH, nil, &inlen,
			(*C.char)(unsafe.Pointer(&w.outBuf)), &outlen)
		if _, err := w.w.Write(w.outBuf[:outlen]); err != nil {
			return err
		}
		if r == C.BZ_STREAM_END {
			return nil
		}
	}
}
