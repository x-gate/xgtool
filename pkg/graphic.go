package pkg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"xgtool/internal"
)

const (
	// GraphicInfoSize is the size of GraphicInfo structure in bytes.
	GraphicInfoSize = 40
)

var (
	// ErrInvalidMagic if graphic header magic is not "RD".
	ErrInvalidMagic = errors.New("invalid magic")
	// ErrDecodeFailed if graphic data decode failed.
	ErrDecodeFailed = errors.New("decode failed")
	// ErrEmptyPalette is returned when Graphic.Image is called but the palette is empty.
	ErrEmptyPalette = errors.New("empty palette")
	// ErrRenderFailed is returned when the Graphic.PaletteData[Graphic.GraphicData[i]] is out of range.
	ErrRenderFailed = errors.New("render failed")
)

// GraphicInfo structure for each graphic info, 40 bytes.
type GraphicInfo struct {
	ID     int32
	Addr   int32
	Len    int32
	OffX   int32
	OffY   int32
	Width  int32
	Height int32
	GridW  byte
	GridH  byte
	Access byte
	_      [5]byte
	MapID  int32
}

// GraphicHeader structure for each graphic header, 16 bytes.
type GraphicHeader struct {
	Magic   [2]byte // "RD" for valid graphic
	Version byte    // 0 for raw data, 1 for encoded data, 2 for raw data with palette, 3 for encoded data with palette
	_       byte    //
	Width   int32   // Width of graphic, it shouldn't be trusted, use GraphicInfo.Width instead.
	Height  int32   // Height of graphic, it shouldn't be trusted, use GraphicInfo.Height instead.
	Len     int32   // Length of graphic data, it shouldn't be trusted, use GraphicInfo.Len instead.
}

// Valid checks if the graphic is valid.
func (gh GraphicHeader) Valid() bool {
	return gh.Magic[0] == 'R' && gh.Magic[1] == 'D'
}

// Graphic stores data for each graphic, not a strict mapping to the file.
type Graphic struct {
	Info        GraphicInfo // Pointer of GraphicInfo, for reverse searching.
	Header      GraphicHeader
	GraphicData []byte        // The decoded (if needed) data from RawData
	PaletteData color.Palette // When Version < 2, set palette data from palette file; otherwise, set palette data from graphic file.
}

// GraphicIndex is a map of Graphic, key is the ID of the graphic.
type GraphicIndex map[int32][]*Graphic

// Find finds graphic data for specific id.
func (idx GraphicIndex) Find(id int32) []*Graphic {
	return idx[id]
}

// First finds the first graphic data for specific id.
func (idx GraphicIndex) First(id int32) *Graphic {
	if g := idx.Find(id); len(g) > 0 {
		return g[0]
	}
	return nil
}

// Load loads graphic data for specific id.
func (idx GraphicIndex) Load(id int32, gf io.ReadSeeker) (err error) {
	for _, g := range idx.Find(id) {
		if err = g.Load(gf); err != nil {
			return
		}
	}

	return
}

// GraphicResource is a map of []*Graphic, key is the ID or MapID of the graphic.
type GraphicResource struct {
	IDx GraphicIndex // Index by GraphicInfo.ID
	MDx GraphicIndex // Index by GraphicInfo.MapID
}

// NewGraphicResource reads graphic info from gif, and returns GraphicResource.
//
// The graphic data is not loaded yet, use GraphicIndex.Load to load graphic data.
func NewGraphicResource(gif io.Reader) (gr GraphicResource, err error) {
	gr.IDx = make(map[int32][]*Graphic)
	gr.MDx = make(map[int32][]*Graphic)

	r := bufio.NewReaderSize(gif, GraphicInfoSize*100)
	for {
		buf := bytes.NewBuffer(make([]byte, GraphicInfoSize))
		if _, err = io.ReadFull(r, buf.Bytes()); err != nil && errors.Is(err, io.EOF) {
			err = nil
			break
		} else if err != nil {
			return
		}

		var gi GraphicInfo
		if err = binary.Read(buf, binary.LittleEndian, &gi); err != nil {
			return
		}

		g := Graphic{Info: gi}
		gr.IDx[gi.ID] = append(gr.IDx[gi.ID], &g)
		if gi.MapID != 0 { // if MapID=0, it's not used in map files
			gr.MDx[gi.MapID] = append(gr.MDx[gi.MapID], &g)
		}
	}

	return
}

// LoadGraphic loads graphic data from graphic file.
func (gi GraphicInfo) LoadGraphic(gf io.ReadSeeker) (g *Graphic, err error) {
	g = new(Graphic)
	g.Info = gi

	if err = g.Load(gf); err != nil {
		return
	}

	return
}

// Load reads from graphic file, and decode if needed
func (g *Graphic) Load(f io.ReadSeeker) (err error) {
	// If GraphicData is not empty, it's already loaded.
	if len(g.GraphicData) != 0 {
		return nil
	}

	if _, err = f.Seek(int64(g.Info.Addr), io.SeekStart); err != nil {
		return
	}

	buf := bytes.NewBuffer(make([]byte, g.Info.Len))
	if _, err = io.ReadFull(f, buf.Bytes()); err != nil {
		return
	}

	if err = binary.Read(buf, binary.LittleEndian, &g.Header); err != nil {
		return
	}

	if !g.Header.Valid() {
		return fmt.Errorf("%w: info=%+v, header=%+v", ErrInvalidMagic, g.Info, g.Header)
	}

	var psz int32
	if g.Header.Version >= 2 {
		if err = binary.Read(buf, binary.LittleEndian, &psz); err != nil {
			return
		}
	}

	var decoded []byte
	if decoded, err = g.decode(buf.Bytes()); err != nil {
		return
	}

	g.GraphicData = decoded[:len(decoded)-int(psz)]
	g.PaletteData, err = NewPaletteFromBytes(decoded[len(decoded)-int(psz):])

	return
}

func (g *Graphic) decode(raw []byte) (decoded []byte, err error) {
	if g.Header.Version&1 == 0 {
		decoded = raw
	} else if decoded, _ = internal.Decode(raw); decoded == nil {
		return
	}

	return
}

// ImgRGBA convert graphic data to image.RGBA
func (g *Graphic) ImgRGBA(p color.Palette) (img *image.RGBA, err error) {
	if len(g.PaletteData) == 0 && len(p) == 0 {
		return nil, ErrEmptyPalette
	} else if len(g.PaletteData) == 0 {
		g.PaletteData = p
	}

	w := int(g.Info.Width)
	h := int(g.Info.Height)
	img = image.NewRGBA(image.Rect(0, 0, w, h))

	for i, pix := range g.GraphicData {
		if int(pix) >= len(g.PaletteData) {
			return nil, fmt.Errorf("%w: info=%+v, header=%+v, g.GraphicData[i]=%d, len(g.PaletteData)=%d", ErrRenderFailed, g.Info, g.Header, pix, len(g.PaletteData))
		}
		img.Set(i%w, h-i/w, g.PaletteData[pix])
	}

	return
}

// ImgPaletted convert graphic data to image.Paletted
func (g *Graphic) ImgPaletted(p color.Palette) (img *image.Paletted, err error) {
	if len(g.PaletteData) == 0 && len(p) == 0 {
		return nil, ErrEmptyPalette
	} else if len(g.PaletteData) == 0 {
		g.PaletteData = p
	}

	w := int(g.Info.Width)
	h := int(g.Info.Height)
	r := image.Rect(0, 0, w, h)
	img = image.NewPaletted(r, g.PaletteData)

	for i, pix := range g.GraphicData {
		// The code is based on image.Paletted.Set() from go standard library.
		// The implementation is very slow because it calls p.Palette.Index(c) for each pixel, but it's not necessary.
		//
		// Ref: https://cs.opensource.google/go/go/+/refs/tags/go1.21.4:src/image/image.go;l=1188
		if !(image.Point{X: i % w, Y: h - i/w}.In(r)) {
			continue
		}
		j := img.PixOffset(i%w, h-i/w)

		img.Pix[j] = pix
	}

	return
}
