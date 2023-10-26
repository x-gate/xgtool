package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	gif2 "image/gif"
	"io"
	"os"
	"testing"
)

func TestMakeAnimeIndex(t *testing.T) {
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

			err := res.OpenAnimeInfo(tc.filename)
			skipIfNotExists(tc.filename, err, t)

			if len(res.AnimeInfoIndex) != tc.expected {
				t.Errorf("expected len(index): %d, got %d", tc.expected, len(res.AnimeInfoIndex))
			}
		})
	}
}

func TestAnimeInfo_LoadAllAnimes(t *testing.T) {
	testcases := []struct {
		name string
		aif  string
		af   string
		gif  string
		gf   string
	}{
		{
			name: "AnimeInfo_4.bin",
			aif:  "../testdata/anime_info/AnimeInfo_4.bin",
			af:   "../testdata/anime/Anime_4.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_66.bin",
			gf:   "../testdata/graphic/Graphic_66.bin",
		},
		{
			name: "AnimeInfoEx_1.bin",
			aif:  "../testdata/anime_info/AnimeInfoEx_1.Bin",
			af:   "../testdata/anime/AnimeEx_1.Bin",
			gif:  "../testdata/graphic_info/GraphicInfoEx_5.bin",
			gf:   "../testdata/graphic/GraphicEx_5.bin",
		},
		{
			name: "AnimeInfoV3_8.bin",
			aif:  "../testdata/anime_info/AnimeInfoV3_8.bin",
			af:   "../testdata/anime/AnimeV3_8.bin",
			gif:  "../testdata/graphic_info/GraphicInfoV3_19.bin",
			gf:   "../testdata/graphic/GraphicV3_19.bin",
		},
		{
			name: "AnimeInfo_PUK2_4.bin",
			aif:  "../testdata/anime_info/AnimeInfo_PUK2_4.bin",
			af:   "../testdata/anime/Anime_PUK2_4.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_PUK2_2.bin",
			gf:   "../testdata/graphic/Graphic_PUK2_2.bin",
		},
		{
			name: "AnimeInfo_PUK3_2.bin",
			aif:  "../testdata/anime_info/AnimeInfo_PUK3_2.bin",
			af:   "../testdata/anime/Anime_PUK3_2.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_PUK3_1.bin",
			gf:   "../testdata/graphic/Graphic_PUK3_1.bin",
		},
		{
			name: "AnimeInfo_Joy_91.bin",
			aif:  "../testdata/anime_info/AnimeInfo_Joy_91.bin",
			af:   "../testdata/anime/Anime_Joy_91.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_Joy_125.bin",
			gf:   "../testdata/graphic/Graphic_Joy_125.bin",
		},
		{
			name: "AnimeInfo_Joy_CH1.bin",
			aif:  "../testdata/anime_info/AnimeInfo_Joy_CH1.Bin",
			af:   "../testdata/anime/Anime_Joy_CH1.Bin",
			gif:  "../testdata/graphic_info/GraphicInfo_Joy_CH1.bin",
			gf:   "../testdata/graphic/Graphic_Joy_CH1.bin",
		},
		{
			name: "AnimeInfo_Joy_EX_146.bin",
			aif:  "../testdata/anime_info/AnimeInfo_Joy_EX_146.bin",
			af:   "../testdata/anime/Anime_Joy_EX_146.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_Joy_EX_152.bin",
			gf:   "../testdata/graphic/Graphic_Joy_EX_152.bin",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res := Resources{}
			defer res.Close()

			var err error
			err = res.OpenAnimeInfo(tc.aif)
			skipIfNotExists(tc.aif, err, t)
			err = res.OpenAnime(tc.af)
			skipIfNotExists(tc.af, err, t)
			err = res.OpenGraphicInfo(tc.gif)
			skipIfNotExists(tc.gif, err, t)
			err = res.OpenGraphic(tc.gf)
			skipIfNotExists(tc.gf, err, t)

			var ai AnimeInfo
			for _, ai = range res.AnimeInfoIndex {
				break
			}

			animes, err := ai.LoadAllAnimes(res.AnimeFile, res.GraphicIDIndex, res.GraphicFile)
			if err != nil {
				t.Logf("%+v", ai)
				t.Fatal(err)
			}

			if len(animes) != int(ai.ActCnt) {
				t.Errorf("expected len(animes): %d, got %d", ai.ActCnt, len(animes))
			}
		})
	}
}

func TestAnime_GIF_1(t *testing.T) {
	testcases := []struct {
		name        string
		animeInfo   string
		anime       string
		graphicInfo string
		graphic     string
		palette     string
	}{
		{
			name:        "AnimeInfo_4.bin",
			animeInfo:   "../testdata/anime_info/AnimeInfo_4.bin",
			anime:       "../testdata/anime/Anime_4.bin",
			graphicInfo: "../testdata/graphic_info/GraphicInfo_66.bin",
			graphic:     "../testdata/graphic/Graphic_66.bin",
			palette:     "../testdata/palette/palet_00.cgp",
		},
		{
			name:        "AnimeInfoEx_1.bin",
			animeInfo:   "../testdata/anime_info/AnimeInfoEx_1.Bin",
			anime:       "../testdata/anime/AnimeEx_1.Bin",
			graphicInfo: "../testdata/graphic_info/GraphicInfoEx_5.bin",
			graphic:     "../testdata/graphic/GraphicEx_5.bin",
			palette:     "../testdata/palette/palet_00.cgp",
		},
		{
			name:        "AnimeInfo_Joy_91.bin",
			animeInfo:   "../testdata/anime_info/AnimeInfo_Joy_91.bin",
			anime:       "../testdata/anime/Anime_Joy_91.bin",
			graphicInfo: "../testdata/graphic_info/GraphicInfo_Joy_125.bin",
			graphic:     "../testdata/graphic/Graphic_Joy_125.bin",
			palette:     "../testdata/palette/palet_00.cgp",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res := Resources{}
			defer res.Close()

			var err error

			err = res.OpenAnimeInfo(tc.animeInfo)
			skipIfNotExists(tc.animeInfo, err, t)
			err = res.OpenAnime(tc.anime)
			skipIfNotExists(tc.anime, err, t)
			err = res.OpenGraphicInfo(tc.graphicInfo)
			skipIfNotExists(tc.graphicInfo, err, t)
			err = res.OpenGraphic(tc.graphic)
			skipIfNotExists(tc.graphic, err, t)
			err = res.OpenPalette(tc.palette)
			skipIfNotExists(tc.palette, err, t)

			ai, err := readAnimeInfo(res.AnimeInfoFile)
			if err != nil {
				t.Fatal(err)
			}

			a, err := ai.LoadAllAnimes(res.AnimeFile, res.GraphicIDIndex, res.GraphicFile)
			if err != nil {
				t.Fatal(err)
			}

			img, err := a[0].GIF(res.Palette)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("%+v", img)

			if len(img.Image) != int(a[0].Header.FrameCnt) {
				t.Errorf("expected len(img.Image): %d, got %d", a[0].Header.FrameCnt, len(img.Image))
			}

			out, err := os.OpenFile(fmt.Sprintf("../output/%s.gif", tc.name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				t.Fatal(err)
			}
			defer out.Close()

			if err = gif2.EncodeAll(out, img); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAnime_GIF_2(t *testing.T) {
	testcases := []struct {
		name               string
		animeInfo          string
		anime              string
		graphicInfo        string
		graphic            string
		paletteGraphicInfo string
		paletteGraphic     string
		palette            string
	}{
		{
			name:               "AnimeInfoV3_8.bin",
			animeInfo:          "../testdata/anime_info/AnimeInfoV3_8.bin",
			anime:              "../testdata/anime/AnimeV3_8.bin",
			graphicInfo:        "../testdata/graphic_info/GraphicInfoV3_19.bin",
			graphic:            "../testdata/graphic/GraphicV3_19.bin",
			paletteGraphicInfo: "../testdata/graphic_info/GraphicInfoV3_19.bin",
			paletteGraphic:     "../testdata/graphic/GraphicV3_19.bin",
		},
		{
			name:               "AnimeInfo_PUK2_4.bin",
			animeInfo:          "../testdata/anime_info/AnimeInfo_PUK2_4.bin",
			anime:              "../testdata/anime/Anime_PUK2_4.bin",
			graphicInfo:        "../testdata/graphic_info/GraphicInfo_PUK2_2.bin",
			graphic:            "../testdata/graphic/Graphic_PUK2_2.bin",
			paletteGraphicInfo: "../testdata/graphic_info/GraphicInfoV3_19.bin",
			paletteGraphic:     "../testdata/graphic/GraphicV3_19.bin",
		},
		{
			name:               "AnimeInfo_PUK3_2.bin",
			animeInfo:          "../testdata/anime_info/AnimeInfo_PUK3_2.bin",
			anime:              "../testdata/anime/Anime_PUK3_2.bin",
			graphicInfo:        "../testdata/graphic_info/GraphicInfo_PUK3_1.bin",
			graphic:            "../testdata/graphic/Graphic_PUK3_1.bin",
			paletteGraphicInfo: "../testdata/graphic_info/GraphicInfoV3_19.bin",
			paletteGraphic:     "../testdata/graphic/GraphicV3_19.bin",
		},
		{
			name:               "AnimeInfo_Joy_CH1.bin",
			animeInfo:          "../testdata/anime_info/AnimeInfo_Joy_CH1.Bin",
			anime:              "../testdata/anime/Anime_Joy_CH1.Bin",
			graphicInfo:        "../testdata/graphic_info/GraphicInfo_Joy_CH1.bin",
			graphic:            "../testdata/graphic/Graphic_Joy_CH1.bin",
			paletteGraphicInfo: "../testdata/graphic_info/GraphicInfo_Joy_CH1.bin",
			paletteGraphic:     "../testdata/graphic/Graphic_Joy_CH1.bin",
		},
		{
			name:               "AnimeInfo_Joy_EX_146.bin",
			animeInfo:          "../testdata/anime_info/AnimeInfo_Joy_EX_146.bin",
			anime:              "../testdata/anime/Anime_Joy_EX_146.bin",
			graphicInfo:        "../testdata/graphic_info/GraphicInfo_Joy_EX_152.bin",
			graphic:            "../testdata/graphic/Graphic_Joy_EX_152.bin",
			paletteGraphicInfo: "../testdata/graphic_info/GraphicInfo_Joy_EX_152.bin",
			paletteGraphic:     "../testdata/graphic/Graphic_Joy_EX_152.bin",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res := Resources{}
			pres := Resources{}
			defer res.Close()
			defer pres.Close()

			var err error

			err = res.OpenAnimeInfo(tc.animeInfo)
			skipIfNotExists(tc.animeInfo, err, t)
			err = res.OpenAnime(tc.anime)
			skipIfNotExists(tc.anime, err, t)
			err = res.OpenGraphicInfo(tc.graphicInfo)
			skipIfNotExists(tc.graphicInfo, err, t)
			err = res.OpenGraphic(tc.graphic)
			skipIfNotExists(tc.graphic, err, t)

			err = pres.OpenGraphicInfo(tc.paletteGraphicInfo)
			skipIfNotExists(tc.paletteGraphicInfo, err, t)
			err = pres.OpenGraphic(tc.paletteGraphic)
			skipIfNotExists(tc.paletteGraphic, err, t)

			ai, err := readAnimeInfo(res.AnimeInfoFile)
			if err != nil {
				t.Fatal(err)
			}

			a, err := ai.LoadAllAnimes(res.AnimeFile, res.GraphicIDIndex, res.GraphicFile)
			if err != nil {
				t.Fatal(err)
			}

			var pg *Graphic
			if _, ok := pres.GraphicMapIndex[ai.ID]; !ok {
				t.Fatalf("gmdx[%d] graphic info not found", ai.ID)
			}
			if pg, err = pres.GraphicMapIndex[ai.ID].LoadGraphic(pres.GraphicFile); err != nil {
				t.Fatal(err)
			}

			img, err := a[0].GIF(pg.PaletteData)
			if err != nil {
				t.Fatal(err)
			}

			if len(img.Image) != int(a[0].Header.FrameCnt) {
				t.Errorf("expected len(img.Image): %d, got %d", a[0].Header.FrameCnt, len(img.Image))
			}

			out, err := os.OpenFile(fmt.Sprintf("../output/%s.gif", tc.name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				t.Fatal(err)
			}
			defer out.Close()

			if err = gif2.EncodeAll(out, img); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func readAnimeInfo(f *os.File) (ai AnimeInfo, err error) {
	buf := bytes.NewBuffer(make([]byte, 12))

	if _, err = io.ReadFull(f, buf.Bytes()); err != nil {
		return
	}

	err = binary.Read(buf, binary.LittleEndian, &ai)

	return
}
