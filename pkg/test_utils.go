package pkg

import (
	"image/color"
	"io"
	"os"
	"testing"
)

type TestRes struct {
	GraphicInfoFileName string
	GraphicInfoFile     *os.File
	GraphicInfoIDIndex  GraphicInfoIndex
	GraphicInfoMapIndex GraphicInfoIndex

	GraphicFileName string
	GraphicFile     *os.File

	PaletteFileName string
	PaletteFile     *os.File
	Palette         color.Palette

	AnimeInfoFileName string
	AnimeInfoFile     *os.File
	AnimeInfoIndex    AnimeInfoIndex

	AnimeFileName string
	AnimeFile     *os.File

	MapFileName string
	MapFile     *os.File
	Map         Map
}

func (r *TestRes) OpenGraphicInfo(gif string) (err error) {
	r.GraphicInfoFileName = gif
	if r.GraphicInfoFile, err = os.Open(gif); err != nil {
		return
	}
	r.GraphicInfoIDIndex, r.GraphicInfoMapIndex, err = MakeGraphicInfoIndexes(r.GraphicInfoFile)
	_, _ = r.GraphicInfoFile.Seek(0, io.SeekStart)

	return
}

func (r *TestRes) OpenGraphic(gf string) (err error) {
	r.GraphicFileName = gf
	r.GraphicFile, err = os.Open(gf)
	_, _ = r.GraphicFile.Seek(0, io.SeekStart)

	return
}

func (r *TestRes) OpenPalette(pf string) (err error) {
	r.PaletteFileName = pf
	if r.PaletteFile, err = os.Open(pf); err != nil {
		return
	}
	r.Palette, err = NewPaletteFromCGP(r.PaletteFile)
	_, _ = r.PaletteFile.Seek(0, io.SeekStart)

	return
}

func (r *TestRes) OpenAnimeInfo(aif string) (err error) {
	r.AnimeInfoFileName = aif
	if r.AnimeInfoFile, err = os.Open(aif); err != nil {
		return
	}
	r.AnimeInfoIndex, err = MakeAnimeInfoIndex(r.AnimeInfoFile)
	_, _ = r.AnimeInfoFile.Seek(0, io.SeekStart)

	return
}

func (r *TestRes) OpenAnime(af string) (err error) {
	r.AnimeFileName = af
	r.AnimeFile, err = os.Open(af)
	_, _ = r.AnimeFile.Seek(0, io.SeekStart)

	return
}

func (r *TestRes) OpenMap(mf string) (err error) {
	r.MapFileName = mf
	if r.MapFile, err = os.Open(mf); err != nil {
		return
	}
	r.Map, err = MakeMap(r.MapFile)
	_, _ = r.MapFile.Seek(0, io.SeekStart)

	return
}

func (r *TestRes) Close() {
	_ = r.GraphicInfoFile.Close()
	_ = r.GraphicFile.Close()
	_ = r.PaletteFile.Close()
	_ = r.AnimeInfoFile.Close()
	_ = r.AnimeFile.Close()
	_ = r.MapFile.Close()
}

func SkipIfNotExists(file string, err error, t *testing.T) {
	if err != nil && os.IsNotExist(err) {
		t.Skipf("skipping test; file %s does not exist", file)
	} else if err != nil {
		t.Fatal(err)
	}
}
