package tmx

// DrawOrder of object group layer
type DrawOrder string

// Encoding of tile layer
type Encoding string

// LayerType of layers
type LayerType string

const (
	TopDown DrawOrder = "topdown"
	Index   DrawOrder = "index"
)

const (
	CSV    Encoding = "csv"
	Base64 Encoding = "base64"
)

const (
	TileLayer   LayerType = "tilelayer"
	ObjectGroup LayerType = "objectgroup"
	ImageLayer  LayerType = "imagelayer"
	Group       LayerType = "group"
)

// Layer presents layers in a map, including tile layers, object groups, image layers and groups.
type Layer struct {
	Chunks           []Chunk    `json:"chunks,omitempty"`                              // Array of chunks (optional). TileLayer only.
	Class            string     `json:"class,omitempty" xml:"class,omitempty"`         // The class of the layer (since 1.9, defaults to “”).
	Compression      string     `json:"compression,omitempty"`                         // `zlib`, `gzip`, `zstd` (since Tiled 1.3) or empty (default). TileLayer only.
	Data             []uint     `json:"data,omitempty"`                                // Array of unsigned int (GIDs) or base64-encoded data. TileLayer only.
	DrawOrder        DrawOrder  `json:"draworder,omitempty"`                           // TopDown (default) or index. ObjectGroup only.
	Encoding         Encoding   `json:"encoding,omitempty"`                            // CSV (default) or Base64. TileLayer only.
	Height           int        `json:"height,omitempty" xml:"height,omitempty"`       // The height of the layer in tiles. Always the same as the map height for fixed-size maps.
	ID               int        `json:"id" xml:"id"`                                   // Unique ID of the layer (defaults to 0, with valid IDs being at least 1). Each layer that added to a map gets a unique id. Even if a layer is deleted, no layer ever gets the same ID. Can not be changed in Tiled. (since Tiled 1.2)
	Image            string     `json:"image,omitempty"`                               // Image used by this layer. ImageLayer only.
	Layers           []Layer    `json:"layers,omitempty"`                              // Array of Layer. Group only.
	Locked           bool       `json:"locked,omitempty"`                              // Whether layer is locked in the editor (default: false). (since Tiled 1.8.2)
	Objects          []Object   `json:"objects,omitempty"`                             // Array of objects. ObjectGroup only.
	Name             string     `json:"name,omitempty" xml:"name,omitempty"`           // The name of the layer. (defaults to “”)
	OffsetX          float64    `json:"offsetx,omitempty" xml:"offsetx,omitempty"`     // Horizontal offset for this layer in pixels. Defaults to 0. (since 0.14)
	OffsetY          float64    `json:"offsety,omitempty" xml:"offsety,omitempty"`     // Vertical offset for this layer in pixels. Defaults to 0. (since 0.14)
	Opacity          float64    `json:"opacity,omitempty" xml:"opacity,omitempty"`     // The opacity of the layer as a value from 0 to 1. Defaults to 1.
	ParallaxX        float64    `json:"parallaxx,omitempty" xml:"parallaxx,omitempty"` // Horizontal parallax factor for this layer. Defaults to 1. (since 1.5)
	ParallaxY        float64    `json:"parallaxy,omitempty" xml:"parallaxy,omitempty"` // Vertical parallax factor for this layer. Defaults to 1. (since 1.5)
	Properties       []Property `json:"properties,omitempty" xml:"property,omitempty"` // Array of Property
	RepeatX          bool       `json:"repeatx,omitempty"`                             // Whether the image drawn by this layer is repeated along the X axis. ImageLayer only. (since Tiled 1.8)
	RepeatY          bool       `json:"repeaty,omitempty"`                             // Whether the image drawn by this layer is repeated along the Y axis. ImageLayer only. (since Tiled 1.8)
	StartX           int        `json:"startx,omitempty"`                              // X coordinate where layer content starts (for infinite maps)
	StartY           int        `json:"starty,omitempty"`                              // Y coordinate where layer content starts (for infinite maps)
	TintColor        string     `json:"tintcolor,omitempty" xml:"tintcolor,omitempty"` // A tint color that is multiplied with any tiles drawn by this layer in #AARRGGBB or #RRGGBB format (optional).
	TransParentColor string     `json:"transparentcolor,omitempty"`                    // Hex-formatted color (#RRGGBB) (optional). ImageLayer only.
	Type             LayerType  `json:"type"`                                          // TileLayer, ObjectGroup, ImageLayer or Group
	Visible          bool       `json:"visible,omitempty" xml:"visible,omitempty"`     // Whether the layer is shown (1) or hidden (0). Defaults to 1.
	Width            int        `json:"width,omitempty" xml:"width,omitempty"`         // The width of the layer in tiles. Always the same as the map width for fixed-size maps.
	X                int        `json:"x" xml:"x"`                                     // The x coordinate of the layer in tiles. Defaults to 0 and can not be changed in Tiled.
	Y                int        `json:"y" xml:"y"`                                     // The y coordinate of the layer in tiles. Defaults to 0 and can not be changed in Tiled.
}

// NewTileLayer creates a Layer with type TileLayer.
func NewTileLayer(name string, id, w, h int) (l Layer) {
	l.Type = TileLayer
	l.ID = id
	l.Name = name
	l.Width = w
	l.Height = h
	l.Visible = true
	l.Opacity = 1.0

	return
}

// NewObjectLayer creates a Layer with type ObjectGroup.
func NewObjectLayer(name string, id int, do DrawOrder) (l Layer) {
	l.Type = ObjectGroup
	l.ID = id
	l.Name = name
	l.DrawOrder = do
	l.Visible = true
	l.Opacity = 1.0

	return
}
