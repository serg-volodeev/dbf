package xbase

import "time"

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
	year := int(h.ModYear) + 1900
	month := time.Month(h.ModMonth)
	day := int(h.ModDay)
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func (h *header) setModDate(d time.Time) {
	h.ModYear = byte(d.Year() - 1900)
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
