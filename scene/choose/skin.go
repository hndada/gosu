package choose

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/scene"
)

const (
	TPS         = scene.TPS
	ScreenSizeX = scene.ScreenSizeX
	ScreenSizeY = scene.ScreenSizeY
)

type Skin struct {
	Type              int
	DefaultBackground draws.Sprite
	DefaultChartPanel draws.Sprite
	ChartBox          draws.Sprite
	ChartLevelBox     draws.Sprite
}

const (
	chartBoxWidth  = 450
	chartBoxHeight = 50
	chartBoxShrink = 0.15 * chartBoxWidth
	chartBoxCount  = ScreenSizeY/chartBoxHeight + 2
)

var (
	DefaultSkin = Skin{Type: mode.Default}
	UserSkin    = Skin{Type: mode.User}
)

func (skin *Skin) Load(fsys fs.FS) {
	skin.DefaultBackground = scene.UserSkin.DefaultBackground
	{
		s := draws.NewSprite(fsys, "interface/box-mask.png")
		s.SetSize(chartBoxWidth, chartBoxHeight)
		x := ScreenSizeX - chartBoxShrink
		y := float64(ScreenSizeY) / 2
		s.Locate(x, y, draws.RightMiddle)
		skin.ChartBox = s
	}
	base := []Skin{{}, DefaultSkin, UserSkin}[skin.Type]
	skin.fillBlank(base)
}
func (skin *Skin) fillBlank(base Skin) {
	if !skin.DefaultBackground.IsValid() {
		skin.DefaultBackground = base.DefaultBackground
	}
	if !skin.DefaultChartPanel.IsValid() {
		skin.DefaultChartPanel = base.DefaultChartPanel
	}
	if !skin.ChartBox.IsValid() {
		skin.ChartBox = base.ChartBox
	}
	if !skin.ChartLevelBox.IsValid() {
		skin.ChartLevelBox = base.ChartLevelBox
	}
}
