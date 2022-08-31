package gosu

// Todo: BarLine color settings
var (
	TPS    = 1000 // TPS should be 1000 or greater.
	Volume = 0.25

	MeterWidth  float64 = 4 // The number of pixels per 1ms.
	MeterHeight float64 = 50
	CursorScale float64 = 0.1

	BgDimness float64 = 0.5

	ChartInfoBoxWidth  float64 = 450
	ChartInfoBoxHeight float64 = 50
	ChartInfoBoxShrink float64 = 0.15

	chartInfoBoxshrink float64 = ChartInfoBoxWidth * ChartInfoBoxShrink

	ScoreScale float64 = 0.65
)

var (
	MusicRoot   = "music"
	WindowSizeX = 1600
	WindowSizeY = 900
)

// Todo: reset all tick-dependent variables.
// They are mostly at drawer.go or play.go, settings.go
// Keyword: TimeToTick
func SetTPS() {

}
