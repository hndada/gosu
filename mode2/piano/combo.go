package piano

import "github.com/hndada/gosu/draws"

type Combo struct {
	Sprites []draws.Sprite
	Combo   int
}

// c stands for component.
func (c Combo) Draw(screen draws.Image) {
	c.ComboSprite.Draw(screen, draws.Op{})
}
