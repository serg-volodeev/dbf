package dbf

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

const (
	dbfId     byte = 0x03
	headerEnd byte = 0x0D

	headerSize = 32
	yearOffset = 1900
)

type header struct {
	Id         byte
	ModYear    byte
	ModMonth   byte
	ModDay     byte
	RecCount   uint32
	DataOffset uint16
	RecSize    uint16
	Filler1    [17]byte
	CP         byte
	Filler2    [2]byte
}

func newHeader() *header {
	h := &header{}
	h.Id = dbfId
	h.setModDate(time.Now())
	return h
}

// Modified date

func (h *header) modDate() time.Time {
	year := int(h.ModYear) + yearOffset
	month := time.Month(h.ModMonth)
	day := int(h.ModDay)
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func (h *header) setModDate(d time.Time) {
	h.ModYear = byte(d.Year() - yearOffset)
	h.ModMonth = byte(d.Month())
	h.ModDay = byte(d.Day())
}

// Field count

func (h *header) fieldCount() int {
	return (int(h.DataOffset) - headerSize - 1) / fieldSize
}

func (h *header) setFieldCount(count int) {
	h.DataOffset = uint16(count*fieldSize + headerSize + 1)
}

// Read/write

func (h *header) read(reader io.Reader) error {
	if err := binary.Read(reader, binary.LittleEndian, h); err != nil {
		return err
	}
	if h.Id != dbfId {
		return fmt.Errorf("not DBF file")
	}
	return nil
}

func (h *header) write(writer io.Writer) error {
	return binary.Write(writer, binary.LittleEndian, h)
}

// Code page

func (h *header) codePage() int {
	return pageByCode(h.CP)
}

func (h *header) setCodePage(cp int) {
	h.CP = codeByPage(cp)
}
