package dbf

import (
	"testing"

	"golang.org/x/text/encoding/charmap"
)

func Test_сodeByPage(t *testing.T) {
	tests := []struct {
		page int
		want byte
	}{
		{page: 866, want: 0x65},
		{page: -1, want: 0x00},
		{page: 1251, want: 0xC9},
	}
	for _, tc := range tests {
		got := codeByPage(tc.page)
		if got != tc.want {
			t.Errorf("codeByPage(%#v), want: %#v, got: %#v", tc.page, tc.want, got)
		}
	}
}

func Test_pageByCode(t *testing.T) {
	tests := []struct {
		code byte
		want int
	}{
		{code: 0x65, want: 866},
		{code: 0x00, want: 0},
		{code: 0xC9, want: 1251},
	}
	for _, tc := range tests {
		got := pageByCode(tc.code)
		if got != tc.want {
			t.Errorf("pageByCode(%#v), want: %#v, got: %#v", tc.code, tc.want, got)
		}
	}
}

func Test_charmapByPage(t *testing.T) {
	tests := []struct {
		page int
		want *charmap.Charmap
	}{
		{page: 866, want: charmap.CodePage866},
		{page: -1, want: nil},
		{page: 1251, want: charmap.Windows1251},
	}
	for _, tc := range tests {
		got := charmapByPage(tc.page)
		if got != tc.want {
			t.Errorf("charmapByPage(%v), want: %v, got: %v", tc.page, tc.want, got)
		}
	}
}
