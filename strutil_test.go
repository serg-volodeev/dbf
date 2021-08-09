package dbf

import (
	"testing"
)

func Test_padRight(t *testing.T) {
	tests := []struct {
		s     string
		width int
		want  string
	}{
		{s: "Abc", width: 6, want: "Abc   "},
		{s: "Abc", width: 2, want: "Abc"},
		{s: "", width: 3, want: "   "},
		{s: "  Abc", width: 7, want: "  Abc  "},
	}
	for _, tc := range tests {
		got := padRight(tc.s, tc.width)
		if got != tc.want {
			t.Errorf("padRight(%#v, %#v), want: %#v, got: %#v", tc.s, tc.width, tc.want, got)
		}
	}
}

func Test_padLeft(t *testing.T) {
	tests := []struct {
		s     string
		width int
		want  string
	}{
		{s: "Abc", width: 6, want: "   Abc"},
		{s: "Abc", width: 2, want: "Abc"},
		{s: "", width: 3, want: "   "},
		{s: "  Abc", width: 7, want: "    Abc"},
	}
	for _, tc := range tests {
		got := padLeft(tc.s, tc.width)
		if got != tc.want {
			t.Errorf("padLeft(%#v, %#v), want: %#v, got: %#v", tc.s, tc.width, tc.want, got)
		}
	}
}

func Test_isASCII(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{s: "Abc", want: true},
		{s: "AbЖc", want: false},
		{s: "", want: true},
	}
	for _, tc := range tests {
		got := isASCII(tc.s)
		if got != tc.want {
			t.Errorf("isASCII(%#v), want: %#v, got: %#v", tc.s, tc.want, got)
		}
	}
}

func Test_trimRight(t *testing.T) {
	tests := []struct {
		b    []byte
		want string
	}{
		{b: []byte("Abc"), want: "Abc"},
		{b: []byte("Abc   "), want: "Abc"},
		{b: []byte(""), want: ""},
		{b: []byte("   "), want: ""},
	}
	for _, tc := range tests {
		got := trimRight(tc.b)
		if got != tc.want {
			t.Errorf("trimRight(%#v), want: %#v, got: %#v", string(tc.b), tc.want, got)
		}
	}
}

func Test_trimLeft(t *testing.T) {
	tests := []struct {
		b    []byte
		want string
	}{
		{b: []byte("Abc"), want: "Abc"},
		{b: []byte("   Abc"), want: "Abc"},
		{b: []byte(""), want: ""},
		{b: []byte("   "), want: ""},
	}
	for _, tc := range tests {
		got := trimLeft(tc.b)
		if got != tc.want {
			t.Errorf("trimLeft(%#v), want: %#v, got: %#v", string(tc.b), tc.want, got)
		}
	}
}
