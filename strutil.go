package dbf

import (
	"strings"
	"unicode"
)

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

func trimLeft(buf []byte) string {
	i := 0
	for ; i < len(buf); i++ {
		if buf[i] != ' ' {
			break
		}
	}
	if i >= len(buf) {
		return ""
	}
	return string(buf[i:])
}

func trimRight(buf []byte) string {
	i := len(buf) - 1
	for ; i >= 0; i-- {
		if buf[i] != ' ' {
			i++
			break
		}
	}
	if i < 0 {
		return ""
	}
	return string(buf[0:i])
}

func isEmpty(buf []byte) bool {
	for i := range buf {
		if buf[i] != ' ' {
			return false
		}
	}
	return true
}
