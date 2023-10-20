package internal

import (
	"image/color"
	"os"
	"xgtool/pkg"
)

type Resources struct {
	InfoFile    *os.File
	IDIndex     pkg.GraphicInfoIndex
	MapIndex    pkg.GraphicInfoIndex
	GraphicFile *os.File
	PaletteFile *os.File
	Palette     color.Palette
	MapFile     *os.File
	Map         pkg.Map
}

func OpenGraphicRes(gif, gf, pf, mf string) (res Resources, err error) {
	if res.InfoFile, err = os.Open(gif); err != nil {
		return res, err
	}
	res.IDIndex, res.MapIndex, err = pkg.MakeGraphicInfoIndexes(res.InfoFile)

	if res.GraphicFile, err = os.Open(gf); err != nil {
		return res, err
	}

	// palette file is optional
	if pf != "" {
		res.PaletteFile, err = os.Open(pf)
		res.Palette, err = pkg.NewPaletteFromCGP(res.PaletteFile)
	}

	// map file is optional
	if mf != "" {
		res.MapFile, err = os.Open(mf)
		res.Map, err = pkg.MakeMap(res.MapFile)
	}

	return
}

func (f Resources) Close() {
	_ = f.InfoFile.Close()
	_ = f.GraphicFile.Close()
	_ = f.PaletteFile.Close()
	_ = f.MapFile.Close()
}
