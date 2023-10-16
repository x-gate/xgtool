package tmx

// ObjectAlignment Controls the alignment for tile objects
type ObjectAlignment string

const (
	Unspecified ObjectAlignment = "unspecified"
	TopLeft     ObjectAlignment = "topleft"
	Top         ObjectAlignment = "top"
	TopRight    ObjectAlignment = "topright"
	Left        ObjectAlignment = "left"
	Center      ObjectAlignment = "center"
	Right       ObjectAlignment = "right"
	BottomLeft  ObjectAlignment = "bottomleft"
	Bottom      ObjectAlignment = "bottom"
	BottomRight ObjectAlignment = "bottomright"
)

// TileSet is a struct that represents the tile set data.
type TileSet struct {
	BackgroundColor  string           `json:"backgroundcolor,omitempty"`  // Hex-formatted color (#RRGGBB or #AARRGGBB) (optional)
	Class            string           `json:"class,omitempty"`            // The class of the TileSet (since 1.9, optional)
	Columns          int              `json:"columns"`                    // The number of tile columns in the TileSet
	FileMode         string           `json:"filemode,omitempty"`         // The fill mode to use when rendering tiles from this TileSet (stretch (default) or preserve-aspect-fit) (since 1.9)
	FirstGID         int              `json:"firstgid"`                   // GID corresponding to the first tile in the set
	Grid             Grid             `json:"grid,omitempty"`             // (optional)
	Image            string           `json:"image,omitempty"`            // Image used for tiles in this set
	ImageHeight      int              `json:"imageheight,omitempty"`      // Height of source image in pixels
	ImageWidth       int              `json:"imagewidth,omitempty"`       // Width of source image in pixels
	Margin           int              `json:"margin"`                     // Buffer between image edge and first tile (pixels)
	Name             string           `json:"name"`                       // Name given to this TileSet
	ObjectAlignment  ObjectAlignment  `json:"objectalignment,omitempty"`  // Alignment to use for tile objects (Unspecified (default), TopLeft, Top, TopRight, Left, Center, Right, BottomLeft, Bottom or BottomRight) (since 1.4)
	Properties       []Property       `json:"properties,omitempty"`       // Array of Property
	Source           string           `json:"source,omitempty"`           // The external file that contains this TileSet data
	Spacing          int              `json:"spacing"`                    // Spacing between adjacent tiles in image (pixels)
	Terrains         []Terrain        `json:"terrains,omitempty"`         // Array of Terrain (optional)
	TileCount        int              `json:"tilecount"`                  // The number of tiles in this TileSet
	TiledVersion     string           `json:"tiledversion"`               // The Tiled version used to save the file
	TileHeight       int              `json:"tileheight"`                 // Maximum height of tiles in this set
	TileOffset       *TileOffset      `json:"tileoffset,omitempty"`       // (optional)
	TileRenderSize   string           `json:"tilerendersize,omitempty"`   // The size to use when rendering tiles from this TileSet on a tile layer (tile (default) or grid) (since 1.9)
	Tiles            []Tile           `json:"tiles,omitempty"`            // Array of Tile (optional)
	TileWidth        int              `json:"tilewidth"`                  // Maximum width of tiles in this set
	Transformations  *Transformations `json:"transformations,omitempty"`  // Allowed transformations (optional)
	TransparentColor string           `json:"transparentcolor,omitempty"` // Hex-formatted color (#RRGGBB) (optional)
	Type             string           `json:"type"`                       // `tileset` (since 1.0)
	Version          string           `json:"version"`                    // The JSON format version (previously a number, saved as string since 1.6)
	WangSets         []WangSet        `json:"wangsets,omitempty"`         // Array of WangSet (since 1.1.5)
}

// Grid is a struct that represents the grid data.
type Grid struct {
	Height      int         `json:"height"`      // Cell height of tile grid
	Width       int         `json:"width"`       // Orientation or Isometric
	Orientation Orientation `json:"orientation"` // Cell width of tile grid
}

// Terrain is a struct that represents the terrain data.
type Terrain struct {
	Name       string
	Properties []Property
	Tile       int
}

// TileOffset is a struct that represents the tile offset data.
type TileOffset struct {
	X int `json:"x"` // Horizontal offset in pixels
	Y int `json:"y"` // Vertical offset in pixels (positive is down)
}

// Tile is a struct that represents the tile data.
type Tile struct {
	Animation   []Frame    `json:"animation,omitempty"`   // Array of Frame
	ID          int        `json:"id"`                    // Local ID of the tile
	Image       string     `json:"image"`                 // Image representing this tile (optional, used for image collection TileSet)
	ImageHeight int        `json:"imageheight,omitempty"` // Height of the tile image in pixels
	ImageWidth  int        `json:"imagewidth,omitempty"`  // Width of the tile image in pixels
	X           int        `json:"x,omitempty"`           // The X position of the sub-rectangle representing this tile (default: 0)
	Y           int        `json:"y,omitempty"`           // The Y position of the sub-rectangle representing this tile (default: 0)
	Width       int        `json:"width,omitempty"`       // The width of the sub-rectangle representing this tile (defaults to the image width)
	Height      int        `json:"height,omitempty"`      // The height of the sub-rectangle representing this tile (defaults to the image height)
	ObjectGroup *Layer     `json:"objectgroup,omitempty"` // Layer with type ObjectGroup, when collision shapes are specified (optional)
	Probability float64    `json:"probability,omitempty"` // Percentage chance this tile is chosen when competing with others in the editor (optional)
	Properties  []Property `json:"properties,omitempty"`  // Array of Properties
	Terrain     []int      `json:"terrain,omitempty"`     // Index of terrain for each corner of tile (optional, replaced by Wang sets since 1.5)
	Type        string     `json:"type,omitempty"`        // The class of the tile (was saved as `class` in 1.9, optional)
}

// Transformations this element is used to describe which transformations can be applied to the tiles (e.g. to extend a Wang set by transforming existing tiles).
type Transformations struct {
	HFlip               bool `json:"hflip,omitempty"`               // Tiles can be flipped horizontally
	VFlip               bool `json:"vflip,omitempty"`               // Tiles can be flipped vertically
	Rotate              bool `json:"rotate,omitempty"`              // Tiles can be rotated in 90-degree increments
	PreferUntransformed bool `json:"preferuntransformed,omitempty"` // Whether untransformed tiles remain preferred, otherwise transformed tiles are used to produce more variations
}

// WangSet contains the list of Wang sets defined for this tileset.
type WangSet struct {
	Class      string      `json:"class,omitempty"`  // The class of the Wang set (since 1.9, optional)
	Colors     []WangColor `json:"colors,omitempty"` // Array of WangColor
	Name       string      `json:"name"`             // Name of the Wang set
	Properties []Property  `json:"properties"`       // Array of Property
	Tile       int         `json:"tile"`             // Local ID of tile representing the Wang set
	Type       string      `json:"type"`             // `corner`, `edge or `mixed` (since 1.5)
	WangTiles  []WangTile  `json:"wangtiles"`        // Array of WangTile
}

// WangColor a color that can be used to define the corner and/or edge of a Wang tile.
type WangColor struct {
	Class       string     `json:"class,omitempty"`       // The class of the Wang color (since 1.9, optional)
	Color       string     `json:"color"`                 // Hex-formatted color (#RRGGBB or #AARRGGBB)
	Name        string     `json:"name"`                  // Name of the Wang color
	Probability float64    `json:"probability,omitempty"` // Probability used when randomizing
	Properties  []Property `json:"properties,omitempty"`  // Array of Property
	Tile        int        `json:"tile"`                  // Local ID of tile representing the Wang color
}

// WangTile defines a Wang tile, by referring to a tile in the tileset and associating it with a certain Wang ID.
type WangTile struct {
	TileID int    `json:"tileid"` // Local ID of tile
	WangID [8]int `json:"wangid"` // Array of Wang color indexes (uchar[8])
}

// Frame is a struct that represents the frame data.
type Frame struct {
	Duration int `json:"duration"` // Frame duration in milliseconds
	TileID   int `json:"tileid"`   // Local tile ID representing this frame
}

// NewTileSet creates a TileSet.
func NewTileSet(name string, fgid int, grid Grid) (ts TileSet) {
	ts.Type = "tileset"
	ts.Name = name
	ts.TiledVersion = TiledVersion
	ts.Version = "0.0"
	ts.Grid = grid
	ts.FirstGID = fgid

	return
}
