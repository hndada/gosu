package choose

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
	defaultskin "github.com/hndada/gosu/skin"
)

const (
	chartBoxWidth  = 450
	chartBoxHeight = 50
	chartBoxshrink = 0.15 * chartBoxWidth
	chartBoxCount  = scene.ScreenSizeY/chartBoxHeight + 2
)

type skinType struct {
	DefaultBackground draws.Sprite
	ChartBox          draws.Sprite
	ChartLevelBox     draws.Sprite
	DefaultChartPanel draws.Sprite
}

var (
	defaultSkin skinType
	Skin        skinType
)

func init() { LoadSkin(defaultskin.FS, scene.LoadSkinDefault) }
func LoadSkin(fsys fs.FS, mode scene.LoadSkinMode) {
	skin := &Skin
	if mode == scene.LoadSkinDefault {
		skin = &defaultSkin
	}
	skin.DefaultBackground = scene.Skin.DefaultBackground
	{
		s := draws.NewSprite(fsys, "interface/box-mask.png")
		s.SetSize(chartBoxWidth, chartBoxHeight)
		x := scene.ScreenSizeX + chartBoxshrink
		y := float64(scene.ScreenSizeY) / 2
		s.Locate(x, y, draws.RightMiddle)
		skin.ChartBox = s
	}
	skin.fillBlank([]skinType{{}, defaultSkin, Skin}[mode])
}
func (skin *skinType) fillBlank(base skinType) {
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
