package tmx

// Property wraps any number of custom properties.
type Property struct {
	Name         string `json:"name" xml:"name"`                                     // Name of the property
	Type         string `json:"type" xml:"type"`                                     // The type of the property. Can be string (default), int, float, bool, color, file, object or class (since 0.16, with color and file added in 0.17, object added in 1.4 and class added in 1.8).
	PropertyType string `json:"propertytype,omitempty" xml:"propertytype,omitempty"` // The name of the custom property type, when applicable (since 1.8).
	Value        any    `json:"value" xml:"value"`                                   // The value of the property. (default string is “”, default number is 0, default boolean is “false”, default color is #00000000, default file is “.” (the current file’s parent directory))
}
