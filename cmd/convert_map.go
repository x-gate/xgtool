package main

import (
	"flag"
	"github.com/rs/zerolog/log"
	"os"
	"xgtool/cmd/internal"
)

type convertMapFlags struct {
	GraphicInfoFile string
	GraphicFile     string
	PaletteFile     string
	MapFile         string
	DryRun          bool
}

func (f *convertMapFlags) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet("dump-graphic", flag.ExitOnError)
	fs.StringVar(&f.GraphicInfoFile, "gif", "", "graphic info file path")
	fs.StringVar(&f.GraphicFile, "gf", "", "graphic file path")
	fs.StringVar(&f.PaletteFile, "pf", "", "palette file path")
	fs.StringVar(&f.MapFile, "mf", "", "map file path")
	fs.BoolVar(&f.DryRun, "dry-run", false, "dump without output files (for testing)")

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
}
