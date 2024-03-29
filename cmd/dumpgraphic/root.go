package dumpgraphic

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"xgtool/pkg"

	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
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
	f   flags
)

// DumpGraphic the entrypoint of "dump-graphic" command
func DumpGraphic(ctx context.Context, args []string) (err error) {
	if err = f.Flags().Parse(args); err != nil {
		return
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

	bar = progressbar.Default(int64(len(res.GraphicResource.IDx)))
	done := make(chan struct{})

	go func() {
		defer close(done)
		for _, gifs := range res.GraphicResource.IDx {
			select {
			case <-ctx.Done():
				return
			default:
				for i, gif := range gifs {
					if err = dumpGraphic(gif.Info, res.GraphicFile, res.Palette, i); err != nil {
						log.Err(err).Send()
						return
					}
				}

				_ = bar.Add(1)
			}

		}
	}()
	<-done

	return nil
}

func dumpGraphic(info pkg.GraphicInfo, gf *os.File, palette color.Palette, serial int) (err error) {
	var g *pkg.Graphic
	g, err = info.LoadGraphic(gf)
	if err != nil && (errors.Is(err, pkg.ErrInvalidMagic) || errors.Is(err, pkg.ErrDecodeFailed)) {
		log.Warn().Msgf("Invalid Graphic: %+v", err)
		return nil
	} else if err != nil {
		return err
	}

	var img image.Image
	if img, err = g.ImgRGBA(palette); err != nil && (errors.Is(err, pkg.ErrRenderFailed) || errors.Is(err, pkg.ErrEmptyPalette)) {
		log.Warn().Msgf("Failed to render: %+v", err)
		return nil
	} else if err != nil {
		return err
	}

	var out *os.File
	if f.dr {
		out, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
	} else {
		out, err = os.OpenFile(fmt.Sprintf("%s/%d-%d.jpg", filepath.Clean(f.outdir), g.Info.ID, serial), os.O_WRONLY|os.O_CREATE, 0644)
	}
	if err != nil {
		log.Err(err).Send()
		return
	}
	defer out.Close()

	if err = jpeg.Encode(out, img, &jpeg.Options{Quality: 75}); err != nil {
		log.Err(err).Send()
		return
	}

	return err
}
