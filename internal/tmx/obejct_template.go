package tmx

// ObjectTemplate is written to its own file and referenced by any instances of that template.
type ObjectTemplate struct {
	Type    string  `json:"type"`              // `template`
	TileSet TileSet `json:"tileset,omitempty"` // External tileset used by the template (optional)
	Object  Object  `json:"object"`            // The object instantiated by this template
}
