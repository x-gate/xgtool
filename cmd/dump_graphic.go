package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"image"
	"image/jpeg"
	"os"
	"sync"
	"xgtool/cmd/internal"
	"xgtool/pkg"
)

var (
	bar *progressbar.ProgressBar
	wg  sync.WaitGroup
)

type dumpGraphicFlags struct {
	GraphicInfoFile string
	GraphicFile     string
	PaletteFile     string
	DryRun          bool
}

func (f *dumpGraphicFlags) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet("dump-graphic", flag.ExitOnError)
	fs.StringVar(&f.GraphicInfoFile, "gif", "", "graphic info file path")
	fs.StringVar(&f.GraphicFile, "gf", "", "graphic file path")
	fs.StringVar(&f.PaletteFile, "pf", "", "palette file path")
	fs.BoolVar(&f.DryRun, "dry-run", false, "dump without output files (for testing)")

	return fs
}

var dgf dumpGraphicFlags

func main() {
	if err := dgf.Flags().Parse(os.Args[1:]); err != nil {
		log.Err(err).Send()
		return
	}

	log.Debug().Msgf("dumpGraphicFlags: %+v", dgf)

	res, err := internal.OpenGraphicRes(dgf.GraphicInfoFile, dgf.GraphicFile, dgf.PaletteFile, "")
	if err != nil {
		log.Err(err).Send()
		return
	}
	defer res.Close()

	log.Debug().Msgf("len(graphicIndex): %d", len(res.IDIndex))

	bar = progressbar.Default(int64(len(res.IDIndex)))

	for _, gif := range res.IDIndex {
		if err = dumpGraphic(gif, res.GraphicFile, res.Palette); err != nil {
			log.Err(err).Send()
			return
		}
		_ = bar.Add(1)
	}

	wg.Wait()

	return
}

func dumpGraphic(info pkg.GraphicInfo, gf *os.File, palette pkg.Palette) error {
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

	go render(g)

	return err
}

func render(g *pkg.Graphic) {
	var err error

	wg.Add(1)
	defer wg.Done()

	var img image.Image
	if img, err = g.ToImage(); err != nil {
		log.Err(err).Send()
		return
	}

	var out *os.File
	if out, err = output(*g.Info); err != nil {
		log.Err(err).Send()
		return
	}
	defer out.Close()

	if err = jpeg.Encode(out, img, &jpeg.Options{Quality: 75}); err != nil {
		log.Err(err).Send()
	}
}

func output(gi pkg.GraphicInfo) (f *os.File, err error) {
	if dgf.DryRun {
		return os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
	}

	return os.OpenFile(fmt.Sprintf("output/%d.jpg", gi.ID), os.O_WRONLY|os.O_CREATE, 0644)
}
