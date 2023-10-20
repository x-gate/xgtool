package pkg

import (
	"bytes"
	"encoding/binary"
	"errors"
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
			aif, err := os.Open(tc.filename)
			if err != nil && errors.Is(err, os.ErrNotExist) {
				t.Skipf("skipping test; file %s does not exist", tc.filename)
			} else if err != nil {
				t.Fatal(err)
			}
			defer aif.Close()

			index, err := MakeAnimeInfoIndex(aif)
			if err != nil {
				t.Fatal(err)
			}

			if len(index) != tc.expected {
				t.Errorf("expected len(index): %d, got %d", tc.expected, len(index))
			}
		})
	}
}

func TestAnimeInfo_LoadAnime(t *testing.T) {
	testcases := []struct {
		infoName        string
		animeName       string
		graphicInfoName string
		expectHeader    animeHeader
	}{
		{
			infoName:        "../testdata/anime_info/AnimeInfo_4.bin",
			animeName:       "../testdata/anime/Anime_4.bin",
			graphicInfoName: "../testdata/graphic_info/GraphicInfo_66.bin",
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
			aif, err := os.Open(tc.infoName)
			if err != nil && errors.Is(err, os.ErrNotExist) {
				t.Skipf("skipping test; file %s does not exist", tc.infoName)
			} else if err != nil {
				t.Fatal(err)
			}
			defer aif.Close()

			af, err := os.Open(tc.animeName)
			if err != nil && errors.Is(err, os.ErrNotExist) {
				t.Skipf("skipping test; file %s does not exist", tc.animeName)
			} else if err != nil {
				t.Fatal(err)
			}
			defer af.Close()

			ai, err := readAnimeInfo(aif)

			gif, err := os.Open("../testdata/graphic_info/GraphicInfo_66.bin")
			if err != nil && errors.Is(err, os.ErrNotExist) {
				t.Skipf("skipping test; file %s does not exist", tc.animeName)
			} else if err != nil {
				t.Fatal(err)
			}
			defer gif.Close()
			gf, err := os.Open("../testdata/graphic/Graphic_66.bin")
			if err != nil && errors.Is(err, os.ErrNotExist) {
				t.Skipf("skipping test; file %s does not exist", tc.animeName)
			} else if err != nil {
				t.Fatal(err)
			}

			idx, _, err := MakeGraphicInfoIndexes(gif)
			if err != nil {
				t.Fatal(err)
			}

			a, err := ai.LoadAnime(af, idx, gf)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expectHeader, a.Header); diff != "" {
				t.Errorf("anime header mismatch (-want +got):\n%s", diff)
			}

			if len(a.Frames) != int(a.Header.FrameCnt) {
				t.Errorf("expected len(a.Frames): %d, got %d", a.Header.FrameCnt, len(a.Frames))
			}

			t.Logf("%+v", a)
		})
	}
}

func TestAnime_GIF(t *testing.T) {
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
			name:        "AnimeInfoV3_8.bin",
			animeInfo:   "../testdata/anime_info/AnimeInfoV3_8.bin",
			anime:       "../testdata/anime/AnimeV3_8.bin",
			graphicInfo: "../testdata/graphic_info/GraphicInfoV3_19.bin",
			graphic:     "../testdata/graphic/GraphicV3_19.bin",
			palette:     "../testdata/palette/palet_00.cgp",
		},
		{
			name:        "AnimeInfo_PUK2_4.bin",
			animeInfo:   "../testdata/anime_info/AnimeInfo_PUK2_4.bin",
			anime:       "../testdata/anime/Anime_PUK2_4.bin",
			graphicInfo: "../testdata/graphic_info/GraphicInfo_PUK2_2.bin",
			graphic:     "../testdata/graphic/Graphic_PUK2_2.bin",
			palette:     "../testdata/palette/palet_00.cgp",
		},
		{
			name:        "AnimeInfo_PUK3_2.bin",
			animeInfo:   "../testdata/anime_info/AnimeInfo_PUK3_2.bin",
			anime:       "../testdata/anime/Anime_PUK3_2.bin",
			graphicInfo: "../testdata/graphic_info/GraphicInfo_PUK3_1.bin",
			graphic:     "../testdata/graphic/Graphic_PUK3_1.bin",
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
		{
			name:        "AnimeInfo_Joy_CH1.bin",
			animeInfo:   "../testdata/anime_info/AnimeInfo_Joy_CH1.bin",
			anime:       "../testdata/anime/Anime_Joy_CH1.bin",
			graphicInfo: "../testdata/graphic_info/GraphicInfo_Joy_CH1.bin",
			graphic:     "../testdata/graphic/Graphic_Joy_CH1.bin",
			palette:     "../testdata/palette/palet_00.cgp",
		},
		{
			name:        "AnimeInfo_Joy_EX_146.bin",
			animeInfo:   "../testdata/anime_info/AnimeInfo_Joy_EX_146.bin",
			anime:       "../testdata/anime/Anime_Joy_EX_146.bin",
			graphicInfo: "../testdata/graphic_info/GraphicInfo_Joy_EX_152.bin",
			graphic:     "../testdata/graphic/Graphic_Joy_EX_152.bin",
			palette:     "../testdata/palette/palet_00.cgp",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			aif, err := os.Open(tc.animeInfo)
			if err != nil && errors.Is(err, os.ErrNotExist) {
				t.Skipf("skipping test; file %s does not exist", tc.animeInfo)
			} else if err != nil {
				t.Fatal(err)
			}
			defer aif.Close()

			af, err := os.Open(tc.anime)
			if err != nil && errors.Is(err, os.ErrNotExist) {
				t.Skipf("skipping test; file %s does not exist", tc.anime)
			} else if err != nil {
				t.Fatal(err)
			}
			defer af.Close()

			gif, err := os.Open(tc.graphicInfo)
			if err != nil && errors.Is(err, os.ErrNotExist) {
				t.Skipf("skipping test; file %s does not exist", tc.graphicInfo)
			} else if err != nil {
				t.Fatal(err)
			}
			defer gif.Close()

			gf, err := os.Open(tc.graphic)
			if err != nil && errors.Is(err, os.ErrNotExist) {
				t.Skipf("skipping test; file %s does not exist", tc.graphic)
			} else if err != nil {
				t.Fatal(err)
			}
			defer gf.Close()

			pf, err := os.Open(tc.palette)
			if err != nil && errors.Is(err, os.ErrNotExist) {
				t.Skipf("skipping test; file %s does not exist", tc.palette)
			} else if err != nil {
				t.Fatal(err)
			}
			defer pf.Close()

			ai, err := readAnimeInfo(aif)
			if err != nil {
				t.Fatal(err)
			}

			gidx, _, err := MakeGraphicInfoIndexes(gif)
			if err != nil {
				t.Fatal(err)
			}

			a, err := ai.LoadAnime(af, gidx, gf)
			if err != nil {
				t.Fatal(err)
			}

			p, err := NewPaletteFromCGP(pf)
			if err != nil {
				t.Fatal(err)
			}

			img, err := a.GIF(p)
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

func readAnimeInfo(f *os.File) (ai AnimeInfo, err error) {
	buf := bytes.NewBuffer(make([]byte, 12))

	if _, err = io.ReadFull(f, buf.Bytes()); err != nil {
		return
	}

	err = binary.Read(buf, binary.LittleEndian, &ai)

	return
}
