package types

// AnimeInfo structure for each anime info in `AnimeInfo*.bin`, 12 bytes.
type AnimeInfo struct {
	ID     int32 // ID, it's the index of the anime in the file, but it could be duplicated.
	Addr   int32
	ActCnt int16
	_      int16
}

// AnimeDataHeader struct for each anime header in `Anime*.bin`, 12 bytes.
type AnimeDataHeader struct {
	Direct   int16
	Action   int16
	Duration int32 // Duration in milliseconds
	FrameCnt int32
}

// AnimeDataHeaderEx struct for each anime header in `Anime*.bin` for version 2 and above (PUK2 and above), 20 bytes.
type AnimeDataHeaderEx struct {
	AnimeDataHeader
	_        int16
	Reversed int16 // Reversed, odd number for reversed, even number for normal
	Sentinel int32 // Sentinel, always 0xFFFF (-1) if it is an extended header
}

// AnimeDataFrame structure for each anime frame in `Anime*.bin`, 10 bytes.
type AnimeDataFrame struct {
	GraphicID int32 // Graphic ID, it's the ID of the graphic in `GraphicInfo*.bin`.
	OffX      int16
	OffY      int16
	Flag      int16
}

// AnimeData structure for each anime data in `Anime*.bin`, 12/18(header) + frame_cnt*10(frames) bytes.
type AnimeData[H AnimeDataHeader | AnimeDataHeaderEx] struct {
	Header H
	Frames []AnimeDataFrame
}
