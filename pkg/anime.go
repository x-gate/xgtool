package pkg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
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
	Header  animeHeader
	Frames  []animeFrame
	Graphic []*Graphic
}

// AnimeIndex is a map is a map of AnimeInfo, key is the ID of the anime.
type AnimeIndex map[int32]AnimeInfo

// MakeAnimeInfoIndex reads anime info from src, and returns AnimeInfoIndex.
func MakeAnimeInfoIndex(src io.Reader) (AnimeIndex, error) {
	index := make(AnimeIndex)

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

// LoadAnime loads anime data from anime file.
func (ai AnimeInfo) LoadAnime(af *os.File, idx GraphicInfoIndex, gf io.ReadSeeker) (a *Anime, err error) {
	a = new(Anime)

	a.Header, err = ai.parseHeader(af)

	buf := bytes.NewBuffer(make([]byte, 10*a.Header.FrameCnt))
	if _, err = io.ReadFull(af, buf.Bytes()); err != nil {
		return
	}

	for i := 0; i < int(a.Header.FrameCnt); i++ {
		var frame animeFrame
		if err = binary.Read(buf, binary.LittleEndian, &frame); err != nil {
			return
		}

		var g *Graphic
		if g, err = idx[frame.GraphicID].LoadGraphic(gf); err != nil {
			return
		}

		a.Frames = append(a.Frames, frame)
		a.Graphic = append(a.Graphic, g)
	}

	return
}

func (ai AnimeInfo) parseHeader(af *os.File) (h animeHeader, err error) {
	h, err = ai.readAnimeHeader(af, 20)
	if err != nil {
		return
	}

	// check if this anime header is extended or not
	// h.Sentinel will be -1 if it's extended
	if h.Sentinel == -1 {
		return
	}

	return ai.readAnimeHeader(af, 12)
}

func (ai AnimeInfo) readAnimeHeader(af *os.File, len int) (h animeHeader, err error) {
	if _, err = af.Seek(int64(ai.Addr), io.SeekStart); err != nil {
		return
	}

	buf := bytes.NewBuffer(make([]byte, len))
	if _, err = io.ReadFull(af, buf.Bytes()); err != nil {
		return
	}

	h = *(*animeHeader)(unsafe.Pointer(&buf.Bytes()[0]))

	return
}
