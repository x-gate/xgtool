package pkg

import (
	"errors"
	"io"
	"os"
	"testing"
)

func TestMakePaletteFromCGP(t *testing.T) {
	palettes, _ := os.ReadDir("../testdata/palette")

	for _, f := range palettes {
		t.Run(f.Name(), func(t *testing.T) {
			res := TestRes{}
			defer res.Close()

			if err := res.OpenPalette("../testdata/palette/" + f.Name()); err != nil {
				t.Fatal(err)
			}

			p, err := NewPaletteFromCGP(res.PaletteFile)
			if err != nil {
				t.Fatal(err)
			}

			if len(p) != 256 {
				t.Errorf("len(p) = %d, want 256", len(p))
			}
		})
	}
}

func TestMakePaletteFromBytes(t *testing.T) {
	testcases := []struct {
		name string
		data []byte
		err  error
	}{
		{
			name: "single color palette",
			data: []byte{0xff, 0xff, 0xff}, // white
			err:  nil,
		},
		{
			name: "invalid length",
			data: []byte{0xff, 0xff, 0xff, 0xff},
			err:  io.ErrUnexpectedEOF,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := NewPaletteFromBytes(tc.data)
			if !errors.Is(err, tc.err) {
				t.Fatal(err)
			}

			if len(p) != len(tc.data)/3 {
				t.Errorf("len(p) = %d, want %d", len(p), len(tc.data)/3)
			}
		})
	}
}
