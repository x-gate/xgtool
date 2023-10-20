package internal

import (
	"image/color"
	"io"
	"os"
	"xgtool/pkg"
)

type Resources struct {
	GraphicInfoFile *os.File
	GraphicIDIndex  pkg.GraphicInfoIndex
	GraphicMapIndex pkg.GraphicInfoIndex

	GraphicFile *os.File

	PaletteFile *os.File
	Palette     color.Palette

	MapFile *os.File
	Map     pkg.Map

	AnimeInfoFile  *os.File
	AnimeInfoIndex pkg.AnimeInfoIndex

	AnimeFile *os.File
}

func (r *Resources) OpenGraphicInfo(gif string) (err error) {
	if r.GraphicInfoFile, err = os.Open(gif); err != nil {
		return
	}

	if r.GraphicIDIndex, r.GraphicMapIndex, err = pkg.MakeGraphicInfoIndexes(r.GraphicInfoFile); err != nil {
		return
	}
	_, _ = r.GraphicInfoFile.Seek(0, io.SeekStart)

	return
}

func (r *Resources) OpenGraphic(gf string) (err error) {
	r.GraphicFile, err = os.Open(gf)

	return
}

func (r *Resources) OpenPalette(pf string) (err error) {
	if r.PaletteFile, err = os.Open(pf); err != nil {
		return
	}
	r.Palette, err = pkg.NewPaletteFromCGP(r.PaletteFile)
	_, _ = r.PaletteFile.Seek(0, io.SeekStart)

	return
}

func (r *Resources) OpenMap(mf string) (err error) {
	if r.MapFile, err = os.Open(mf); err != nil {
		return
	}
	r.Map, err = pkg.MakeMap(r.MapFile)
	_, _ = r.MapFile.Seek(0, io.SeekStart)

	return
}

func (r *Resources) Close() {
	_ = r.GraphicInfoFile.Close()
	_ = r.GraphicFile.Close()
	_ = r.PaletteFile.Close()
	_ = r.MapFile.Close()
	_ = r.AnimeInfoFile.Close()
	_ = r.AnimeFile.Close()
}
