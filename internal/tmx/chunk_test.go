package tmx

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecodeChunk(t *testing.T) {
	data, _ := os.ReadFile("testdata/chunk.json")
	expect := Chunk{
		Data:   []uint{1, 2, 1, 2, 3, 1, 3, 1, 2, 2, 3, 3, 4, 4, 4, 1},
		Height: 16,
		Width:  16,
		X:      0,
		Y:      -16,
	}

	var c Chunk
	if err := json.Unmarshal(data, &c); err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expect, c); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
