package draws

import "image/color"

type Color struct{ color.Color }

func NewColor(c color.Color) Color { return Color{c} }
func (c Color) IsEmpty() bool      { return c.Color == nil }
func (c Color) Size() Vector2      { return Vector2{100, 100} }
