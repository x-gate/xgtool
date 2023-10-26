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
	BackgroundColor  string      `json:"backgroundcolor,omitempty" xml:"backgroundcolor,omitempty"`   // The background color of the map. (optional, may include alpha value since 0.15 in the form #AARRGGBB. Defaults to fully transparent.)
	Class            string      `json:"class,omitempty" xml:"class,omitempty"`                       // The class of this map (since 1.9, defaults to "").
	CompressionLevel int         `json:"compressionlevel,omitempty" xml:"compressionlevel,omitempty"` // The compression level to use for tile layer data (defaults to -1, which means to use the algorithm default).
	Height           int         `json:"height" xml:"height"`                                         // The map height in tiles.
	HexSideLength    int         `json:"hexsidelength,omitempty" xml:"hexsidelength,omitempty"`       // Only for Hexagonal maps. Determines the width or height (depending on the staggered axis) of the tile’s edge, in pixels.
	Infinite         bool        `json:"infinite,omitempty" xml:"infinite,omitempty"`                 // Whether this map is infinite. An infinite map has no fixed size and can grow in all directions. Its layer data is stored in chunks. (0 for false, 1 for true, defaults to 0)
	Layers           []Layer     `json:"layers" xml:"layer"`                                          // Array of Layers
	NextLayerID      int         `json:"nextlayerid,omitempty" xml:"nextlayerid,omitempty"`           // Stores the next available ID for new layers. This number is stored to prevent reuse of the same ID after layers have been removed. (since 1.2) (defaults to the highest layer id in the file + 1)
	NextObjectID     int         `json:"nextobjectid,omitempty" xml:"nextobjectid,omitempty"`         // Stores the next available ID for new objects. This number is stored to prevent reuse of the same ID after objects have been removed. (since 0.11) (defaults to the highest object id in the file + 1)
	Orientation      Orientation `json:"orientation" xml:"orientation"`                               // Map orientation. Tiled supports Orthogonal, Isometric, Staggered and Hexagonal (since 0.11).
	ParallaxOriginX  float64     `json:"parallaxoriginx,omitempty" xml:"parallaxoriginx,omitempty"`   // X coordinate of the parallax origin in pixels (defaults to 0). (since 1.8)
	ParallaxOriginY  float64     `json:"parallaxoriginy,omitempty" xml:"parallaxoriginy,omitempty"`   // Y coordinate of the parallax origin in pixels (defaults to 0). (since 1.8)
	Properties       []Property  `json:"properties,omitempty" xml:"properties,omitempty"`             // Array of Property
	RenderOrder      RenderOrder `json:"renderorder,omitempty" xml:"renderorder,omitempty"`           // The order in which tiles on tile layers are rendered. Valid values are RightDown (the default), RightUp, LeftDown and LeftUp. In all cases, the map is drawn row-by-row. (only supported for orthogonal maps at the moment)
	StaggerAxis      string      `json:"staggeraxis,omitempty" xml:"staggeraxis,omitempty"`           // For Staggered and Hexagonal maps, determines which axis (“x” or “y”) is staggered. (since 0.11)
	StaggerIndex     string      `json:"staggerindex,omitempty" xml:"staggerindex,omitempty"`         // For Staggered and Hexagonal maps, determines whether the “even” or “odd” indexes along the staggered axis are shifted. (since 0.11)
	TiledVersion     string      `json:"tiledversion" xml:"tiledversion"`                             // The Tiled version used to save the file (since Tiled 1.0.1). Maybe a date (for snapshot builds). (optional)
	TileHeight       int         `json:"tileheight" xml:"tileheight"`                                 // The height of a tile.
	TileSets         []TileSet   `json:"tilesets" xml:"tileset"`                                      // Array of TileSet
	TileWidth        int         `json:"tilewidth" xml:"tilewidth"`                                   // The width of a tile.
	Type             string      `json:"type"`                                                        // `map` (since 1.0)
	Version          string      `json:"version" xml:"version"`                                       // The TMX format version. Was “1.0” so far, and will be incremented to match minor Tiled releases.
	Width            int         `json:"width" xml:"width"`                                           // The map width in tiles.
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
