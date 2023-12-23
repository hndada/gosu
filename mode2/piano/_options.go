// Stage is a virtual component: no resources nor components.
type Options struct {
	Screen mode.ScreenOpts
	Music  mode.MusicOpts
	Sound  mode.SoundOpts

	// Common options
	KeyCount    int
	KeyMappings map[int][]string
	KeyOrders   map[int][]KeyKind
	Scratches   map[int]Scratch
	KeyRWs      [4]float64
	StageRWs    map[int]float64
	StageRX     float64
	BaselineRY  float64

	// Screen mode.ScreenOpts
	screen draws.WHXY
	stage  draws.WHXY
	key    draws.WHXY

	Field    FieldOpts
	Hint     HintOpts
	Bar      BarOpts
	Note     NoteOpts
	Light    LightOpts
	Judgment JudgmentOpts
	Combo    mode.ComboOpts
	Score    mode.ScoreOpts
}