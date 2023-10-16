package pkg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var (
	// ErrInvalidMagic if graphic header magic is not "RD".
	ErrInvalidMagic = errors.New("invalid magic")
	// ErrDecodeFailed if graphic data decode failed.
	ErrDecodeFailed = errors.New("decode failed")
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

// graphicHeader structure for each graphic header, 16 bytes.
type graphicHeader struct {
	Magic   [2]byte // "RD" for valid graphic
	Version byte    // 0 for raw data, 1 for encoded data, 2 for raw data with palette, 3 for encoded data with palette
	_       byte    //
	Width   int32   // Width of graphic, it shouldn't be trusted, use GraphicInfo.Width instead.
	Height  int32   // Height of graphic, it shouldn't be trusted, use GraphicInfo.Height instead.
	Len     int32   // Length of graphic data, it shouldn't be trusted, use GraphicInfo.Len instead.
}

// Graphic stores data for each graphic, not a strict mapping to the file.
type Graphic struct {
	Info        *GraphicInfo // Pointer of GraphicInfo, for reverse searching.
	Header      graphicHeader
	RawData     []byte  // The raw data which read from graphic file.
	GraphicData []byte  // The decoded (if needed) data from RawData
	PaletteLen  int32   // When Version >= 2, read this field from graphic file, it couldn't be set by direct set palette data.
	PaletteData Palette // When Version < 2, set palette data from palette file; otherwise, set palette data from graphic file.
}

// GraphicInfoIndex is a map of GraphicInfo, key is GraphicInfo.ID or GraphicInfo.MapID.
type GraphicInfoIndex map[int32]GraphicInfo

// MakeGraphicInfoIndexes reads graphic info from src, and returns two GraphicInfoIndex,
// first is indexed by GraphicInfo.ID, second is indexed by GraphicInfo.MapID.
func MakeGraphicInfoIndexes(gif io.Reader) (idx, mapIdx GraphicInfoIndex, err error) {
	idx = make(GraphicInfoIndex)
	mapIdx = make(GraphicInfoIndex)

	r := bufio.NewReader(gif)
	for {
		buf := bytes.NewBuffer(make([]byte, 40))
		if _, err = io.ReadFull(r, buf.Bytes()); err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, nil, err
		}

		var info GraphicInfo
		if err = binary.Read(buf, binary.LittleEndian, &info); err != nil {
			return nil, nil, err
		}

		idx[info.ID] = info
		if info.MapID != 0 { // there are a lot of graphic info with MapID=0, but they are not used in map files
			mapIdx[info.MapID] = info
		}
	}

	return
}

// LoadGraphic loads graphic data from graphic file.
func (gi GraphicInfo) LoadGraphic(gf io.ReadSeeker) (g *Graphic, err error) {
	g = new(Graphic)
	g.Info = &gi

	if err = g.readGraphic(gf, int64(gi.Addr), int64(gi.Len)); err != nil {
		return
	}

	if err = g.decode(); err != nil {
		return
	}

	return
}

func (g *Graphic) readGraphic(f io.ReadSeeker, offset, len int64) (err error) {
	if _, err = f.Seek(offset, io.SeekStart); err != nil {
		return
	}

	buf := bytes.NewBuffer(make([]byte, len))
	if _, err = io.ReadFull(f, buf.Bytes()); err != nil {
		return
	}

	if err = binary.Read(buf, binary.LittleEndian, &g.Header); err != nil {
		return
	}

	if g.Header.Magic[0] != 'R' || g.Header.Magic[1] != 'D' {
		return fmt.Errorf("%w: info=%+v, header=%+v", ErrInvalidMagic, g.Info, g.Header)
	}

	if g.Header.Version >= 2 {
		if err = binary.Read(buf, binary.LittleEndian, &g.PaletteLen); err != nil {
			return
		}
	}

	g.RawData = buf.Bytes()

	return
}

func (g *Graphic) decode() (err error) {
	var decoded []byte

	if g.Header.Version&1 == 0 {
		decoded = g.RawData
	} else if decoded, err = Decode(g.RawData); err != nil {
		return fmt.Errorf("%w: info=%+v, header=%+v", ErrDecodeFailed, g.Info, g.Header)
	}

	g.GraphicData = decoded[:len(decoded)-int(g.PaletteLen)]
	g.PaletteData, err = NewPaletteFromBytes(decoded[len(decoded)-int(g.PaletteLen):])

	return
}
