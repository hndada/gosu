package choose

import (
	"image/color"

	"github.com/hndada/gosu/assets"
	"github.com/hndada/gosu/draws"
)

// Load box-mask.png from assets
var boxMask = draws.LoadSprite(assets.FS, "box-mask.png") // interface/

// scene/asset.go (BaseScene)
// type BaseScene = Scene
func NewSearchDrawer(query *string) SearchDrawer {
	const (
		x = ScreenSizeX - RowWidth
		y = 25
	)
	i := draws.NewImage(RowWidth, 50)
	i.Fill(color.NRGBA{153, 217, 234, 192})
	s := draws.NewSprite(i)
	s.Locate(x, y, draws.LeftTop)
}
