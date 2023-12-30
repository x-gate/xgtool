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
	"unsafe"
)

const (
	// AnimeInfoSize every block of anime info is 12 bytes
	AnimeInfoSize = 12

	// AnimeFrameSize every block of anime frame is 10 bytes
	AnimeFrameSize = 10
)

// AnimeID is the ID of an anime, from anime info.
type AnimeID int32

// ActionID is the ID of an action, from anime header
type ActionID int16

type animeInfo struct {
	ID     AnimeID
	Addr   int32
	ActCnt int16
	_      int16
}

type animeHeader struct {
	Direct   int16
	Action   ActionID
	Duration int32
	FrameCnt int32
	_        int16 // v3 only
	Reversed int16 // v3 only
	Sentinel int32 // v3 only
}

type animeFrameData struct {
	GraphicID int32
	OffX      int16
	OffY      int16
	Flag      int16
}

type animeFrame struct {
	Data    animeFrameData
	Graphic *Graphic
}

// Anime is a collection of frames.
type Anime struct {
	Index  AnimeIndex // point to the index of this anime
	Header animeHeader
	Frames []animeFrame
}

// AnimeIndex built from anime info, Animes are grouped by ActionID.
type AnimeIndex struct {
	Info   animeInfo
	Animes map[ActionID][]Anime
}

// AnimeResource is a map of AnimeIndex, key is AnimeID in anime info.
type AnimeResource map[AnimeID]AnimeIndex

// NewAnimeResource creates a new AnimeResource from anime info file.
//
// The anime data hasn't been loaded yet, you need to call AnimeIndex.Load for each AnimeIndex.
func NewAnimeResource(aif io.Reader) (ar AnimeResource, err error) {
	ar = make(AnimeResource)

	r := bufio.NewReaderSize(aif, AnimeInfoSize*100)
	for {
		buf := bytes.NewBuffer(make([]byte, AnimeInfoSize))
		if _, err = io.ReadFull(r, buf.Bytes()); err != nil && errors.Is(err, io.EOF) {
			err = nil
			break
		} else if err != nil {
			return
		}

		var ai animeInfo
		if err = binary.Read(buf, binary.LittleEndian, &ai); err != nil {
			return
		}

		ar[ai.ID] = AnimeIndex{Info: ai, Animes: make(map[ActionID][]Anime)}
	}

	return
}

// Load loads anime data from anime file.
func (aidx AnimeIndex) Load(af io.ReadSeeker, gr GraphicResource) (err error) {
	if _, err = af.Seek(int64(aidx.Info.Addr), io.SeekStart); err != nil {
		return
	}

	hsz := getHeaderSize(af)

	if _, err = af.Seek(int64(aidx.Info.Addr), io.SeekStart); err != nil {
		return
	}

	for i := 0; i < int(aidx.Info.ActCnt); i++ {
		var a Anime

		a.Index = aidx
		if a.Header, err = a.readHeader(af, hsz); err != nil {
			return
		}
		if a.Frames, err = a.readFrames(af, int(a.Header.FrameCnt), gr); err != nil {
			return
		}

		aidx.Animes[a.Header.Action] = append(aidx.Animes[a.Header.Action], a)
	}

	return
}

func (a Anime) readHeader(af io.Reader, sz int) (h animeHeader, err error) {
	buf := bytes.NewBuffer(make([]byte, sz))
	if _, err = io.ReadFull(af, buf.Bytes()); err != nil {
		return
	}

	h = *(*animeHeader)(unsafe.Pointer(&buf.Bytes()[0]))
	if sz == 12 {
		h.Reversed = 0
		h.Sentinel = 0
	}

	return
}

func (a Anime) readFrames(af io.Reader, cnt int, gr GraphicResource) (f []animeFrame, err error) {
	f = make([]animeFrame, 0, cnt)

	buf := bytes.NewBuffer(make([]byte, AnimeFrameSize*cnt))
	if _, err = io.ReadFull(af, buf.Bytes()); err != nil {
		return
	}

	for i := 0; i < cnt; i++ {
		var fd animeFrameData
		if err = binary.Read(buf, binary.LittleEndian, &fd); err != nil {
			return
		}

		f = append(f, animeFrame{Data: fd, Graphic: gr.IDx.First(fd.GraphicID)})
	}

	return
}

// GIF creates a gif from anime frames with given palette
func (a Anime) GIF(gf io.ReadSeeker, p color.Palette) (img *gif.GIF, err error) {
	img = new(gif.GIF)

	var w, h int
	for _, f := range a.Frames {
		if err = f.Graphic.Load(gf); err != nil {
			return
		}

		w = max(w, int(f.Graphic.Header.Width))
		h = max(h, int(f.Graphic.Header.Height))

		var i *image.Paletted
		if i, err = f.Graphic.ImgPaletted(p); err != nil {
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

func getHeaderSize(af io.Reader) (sz int) {
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
