package dumpanime

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"image/color"
	"image/gif"
	"os"
	"path/filepath"
	"xgtool/pkg"

	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
)

var errPaletteNotFound = errors.New("palette not found")

type flags struct {
	aif    string
	af     string
	gif    string
	gf     string
	pgif   string
	pgf    string
	pf     string
	outdir string
	dr     bool // dry-run
}

func (f *flags) Flags() (fs *flag.FlagSet) {
	fs = flag.NewFlagSet("dump-anime", flag.ExitOnError)
	fs.StringVar(&f.aif, "aif", "", "anime info file path")
	fs.StringVar(&f.af, "af", "", "anime file path")
	fs.StringVar(&f.gif, "gif", "", "graphic info file path")
	fs.StringVar(&f.gf, "gf", "", "graphic file path")
	fs.StringVar(&f.pgif, "pgif", "", "palette graphic info file path")
	fs.StringVar(&f.pgf, "pgf", "", "palette graphic file path")
	fs.StringVar(&f.pf, "pf", "", "palette file path")
	fs.StringVar(&f.outdir, "o", "output", "output directory")
	fs.BoolVar(&f.dr, "dry-run", false, "dump without output files (for testing)")

	return
}

var (
	bar *progressbar.ProgressBar
	f   flags
)

// DumpAnime the entrypoint of "dump-anime" command
func DumpAnime(ctx context.Context, args []string) (err error) {
	if err = f.Flags().Parse(args); err != nil {
		return
	}

	res := pkg.Resources{}
	pres := pkg.Resources{}
	defer res.Close()
	defer pres.Close()

	if err = res.OpenAnimeResource(f.aif); err != nil {
		return
	}
	if err = res.OpenAnime(f.af); err != nil {
		return
	}
	if err = res.OpenGraphicResource(f.gif); err != nil {
		return
	}
	if err = res.OpenGraphic(f.gf); err != nil {
		return
	}
	if err = res.OpenPalette(f.pf); err != nil {
		return
	}
	if f.pgif != "" {
		if err = pres.OpenGraphicResource(f.pgif); err != nil {
			return
		}
	}
	if f.pgf != "" {
		if err = pres.OpenGraphic(f.pgf); err != nil {
			return
		}
	}

	bar = progressbar.Default(int64(len(res.AnimeResource)))
	done := make(chan struct{})

	go func() {
		defer close(done)
		for _, ai := range res.AnimeResource {
			select {
			case <-ctx.Done():
				return
			default:
				var p color.Palette
				if p, err = palette(res, pres, ai); err != nil {
					return
				}
				if err = dumpAnime(ai, res.AnimeFile, res.GraphicResource, res.GraphicFile, p); err != nil {
					log.Err(err).Send()
				}
				_ = bar.Add(1)
			}
		}
	}()

	<-done

	return
}

func palette(res pkg.Resources, pres pkg.Resources, ai pkg.AnimeIndex) (p color.Palette, err error) {
	// use hidden palette
	if len(pres.GraphicResource.MDx) > 0 {
		if _, ok := pres.GraphicResource.MDx[int32(ai.Info.ID)]; ok {
			var pg *pkg.Graphic
			if err = pres.GraphicResource.MDx[int32(ai.Info.ID)][0].Load(pres.GraphicFile); err != nil {
				return nil, err
			}

			return pg.PaletteData, nil
		}

		log.Debug().Msgf("hidden palette not found: %+v", ai)
	}

	// use cgp palette
	if len(res.Palette) > 0 {
		return res.Palette, nil
	}

	return nil, fmt.Errorf("%w: %d", errPaletteNotFound, ai.Info.ID)
}

func dumpAnime(ai pkg.AnimeIndex, af *os.File, gr pkg.GraphicResource, gf *os.File, p color.Palette) (err error) {
	if err = ai.Load(af, gr); err != nil {
		return
	}

	for i, animes := range ai.Animes {
		for _, a := range animes {
			var img *gif.GIF
			if img, err = a.GIF(gf, p); err != nil {
				log.Err(err).Msgf("anime: %+v", ai.Info)
				return
			}

			var out *os.File
			if f.dr {
				out, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
			} else {
				out, err = os.OpenFile(fmt.Sprintf("%s/%d-%d.gif", filepath.Clean(f.outdir), ai.Info.ID, i), os.O_WRONLY|os.O_CREATE, 0644)
			}
			if err != nil {
				log.Err(err).Msgf("anime: %+v", ai.Info)
				return
			}

			if err = gif.EncodeAll(out, img); err != nil {
				log.Err(err).Msgf("anime: %+v", ai.Info)
				return
			}
		}
	}

	return
}
