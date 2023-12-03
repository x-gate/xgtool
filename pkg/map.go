package pkg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"xgtool/internal/mat"
	"xgtool/internal/tmx"
)

type mapHeader struct {
	Magic  [12]byte
	Width  int32
	Height int32
}

// Map is a struct that represents the CrossGate map data.
type Map struct {
	Header mapHeader
	Ground []uint16
	Object []uint16
	Meta   []uint16
}

// The base unit of CrossGate data.
const (
	TileWidth  = 64
	TileHeight = 47
)

// MakeMap make a Map from Crossgate map file.
func MakeMap(f io.Reader) (m Map, err error) {
	if m.Header, err = readHeader(f); err != nil {
		return
	}

	m.Ground, err = readBlock(f, int(m.Header.Width*m.Header.Height*2))
	if err != nil {
		return
	}

	m.Object, err = readBlock(f, int(m.Header.Width*m.Header.Height*2))
	if err != nil {
		return
	}

	m.Meta, err = readBlock(f, int(m.Header.Width*m.Header.Height*2))

	return
}

func readHeader(f io.Reader) (h mapHeader, err error) {
	buf := bytes.NewBuffer(make([]byte, 20))
	if _, err = io.ReadFull(f, buf.Bytes()); err != nil {
		return
	}

	if err = binary.Read(buf, binary.LittleEndian, &h); err != nil {
		return
	}

	if !(h.Magic[0] == 'M' || h.Magic[1] == 'A' || h.Magic[2] == 'P') {
		err = fmt.Errorf("%w: header=%+v", ErrInvalidMagic, h)
	}

	return
}

func readBlock(f io.Reader, len int) (b []uint16, err error) {
	b = make([]uint16, 0, len/2)

	buf := bytes.NewBuffer(make([]byte, len))
	if _, err = io.ReadFull(f, buf.Bytes()); err != nil && !errors.Is(err, io.EOF) {
		return
	}

	for i := 0; i < len/2; i++ {
		var data uint16
		if err = binary.Read(buf, binary.LittleEndian, &data); err != nil {
			return
		}

		b = append(b, data)
	}

	return
}

// TiledMap convert the Map to a tmx.Map
func (m Map) TiledMap(index GraphicInfoIndex, gf io.ReadSeeker, p color.Palette, outdir string) (tiled tmx.Map, err error) {
	tiled = tmx.NewMap(
		// reverse width and height, because the map will be rotated -90 degrees
		int(m.Header.Height),
		int(m.Header.Width),
		tmx.Isometric,
		tmx.LeftUp,
	)
	tiled.Layers = make([]tmx.Layer, 0, 2)
	tiled.TileSets = make([]tmx.TileSet, 0, 2)

	var gid int
	if gid, err = m.setGround(&tiled, index, gf, p, outdir); err != nil {
		return
	}
	if err = m.setObject(&tiled, gid, index, gf, p, outdir); err != nil {
		return
	}

	tiled.TileWidth = TileWidth
	tiled.TileHeight = TileHeight

	return
}

func (m Map) setGround(tiled *tmx.Map, index GraphicInfoIndex, gf io.ReadSeeker, p color.Palette, outdir string) (gid int, err error) {
	var layer tmx.Layer
	if layer, err = m.buildGroundLayer(); err != nil {
		return
	}
	tiled.Layers = append(tiled.Layers, layer)

	var tileset tmx.TileSet
	if tileset, err = m.buildTileSet("ground", &gid, m.Ground, index, gf, p, outdir); err != nil {
		return
	}
	tiled.TileSets = append(tiled.TileSets, tileset)

	return
}

func (m Map) setObject(tiled *tmx.Map, gid int, index GraphicInfoIndex, gf io.ReadSeeker, p color.Palette, outdir string) (err error) {
	var layer tmx.Layer
	if layer, err = m.buildObjectLayer(index, gid); err != nil {
		return
	}
	tiled.Layers = append(tiled.Layers, layer)

	var tileset tmx.TileSet
	if tileset, err = m.buildTileSet("object", &gid, m.Object, index, gf, p, outdir); err != nil {
		return
	}
	tiled.TileSets = append(tiled.TileSets, tileset)

	return
}

func (m Map) buildGroundLayer() (layer tmx.Layer, err error) {
	layer = tmx.NewTileLayer(
		"ground",
		1,
		// reverse width and height, because the map will be rotated -90 degrees
		int(m.Header.Height),
		int(m.Header.Width),
	)
	tile := m.Ground

	// The generated map needs to be rotated 90 degrees clockwise.
	matrix, err := mat.NewMatrix(tile, int(m.Header.Width), int(m.Header.Height))
	if err != nil {
		return
	}
	tile = matrix.Rotate().ToUint16Array()

	for _, t := range tile {
		layer.Data = append(layer.Data, uint(t))
	}

	return
}

func (m Map) buildObjectLayer(index GraphicInfoIndex, fgid int) (layer tmx.Layer, err error) {
	layer = tmx.NewObjectLayer("object", 2, tmx.TopDown)

	for i, t := range m.Object {
		if t == 0 {
			continue
		}
		if _, ok := index[int32(t)]; !ok {
			continue
		}

		gi := index[int32(t)]
		obj := tmx.NewObject(int(t)+fgid, int(t), float64(gi.Width), float64(gi.Height))
		obj.X, obj.Y = objectCoordinate(int32(i), m.Header.Width, gi.Width, gi.Height, gi.OffX, gi.OffY)

		layer.Objects = append(layer.Objects, obj)
	}

	return
}

// objectCoordinate get the (X, Y) from the graphic offset and the map width.
//
// Ref: https://github.com/x-gate/CrossGateRemastered/blob/master/toolchain/hackMap/getCGMap.cpp#L304-L374
func objectCoordinate(i, mapWidth, w, h, offX, offY int32) (x, y float64) {
	row, col := i/mapWidth, i%mapWidth

	offsetX := (float64(offX) + (float64(w) / 2.0)) / TileWidth * TileHeight
	offsetY := float64(offY) + float64(h) - TileHeight/2.0

	x = float64(row+1)*TileHeight + offsetY + offsetX
	y = float64(mapWidth-col)*TileHeight + offsetY - offsetX

	return
}

func (m Map) buildTileSet(name string, fgid *int, tiles []uint16, index GraphicInfoIndex, gf io.ReadSeeker, p color.Palette, outdir string) (ts tmx.TileSet, err error) {
	*fgid++
	ts = tmx.NewTileSet(name, *fgid, tmx.Grid{Orientation: tmx.Orthogonal, Width: 1, Height: 1})

	mapping := make(map[uint16]GraphicInfo)
	for i, t := range tiles {
		// When the tile is 0, means empty.
		if t == 0 {
			continue
		}
		// When the tile is not found in the index, mark it as empty.
		if _, ok := index[int32(t)]; !ok {
			tiles[i] = 0
			continue
		}
		// When the tile is found in the map, ignore it because it has been processed.
		if _, ok := mapping[t]; ok {
			continue
		}

		mapping[t] = index[int32(t)]
		if t > uint16(*fgid) {
			*fgid = int(t)
		}

		if err = render(mapping[t], gf, p, outdir); err != nil {
			return
		}
	}

	for _, v := range mapping {
		ts.TileCount++
		ts.Tiles = append(ts.Tiles, tmx.Tile{
			ID:          int(v.MapID) - 1,
			Image:       fmt.Sprintf("%d.png", v.MapID),
			ImageWidth:  int(v.Width),
			ImageHeight: int(v.Height),
		})
	}
	ts.TileWidth, ts.TileHeight = TileWidth, TileHeight

	return
}

func render(info GraphicInfo, gf io.ReadSeeker, p color.Palette, outdir string) (err error) {
	var g *Graphic
	if g, err = info.LoadGraphic(gf); err != nil {
		return
	}

	if len(g.PaletteData) == 0 {
		g.PaletteData = p
	}

	var img image.Image
	if img, err = g.ImgRGBA(p); err != nil {
		return
	}

	var out *os.File
	if out, err = os.OpenFile(
		fmt.Sprintf("%s/%d.png", filepath.Clean(outdir), info.MapID),
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644,
	); err != nil {
		return
	}

	return png.Encode(out, img)
}
