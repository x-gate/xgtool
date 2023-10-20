package dumpgraphic

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"sync"
	"xgtool/pkg"
)

type flags struct {
	gif    string
	gf     string
	pf     string
	outdir string
	dr     bool // dry-run
}

func (f *flags) Flags() (fs *flag.FlagSet) {
	fs = flag.NewFlagSet("dump-graphic", flag.ExitOnError)
	fs.StringVar(&f.gif, "gif", "", "graphic info file path")
	fs.StringVar(&f.gf, "gf", "", "graphic file path")
	fs.StringVar(&f.pf, "pf", "", "palette file path")
	fs.StringVar(&f.outdir, "o", "output", "output directory")
	fs.BoolVar(&f.dr, "dry-run", false, "dump without output files (for testing)")

	return
}

var (
	bar *progressbar.ProgressBar
	wg  sync.WaitGroup
	f   flags
)

// DumpGraphic the entrypoint of "dump-graphic" command
func DumpGraphic(ctx context.Context, args []string) (err error) {
	if err = f.Flags().Parse(args); err != nil {
		return
	}

	res := pkg.Resources{}
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

	bar = progressbar.Default(int64(len(res.GraphicIDIndex)))

	for _, gif := range res.GraphicIDIndex {
		if err = dumpGraphic(gif, res.GraphicFile, res.Palette); err != nil {
			log.Err(err).Send()
			return
		}
		_ = bar.Add(1)
	}

	wg.Wait()

	return nil
}

func dumpGraphic(info pkg.GraphicInfo, gf *os.File, palette color.Palette) error {
	g, err := info.LoadGraphic(gf)
	if err != nil && (errors.Is(err, pkg.ErrInvalidMagic) || errors.Is(err, pkg.ErrDecodeFailed)) {
		log.Warn().Msgf("Invalid Graphic: %+v", err)
		return nil
	} else if err != nil {
		return err
	}

	if len(g.PaletteData) == 0 {
		if len(palette) == 0 {
			return pkg.ErrEmptyPalette
		}
		g.SetPalette(palette)
	}

	go func() {
		wg.Add(1)
		defer wg.Done()

		var img image.Image
		if img, err = g.ImgRGBA(); err != nil {
			log.Err(err).Send()
			return
		}

		var out *os.File
		if f.dr {
			out, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
		} else {
			out, err = os.OpenFile(fmt.Sprintf("%s/%d.jpg", filepath.Clean(f.outdir), g.Info.ID), os.O_WRONLY|os.O_CREATE, 0644)
		}
		if err != nil {
			log.Err(err).Send()
		}
		defer out.Close()

		if err = jpeg.Encode(out, img, &jpeg.Options{Quality: 75}); err != nil {
			log.Err(err).Send()
		}
	}()

	return err
}
