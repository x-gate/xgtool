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
	Chunks           []Chunk    `json:"chunks,omitempty"`           // Array of chunks (optional). TileLayer only.
	Class            string     `json:"class,omitempty"`            // The class of the layer (since 1.9, optional)
	Compression      string     `json:"compression,omitempty"`      // `zlib`, `gzip`, `zstd` (since Tiled 1.3) or empty (default). TileLayer only.
	Data             []uint     `json:"data,omitempty"`             // Array of unsigned int (GIDs) or base64-encoded data. TileLayer only.
	DrawOrder        DrawOrder  `json:"draworder,omitempty"`        // TopDown (default) or index. ObjectGroup only.
	Encoding         Encoding   `json:"encoding,omitempty"`         // CSV (default) or Base64. TileLayer only.
	Height           int        `json:"height,omitempty"`           // Row count. Same as map height for fixed-size maps. TileLayer only.
	ID               int        `json:"id"`                         // Incremental ID - unique across all layers
	Image            string     `json:"image,omitempty"`            // Image used by this layer. ImageLayer only.
	Layers           []Layer    `json:"layers,omitempty"`           // Array of Layer. Group only.
	Locked           bool       `json:"locked,omitempty"`           // Whether layer is locked in the editor (default: false). (since Tiled 1.8.2)
	Name             string     `json:"name"`                       // Name assigned to this layer
	Objects          []Object   `json:"objects,omitempty"`          // Array of objects. ObjectGroup only.
	OffsetX          float64    `json:"offsetx,omitempty"`          // Horizontal layer offset in pixels (default: 0)
	OffsetY          float64    `json:"offsety,omitempty"`          // Vertical layer offset in pixels (default: 0)
	Opacity          float64    `json:"opacity"`                    // Value between 0 and 1
	ParallaxX        float64    `json:"parallaxx,omitempty"`        // Horizontal parallax factor for this layer (default: 1). (since Tiled 1.5)
	ParallaxY        float64    `json:"parallaxy,omitempty"`        // Vertical parallax factor for this layer (default: 1). (since Tiled 1.5)
	Properties       []Property `json:"properties,omitempty"`       // Array of Property
	RepeatX          bool       `json:"repeatx,omitempty"`          // Whether the image drawn by this layer is repeated along the X axis. ImageLayer only. (since Tiled 1.8)
	RepeatY          bool       `json:"repeaty,omitempty"`          // Whether the image drawn by this layer is repeated along the Y axis. ImageLayer only. (since Tiled 1.8)
	StartX           int        `json:"startx,omitempty"`           // X coordinate where layer content starts (for infinite maps)
	StartY           int        `json:"starty,omitempty"`           // Y coordinate where layer content starts (for infinite maps)
	TintColor        string     `json:"tintcolor,omitempty"`        // Hex-formatted [tint color](https://doc.mapeditor.org/en/stable/manual/layers/#tint-color) (#RRGGBB or #AARRGGBB) that is multiplied with any graphics drawn by this layer or any child layers (optional).
	TransParentColor string     `json:"transparentcolor,omitempty"` // Hex-formatted color (#RRGGBB) (optional). ImageLayer only.
	Type             LayerType  `json:"type"`                       // TileLayer, ObjectGroup, ImageLayer or Group
	Visible          bool       `json:"visible"`                    // Whether layer is shown or hidden in editor
	Width            int        `json:"width,omitempty"`            // Column count. Same as map width for fixed-size maps. TileLayer only.
	X                int        `json:"x"`                          // Horizontal layer offset in tiles. Always 0.
	Y                int        `json:"y"`                          // Vertical layer offset in tiles. Always 0.
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
