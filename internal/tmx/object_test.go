package tmx

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecodeObject(t *testing.T) {
	testcases := []struct {
		name   string
		data   string
		expect Object
	}{
		{
			name: "object",
			data: "testdata/object_object.json",
			expect: Object{
				GID:    5,
				Height: 0,
				ID:     1,
				Name:   "villager",
				Properties: []Property{
					{
						Name:  "hp",
						Type:  "int",
						Value: float64(12),
					},
				},
				Rotation: 0,
				Type:     "npc",
				Visible:  true,
				Width:    0,
				X:        32,
				Y:        32,
			},
		},
		{
			name: "ellipse",
			data: "testdata/object_ellipse.json",
			expect: Object{
				Ellipse:  true,
				Height:   152,
				ID:       13,
				Name:     "",
				Rotation: 0,
				Type:     "",
				Visible:  true,
				Width:    248,
				X:        560,
				Y:        808,
			},
		},
		{
			name: "rectangle",
			data: "testdata/object_rectangle.json",
			expect: Object{
				Height:   184,
				ID:       14,
				Name:     "",
				Rotation: 0,
				Type:     "",
				Visible:  true,
				Width:    368,
				X:        576,
				Y:        584,
			},
		},
		{
			name: "point",
			data: "testdata/object_point.json",
			expect: Object{
				Height:   0,
				ID:       20,
				Name:     "",
				Point:    true,
				Rotation: 0,
				Type:     "",
				Visible:  true,
				Width:    0,
				X:        220,
				Y:        350,
			},
		},
		{
			name: "polygon",
			data: "testdata/object_polygon.json",
			expect: Object{
				Height: 0,
				ID:     15,
				Name:   "",
				Polygon: []Point{
					{0, 0},
					{152, 88},
					{136, -128},
					{80, -280},
					{16, -288},
				},
				Rotation: 0,
				Type:     "",
				Visible:  true,
				Width:    0,
				X:        -176,
				Y:        432,
			},
		},
		{
			name: "polyline",
			data: "testdata/object_polyline.json",
			expect: Object{
				Height: 0,
				ID:     16,
				Polyline: []Point{
					{0, 0},
					{248, -32},
					{376, 72},
					{544, 288},
					{656, 120},
					{512, 0},
				},
				Rotation: 0,
				Type:     "",
				Visible:  true,
				Width:    0,
				X:        240,
				Y:        88,
			},
		},
		{
			name: "text",
			data: "testdata/object_text.json",
			expect: Object{
				Height: 19,
				ID:     15,
				Name:   "",
				Text: &Text{
					Text: "Hello World",
					Wrap: true,
				},
				Rotation: 0,
				Type:     "",
				Visible:  true,
				Width:    248,
				X:        48,
				Y:        136,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			data, _ := os.ReadFile(tc.data)

			var object Object
			if err := json.Unmarshal(data, &object); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expect, object); diff != "" {
				t.Errorf("object mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
