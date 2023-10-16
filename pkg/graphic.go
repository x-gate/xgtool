package pkg

// GraphicInfo structure for each graphic info, 40 bytes.
type GraphicInfo struct {
	ID     int32
	Addr   int32
	Len    int32
	OffX   int32
	OffY   int32
	Width  int32
	Height int32
	GridW  byte
	GridH  byte
	Access byte
	_      [5]byte
	MapID  int32
}

// graphicHeader structure for each graphic header, 16 bytes.
type graphicHeader struct {
	Magic   [2]byte // "RD" for valid graphic
	Version byte    // 0 for raw data, 1 for encoded data, 2 for raw data with palette, 3 for encoded data with palette
	_       byte    //
	Width   int32   // Width of graphic, it shouldn't be trusted, use GraphicInfo.Width instead.
	Height  int32   // Height of graphic, it shouldn't be trusted, use GraphicInfo.Height instead.
	Len     int32   // Length of graphic data, it shouldn't be trusted, use GraphicInfo.Len instead.
}

// Graphic stores data for each graphic, not a strict mapping to the file.
type Graphic struct {
	Info        *GraphicInfo // Pointer of GraphicInfo, for reverse searching.
	Header      graphicHeader
	RawData     []byte  // The raw data which read from graphic file.
	GraphicData []byte  // The decoded (if needed) data from RawData
	PaletteLen  int32   // When Version >= 2, read this field from graphic file, it couldn't be set by direct set palette data.
	PaletteData Palette // When Version < 2, set palette data from palette file; otherwise, set palette data from graphic file.
}

// GraphicInfoIndex is a map of GraphicInfo, key is GraphicInfo.ID or GraphicInfo.MapID.
type GraphicInfoIndex map[int]GraphicInfo
