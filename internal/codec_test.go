package internal

import (
	"bytes"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecode(t *testing.T) {
	testcases := []struct {
		name    string
		encoded []byte
		decoded []byte
	}{
		{
			name:    "read 1 bytes",
			encoded: []byte{0x01, 0xaa},
			decoded: []byte{0xaa},
		},
		{
			name:    "read 4078 bytes",
			encoded: append([]byte{0x1f, 0xee}, bytes.Repeat([]byte{0xaa}, 0x0f*0x100+0xee)...),
			decoded: bytes.Repeat([]byte{0xaa}, 0x0f*0x100+0xee),
		},
		{
			name:    "read 74291 bytes",
			encoded: append([]byte{0x21, 0x22, 0x33}, bytes.Repeat([]byte{0xaa}, 0x01*0x10000+0x22*0x100+0x33)...),
			decoded: bytes.Repeat([]byte{0xaa}, 0x01*0x10000+0x22*0x100+0x33),
		},
		{
			name:    "repeat 1 byte 2 times",
			encoded: []byte{0x82, 0xaa},
			decoded: []byte{0xaa, 0xaa},
		},
		{
			name:    "repeat 1 byte 4078 times",
			encoded: []byte{0x9f, 0xaa, 0xee},
			decoded: bytes.Repeat([]byte{0xaa}, 0x0f*0x100+0xee),
		},
		{
			name:    "repeat 1 byte 74291 times",
			encoded: []byte{0xa1, 0xaa, 0x22, 0x33},
			decoded: bytes.Repeat([]byte{0xaa}, 0x01*0x10000+0x22*0x100+0x33),
		},
		{
			name:    "repeat 1 alpha byte",
			encoded: []byte{0xc1},
			decoded: []byte{0x00},
		},
		{
			name:    "repeat 4078 alpha bytes",
			encoded: []byte{0xdf, 0xee},
			decoded: bytes.Repeat([]byte{0x00}, 0x0f*0x100+0xee),
		},
		{
			name:    "repeat 74291 alpha bytes",
			encoded: []byte{0xe1, 0x22, 0x33},
			decoded: bytes.Repeat([]byte{0x00}, 0x01*0x10000+0x22*0x100+0x33),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			decoded, err := Decode(tc.encoded)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.decoded, decoded); diff != "" {
				t.Errorf("decoded mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDecode_InvalidFlag(t *testing.T) {
	testcases := []struct {
		name string
		data []byte
	}{
		{
			name: "invalid flag 0x3?",
			data: []byte{0x31},
		},
		{
			name: "invalid flag 0x4?",
			data: []byte{0x42},
		},
		{
			name: "invalid flag 0x5?",
			data: []byte{0x53},
		},
		{
			name: "invalid flag 0x6?",
			data: []byte{0x64},
		},
		{
			name: "invalid flag 0x7?",
			data: []byte{0x75},
		},
		{
			name: "invalid flag 0xb?",
			data: []byte{0xb6},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Decode(tc.data)

			if !errors.Is(err, ErrInvalidFlag) {
				t.Errorf("expected error: %v, got: %v", ErrInvalidFlag, err)
			}
		})
	}
}

func BenchmarkDecode_RepeatSingleByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Decode([]byte{0xaf, 0xaa, 0xff, 0xff})
	}
}

func BenchmarkDecode_RepeatAlphaByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Decode([]byte{0xef, 0xff, 0xff})
	}
}
