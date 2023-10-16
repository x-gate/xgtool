package tmx

// Chunk are used to store the tile layer data for infinite maps.
type Chunk struct {
	Data   []uint `json:"data"`   // Array of unsigned int (GIDs) or base64-encoded data
	Height int    `json:"height"` // Height in tiles
	Width  int    `json:"width"`  // Width in tiles
	X      int    `json:"x"`      // X coordinate in tiles
	Y      int    `json:"y"`      // Y coordinate in tiles
}
