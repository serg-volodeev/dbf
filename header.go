package xbase

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
	// h.setModDate(time.Now())
	return h
}
