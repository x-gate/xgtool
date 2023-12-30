package pkg

import (
	"fmt"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
	"image/gif"
	"os"
	"testing"
)

func TestNewAnimeResource(t *testing.T) {
	testcases := []struct {
		filename string
		expected int
	}{
		{"../testdata/anime_info/AnimeInfo_4.bin", 806},
		{"../testdata/anime_info/AnimeInfoEx_1.bin", 827},
		{"../testdata/anime_info/AnimeInfoV3_8.bin", 342},
		{"../testdata/anime_info/AnimeInfo_PUK2_4.bin", 343},
		{"../testdata/anime_info/AnimeInfo_PUK3_2.bin", 149},
		{"../testdata/anime_info/AnimeInfo_Joy_91.bin", 785},
		{"../testdata/anime_info/AnimeInfo_Joy_CH1.Bin", 68},
		{"../testdata/anime_info/AnimeInfo_Joy_EX_146.bin", 569},
	}

	for _, tc := range testcases {
		t.Run(tc.filename, func(t *testing.T) {
			res := Resources{}
			defer res.Close()

			err := res.OpenAnimeResource(tc.filename)
			skipIfNotExists(tc.filename, err, t)

			if len(res.AnimeResource) != tc.expected {
				t.Errorf("expected len(index): %d, got %d", tc.expected, len(res.AnimeResource))
			}
		})
	}
}

func TestAnimeIndex_Load(t *testing.T) {
	testcases := []struct {
		name string
		aif  string
		af   string
		gif  string
	}{
		{
			name: "AnimeInfo_4.bin",
			aif:  "../testdata/anime_info/AnimeInfo_4.bin",
			af:   "../testdata/anime/Anime_4.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_66.bin",
		},
		{
			name: "AnimeInfoEx_1.bin",
			aif:  "../testdata/anime_info/AnimeInfoEx_1.Bin",
			af:   "../testdata/anime/AnimeEx_1.Bin",
			gif:  "../testdata/graphic_info/GraphicInfoEx_5.bin",
		},
		{
			name: "AnimeInfoV3_8.bin",
			aif:  "../testdata/anime_info/AnimeInfoV3_8.bin",
			af:   "../testdata/anime/AnimeV3_8.bin",
			gif:  "../testdata/graphic_info/GraphicInfoV3_19.bin",
		},
		{
			name: "AnimeInfo_PUK2_4.bin",
			aif:  "../testdata/anime_info/AnimeInfo_PUK2_4.bin",
			af:   "../testdata/anime/Anime_PUK2_4.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_PUK2_2.bin",
		},
		{
			name: "AnimeInfo_PUK3_2.bin",
			aif:  "../testdata/anime_info/AnimeInfo_PUK3_2.bin",
			af:   "../testdata/anime/Anime_PUK3_2.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_PUK3_1.bin",
		},
		{
			name: "AnimeInfo_Joy_91.bin",
			aif:  "../testdata/anime_info/AnimeInfo_Joy_91.bin",
			af:   "../testdata/anime/Anime_Joy_91.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_Joy_125.bin",
		},
		{
			name: "AnimeInfo_Joy_CH1.bin",
			aif:  "../testdata/anime_info/AnimeInfo_Joy_CH1.Bin",
			af:   "../testdata/anime/Anime_Joy_CH1.Bin",
			gif:  "../testdata/graphic_info/GraphicInfo_Joy_CH1.bin",
		},
		{
			name: "AnimeInfo_Joy_EX_146.bin",
			aif:  "../testdata/anime_info/AnimeInfo_Joy_EX_146.bin",
			af:   "../testdata/anime/Anime_Joy_EX_146.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_Joy_EX_152.bin",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res := Resources{}
			defer res.Close()

			var err error
			err = res.OpenAnimeResource(tc.aif)
			err = res.OpenAnime(tc.af)
			err = res.OpenGraphicResource(tc.gif)
			skipIfNotExists(tc.gif, err, t)

			var aidx AnimeIndex
			for _, aidx = range res.AnimeResource {
				break
			}

			err = aidx.Load(res.AnimeFile, res.GraphicResource)
			if err != nil {
				t.Logf("%+v", aidx)
				t.Fatal(err)
			}

			sum := lo.SumBy(maps.Values(aidx.Animes), func(animes []Anime) int { return len(animes) })
			if sum != int(aidx.Info.ActCnt) {
				t.Errorf("expected sum: %d, got %d", aidx.Info.ActCnt, sum)
			}
		})
	}
}

func TestAnime_GIF_ExternalPalette(t *testing.T) {
	testcases := []struct {
		name string
		aif  string
		af   string
		gif  string
		gf   string
		pf   string
	}{
		{
			name: "AnimeInfo_4.bin",
			aif:  "../testdata/anime_info/AnimeInfo_4.bin",
			af:   "../testdata/anime/Anime_4.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_66.bin",
			gf:   "../testdata/graphic/Graphic_66.bin",
			pf:   "../testdata/palette/palet_00.cgp",
		},
		{
			name: "AnimeInfoEx_1.bin",
			aif:  "../testdata/anime_info/AnimeInfoEx_1.Bin",
			af:   "../testdata/anime/AnimeEx_1.Bin",
			gif:  "../testdata/graphic_info/GraphicInfoEx_5.bin",
			gf:   "../testdata/graphic/GraphicEx_5.bin",
			pf:   "../testdata/palette/palet_00.cgp",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res := Resources{}
			defer res.Close()

			var err error
			err = res.OpenAnimeResource(tc.aif)
			err = res.OpenAnime(tc.af)
			err = res.OpenGraphicResource(tc.gif)
			err = res.OpenGraphic(tc.gf)
			err = res.OpenPalette(tc.pf)
			skipIfNotExists(tc.pf, err, t)

			var aidx AnimeIndex
			for _, aidx = range res.AnimeResource {
				break
			}
			_ = aidx.Load(res.AnimeFile, res.GraphicResource)

			var a Anime
			for _, a = range aidx.Animes[0] {
				break
			}

			var img *gif.GIF
			if img, err = a.GIF(res.GraphicFile, res.Palette); err != nil {
				t.Fatal(err)
			}

			t.Logf("%+v", img)

			if len(img.Image) != int(a.Header.FrameCnt) {
				t.Errorf("expected len(img.Image): %d, got %d", a.Header.FrameCnt, len(img.Image))
			}

			out, _ := os.OpenFile(fmt.Sprintf("../output/%s.gif", tc.name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			defer out.Close()

			if err = gif.EncodeAll(out, img); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAnime_GIF_InternalPalette(t *testing.T) {
	testcases := []struct {
		name string
		aif  string
		af   string
		gif  string
		gf   string
		pgif string
		pgf  string
	}{
		{
			name: "AnimeInfoV3_8.bin",
			aif:  "../testdata/anime_info/AnimeInfoV3_8.bin",
			af:   "../testdata/anime/AnimeV3_8.bin",
			gif:  "../testdata/graphic_info/GraphicInfoV3_19.bin",
			gf:   "../testdata/graphic/GraphicV3_19.bin",
			pgif: "../testdata/graphic_info/GraphicInfoV3_19.bin",
			pgf:  "../testdata/graphic/GraphicV3_19.bin",
		},
		{
			name: "AnimeInfo_PUK2_4.bin",
			aif:  "../testdata/anime_info/AnimeInfo_PUK2_4.bin",
			af:   "../testdata/anime/Anime_PUK2_4.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_PUK2_2.bin",
			gf:   "../testdata/graphic/Graphic_PUK2_2.bin",
			pgif: "../testdata/graphic_info/GraphicInfoV3_19.bin",
			pgf:  "../testdata/graphic/GraphicV3_19.bin",
		},
		{
			name: "AnimeInfo_PUK3_2.bin",
			aif:  "../testdata/anime_info/AnimeInfo_PUK3_2.bin",
			af:   "../testdata/anime/Anime_PUK3_2.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_PUK3_1.bin",
			gf:   "../testdata/graphic/Graphic_PUK3_1.bin",
			pgif: "../testdata/graphic_info/GraphicInfoV3_19.bin",
			pgf:  "../testdata/graphic/GraphicV3_19.bin",
		},
		{
			name: "AnimeInfo_Joy_91.bin",
			aif:  "../testdata/anime_info/AnimeInfo_Joy_91.bin",
			af:   "../testdata/anime/Anime_Joy_91.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_Joy_125.bin",
			gf:   "../testdata/graphic/Graphic_Joy_125.bin",
			pgif: "../testdata/graphic_info/GraphicInfo_Joy_125.bin",
			pgf:  "../testdata/graphic/Graphic_Joy_125.bin",
		},
		{
			name: "AnimeInfo_Joy_CH1.bin",
			aif:  "../testdata/anime_info/AnimeInfo_Joy_CH1.Bin",
			af:   "../testdata/anime/Anime_Joy_CH1.Bin",
			gif:  "../testdata/graphic_info/GraphicInfo_Joy_CH1.bin",
			gf:   "../testdata/graphic/Graphic_Joy_CH1.bin",
			pgif: "../testdata/graphic_info/GraphicInfo_Joy_CH1.bin",
			pgf:  "../testdata/graphic/Graphic_Joy_CH1.bin",
		},
		{
			name: "AnimeInfo_Joy_EX_146.bin",
			aif:  "../testdata/anime_info/AnimeInfo_Joy_EX_146.bin",
			af:   "../testdata/anime/Anime_Joy_EX_146.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_Joy_EX_152.bin",
			gf:   "../testdata/graphic/Graphic_Joy_EX_152.bin",
			pgif: "../testdata/graphic_info/GraphicInfo_Joy_EX_152.bin",
			pgf:  "../testdata/graphic/Graphic_Joy_EX_152.bin",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res := Resources{}
			pres := Resources{}
			defer res.Close()
			defer pres.Close()

			var err error
			err = res.OpenAnimeResource(tc.aif)
			err = res.OpenAnime(tc.af)
			err = res.OpenGraphicResource(tc.gif)
			err = res.OpenGraphic(tc.gf)
			err = pres.OpenGraphicResource(tc.pgif)
			err = pres.OpenGraphic(tc.pgf)
			skipIfNotExists(tc.pgf, err, t)

			var aidx AnimeIndex
			for _, aidx = range res.AnimeResource {
				break
			}
			_ = aidx.Load(res.AnimeFile, res.GraphicResource)

			var a Anime
			for _, a = range aidx.Animes[0] {
				break
			}

			if err = pres.GraphicResource.MDx.Load(int32(aidx.Info.ID), pres.GraphicFile); err != nil {
				t.Fatal(err)
			}
			pg := pres.GraphicResource.MDx.First(int32(aidx.Info.ID))
			if pg == nil {
				t.Logf("palette graphic (id: %d) is nil", aidx.Info.ID)
				return
			}

			var img *gif.GIF
			if img, err = a.GIF(res.GraphicFile, pg.PaletteData); err != nil {
				t.Fatal(err)
			}

			t.Logf("%+v", img)

			if len(img.Image) != int(a.Header.FrameCnt) {
				t.Errorf("expected len(img.Image): %d, got %d", a.Header.FrameCnt, len(img.Image))
			}

			out, _ := os.OpenFile(fmt.Sprintf("../output/%s.gif", tc.name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			defer out.Close()

			if err = gif.EncodeAll(out, img); err != nil {
				t.Fatal(err)
			}
		})
	}
}
