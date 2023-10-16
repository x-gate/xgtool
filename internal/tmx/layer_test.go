package tmx

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecodeLayer(t *testing.T) {
	var testcases = []struct {
		name   string
		data   string
		expect Layer
	}{
		{
			name: "tile layer",
			data: "testdata/layer_tile.json",
			expect: Layer{
				Data:    []uint{1, 2, 1, 2, 3, 1, 3, 1, 2, 2, 3, 3, 4, 4, 4, 1},
				Height:  4,
				Name:    "ground",
				Opacity: 1,
				Properties: []Property{
					{
						Name:  "tileLayerProp",
						Type:  "int",
						Value: float64(1),
					},
				},
				Type:    "tilelayer",
				Visible: true,
				Width:   4,
				X:       0,
				Y:       0,
			},
		},
		{
			name: "object layer",
			data: "testdata/layer_object.json",
			expect: Layer{
				DrawOrder: TopDown,
				Height:    0,
				Name:      "people",
				Objects:   []Object{},
				Opacity:   1,
				Properties: []Property{
					{
						Name:  "layerProp1",
						Type:  "string",
						Value: "someStringValue",
					},
				},
				Type:    "objectgroup",
				Visible: true,
				Width:   0,
				X:       0,
				Y:       0,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			data, _ := os.ReadFile(tc.data)
			var l Layer
			if err := json.Unmarshal(data, &l); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expect, l); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
