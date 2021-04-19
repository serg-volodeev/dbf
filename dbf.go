// Package dbf reads and writes DBF files.
// The api of the dbf package is similar to the csv package from the standard library.
package dbf

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/encoding/charmap"
)

const (
	dbfId     byte = 0x03
	headerEnd byte = 0x0D
	fileEnd   byte = 0x1A
)

const (
	fieldSize  = 32
	headerSize = 32
)

type cPage struct {
	code byte
	page int
	cm   *charmap.Charmap
}

var cPages = []cPage{
	{code: 0x01, page: 437, cm: charmap.CodePage437},  // US MS-DOS
	{code: 0x02, page: 850, cm: charmap.CodePage850},  // International MS-DOS
	{code: 0x03, page: 1252, cm: charmap.Windows1252}, // Windows ANSI
	{code: 0x04, page: 10000, cm: charmap.Macintosh},  // Standard Macintosh
	{code: 0x64, page: 852, cm: charmap.CodePage852},  // Easern European MS-DOS
	{code: 0x65, page: 866, cm: charmap.CodePage866},  // Russian MS-DOS
	{code: 0x66, page: 865, cm: charmap.CodePage865},  // Nordic MS-DOS

	// Not found in package charmap
	// 0x67	Codepage 861 Icelandic MS-DOS
	// 0x68	Codepage 895 Kamenicky (Czech) MS-DOS
	// 0x69	Codepage 620 Mazovia (Polish) MS-DOS
	// 0x6A	Codepage 737 Greek MS-DOS (437G)
	// 0x6B	Codepage 857 Turkish MS-DOS
	// 0x78	Codepage 950 Chinese (Hong Kong SAR, Taiwan) Windows
	// 0x79	Codepage 949 Korean Windows
	// 0x7A	Codepage 936 Chinese (PRC, Singapore) Windows
	// 0x7B	Codepage 932 Japanese Windows
	// 0x7C	Codepage 874 Thai Windows

	{code: 0x7D, page: 1255, cm: charmap.Windows1255},        // Hebrew Windows
	{code: 0x7E, page: 1256, cm: charmap.Windows1256},        // Arabic Windows
	{code: 0x96, page: 10007, cm: charmap.MacintoshCyrillic}, // Russian MacIntosh

	// Not found in package charmap
	// 0x97	Codepage 10029 MacIntosh EE
	// 0x98	Codepage 10006 Greek MacIntosh

	{code: 0xC8, page: 1250, cm: charmap.Windows1250}, // Eastern European Windows
	{code: 0xC9, page: 1251, cm: charmap.Windows1251}, // Russian Windows
	{code: 0xCA, page: 1254, cm: charmap.Windows1254}, // Turkish Windows
	{code: 0xCB, page: 1253, cm: charmap.Windows1253}, // Greek Windows
}

func charmapByPage(page int) *charmap.Charmap {
	for i := range cPages {
		if cPages[i].page == page {
			return cPages[i].cm
		}
	}
	return nil
}

func codeByPage(page int) byte {
	for i := range cPages {
		if cPages[i].page == page {
			return cPages[i].code
		}
	}
	return 0
}

func pageByCode(code byte) int {
	for i := range cPages {
		if cPages[i].code == code {
			return cPages[i].page
		}
	}
	return 0
}

// String utils

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func padLeft(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-len(s)) + s
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// Convert interface{}

func interfaceToString(v interface{}) (string, error) {
	result, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("error convert %v to string", v)
	}
	return result, nil
}

func interfaceToBool(v interface{}) (bool, error) {
	result, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("error convert %v to bool", v)
	}
	return result, nil
}

func interfaceToInt(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	}
	return 0, fmt.Errorf("error convert %v to int", value)
}

func interfaceToFloat(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return float64(v), nil
	}
	return 0, fmt.Errorf("error convert %v to float", value)
}

func interfaceToDate(v interface{}) (time.Time, error) {
	result, ok := v.(time.Time)
	if !ok {
		return result, fmt.Errorf("error convert %v to date", v)
	}
	return result, nil
}
