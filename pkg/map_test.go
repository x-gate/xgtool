package pkg

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestMakeMap(t *testing.T) {
	maps, _ := os.ReadDir("../testdata/map")

	for _, f := range maps {
		t.Run(f.Name(), func(t *testing.T) {
			file, err := os.Open("../testdata/map/" + f.Name())
			if err != nil {
				t.Fatal(err)
			}
			defer file.Close()

			m, err := MakeMap(file)
			if err != nil {
				t.Fatal(err)
			}

			if len(m.Ground) != int(m.Header.Width*m.Header.Height) {
				t.Errorf("len(m.Ground) = %d, want %d", len(m.Ground), int(m.Header.Width*m.Header.Height))
			}

			if len(m.Object) != int(m.Header.Width*m.Header.Height) {
				t.Errorf("len(m.Object) = %d, want %d", len(m.Object), int(m.Header.Width*m.Header.Height))
			}

			if len(m.Meta) != int(m.Header.Width*m.Header.Height) {
				t.Errorf("len(m.Meta) = %d, want %d", len(m.Meta), int(m.Header.Width*m.Header.Height))
			}

			t.Logf("%+v", m.Header)
		})
	}
}

func TestMap_TiledMap(t *testing.T) {
	gif, err := os.Open("../testdata/graphic_info/GraphicInfo_66.bin")
	if errors.Is(err, os.ErrNotExist) {
		t.Skipf("skipping test; file %s does not exist", "../testdata/graphic_info/GraphicInfo_66.bin")
	}
	defer gif.Close()
	_, mi, _ := MakeGraphicInfoIndexes(gif)

	gf, _ := os.Open("../testdata/graphic/Graphic_66.bin")
	defer gf.Close()

	pf, _ := os.Open("../testdata/palette/palet_00.cgp")
	defer pf.Close()
	p, _ := NewPaletteFromCGP(pf)

	mf, _ := os.Open("../testdata/map/1091.dat")
	defer mf.Close()
	m, _ := MakeMap(mf)

	tm, err := m.TiledMap(mi, gf, p, "../output/")
	if err != nil {
		t.Fatal(err)
	}

	out, _ := os.OpenFile("../output/map.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer out.Close()

	_, err = json.Marshal(tm)
	if err != nil {
		t.Fatal(err)
	}
}
