package pkg

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMakeGraphicIndex(t *testing.T) {
	testcases := []struct {
		infoName string
		expected [2]int // [0] = len(idIndex), [1] = len(mapIndex)
	}{

		{"../testdata/graphic_info/GraphicInfo_66.bin", [...]int{252788, 21209}},
		{"../testdata/graphic_info/GraphicInfoEx_5.bin", [...]int{343869, 7390}},
		{"../testdata/graphic_info/GraphicInfoV3_19.bin", [...]int{20024, 2672}},
		{"../testdata/graphic_info/GraphicInfo_PUK2_2.bin", [...]int{11033, 4032}},
		{"../testdata/graphic_info/GraphicInfo_PUK3_1.bin", [...]int{4592, 162}},
		{"../testdata/graphic_info/GraphicInfo_Joy_125.bin", [...]int{493880, 5250}},
		{"../testdata/graphic_info/GraphicInfo_Joy_CH1.bin", [...]int{53541, 268}},
		{"../testdata/graphic_info/GraphicInfo_Joy_EX_152.bin", [...]int{199515, 810}},
	}

	for _, tc := range testcases {
		t.Run(tc.infoName, func(t *testing.T) {
			res := Resources{}
			defer res.Close()

			err := res.OpenGraphicInfo(tc.infoName)
			skipIfNotExists(tc.infoName, err, t)

			if len(res.GraphicIDIndex) != tc.expected[0] {
				t.Errorf("expected len(index): %d, got %d", tc.expected, len(res.GraphicIDIndex))
			}
			if len(res.GraphicMapIndex) != tc.expected[1] {
				t.Errorf("expected len(index): %d, got %d", tc.expected, len(res.GraphicMapIndex))
			}
		})
	}
}

func TestNewGraphicResource(t *testing.T) {
	testcases := []struct {
		infoName string
		expected [2]int // [0] = len(gres.idx), [1] = len(gres.mdx)
	}{
		{
			infoName: "../testdata/graphic_info/GraphicInfo_66.bin",
			expected: [...]int{252788, 21209},
		},
		{
			infoName: "../testdata/graphic_info/GraphicInfoEx_5.bin",
			expected: [...]int{343869, 7390},
		},
		{
			infoName: "../testdata/graphic_info/GraphicInfoV3_19.bin",
			expected: [...]int{20024, 2672},
		},
		{
			infoName: "../testdata/graphic_info/GraphicInfo_PUK2_2.bin",
			expected: [...]int{11033, 4032},
		},
		{
			infoName: "../testdata/graphic_info/GraphicInfo_PUK3_1.bin",
			expected: [...]int{4592, 162},
		},
		{
			infoName: "../testdata/graphic_info/GraphicInfo_Joy_125.bin",
			expected: [...]int{493880, 5250},
		},
		{
			infoName: "../testdata/graphic_info/GraphicInfo_Joy_CH1.bin",
			expected: [...]int{53541, 268},
		},
		{
			infoName: "../testdata/graphic_info/GraphicInfo_Joy_EX_152.bin",
			expected: [...]int{199515, 810},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.infoName, func(t *testing.T) {
			res := Resources{}
			defer res.Close()

			var err error
			err = res.OpenGraphicInfo(tc.infoName)
			skipIfNotExists(tc.infoName, err, t)

			gres, err := NewGraphicResource(res.GraphicInfoFile)
			if err != nil {
				t.Fatal(err)
			}

			if len(gres.idx) != tc.expected[0] {
				t.Errorf("expected len(gres.idx): %d, got %d", tc.expected[0], len(gres.idx))
			}
			if len(gres.mdx) != tc.expected[1] {
				t.Errorf("expected len(gres.mdx): %d, got %d", tc.expected[1], len(gres.mdx))
			}
		})
	}
}

func TestGraphicInfo_LoadGraphic(t *testing.T) {
	testcases := []struct {
		infoName           string
		graphicName        string
		expectedHeader     GraphicHeader
		expectedGraphicLen int
		expectedPaletteLen int
	}{
		{
			infoName:    "../testdata/graphic_info/GraphicInfo_66.bin",
			graphicName: "../testdata/graphic/Graphic_66.bin",
			expectedHeader: GraphicHeader{
				Magic:   [2]byte{'R', 'D'},
				Version: 1,
				Width:   64,
				Height:  47,
				Len:     424,
			},
			expectedGraphicLen: 3008,
			expectedPaletteLen: 0,
		},
		{
			infoName:    "../testdata/graphic_info/GraphicInfoEx_5.bin",
			graphicName: "../testdata/graphic/GraphicEx_5.bin",
			expectedHeader: GraphicHeader{
				Magic:   [2]byte{'R', 'D'},
				Version: 1,
				Width:   64,
				Height:  47,
				Len:     1697,
			},
			expectedGraphicLen: 3008,
			expectedPaletteLen: 0,
		},
		{
			infoName:    "../testdata/graphic_info/GraphicInfoV3_19.bin",
			graphicName: "../testdata/graphic/GraphicV3_19.bin",
			expectedHeader: GraphicHeader{
				Magic:   [2]byte{'R', 'D'},
				Version: 1,
				Width:   228,
				Height:  165,
				Len:     18895,
			},
			expectedGraphicLen: 37620,
			expectedPaletteLen: 0,
		},
		{
			infoName:    "../testdata/graphic_info/GraphicInfo_PUK2_2.bin",
			graphicName: "../testdata/graphic/Graphic_PUK2_2.bin",
			expectedHeader: GraphicHeader{
				Magic:   [2]byte{'R', 'D'},
				Version: 3,
				Width:   640,
				Height:  480,
				Len:     2012,
			},
			expectedGraphicLen: 307200,
			expectedPaletteLen: 768 / 3,
		},
		{
			infoName:    "../testdata/graphic_info/GraphicInfo_PUK3_1.bin",
			graphicName: "../testdata/graphic/Graphic_PUK3_1.bin",
			expectedHeader: GraphicHeader{
				Magic:   [2]byte{'R', 'D'},
				Version: 3,
				Width:   548,
				Height:  450,
				Len:     107742,
			},
			expectedGraphicLen: 246600,
			expectedPaletteLen: 768 / 3,
		},
		{
			infoName:    "../testdata/graphic_info/GraphicInfo_Joy_125.bin",
			graphicName: "../testdata/graphic/Graphic_Joy_125.bin",
			expectedHeader: GraphicHeader{
				Magic:   [2]byte{'R', 'D'},
				Version: 3,
				Width:   80,
				Height:  15,
				Len:     563,
			},
			expectedGraphicLen: 1200,
			expectedPaletteLen: 63 / 3,
		},
		{
			infoName:    "../testdata/graphic_info/GraphicInfo_Joy_CH1.bin",
			graphicName: "../testdata/graphic/Graphic_Joy_CH1.bin",
			expectedHeader: GraphicHeader{
				Magic:   [2]byte{'R', 'D'},
				Version: 3,
				Width:   88,
				Height:  149,
				Len:     6545,
			},
			expectedGraphicLen: 13112,
			expectedPaletteLen: 768 / 3,
		},
		{
			infoName:    "../testdata/graphic_info/GraphicInfo_Joy_EX_152.bin",
			graphicName: "../testdata/graphic/Graphic_Joy_EX_152.bin",
			expectedHeader: GraphicHeader{
				Magic:   [2]byte{'R', 'D'},
				Version: 3,
				Width:   88,
				Height:  149,
				Len:     6545,
			},
			expectedGraphicLen: 13112,
			expectedPaletteLen: 768 / 3,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.graphicName, func(t *testing.T) {
			res := Resources{}
			defer res.Close()

			var err error
			err = res.OpenGraphicInfo(tc.infoName)
			skipIfNotExists(tc.infoName, err, t)
			err = res.OpenGraphic(tc.graphicName)
			skipIfNotExists(tc.graphicName, err, t)

			gi, err := readGraphicInfo(res.GraphicInfoFile)
			if err != nil {
				t.Fatal(err)
			}

			g, err := gi.LoadGraphic(res.GraphicFile)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expectedHeader, g.Header); diff != "" {
				t.Errorf("graphic header mismatch (-want +got):\n%s", diff)
			}
			if len(g.GraphicData) != tc.expectedGraphicLen {
				t.Errorf("expected len(g.GraphicData): %d, got %d", tc.expectedGraphicLen, len(g.GraphicData))
			}
			if len(g.PaletteData) != tc.expectedPaletteLen {
				t.Errorf("expected len(g.PaletteLen): %d, got %d", tc.expectedPaletteLen, len(g.PaletteData))
			}

			t.Logf("header: %+v, len(graphic): %d, len(palette): %d", g.Header, len(g.GraphicData), len(g.PaletteData))
		})
	}
}

func TestGraphicResource_LoadGraphic(t *testing.T) {
	const GraphicInfoFile = "../testdata/graphic_info/GraphicInfo_66.bin"
	const GraphicFile = "../testdata/graphic/Graphic_66.bin"
	const PaletteFile = "../testdata/palette/palet_00.cgp"

	res := Resources{}
	defer res.Close()

	var err error
	err = res.OpenGraphicInfo(GraphicInfoFile)
	skipIfNotExists(GraphicInfoFile, err, t)
	err = res.OpenGraphic(GraphicFile)
	skipIfNotExists(GraphicFile, err, t)
	err = res.OpenPalette(PaletteFile)
	skipIfNotExists(PaletteFile, err, t)

	gres, _ := NewGraphicResource(res.GraphicInfoFile)
	gres.LoadGraphic(res.GraphicFile, res.Palette)

	if len(gres.idx[0][0].GraphicData) != 3008 {
		t.Errorf("expected len(gres.idx[0][0].GraphicData): %d, got %d", 3008, len(gres.idx[0][0].GraphicData))
	}
}

func TestGraphic_Img(t *testing.T) {
	testcases := []struct {
		name string
		gif  string
		gf   string
		pf   string
	}{
		{
			name: "GraphicInfo_66.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_66.bin",
			gf:   "../testdata/graphic/Graphic_66.bin",
			pf:   "../testdata/palette/palet_00.cgp",
		},
		{
			name: "GraphicInfoEx_5.bin",
			gif:  "../testdata/graphic_info/GraphicInfoEx_5.bin",
			gf:   "../testdata/graphic/GraphicEx_5.bin",
			pf:   "../testdata/palette/palet_00.cgp",
		},
		{
			name: "GraphicInfoV3_19.bin",
			gif:  "../testdata/graphic_info/GraphicInfoV3_19.bin",
			gf:   "../testdata/graphic/GraphicV3_19.bin",
			pf:   "../testdata/palette/palet_00.cgp",
		},
		{
			name: "GraphicInfo_PUK2_2.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_PUK2_2.bin",
			gf:   "../testdata/graphic/Graphic_PUK2_2.bin",
		},
		{
			name: "GraphicInfo_PUK3_1.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_PUK3_1.bin",
			gf:   "../testdata/graphic/Graphic_PUK3_1.bin",
		},
		{
			name: "GraphicInfo_Joy_125.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_Joy_125.bin",
			gf:   "../testdata/graphic/Graphic_Joy_125.bin",
		},
		{
			name: "GraphicInfo_Joy_CH1.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_Joy_CH1.bin",
			gf:   "../testdata/graphic/Graphic_Joy_CH1.bin",
		},
		{
			name: "GraphicInfo_Joy_EX_152.bin",
			gif:  "../testdata/graphic_info/GraphicInfo_Joy_EX_152.bin",
			gf:   "../testdata/graphic/Graphic_Joy_EX_152.bin",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res := Resources{}
			defer res.Close()

			var err error
			err = res.OpenGraphicInfo(tc.gif)
			skipIfNotExists(tc.gif, err, t)
			err = res.OpenGraphic(tc.gf)
			skipIfNotExists(tc.gf, err, t)

			gi, err := readGraphicInfo(res.GraphicInfoFile)
			if err != nil {
				t.Fatal(err)
			}

			g, err := gi.LoadGraphic(res.GraphicFile)
			if err != nil {
				t.Fatal(err)
			}

			if tc.pf != "" {
				err = res.OpenPalette(tc.pf)
				g.PaletteData = res.Palette
			}

			_, err = g.ImgRGBA()
			if err != nil {
				t.Fatal(err)
			}

			_, err = g.ImgPaletted()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func BenchmarkGraphic_ImgRGBA(b *testing.B) {
	res := Resources{}
	defer res.Close()

	_ = res.OpenGraphicInfo("../testdata/graphic_info/GraphicInfo_66.bin")
	_ = res.OpenGraphic("../testdata/graphic/Graphic_66.bin")

	gi, _ := readGraphicInfo(res.GraphicInfoFile)
	g, _ := gi.LoadGraphic(res.GraphicFile)

	_ = res.OpenPalette("../testdata/palette/palet_00.cgp")
	g.PaletteData = res.Palette

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = g.ImgRGBA()
	}
}

func BenchmarkGraphic_ImgPaletted(b *testing.B) {
	res := Resources{}
	defer res.Close()

	_ = res.OpenGraphicInfo("../testdata/graphic_info/GraphicInfo_66.bin")
	_ = res.OpenGraphic("../testdata/graphic/Graphic_66.bin")

	gi, _ := readGraphicInfo(res.GraphicInfoFile)
	g, _ := gi.LoadGraphic(res.GraphicFile)

	_ = res.OpenPalette("../testdata/palette/palet_00.cgp")
	g.PaletteData = res.Palette

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = g.ImgPaletted()
	}
}

func readGraphicInfo(f *os.File) (gi GraphicInfo, err error) {
	buf := bytes.NewBuffer(make([]byte, 40))

	if _, err = io.ReadFull(f, buf.Bytes()); err != nil {
		return
	}

	err = binary.Read(buf, binary.LittleEndian, &gi)

	return
}
