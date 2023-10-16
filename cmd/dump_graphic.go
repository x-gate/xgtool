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
	"xgtool/pkg"
)

var bar *progressbar.ProgressBar
var wg sync.WaitGroup
var flags dumpGraphicFlags

func main() {
	if err := flags.Flags().Parse(os.Args[1:]); err != nil {
		log.Err(err).Send()
		return
	}

	log.Debug().Msgf("dumpGraphicFlags: %+v", flags)

	files, err := openGraphicFiles(flags)
	if err != nil {
		log.Err(err).Send()
		return
	}
	defer files.Close()

	var palette pkg.Palette
	if files.Palette != nil {
		palette, err = pkg.NewPaletteFromCGP(files.Palette)
		if err != nil {
			log.Err(err).Send()
			return
		}
	}

	graphicIndex, _, err := pkg.MakeGraphicInfoIndexes(files.Info)
	if err != nil {
		log.Err(err).Send()
		return
	}

	log.Debug().Msgf("len(graphicIndex): %d", len(graphicIndex))

	bar = progressbar.Default(int64(len(graphicIndex)))

	for _, gif := range graphicIndex {
		if err = dumpGraphic(gif, files.Graphic, palette); err != nil {
			log.Err(err).Send()
			return
		}
		_ = bar.Add(1)
	}

	wg.Wait()

	return
}

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

type graphicFiles struct {
	Info    *os.File
	Graphic *os.File
	Palette *os.File
}

func openGraphicFiles(flags dumpGraphicFlags) (files graphicFiles, err error) {
	if files.Info, err = os.Open(flags.GraphicInfoFile); err != nil {
		return files, err
	}
	if files.Graphic, err = os.Open(flags.GraphicFile); err != nil {
		return files, err
	}

	// palette file is optional
	if flags.PaletteFile == "" {
		return
	}

	files.Palette, err = os.Open(flags.PaletteFile)

	return
}

func (f graphicFiles) Close() {
	_ = f.Info.Close()
	_ = f.Graphic.Close()
	_ = f.Palette.Close()
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

	go func() {
		wg.Add(1)
		defer wg.Done()

		var img image.Image
		if img, err = g.ToImage(); err != nil {
			log.Err(err).Send()
			return
		}

		var out *os.File
		if flags.DryRun {
			out, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
		} else {
			out, err = os.OpenFile(fmt.Sprintf("output/%d.jpg", g.Info.ID), os.O_WRONLY|os.O_CREATE, 0644)
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
