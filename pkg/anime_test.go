package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	gif2 "image/gif"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
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
			res := TestRes{}
			defer res.Close()

			err := res.OpenAnimeInfo(tc.filename)
			SkipIfNotExists(tc.filename, err, t)

			if len(res.AnimeInfoIndex) != tc.expected {
				t.Errorf("expected len(index): %d, got %d", tc.expected, len(res.AnimeInfoIndex))
			}
		})
	}
}

func TestAnimeInfo_LoadAnime(t *testing.T) {
	testcases := []struct {
		infoName        string
		animeName       string
		graphicInfoName string
		graphicName     string
		expectHeader    animeHeader
	}{
		{
			infoName:        "../testdata/anime_info/AnimeInfo_4.bin",
			animeName:       "../testdata/anime/Anime_4.bin",
			graphicInfoName: "../testdata/graphic_info/GraphicInfo_66.bin",
			graphicName:     "../testdata/graphic/Graphic_66.bin",
			expectHeader: animeHeader{
				Direct:   0,
				Action:   5,
				Duration: 1500,
				FrameCnt: 26,
			},
		},
		{
			infoName:        "../testdata/anime_info/AnimeInfoEx_1.bin",
			animeName:       "../testdata/anime/AnimeEx_1.bin",
			graphicInfoName: "../testdata/graphic_info/GraphicInfoEx_5.bin",
			graphicName:     "../testdata/graphic/GraphicEx_5.bin",
			expectHeader: animeHeader{
				Direct:   0,
				Action:   5,
				Duration: 1000,
				FrameCnt: 8,
			},
		},
		{
			infoName:        "../testdata/anime_info/AnimeInfoV3_8.bin",
			animeName:       "../testdata/anime/AnimeV3_8.bin",
			graphicInfoName: "../testdata/graphic_info/GraphicInfoV3_19.bin",
			graphicName:     "../testdata/graphic/GraphicV3_19.bin",
			expectHeader: animeHeader{
				Direct:   0,
				Action:   5,
				Duration: 700,
				FrameCnt: 10,
				Reversed: 4,
				Sentinel: -1,
			},
		},
		{
			infoName:        "../testdata/anime_info/AnimeInfo_PUK2_4.bin",
			animeName:       "../testdata/anime/Anime_PUK2_4.bin",
			graphicInfoName: "../testdata/graphic_info/GraphicInfo_PUK2_2.bin",
			graphicName:     "../testdata/graphic/Graphic_PUK2_2.bin",
			expectHeader: animeHeader{
				Direct:   0,
				Action:   5,
				Duration: 1200,
				FrameCnt: 22,
				Reversed: 0,
				Sentinel: -1,
			},
		},
		{
			infoName:        "../testdata/anime_info/AnimeInfo_PUK3_2.bin",
			animeName:       "../testdata/anime/Anime_PUK3_2.bin",
			graphicInfoName: "../testdata/graphic_info/GraphicInfo_PUK3_1.bin",
			graphicName:     "../testdata/graphic/Graphic_PUK3_1.bin",
			expectHeader: animeHeader{
				Direct:   0,
				Action:   0,
				Duration: 700,
				FrameCnt: 21,
				Reversed: 0,
				Sentinel: -1,
			},
		},
		{
			infoName:        "../testdata/anime_info/AnimeInfo_Joy_91.bin",
			animeName:       "../testdata/anime/Anime_Joy_91.bin",
			graphicInfoName: "../testdata/graphic_info/GraphicInfo_Joy_125.bin",
			graphicName:     "../testdata/graphic/Graphic_Joy_125.bin",
			expectHeader: animeHeader{
				Direct:   0,
				Action:   0,
				Duration: 2000,
				FrameCnt: 17,
				Reversed: 0,
				Sentinel: -1,
			},
		},
		{
			infoName:        "../testdata/anime_info/AnimeInfo_Joy_CH1.Bin",
			animeName:       "../testdata/anime/Anime_Joy_CH1.bin",
			graphicInfoName: "../testdata/graphic_info/GraphicInfo_Joy_CH1.bin",
			graphicName:     "../testdata/graphic/Graphic_Joy_CH1.bin",
			expectHeader: animeHeader{
				Direct:   0,
				Action:   5,
				Duration: 1000,
				FrameCnt: 8,
				Reversed: 0,
				Sentinel: -1,
			},
		},
		{
			infoName:        "../testdata/anime_info/AnimeInfo_Joy_EX_146.bin",
			animeName:       "../testdata/anime/Anime_Joy_EX_146.bin",
			graphicInfoName: "../testdata/graphic_info/GraphicInfo_Joy_EX_152.bin",
			graphicName:     "../testdata/graphic/Graphic_Joy_EX_152.bin",
			expectHeader: animeHeader{
				Direct:   0,
				Action:   5,
				Duration: 1000,
				FrameCnt: 8,
				Reversed: 0,
				Sentinel: -1,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.animeName, func(t *testing.T) {
			res := TestRes{}
			defer res.Close()

			var err error

			err = res.OpenAnimeInfo(tc.infoName)
			SkipIfNotExists(tc.infoName, err, t)
			err = res.OpenAnime(tc.animeName)
			SkipIfNotExists(tc.animeName, err, t)
			err = res.OpenGraphicInfo(tc.graphicInfoName)
			SkipIfNotExists(tc.graphicInfoName, err, t)
			err = res.OpenGraphic(tc.graphicName)
			SkipIfNotExists(tc.graphicName, err, t)

			ai, err := readAnimeInfo(res.AnimeInfoFile)

			a, err := ai.LoadAnime(res.AnimeFile, res.GraphicInfoIDIndex, res.GraphicFile)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expectHeader, a.Header); diff != "" {
				t.Errorf("anime header mismatch (-want +got):\n%s", diff)
			}

			if len(a.Frames) != int(a.Header.FrameCnt) {
				t.Errorf("expected len(a.Frames): %d, got %d", a.Header.FrameCnt, len(a.Frames))
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
			res := TestRes{}
			defer res.Close()

			var err error

			err = res.OpenAnimeInfo(tc.animeInfo)
			SkipIfNotExists(tc.animeInfo, err, t)
			err = res.OpenAnime(tc.anime)
			SkipIfNotExists(tc.anime, err, t)
			err = res.OpenGraphicInfo(tc.graphicInfo)
			SkipIfNotExists(tc.graphicInfo, err, t)
			err = res.OpenGraphic(tc.graphic)
			SkipIfNotExists(tc.graphic, err, t)
			err = res.OpenPalette(tc.palette)
			SkipIfNotExists(tc.palette, err, t)

			ai, err := readAnimeInfo(res.AnimeInfoFile)
			if err != nil {
				t.Fatal(err)
			}

			a, err := ai.LoadAnime(res.AnimeFile, res.GraphicInfoIDIndex, res.GraphicFile)
			if err != nil {
				t.Fatal(err)
			}

			img, err := a.GIF(res.Palette)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("%+v", img)

			if len(img.Image) != int(a.Header.FrameCnt) {
				t.Errorf("expected len(img.Image): %d, got %d", a.Header.FrameCnt, len(img.Image))
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
			res := TestRes{}
			pres := TestRes{}
			defer res.Close()
			defer pres.Close()

			var err error

			err = res.OpenAnimeInfo(tc.animeInfo)
			SkipIfNotExists(tc.animeInfo, err, t)
			err = res.OpenAnime(tc.anime)
			SkipIfNotExists(tc.anime, err, t)
			err = res.OpenGraphicInfo(tc.graphicInfo)
			SkipIfNotExists(tc.graphicInfo, err, t)
			err = res.OpenGraphic(tc.graphic)
			SkipIfNotExists(tc.graphic, err, t)

			err = pres.OpenGraphicInfo(tc.paletteGraphicInfo)
			SkipIfNotExists(tc.paletteGraphicInfo, err, t)
			err = pres.OpenGraphic(tc.paletteGraphic)
			SkipIfNotExists(tc.paletteGraphic, err, t)

			ai, err := readAnimeInfo(res.AnimeInfoFile)
			if err != nil {
				t.Fatal(err)
			}

			a, err := ai.LoadAnime(res.AnimeFile, res.GraphicInfoIDIndex, res.GraphicFile)
			if err != nil {
				t.Fatal(err)
			}

			var pg *Graphic
			if _, ok := pres.GraphicInfoMapIndex[ai.ID]; !ok {
				t.Fatalf("gmdx[%d] graphic info not found", ai.ID)
			}
			if pg, err = pres.GraphicInfoMapIndex[ai.ID].LoadGraphic(pres.GraphicFile); err != nil {
				t.Fatal(err)
			}

			img, err := a.GIF(pg.PaletteData)
			if err != nil {
				t.Fatal(err)
			}

			if len(img.Image) != int(a.Header.FrameCnt) {
				t.Errorf("expected len(img.Image): %d, got %d", a.Header.FrameCnt, len(img.Image))
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
