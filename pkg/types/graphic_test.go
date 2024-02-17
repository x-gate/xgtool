package types

import (
	"encoding/binary"
	"errors"
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

func TestGraphicInfo(t *testing.T) {
	const path = testdata + "graphic_info/"

	testcases := []struct {
		filename string
		expected GraphicInfo
	}{
		{
			filename: "GraphicInfo_66.bin",
			expected: GraphicInfo{Len: 424, OffX: -32, OffY: -24, Width: 64, Height: 47, GridW: 1, GridH: 1, Access: 1, MapID: 999},
		},
		{
			filename: "GraphicInfoEx_5.bin",
			expected: GraphicInfo{Len: 1697, OffX: -32, OffY: -24, Width: 64, Height: 47, GridW: 1, GridH: 1, MapID: 223018},
		},
		{
			filename: "GraphicInfoV3_19.bin",
			expected: GraphicInfo{Len: 18895, OffX: -127, OffY: -7, Width: 228, Height: 165, GridW: 1, GridH: 1},
		},
		{
			filename: "GraphicInfo_PUK2_2.bin",
			expected: GraphicInfo{Len: 2012, OffX: -320, OffY: -240, Width: 640, Height: 480, GridW: 1, GridH: 1, MapID: 300000},
		},
		{
			filename: "GraphicInfo_PUK3_1.bin",
			expected: GraphicInfo{Len: 107742, OffX: -312, OffY: -225, Width: 548, Height: 450, GridW: 1, GridH: 1, MapID: 301114},
		},
		{
			filename: "GraphicInfo_Joy_125.bin",
			expected: GraphicInfo{Len: 563, OffX: -40, OffY: -8, Width: 80, Height: 15, GridW: 1, GridH: 1, MapID: 243021},
		},
		{
			filename: "GraphicInfo_Joy_CH1.bin",
			expected: GraphicInfo{Len: 6545, OffX: -37, OffY: -135, Width: 88, Height: 149, GridW: 1, GridH: 1, MapID: 104854},
		},
		{
			filename: "GraphicInfo_Joy_EX_152.bin",
			expected: GraphicInfo{Len: 6545, OffX: -37, OffY: -135, Width: 88, Height: 149, GridW: 1, GridH: 1, MapID: 104854},
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

			var gi GraphicInfo
			if err := binary.Read(f, binary.LittleEndian, &gi); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expected, gi); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGraphicHeader(t *testing.T) {
	const path = testdata + "graphic/"

	testcases := []struct {
		filename string
		expected GraphicDataHeader
	}{
		{
			filename: "Graphic_66.bin",
			expected: GraphicDataHeader{Magic: [2]byte{'R', 'D'}, Version: 1, Width: 64, Height: 47, Len: 424},
		},
		{
			filename: "GraphicEx_5.bin",
			expected: GraphicDataHeader{Magic: [2]byte{'R', 'D'}, Version: 1, Width: 64, Height: 47, Len: 1697},
		},
		{
			filename: "GraphicV3_19.bin",
			expected: GraphicDataHeader{Magic: [2]byte{'R', 'D'}, Version: 1, Width: 228, Height: 165, Len: 18895},
		},
		{
			filename: "Graphic_PUK2_2.bin",
			expected: GraphicDataHeader{Magic: [2]byte{'R', 'D'}, Version: 3, Width: 640, Height: 480, Len: 2012},
		},
		{
			filename: "Graphic_PUK3_1.bin",
			expected: GraphicDataHeader{Magic: [2]byte{'R', 'D'}, Version: 3, Width: 548, Height: 450, Len: 107742},
		},
		{
			filename: "Graphic_Joy_125.bin",
			expected: GraphicDataHeader{Magic: [2]byte{'R', 'D'}, Version: 3, Width: 80, Height: 15, Len: 563},
		},
		{
			filename: "Graphic_Joy_CH1.bin",
			expected: GraphicDataHeader{Magic: [2]byte{'R', 'D'}, Version: 3, Width: 88, Height: 149, Len: 6545},
		},
		{
			filename: "Graphic_Joy_EX_152.bin",
			expected: GraphicDataHeader{Magic: [2]byte{'R', 'D'}, Version: 3, Width: 88, Height: 149, Len: 6545},
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

			var gh GraphicDataHeader
			if err := binary.Read(f, binary.LittleEndian, &gh); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expected, gh); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
