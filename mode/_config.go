package mode

import "github.com/hndada/gosu/draws"

type Config interface {
	ScreenSize() draws.Vector2
	ScoreScale() float64
}
