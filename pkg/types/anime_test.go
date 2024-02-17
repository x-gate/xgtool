package types

import (
	"encoding/binary"
	"errors"
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

func TestAnimeInfo(t *testing.T) {
	const path = testdata + "/anime_info/"

	testcases := []struct {
		filename string
		expected AnimeInfo
	}{
		{
			filename: "AnimeInfo_4.bin",
			expected: AnimeInfo{ID: 100000, ActCnt: 160},
		},
		{
			filename: "AnimeInfoEx_1.Bin",
			expected: AnimeInfo{ID: 105000, ActCnt: 160},
		},
		{
			filename: "AnimeInfoV3_8.bin",
			expected: AnimeInfo{ID: 110350, ActCnt: 152},
		},
		{
			filename: "AnimeInfo_PUK2_4.bin",
			expected: AnimeInfo{ID: 110300, ActCnt: 64},
		},
		{
			filename: "AnimeInfo_PUK3_2.bin",
			expected: AnimeInfo{ID: 110505, ActCnt: 64},
		},
		{
			filename: "AnimeInfo_Joy_91.bin",
			expected: AnimeInfo{ID: 104000, ActCnt: 1},
		},
		{
			filename: "AnimeInfo_Joy_CH1.Bin",
			expected: AnimeInfo{ID: 104854, ActCnt: 72},
		},
		{
			filename: "AnimeInfo_Joy_EX_146.bin",
			expected: AnimeInfo{ID: 104854, ActCnt: 72},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.filename, func(t *testing.T) {
			f, err := os.Open(path + tc.filename)
			if errors.Is(err, os.ErrNotExist) {
				t.Skipf("skip %s: not found", tc.filename)
			} else if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = f.Close() }()

			var ai AnimeInfo
			if err := binary.Read(f, binary.LittleEndian, &ai); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expected, ai); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAnimeDataHeader(t *testing.T) {
	const path = testdata + "/anime/"

	testcases := []struct {
		filename string
		expected AnimeDataHeader
	}{
		{
			filename: "Anime_4.bin",
			expected: AnimeDataHeader{Action: 5, Duration: 1500, FrameCnt: 26},
		},
		{
			filename: "AnimeEx_1.Bin",
			expected: AnimeDataHeader{Action: 5, Duration: 1000, FrameCnt: 8},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.filename, func(t *testing.T) {
			f, err := os.Open(path + tc.filename)
			if errors.Is(err, os.ErrNotExist) {
				t.Skipf("skip %s: not found", tc.filename)
			} else if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = f.Close() }()

			var ah AnimeDataHeader
			if err := binary.Read(f, binary.LittleEndian, &ah); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expected, ah); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAnimeDataHeaderEx(t *testing.T) {
	const path = testdata + "/anime/"

	testcases := []struct {
		filename string
		expected AnimeDataHeaderEx
	}{
		{
			filename: "AnimeV3_8.bin",
			expected: AnimeDataHeaderEx{AnimeDataHeader: AnimeDataHeader{Action: 5, Duration: 700, FrameCnt: 10}, Reversed: 4, Sentinel: -1},
		},
		{
			filename: "Anime_PUK2_4.bin",
			expected: AnimeDataHeaderEx{AnimeDataHeader: AnimeDataHeader{Action: 5, Duration: 1200, FrameCnt: 22}, Sentinel: -1},
		},
		{
			filename: "Anime_PUK3_2.bin",
			expected: AnimeDataHeaderEx{AnimeDataHeader: AnimeDataHeader{Duration: 700, FrameCnt: 21}, Sentinel: -1},
		},
		{
			filename: "Anime_Joy_91.bin",
			expected: AnimeDataHeaderEx{AnimeDataHeader: AnimeDataHeader{Duration: 2000, FrameCnt: 17}, Sentinel: -1},
		},
		{
			filename: "Anime_Joy_CH1.Bin",
			expected: AnimeDataHeaderEx{AnimeDataHeader: AnimeDataHeader{Action: 5, Duration: 1000, FrameCnt: 8}, Sentinel: -1},
		},
		{
			filename: "Anime_Joy_EX_146.bin",
			expected: AnimeDataHeaderEx{AnimeDataHeader: AnimeDataHeader{Action: 5, Duration: 1000, FrameCnt: 8}, Sentinel: -1},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.filename, func(t *testing.T) {
			f, err := os.Open(path + tc.filename)
			if errors.Is(err, os.ErrNotExist) {
				t.Skipf("skip %s: not found", tc.filename)
			} else if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = f.Close() }()

			var ah AnimeDataHeaderEx
			if err := binary.Read(f, binary.LittleEndian, &ah); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expected, ah); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAnimeDataFrame(t *testing.T) {
	const path = testdata + "/anime/"

	type testcase struct {
		filename string
		expected AnimeDataFrame
	}
	f := func(offset int64, tc testcase) func(t *testing.T) {
		return func(t *testing.T) {
			f, err := os.Open(path + tc.filename)
			if errors.Is(err, os.ErrNotExist) {
				t.Skipf("skip %s: not found", tc.filename)
			} else if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = f.Close() }()

			if _, err = f.Seek(offset, 0); err != nil {
				t.Fatal(err)
			}

			var ad AnimeDataFrame
			if err := binary.Read(f, binary.LittleEndian, &ad); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expected, ad); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		}
	}

	testcases := []testcase{
		{
			filename: "Anime_4.bin",
			expected: AnimeDataFrame{GraphicID: 19642},
		},
		{
			filename: "AnimeEx_1.Bin",
			expected: AnimeDataFrame{GraphicID: 7100},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.filename, f(12, tc))
	}

	testcases = []testcase{
		{
			filename: "AnimeV3_8.bin",
			expected: AnimeDataFrame{GraphicID: 3365},
		},
		{
			filename: "Anime_PUK2_4.bin",
			expected: AnimeDataFrame{GraphicID: 5095},
		},
		{
			filename: "Anime_PUK3_2.bin",
			expected: AnimeDataFrame{GraphicID: 422},
		},
		{
			filename: "Anime_Joy_91.bin",
			expected: AnimeDataFrame{GraphicID: 627},
		},
		{
			filename: "Anime_Joy_CH1.Bin",
			expected: AnimeDataFrame{GraphicID: 2},
		},
		{
			filename: "Anime_Joy_EX_146.bin",
			expected: AnimeDataFrame{GraphicID: 2},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.filename, f(20, tc))
	}
}
