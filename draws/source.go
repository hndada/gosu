package draws

import "github.com/hajimehoshi/ebiten/v2"

type Op = ebiten.DrawImageOptions

// Image, Text, and Blank implement Source.
type Source interface {
	IsValid() bool
	Size() Vector2
	Draw(Image, Op)
}

// Blank is for wrapping Sprite with specific Outer size.
type Blank struct {
	Size_ Vector2
}

func (b Blank) IsValid() bool         { return true }
func (b Blank) Size() Vector2         { return b.Size_ }
func (b Blank) Draw(dst Image, op Op) {}
