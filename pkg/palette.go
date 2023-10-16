package pkg

import "image/color"

// Palette is a collection of colors, it usually 768 bytes (256 * 3), but not always.
type Palette []color.Color
