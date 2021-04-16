package xbase

import (
	"bufio"
	"io"
)

type Reader struct {
	// header  *header
	// fields  []*field
	reader *bufio.Reader
	// buf     []byte
	// recNo   int64
	// decoder *encoding.Decoder
}

func NewReader(r io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(r)}
}
