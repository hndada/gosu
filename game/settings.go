package game

var (
	MusicRoot   = "music"
	WindowSizeX = 1600
	WindowSizeY = 900
)
var (
	// TPS supposed to be multiple of 1000, since only one speed value
	// goes passed per Update, while unit of TransPoint's time is 1ms.
	// TPS affects only on Update(), not on Draw().
	TPS int = 1000 // TPS should be 1000 or greater.

	CursorScale        float64 = 0.1
	ChartInfoBoxWidth  float64 = 450
	ChartInfoBoxHeight float64 = 50
	ChartInfoBoxShrink float64 = 0.15
	chartInfoBoxshrink float64 = ChartInfoBoxWidth * ChartInfoBoxShrink
	chartItemBoxCount  int     = int(ScreenSizeY/ChartInfoBoxHeight) + 2 // Gives some margin.

	ScoreScale    float64 = 0.65
	ScoreDigitGap float64 = 0
	MeterWidth    float64 = 4 // The number of pixels per 1ms.
	MeterHeight   float64 = 50

	MusicVolume          float64 = 0.25
	EffectVolume         float64 = 0.25
	BackgroundBrightness float64 = 0.6

	Offset int64 = -135 //-65
)

// Todo: reset all tick-dependent variables.
// They are mostly at drawer.go or play.go, settings.go
// Keyword: ToTick
func SetTPS() {}