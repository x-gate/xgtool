package tmx

// TiledVersion define the version of TMX format.
const TiledVersion = "1.10"

// Orientation define the orientation of the map: allows "orthogonal", "isometric", "staggered" or "hexagonal".
type Orientation string

const (
	Orthogonal Orientation = "orthogonal"
	Isometric  Orientation = "isometric"
	Staggered  Orientation = "staggered"
	Hexagonal  Orientation = "hexagonal"
)

// RenderOrder define the render order of the map: allows "right-down", "right-up", "left-down" or "left-up".
type RenderOrder string

const (
	RightDown RenderOrder = "right-down"
	RightUp   RenderOrder = "right-up"
	LeftDown  RenderOrder = "left-down"
	LeftUp    RenderOrder = "left-up"
)

// Map The entrypoint (root) of TMX format.
type Map struct {
	BackgroundColor  string      `json:"backgroundcolor,omitempty"`  // Hex-formatted color (#RRGGBB or #AARRGGBB) (optional)
	Class            string      `json:"class,omitempty"`            // The class of the map (since 1.9, optional)
	CompressionLevel int         `json:"compressionlevel,omitempty"` // The compression level to use for tile layer data (defaults to -1, which means to use the algorithm default)
	Height           int         `json:"height"`                     // Number of tile rows
	HexSideLength    int         `json:"hexsidelength,omitempty"`    // Length of the side of a hex tile in pixels (hexagonal maps only)
	Infinite         bool        `json:"infinite"`                   // Whether the map has infinite dimensions
	Layers           []Layer     `json:"layers"`                     // Array of Layers
	NextLayerID      int         `json:"nextlayerid,omitempty"`      // Auto-increments for each layer
	NextObjectID     int         `json:"nextobjectid,omitempty"`     // Auto-increments for each placed object
	Orientation      Orientation `json:"orientation"`                // Orthogonal, Isometric, Staggered or Hexagonal
	ParallaxOriginX  float64     `json:"parallaxoriginx,omitempty"`  // X coordinate of the parallax origin in pixels (since 1.8, default: 0)
	ParallaxOriginY  float64     `json:"parallaxoriginy,omitempty"`  // Y coordinate of the parallax origin in pixels (since 1.8, default: 0)
	Properties       []Property  `json:"properties,omitempty"`       // Array of Property
	RenderOrder      RenderOrder `json:"renderorder"`                // RightDown (the default), RightUp, LeftDown or LeftUp (currently only supported for orthogonal maps)
	StaggerAxis      string      `json:"staggeraxis,omitempty"`      // x or y (staggered / hexagonal maps only)
	StaggerIndex     string      `json:"staggerindex,omitempty"`     // odd or even (staggered / hexagonal maps only)
	TiledVersion     string      `json:"tiledversion"`               // The Tiled version used to save the file
	TileHeight       int         `json:"tileheight"`                 // Map grid height
	TileSets         []TileSet   `json:"tilesets"`                   // Array of TileSet
	TileWidth        int         `json:"tilewidth"`                  // Map grid width
	Type             string      `json:"type"`                       // `map` (since 1.0)
	Version          string      `json:"version"`                    // The JSON format version (previously a number, saved as string since 1.6)
	Width            int         `json:"width"`                      // Number of tile columns
}

// NewMap creates a Map.
func NewMap(w, h int, or Orientation, ro RenderOrder) (m Map) {
	m.Type = "map"
	// reverse width and height, because the map will be rotated -90 degrees
	m.Width = w
	m.Height = h
	m.Orientation = or
	m.RenderOrder = ro
	m.TiledVersion = TiledVersion
	m.Version = "0.0"
	return
}
