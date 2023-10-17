package pkg

import (
	"bytes"
	"encoding/binary"
	"errors"
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
		infoName     string
		animeName    string
		expectHeader animeHeader
	}{
		{
			infoName:  "../testdata/anime_info/AnimeInfo_4.bin",
			animeName: "../testdata/anime/Anime_4.bin",
			expectHeader: animeHeader{
				Direct:   0,
				Action:   5,
				Duration: 1500,
				FrameCnt: 26,
			},
		},
		{
			infoName:  "../testdata/anime_info/AnimeInfoEx_1.bin",
			animeName: "../testdata/anime/AnimeEx_1.bin",
			expectHeader: animeHeader{
				Direct:   0,
				Action:   5,
				Duration: 1000,
				FrameCnt: 8,
			},
		},
		{
			infoName:  "../testdata/anime_info/AnimeInfoV3_8.bin",
			animeName: "../testdata/anime/AnimeV3_8.bin",
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
			infoName:  "../testdata/anime_info/AnimeInfo_PUK2_4.bin",
			animeName: "../testdata/anime/Anime_PUK2_4.bin",
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
			infoName:  "../testdata/anime_info/AnimeInfo_PUK3_2.bin",
			animeName: "../testdata/anime/Anime_PUK3_2.bin",
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
			infoName:  "../testdata/anime_info/AnimeInfo_Joy_91.bin",
			animeName: "../testdata/anime/Anime_Joy_91.bin",
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
			infoName:  "../testdata/anime_info/AnimeInfo_Joy_CH1.Bin",
			animeName: "../testdata/anime/Anime_Joy_CH1.bin",
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
			infoName:  "../testdata/anime_info/AnimeInfo_Joy_EX_146.bin",
			animeName: "../testdata/anime/Anime_Joy_EX_146.bin",
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

			a, err := ai.LoadAnime(af)
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

func readAnimeInfo(f *os.File) (ai AnimeInfo, err error) {
	buf := bytes.NewBuffer(make([]byte, 12))

	if _, err = io.ReadFull(f, buf.Bytes()); err != nil {
		return
	}

	err = binary.Read(buf, binary.LittleEndian, &ai)

	return
}
