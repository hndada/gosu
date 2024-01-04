package draws

import "image/color"

// Filler can realize Background and Shadow.
// Maybe Border too.
// By introducing an image, API becomes much simpler than Web's.
type Color struct{ color.Color }

func (c Color) IsEmpty() bool { return c.Color == nil }
func (c Color) Size() Vector2 { return Vector2{} }
