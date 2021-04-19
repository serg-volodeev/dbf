// Package dbf reads and writes DBF files.
// The api of the dbf package is similar to the csv package from the standard library.
package dbf

const (
	dbfId     byte = 0x03
	headerEnd byte = 0x0D
	fileEnd   byte = 0x1A
)

const (
	fieldSize  = 32
	headerSize = 32
)
