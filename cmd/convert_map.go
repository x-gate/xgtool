package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"xgtool/cmd/internal"
)

type convertMapFlags struct {
	GraphicInfoFile string
	GraphicFile     string
	PaletteFile     string
	MapFile         string
	OutDir          string
	OutName         string
}

func (f *convertMapFlags) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet("dump-graphic", flag.ExitOnError)
	fs.StringVar(&f.GraphicInfoFile, "gif", "", "graphic info file path")
	fs.StringVar(&f.GraphicFile, "gf", "", "graphic file path")
	fs.StringVar(&f.PaletteFile, "pf", "", "palette file path")
	fs.StringVar(&f.MapFile, "mf", "", "map file path")
	fs.StringVar(&f.OutDir, "o", "output", "output directory")
	fs.StringVar(&f.OutName, "n", "map", "output file name")

	return fs
}

var cmf convertMapFlags

func main() {
	if err := cmf.Flags().Parse(os.Args[1:]); err != nil {
		log.Err(err).Send()
		return
	}

	log.Debug().Msgf("dumpGraphicFlags: %+v", cmf)

	res, err := internal.OpenGraphicRes(cmf.GraphicInfoFile, cmf.GraphicFile, cmf.PaletteFile, cmf.MapFile)
	if err != nil {
		log.Err(err).Send()
		return
	}
	defer res.Close()

	tm, err := res.Map.TiledMap(res.MapIndex, res.GraphicFile, res.Palette, cmf.OutDir)
	if err != nil {
		log.Err(err).Send()
	}

	out, err := json.Marshal(tm)
	if err != nil {
		log.Err(err).Send()
	}

	if err = os.WriteFile(
		fmt.Sprintf("%s/%s.json", filepath.Clean(cmf.OutDir), filepath.Base(cmf.OutName)), out, 0644,
	); err != nil {
		log.Err(err).Send()
	}
}
