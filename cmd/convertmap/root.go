package convertmap

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"xgtool/internal/tmx"
	"xgtool/pkg"

	"github.com/rs/zerolog/log"
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

// ConvertMap the entrypoint of "convert-map" command
func ConvertMap(ctx context.Context, args []string) (err error) {
	if err = f.Flags().Parse(args); err != nil {
		return
	}
	if f.dr {
		f.outdir = os.TempDir()
	}

	res := pkg.Resources{}
	defer res.Close()

	if err = res.OpenGraphicResource(f.gif); err != nil {
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
		res.GraphicResource.MDx,
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
