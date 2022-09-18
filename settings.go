package gosu

var (
	MusicRoot   = "music"
	WindowSizeX = 1600
	WindowSizeY = 900
)

// TPS supposed to be multiple of 1000, since only one speed value
// goes passed per Update, while unit of TransPoint's time is 1ms.
// TPS affects only on Update(), not on Draw().
// Todo: BarLine color settings
var (
	TPS         int  = 1000 // TPS should be 1000 or greater.
	VsyncSwitch bool = true

	MusicVolume  float64 = 0.25
	EffectVolume float64 = 0.25

	CursorScale float64 = 0.1

	ChartInfoBoxWidth  float64 = 450
	ChartInfoBoxHeight float64 = 50
	ChartInfoBoxShrink float64 = 0.15
	chartInfoBoxshrink float64 = ChartInfoBoxWidth * ChartInfoBoxShrink

	BackgroundDimness float64 = 0.5
	ScoreScale        float64 = 0.65
	ScoreDigitGap     float64 = 0
	MeterWidth        float64 = 4 // The number of pixels per 1ms.
	MeterHeight       float64 = 50
)

// Todo: reset all tick-dependent variables.
// They are mostly at drawer.go or play.go, settings.go
// Keyword: TimeToTick
func SetTPS() {}
