package piano

import "github.com/hndada/gosu/draws"

type StageComponent struct {
	fieldSprite draws.Sprite
	hintSprite  draws.Sprite
}

type ChartComponent struct {
	barSprite draws.Sprite
	// longNoteBodySprites [][4]draws.Animation
	// noteSprites         [][2]draws.Sprite
}
