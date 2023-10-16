package tmx

// Property wraps any number of custom properties.
type Property struct {
	Name         string `json:"name"`                   // Name of the property
	Type         string `json:"type"`                   // Type of the property (string (default), int, float, bool, color, file, object or class (since 0.16, with color and file added in 0.17, object added in 1.4 and class added in 1.8))
	PropertyType string `json:"propertytype,omitempty"` // Name of the custom property type, when applicable (since 1.8)
	Value        any    `json:"value"`                  // Value of the property
}
