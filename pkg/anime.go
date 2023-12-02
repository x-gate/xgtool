package pkg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"image/gif"
	"io"
	"os"
	"unsafe"
)

// AnimeInfo structure for each anime info, 12 bytes
type AnimeInfo struct {
	ID     int32
	Addr   int32
	ActCnt int16
	_      int16
}

// animeHeader structure for each anime header, 12 bytes for <= V2, 20 bytes > V3.
type animeHeader struct {
	Direct   int16
	Action   int16
	Duration int32
	FrameCnt int32
	_        int16 // v3 only
	Reversed int16 // v3 only
	Sentinel int32 // v3 only
}

// animeFrame structure for each anime frame, 10 bytes.
type animeFrame struct {
	GraphicID int32
	OffX      int16
	OffY      int16
	Flag      int16
}

// Anime stores data for each anime, not a strict mapping to the file.
type Anime struct {
	Info    AnimeInfo
	Header  animeHeader
	Frames  []animeFrame
	Graphic []*Graphic
}

// AnimeInfoIndex is a map is a map of AnimeInfo, key is the ID of the anime.
type AnimeInfoIndex map[int32]AnimeInfo

// MakeAnimeInfoIndex reads anime info from src, and returns AnimeInfoIndex.
func MakeAnimeInfoIndex(src io.Reader) (AnimeInfoIndex, error) {
	index := make(AnimeInfoIndex)

	r := bufio.NewReader(src)
	for {
		buf := bytes.NewBuffer(make([]byte, 12))
		if _, err := io.ReadFull(r, buf.Bytes()); err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		var info AnimeInfo
		if err := binary.Read(buf, binary.LittleEndian, &info); err != nil {
			return nil, err
		}

		index[info.ID] = info
	}

	return index, nil
}

// LoadAllAnimes loads all animes (with all directions and actions) from anime file.
func (ai AnimeInfo) LoadAllAnimes(af *os.File, idx GraphicInfoIndex, gf io.ReadSeeker) (animes []*Anime, err error) {
	animes = make([]*Anime, 0, ai.ActCnt)

	if _, err = af.Seek(int64(ai.Addr), io.SeekStart); err != nil {
		return
	}

	headerSize := getHeaderSize(af)

	if _, err = af.Seek(int64(ai.Addr), io.SeekStart); err != nil {
		return
	}

	for i := 0; i < int(ai.ActCnt); i++ {
		a := new(Anime)

		if a.Header, err = ai.readAnimeHeader(af, headerSize); err != nil {
			return
		}
		if a.Frames, a.Graphic, err = ai.readAnimeFrames(af, int(a.Header.FrameCnt), idx, gf); err != nil {
			return
		}

		animes = append(animes, a)
	}

	return
}

func (ai AnimeInfo) readAnimeHeader(af io.Reader, len int) (h animeHeader, err error) {
	buf := bytes.NewBuffer(make([]byte, len))
	if _, err = io.ReadFull(af, buf.Bytes()); err != nil {
		return
	}

	h = *(*animeHeader)(unsafe.Pointer(&buf.Bytes()[0]))
	if len == 12 {
		h.Reversed = 0
		h.Sentinel = 0
	}

	return
}

func (ai AnimeInfo) readAnimeFrames(af io.Reader, cnt int, idx GraphicInfoIndex, gf io.ReadSeeker) (f []animeFrame, g []*Graphic, err error) {
	f = make([]animeFrame, 0, cnt)
	g = make([]*Graphic, 0, cnt)

	buf := bytes.NewBuffer(make([]byte, 10*cnt))
	if _, err = io.ReadFull(af, buf.Bytes()); err != nil {
		return
	}

	for i := 0; i < cnt; i++ {
		var frame animeFrame
		if err = binary.Read(buf, binary.LittleEndian, &frame); err != nil {
			return
		}

		var graphic *Graphic
		if graphic, err = idx[frame.GraphicID].LoadGraphic(gf); err != nil {
			return
		}

		f = append(f, frame)
		g = append(g, graphic)
	}

	return
}

// GIF returns a gif.GIF from the anime.
func (a Anime) GIF(p color.Palette) (img *gif.GIF, err error) {
	img = new(gif.GIF)

	var w, h int
	for _, g := range a.Graphic {
		if len(g.PaletteData) == 0 {
			g.PaletteData = p
		}

		w = max(w, int(g.Header.Width))
		h = max(h, int(g.Header.Height))

		var i *image.Paletted
		if i, err = g.ImgPaletted(); err != nil {
			return
		}

		img.Image = append(img.Image, i)
		img.Delay = append(img.Delay, int(a.Header.Duration)/int(a.Header.FrameCnt)/10)
		img.Disposal = append(img.Disposal, gif.DisposalBackground)
	}

	img.Config = image.Config{
		Width:  w,
		Height: h,
	}

	return
}

func getHeaderSize(af *os.File) (sz int) {
	buf := bytes.NewBuffer(make([]byte, 20))
	if _, err := io.ReadFull(af, buf.Bytes()); err != nil {
		return
	}

	h := *(*animeHeader)(unsafe.Pointer(&buf.Bytes()[0]))

	// check if this anime header is extended or not
	// h.Sentinel will be -1 if it's extended
	if h.Sentinel == -1 {
		return 20
	}

	return 12
}
