package types

// MapHeader is the header of a map file, 20 bytes.
type MapHeader struct {
	Magic  [12]byte // "MAP" + (0x1) + (0x0*8)
	Width  int32
	Height int32
}

// MapData structure for each map in `map/[0-9]+.dat`, 20 + width*height*2 (ground) + width*height*2 (object) + width*height*2 (metadata) bytes.
type MapData struct {
	MapHeader
	Ground   []uint16 // Ground layer, 2 bytes per tile
	Object   []uint16 // Object layer, 2 bytes per tile
	Metadata []uint16 // Metadata layer, 2 bytes per tile
}
