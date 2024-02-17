package types

// MapHeader is the header of a map file, 20 bytes.
type MapHeader struct {
	Magic  [12]byte
	Width  int32
	Height int32
}

// Map structure for each map in `map/[0-9]+.dat`, 20 + width*height*2 (ground) + width*height*2 (object) + width*height*2 (metadata) bytes.
type Map struct {
	MapHeader
	Ground   []uint16
	Object   []uint16
	Metadata []uint16
}
