package draws

// Blank is for wrapping Sprite with specific Outer size.
type Blank struct {
	Size_ Vector2
}

func (b Blank) IsValid() bool         { return true }
func (b Blank) Size() Vector2         { return b.Size_ }
func (b Blank) Draw(dst Image, op Op) {}
