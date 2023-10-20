package convert_map

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"xgtool/cmd/internal"
	"xgtool/internal/tmx"
)

type flags struct {
	gif    string
	gf     string
	pf     string
	mf     string
	outdir string
	outmap string
	dr     bool // dry-run
}

func (f *flags) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet("dump-graphic", flag.ExitOnError)
	fs.StringVar(&f.gif, "gif", "", "graphic info file path")
	fs.StringVar(&f.gf, "gf", "", "graphic file path")
	fs.StringVar(&f.pf, "pf", "", "palette file path")
	fs.StringVar(&f.mf, "mf", "", "map file path")
	fs.StringVar(&f.outdir, "o", "output", "output directory")
	fs.StringVar(&f.outmap, "n", "map.json", "output file name")
	fs.BoolVar(&f.dr, "dry-run", false, "dump without output files (for testing)")

	return fs
}

var (
	f flags
)

func ConvertMap(ctx context.Context, args []string) (err error) {
	if err = f.Flags().Parse(args); err != nil {
		return
	}
	if f.dr {
		f.outdir = os.TempDir()
	}

	res := internal.Resources{}
	defer res.Close()

	if err = res.OpenGraphicInfo(f.gif); err != nil {
		return
	}
	if err = res.OpenGraphic(f.gf); err != nil {
		return
	}
	if err = res.OpenPalette(f.pf); err != nil {
		return
	}
	if err = res.OpenMap(f.mf); err != nil {
		return
	}

	var tm tmx.Map
	if tm, err = res.Map.TiledMap(
		res.GraphicMapIndex,
		res.GraphicFile,
		res.Palette,
		f.outdir,
	); err != nil {
		return
	}

	var out []byte
	if out, err = json.Marshal(tm); err != nil {
		log.Err(err).Send()
		return
	}

	if err = os.WriteFile(
		fmt.Sprintf("%s/%s", filepath.Clean(f.outdir), f.outmap),
		out,
		0644,
	); err != nil {
		log.Err(err).Send()
		return
	}

	return
}
