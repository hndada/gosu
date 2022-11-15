package options

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
)

// All settings (or options) can be configured on this scene.
type Option struct {
	Target interface{}
}

var options2 = map[string]draws.Vector2{
	"800x600":    draws.Vec2(800, 600),
	"1600x900":   draws.Vec2(1600, 900),
	"Fullscreen": draws.IntVec2(ebiten.ScreenSizeInFullscreen()),
}

func set() {
	s := scene.Settings{}
	current := s.Current()
	current.(scene.Settings).WindowSizeX = 800
}

// type KeyVector2 struct {
// 	Key string
// 	draws.Vector2
// }

// var options = []KeyVector2{
// 	{
// 		Key:     "800x600",
// 		Vector2: draws.Vec2(800, 600),
// 	},
// 	{
// 		Key:     "1600x900",
// 		Vector2: draws.Vec2(1600, 900),
// 	},
// 	{
// 		Key:     "Fullscreen",
// 		Vector2: draws.IntVec2(ebiten.ScreenSizeInFullscreen()),
// 	},
// }
