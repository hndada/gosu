package choose

import (
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
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
