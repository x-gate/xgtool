package pkg

import (
	"bytes"
	"errors"
	"image/color"
	"io"
)

// Transparent is a color.RGBA with all fields set to 0.
var Transparent = color.RGBA{}

// Palette is a collection of colors, it usually 768 bytes (256 * 3), but not always.
type Palette []color.Color

// NewPaletteFromBytes make palette from bytes.
func NewPaletteFromBytes(b []byte) (p Palette, err error) {
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
