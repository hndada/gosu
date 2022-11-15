package play

var (
	scoreScale    float64
	scoreDigitGap float64
	meterWidth    float64 // number of pixels per 1ms
	meterHeight   float64
	offset        int64
)

type Settings struct {
	ScoreScale    float64
	ScoreDigitGap float64
	MeterWidth    float64
	MeterHeight   float64
	Offset        int64
}

func (Settings) Default() Settings {
	return Settings{
		ScoreScale:    0.65,
		ScoreDigitGap: 0,
		MeterWidth:    4,
		MeterHeight:   50,
		Offset:        -135, // -65
	}
}
func (Settings) Current() Settings {
	return Settings{
		ScoreScale:    scoreScale,
		ScoreDigitGap: scoreDigitGap,
		MeterWidth:    meterWidth,
		MeterHeight:   meterHeight,
		Offset:        offset,
	}
}
func (Settings) Set(s Settings) {
	scoreScale = s.ScoreScale
	scoreDigitGap = s.ScoreDigitGap
	meterWidth = s.MeterWidth
	meterHeight = s.MeterHeight
	offset = s.Offset
}
