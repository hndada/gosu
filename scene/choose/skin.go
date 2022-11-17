package choose

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
	defaultskin "github.com/hndada/gosu/skin"
)

// There are ChartPanel, ChartBox, ChartLevelBox.
type Skin struct {
	DefaultBackground draws.Sprite
	DefaultChartPanel draws.Sprite
	ChartBox          draws.Sprite
	ChartLevelBox     draws.Sprite
	Sounds
}

var (
	defaultSkin Skin
	skin        Skin
)

func initSkin() {
	LoadSkin(defaultskin.FS, scene.LoadSkinDefault)
}
func DefaultSkin() Skin { return defaultSkin }
func CurrentSkin() Skin { return skin }

func LoadSkin(fsys fs.FS, mode scene.LoadSkinMode) {
	// {
	// 	s := draws.NewSprite(fsys, "box-mask.png")
	// 	s.SetSize(ChartInfoBoxWidth, ChartInfoBoxHeight)
	// 	s.Locate(ScreenSizeX+chartInfoBoxshrink, ScreenSizeY/2, draws.RightMiddle)
	// 	ChartItemBoxSprite = s
	// }
}
