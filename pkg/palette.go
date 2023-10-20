package pkg

import (
	"bytes"
	"errors"
	"image/color"
	"io"
)

// CGPSize reads bytes from CGP file, (256 colors - 32 default colors) * 3 bytes per color.
const CGPSize = (256 - 32) * 3

// Transparent is a color.RGBA with all fields set to 0.
var Transparent = color.RGBA{}

var prefix = [...]color.Color{
	color.RGBA{B: 0x00, G: 0x00, R: 0x00, A: 0x00}, // RGB(0, 0, 0) is a transparent color for CrossGate, set Alpha to 0 for transparent.
	color.RGBA{B: 0x00, G: 0x00, R: 0x80, A: 0xff},
	color.RGBA{B: 0x00, G: 0x80, R: 0x00, A: 0xff},
	color.RGBA{B: 0x00, G: 0x80, R: 0x80, A: 0xff},
	color.RGBA{B: 0x80, G: 0x00, R: 0x00, A: 0xff},
	color.RGBA{B: 0x80, G: 0x00, R: 0x80, A: 0xff},
	color.RGBA{B: 0x80, G: 0x80, R: 0x00, A: 0xff},
	color.RGBA{B: 0xc0, G: 0xc0, R: 0xc0, A: 0xff},
	color.RGBA{B: 0xc0, G: 0xdc, R: 0xc0, A: 0xff},
	color.RGBA{B: 0xf0, G: 0xca, R: 0xa6, A: 0xff},
	color.RGBA{B: 0x00, G: 0x00, R: 0xde, A: 0xff},
	color.RGBA{B: 0x00, G: 0x5f, R: 0xff, A: 0xff},
	color.RGBA{B: 0xa0, G: 0xff, R: 0xff, A: 0xff},
	color.RGBA{B: 0xd2, G: 0x5f, R: 0x00, A: 0xff},
	color.RGBA{B: 0xff, G: 0xd2, R: 0x50, A: 0xff},
	color.RGBA{B: 0x28, G: 0xe1, R: 0x28, A: 0xff},
}

var suffix = [...]color.Color{
	color.RGBA{B: 0x96, G: 0xc3, R: 0xf5, A: 0xff},
	color.RGBA{B: 0x5f, G: 0xa0, R: 0x1e, A: 0xff},
	color.RGBA{B: 0x46, G: 0x7d, R: 0xc3, A: 0xff},
	color.RGBA{B: 0x1e, G: 0x55, R: 0x9b, A: 0xff},
	color.RGBA{B: 0x37, G: 0x41, R: 0x46, A: 0xff},
	color.RGBA{B: 0x1e, G: 0x23, R: 0x28, A: 0xff},
	color.RGBA{B: 0xf0, G: 0xfb, R: 0xff, A: 0xff},
	color.RGBA{B: 0xa5, G: 0x6e, R: 0x3a, A: 0xff},
	color.RGBA{B: 0x80, G: 0x80, R: 0x80, A: 0xff},
	color.RGBA{B: 0x00, G: 0x00, R: 0xff, A: 0xff},
	color.RGBA{B: 0x00, G: 0xff, R: 0x00, A: 0xff},
	color.RGBA{B: 0x00, G: 0xff, R: 0xff, A: 0xff},
	color.RGBA{B: 0xff, G: 0x00, R: 0x00, A: 0xff},
	color.RGBA{B: 0xff, G: 0x80, R: 0xff, A: 0xff},
	color.RGBA{B: 0xff, G: 0xff, R: 0x00, A: 0xff},
	color.RGBA{B: 0xff, G: 0xff, R: 0xff, A: 0xff},
}

// NewPaletteFromCGP make palette from CGP file.
func NewPaletteFromCGP(r io.Reader) (p color.Palette, err error) {
	buf := bytes.NewBuffer(make([]byte, CGPSize))
	if _, err = io.ReadFull(r, buf.Bytes()); err != nil {
		return
	}

	if p, err = NewPaletteFromBytes(buf.Bytes()); err != nil {
		return
	}

	p = append(prefix[:], p...)
	p = append(p, suffix[:]...)

	return
}

// NewPaletteFromBytes make palette from bytes.
func NewPaletteFromBytes(b []byte) (p color.Palette, err error) {
	buf := bytes.NewBuffer(b)
	tmp := make([]byte, 3)

	for _, err = io.ReadFull(buf, tmp); err == nil; _, err = io.ReadFull(buf, tmp) {
		if tmp[0] == 0 && tmp[1] == 0 && tmp[2] == 0 {
			p = append(p, Transparent)
		} else {
			p = append(p, color.RGBA{B: tmp[0], G: tmp[1], R: tmp[2], A: 0xff})
		}
	}

	if errors.Is(err, io.EOF) {
		err = nil
	}

	return
}
