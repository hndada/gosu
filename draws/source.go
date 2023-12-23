package draws

import "github.com/hajimehoshi/ebiten/v2/colorm"

// type Op = ebiten.DrawImageOptions
type Op struct {
	colorm.DrawImageOptions
	ColorM colorm.ColorM
}

// Image, Text, and Blank implement Source.
type Source interface {
	Size() Vector2
	Draw(Image, Op)
	IsEmpty() bool // Whether the source is nil.
}

// // Blank is for wrapping Sprite with specific Outer size.
// type Blank struct{ Size_ Vector2 }

// func (b Blank) Size() Vector2         { return b.Size_ }
// func (b Blank) Draw(dst Image, op Op) {}
// func (b Blank) IsEmpty() bool         { return false }
