package tmx

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecodeTileSet(t *testing.T) {
	data, _ := os.ReadFile("testdata/tileset.json")
	expect := TileSet{
		Columns:     19,
		FirstGID:    1,
		Image:       "../image/fishbaddie_parts.png",
		ImageHeight: 480,
		ImageWidth:  640,
		Margin:      3,
		Name:        "",
		Properties: []Property{
			{
				Name:  "myProperty1",
				Type:  "string",
				Value: "myProperty1_value",
			},
		},
		Spacing:    1,
		TileCount:  266,
		TileHeight: 32,
		TileWidth:  32,
	}

	var tileset TileSet
	err := json.Unmarshal(data, &tileset)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, tileset); diff != "" {
		t.Errorf("tileset mismatch (-want +got):\n%s", diff)
	}
}

func TestDecodeTile(t *testing.T) {
	data, _ := os.ReadFile("testdata/tile.json")
	expect := Tile{
		ID: 11,
		Properties: []Property{
			{
				Name:  "myProperty2",
				Type:  "string",
				Value: "myProperty2_value",
			},
		},
		Terrain: []int{0, 1, 0, 1},
	}

	var tile Tile
	err := json.Unmarshal(data, &tile)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, tile); diff != "" {
		t.Errorf("tile mismatch (-want +got):\n%s", diff)
	}
}

func TestDecodeTerrain(t *testing.T) {
	data, _ := os.ReadFile("testdata/terrain.json")
	expect := Terrain{
		Name: "ground",
		Tile: 0,
	}

	var terrain Terrain
	if err := json.Unmarshal(data, &terrain); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, terrain); diff != "" {
		t.Errorf("terrain mismatch (-want +got):\n%s", diff)
	}
}

func TestDecodeWangColor(t *testing.T) {
	data, _ := os.ReadFile("testdata/wangcolor.json")
	expect := WangColor{
		Color:       "#d31313",
		Name:        "Rails",
		Probability: 1,
		Tile:        18,
	}

	var wangcolor WangColor
	if err := json.Unmarshal(data, &wangcolor); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, wangcolor); diff != "" {
		t.Errorf("wangcolor mismatch (-want +got):\n%s", diff)
	}
}

func TestDecodeWangTile(t *testing.T) {
	data, _ := os.ReadFile("testdata/wangtile.json")
	expect := WangTile{
		TileID: 0,
		WangID: [8]int{2, 0, 1, 0, 1, 0, 2, 0},
	}

	var wangtile WangTile
	if err := json.Unmarshal(data, &wangtile); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, wangtile); diff != "" {
		t.Errorf("wangtile mismatch (-want +got):\n%s", diff)
	}
}
