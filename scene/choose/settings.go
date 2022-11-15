package choose

import "github.com/hndada/gosu/scene"

const ()

var (
	BackgroundBrightness = 0.6

	ChartInfoBoxWidth  float64 = 450
	ChartInfoBoxHeight float64 = 50
	ChartInfoBoxShrink float64 = 0.15

	chartInfoBoxshrink float64 = ChartInfoBoxWidth * ChartInfoBoxShrink
	chartItemBoxCount  int     = int(scene.ScreenSizeY/ChartInfoBoxHeight) + 2 // Gives some margin.
)
