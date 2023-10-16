package tmx

// Object is a struct that represents the object data.
type Object struct {
	Ellipse    bool       `json:"ellipse,omitempty"`    // Used to mark an object as an ellipse
	GID        int        `json:"gid"`                  // Global tile ID, only if object represents a tile
	Height     float64    `json:"height"`               // Height in pixels.
	ID         int        `json:"id"`                   // Incremental ID, unique across all objects
	Name       string     `json:"name"`                 // String assigned to name field in editor
	Point      bool       `json:"point,omitempty"`      // Used to mark an object as a point
	Polygon    []Point    `json:"polygon,omitempty"`    // Array of Point, in case the object is a polygon
	Polyline   []Point    `json:"polyline,omitempty"`   // Array of Point, in case the object is a polyline
	Properties []Property `json:"properties,omitempty"` // Array of Property
	Rotation   float64    `json:"rotation,omitempty"`   // Angle in degrees clockwise
	Template   string     `json:"template,omitempty"`   // Reference to a template file, in case object is a template instance
	Text       *Text      `json:"text,omitempty"`       // Only used for text objects
	Type       string     `json:"type"`                 // The class of the object (was saved as class in 1.9, optional)
	Visible    bool       `json:"visible"`              // Whether object is shown in editor.
	Width      float64    `json:"width"`                // Width in pixels.
	X          float64    `json:"x"`                    // X coordinate in pixels
	Y          float64    `json:"y"`                    // Y coordinate in pixels
}

// Text is a struct that represents the text data.
type Text struct {
	Bold       bool   `json:"bold,omitempty"`      // Whether to use a bold font (default: false)
	Color      string `json:"color"`               // Hex-formatted color (#RRGGBB or #AARRGGBB) (default: #000000)
	FontFamily string `json:"fontfamily"`          // Font family (default: sans-serif)
	HAlign     string `json:"halign"`              // Horizontal alignment (center, right, justify or left (default))
	Italic     bool   `json:"italic,omitempty"`    // Whether to use an italic font (default: false)
	Kerning    bool   `json:"kerning,omitempty"`   // Whether to use kerning when placing characters (default: true)
	PixelSize  int    `json:"pixelsize"`           // Pixel size of the font (default: 16)
	Strikeout  bool   `json:"strikeout,omitempty"` // Whether to strike out the text (default: false)
	Text       string `json:"text"`                // String of text
	Underline  bool   `json:"underline,omitempty"` // Whether to underline the text (default: false)
	VAlign     string `json:"valign"`              // Vertical alignment (center, bottom or top (default))
	Wrap       bool   `json:"wrap,omitempty"`      // Whether the text is wrapped within the object bounds (default: false)
}

// Point is a struct that represents the point data.
type Point struct {
	X float64 `json:"x"` // X coordinate in pixels
	Y float64 `json:"y"` // Y coordinate in pixels
}

// NewObject creates an Object.
func NewObject(gid, id int, w, h float64) (o Object) {
	o.GID = gid
	o.ID = id
	o.Width = w
	o.Height = h
	o.Visible = true

	return
}
