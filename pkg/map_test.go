package pkg

import (
	"encoding/json"
	"os"
	"testing"
)

func TestMakeMap(t *testing.T) {
	maps, _ := os.ReadDir("../testdata/map")

	for _, f := range maps {
		t.Run(f.Name(), func(t *testing.T) {
			res := Resources{}
			defer res.Close()

			if err := res.OpenMap("../testdata/map/" + f.Name()); err != nil {
				t.Fatal(err)
			}

			if len(res.Map.Ground) != int(res.Map.Header.Width*res.Map.Header.Height) {
				t.Errorf("len(m.Ground) = %d, want %d", len(res.Map.Ground), int(res.Map.Header.Width*res.Map.Header.Height))
			}

			if len(res.Map.Object) != int(res.Map.Header.Width*res.Map.Header.Height) {
				t.Errorf("len(m.Object) = %d, want %d", len(res.Map.Object), int(res.Map.Header.Width*res.Map.Header.Height))
			}

			if len(res.Map.Meta) != int(res.Map.Header.Width*res.Map.Header.Height) {
				t.Errorf("len(m.Meta) = %d, want %d", len(res.Map.Meta), int(res.Map.Header.Width*res.Map.Header.Height))
			}
		})
	}
}

func TestMap_TiledMap(t *testing.T) {
	res := Resources{}
	defer res.Close()

	var err error
	err = res.OpenGraphicResource("../testdata/graphic_info/GraphicInfo_66.bin")
	err = res.OpenGraphic("../testdata/graphic/Graphic_66.bin")
	err = res.OpenPalette("../testdata/palette/palet_00.cgp")
	err = res.OpenMap("../testdata/map/1091.dat")
	skipIfNotExists("../testdata/map/1091.dat", err, t)

	tm, err := res.Map.TiledMap(res.GraphicResource.MDx, res.GraphicFile, res.Palette, "../output/")
	if err != nil {
		t.Fatal(err)
	}

	_, err = json.Marshal(tm)
	if err != nil {
		t.Fatal(err)
	}
}
