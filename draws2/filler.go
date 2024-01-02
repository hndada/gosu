package draws

import "image/color"

// Filler can realize Background and Shadow.
// Maybe Border too.
// By introducing an image, API becomes much simpler than Web's.
type Filler struct{ Color color.Color }

func (f Filler) IsEmpty() bool { return f.Color == nil }
func (f Filler) Size() Vector2 { return Vector2{} }
