package tmx

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecodeMap(t *testing.T) {
	data, _ := os.ReadFile("testdata/map.json")
	expect := Map{
		BackgroundColor: "#656667",
		Height:          4,
		Layers:          []Layer{},
		NextObjectID:    1,
		Orientation:     Orthogonal,
		Properties: []Property{
			{
				Name:  "mapProperty1",
				Type:  "string",
				Value: "one",
			},
			{
				Name:  "mapProperty2",
				Type:  "string",
				Value: "two",
			},
		},
		RenderOrder:  RightDown,
		TileHeight:   32,
		TileSets:     []TileSet{},
		TileWidth:    32,
		Version:      "1",
		TiledVersion: "1.0.3",
		Width:        4,
	}

	var m Map
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, m); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
