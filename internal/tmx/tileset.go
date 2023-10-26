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
	BackgroundColor  string           `json:"backgroundcolor,omitempty"`                                 // Hex-formatted color (#RRGGBB or #AARRGGBB) (optional)
	Class            string           `json:"class,omitempty" xml:"class,omitempty"`                     // The class of this tileset (since 1.9, defaults to “”).
	Columns          int              `json:"columns" xml:"columns"`                                     // The number of tile columns in the tileset. For image collection tilesets it is editable and is used when displaying the tileset. (since 0.15)
	FileMode         string           `json:"filemode,omitempty" xml:"fillmode,omitempty"`               // The fill mode to use when rendering tiles from this tileset. Valid values are stretch (the default) and preserve-aspect-fit. Only relevant when the tiles are not rendered at their native size, so this applies to resized tile objects or in combination with tilerendersize set to grid. (since 1.9)
	FirstGID         int              `json:"firstgid" xml:"firstgid"`                                   // The first global tile ID of this tileset (this global ID maps to the first tile in this tileset).
	Grid             Grid             `json:"grid,omitempty" xml:"grid,omitempty"`                       // (optional)
	Image            string           `json:"image,omitempty"`                                           // Image used for tiles in this set
	ImageHeight      int              `json:"imageheight,omitempty"`                                     // Height of source image in pixels
	ImageWidth       int              `json:"imagewidth,omitempty"`                                      // Width of source image in pixels
	Margin           int              `json:"margin,omitempty" xml:"margin,omitempty"`                   // The margin around the tiles in this tileset (applies to the tileset image, defaults to 0). Irrelevant for image collection tilesets.
	Name             string           `json:"name" xml:"name"`                                           //  The name of this tileset.
	ObjectAlignment  ObjectAlignment  `json:"objectalignment,omitempty" xml:"objectalignment,omitempty"` // Controls the alignment for tile objects. Valid values are unspecified, topleft, top, topright, left, center, right, bottomleft, bottom and bottomright. The default value is unspecified, for compatibility reasons. When unspecified, tile objects use bottomleft in orthogonal mode and bottom in isometric mode. (since 1.4)
	Properties       []Property       `json:"properties,omitempty" xml:"property,omitempty"`             // Array of Property
	Source           string           `json:"source,omitempty" xml:"source,omitempty"`                   //  If this tileset is stored in an external TSX (Tile Set XML) file, this attribute refers to that file. That TSX file has the same structure as the <tileset> element described here. (There is the firstgid attribute missing and this source attribute is also not there. These two attributes are kept in the TMX map, since they are map specific.)
	Spacing          int              `json:"spacing,omitempty" xml:"spacing,omitempty"`                 // The spacing in pixels between the tiles in this tileset (applies to the tileset image, defaults to 0). Irrelevant for image collection tilesets.
	Terrains         []Terrain        `json:"terrains,omitempty" xml:"terraintypes,omitempty"`           // Array of Terrain (optional)
	TileCount        int              `json:"tilecount" xml:"tilecount"`                                 // The number of tiles in this tileset (since 0.13). Note that there can be tiles with a higher ID than the tile count, in case the tileset is an image collection from which tiles have been removed.
	TiledVersion     string           `json:"tiledversion"`                                              // The Tiled version used to save the file
	TileHeight       int              `json:"tileheight" xml:"tileheight"`                               // The (maximum) height of the tiles in this tileset. Irrelevant for image collection tilesets, but stores the maximum tile height.
	TileOffset       *TileOffset      `json:"tileoffset,omitempty" xml:"tileoffset,omitempty"`           // (optional)
	TileRenderSize   string           `json:"tilerendersize,omitempty" xml:"tilerendersize,omitempty"`   // The size to use when rendering tiles from this tileset on a tile layer. Valid values are tile (the default) and grid. When set to grid, the tile is drawn at the tile grid size of the map. (since 1.9)
	Tiles            []Tile           `json:"tiles,omitempty" xml:"tile"`                                // Array of Tile (optional)
	TileWidth        int              `json:"tilewidth" xml:"tilewidth"`                                 // The (maximum) width of the tiles in this tileset. Irrelevant for image collection tilesets, but stores the maximum tile width.
	Transformations  *Transformations `json:"transformations,omitempty" xml:"transformations,omitempty"` // Allowed transformations (optional)
	TransparentColor string           `json:"transparentcolor,omitempty"`                                // Hex-formatted color (#RRGGBB) (optional)
	Type             string           `json:"type"`                                                      // `tileset` (since 1.0)
	Version          string           `json:"version"`                                                   // The JSON format version (previously a number, saved as string since 1.6)
	WangSets         []WangSet        `json:"wangsets,omitempty" xml:"wangsets,omitempty"`               // Array of WangSet (since 1.1.5)
}

// Grid is a struct that represents the grid data.
type Grid struct {
	Height      int         `json:"height" xml:"height"`           // Cell height of tile grid
	Width       int         `json:"width" xml:"width"`             // Orientation or Isometric
	Orientation Orientation `json:"orientation" xml:"orientation"` // Cell width of tile grid
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
