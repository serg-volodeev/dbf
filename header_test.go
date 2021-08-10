package dbf

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

var headerBytes = []byte{0x3, 0x1e, 0x2, 0x14, 0x3, 0, 0, 0, 0xc1, 0, 0x27,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x65, 0, 0}

func currentDate() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func Test_newHeader(t *testing.T) {
	h := newHeader()

	tpl := "newHeader(): %s: want: %v, got: %v"

	if h.Id != 0x03 {
		t.Errorf(tpl, "h.Id", 0x03, h.Id)
	}
	if h.RecCount != 0 {
		t.Errorf(tpl, "h.RecCount", 0, h.RecCount)
	}
	if h.fieldCount() != 0 {
		t.Errorf(tpl, "h.fieldCount()", 0, h.fieldCount())
	}
	if h.RecSize != 0 {
		t.Errorf(tpl, "h.RecSize", 0, h.RecSize)
	}
	if h.codePage() != 0 {
		t.Errorf(tpl, "h.codePage()", 0, h.codePage())
	}
	d := currentDate()
	if h.modDate() != d {
		t.Errorf(tpl, "h.modDate()", d, h.modDate())
	}
}

func Test_header_write(t *testing.T) {
	h := newHeader()
	h.RecCount = uint32(3)
	h.RecSize = uint16(39)
	h.setFieldCount(5)
	h.setModDate(time.Date(1930, 2, 20, 0, 0, 0, 0, time.UTC))
	h.setCodePage(866)

	buf := bytes.NewBuffer(nil)
	err := h.write(buf)

	if err != nil {
		t.Errorf("header.write(): %v", err)
	}
	if !reflect.DeepEqual(buf.Bytes(), headerBytes) {
		t.Errorf("header.write():\nwant: %#v\ngot : %#v", headerBytes, buf.Bytes())
	}
}

func Test_header_read(t *testing.T) {
	h := &header{}
	h.read(bytes.NewReader(headerBytes))

	tpl := "header.read(): %s: want: %v, got: %v"

	if h.Id != 0x03 {
		t.Errorf(tpl, "h.Id", 0x03, h.Id)
	}
	if h.RecCount != 3 {
		t.Errorf(tpl, "h.RecCount", 3, h.RecCount)
	}
	if h.fieldCount() != 5 {
		t.Errorf(tpl, "h.fieldCount()", 5, h.fieldCount())
	}
	if h.RecSize != 39 {
		t.Errorf(tpl, "h.RecSize", 39, h.RecSize)
	}
	if h.codePage() != 866 {
		t.Errorf(tpl, "h.codePage()", 866, h.codePage())
	}
	d := time.Date(1930, 2, 20, 0, 0, 0, 0, time.UTC)
	if h.modDate() != d {
		t.Errorf(tpl, "h.modDate()", d, h.modDate())
	}
}

func Test_header_read_notDBF(t *testing.T) {
	b := make([]byte, headerSize)
	b[0] = 0x05 // valid 0x03

	h := &header{}
	err := h.read(bytes.NewReader(b))

	if err == nil {
		t.Errorf("header.read(): error not DBF expected")
	}
}

func Test_header_setCodePage(t *testing.T) {
	h := &header{}
	h.setCodePage(866)

	if h.CP != 0x65 {
		t.Errorf("header.setCodePage(866): h.CP: want: %v, got: %v", 0x65, h.CP)
	}
	if h.codePage() != 866 {
		t.Errorf("header.setCodePage(866): h.codePage(): want: %v, got: %v", 866, h.codePage())
	}
}
