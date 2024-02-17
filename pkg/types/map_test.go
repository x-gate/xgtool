package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/google/go-cmp/cmp"
	"io"
	"os"
	"testing"
)

func TestMapHeader(t *testing.T) {
	const path = testdata + "/map/"

	testcases := []struct {
		filename string
		expected MapHeader
	}{
		{
			filename: "100.dat", // 芙蕾雅島
			expected: MapHeader{Magic: [12]byte{'M', 'A', 'P', 0x1}, Width: 840, Height: 610},
		},
		{
			filename: "1000.dat", // 法蘭城
			expected: MapHeader{Magic: [12]byte{'M', 'A', 'P', 0x1}, Width: 300, Height: 300},
		},
		{
			filename: "2000.dat", // 伊爾村
			expected: MapHeader{Magic: [12]byte{'M', 'A', 'P', 0x1}, Width: 90, Height: 120},
		},
		{
			filename: "300.dat", // 索奇亞島
			expected: MapHeader{Magic: [12]byte{'M', 'A', 'P', 0x1}, Width: 800, Height: 532},
		},
		{
			filename: "3000.dat", // 加納村
			expected: MapHeader{Magic: [12]byte{'M', 'A', 'P', 0x1}, Width: 100, Height: 100},
		},
		{
			filename: "3200.dat", // 奇利村
			expected: MapHeader{Magic: [12]byte{'M', 'A', 'P', 0x1}, Width: 150, Height: 150},
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

			var mh MapHeader
			if err := binary.Read(f, binary.LittleEndian, &mh); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expected, mh); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestMapData(t *testing.T) {
	const path = testdata + "/map/"

	maps, _ := os.ReadDir(path)
	for _, f := range maps {
		t.Run(f.Name(), func(t *testing.T) {
			data, err := os.ReadFile(path + f.Name())
			if err != nil {
				t.Fatal(err)
			}

			reader := bytes.NewReader(data)

			var mh MapHeader
			if err := binary.Read(reader, binary.LittleEndian, &mh); err != nil {
				t.Fatal(err)
			}

			buf := make([]byte, mh.Width*mh.Height*2)
			for range 3 {
				if _, err := io.ReadFull(reader, buf); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
