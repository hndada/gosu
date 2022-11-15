package piano

import "github.com/hndada/gosu/input"

// Logical size of in-game screen.
const (
	ScreenSizeX = mode.ScreenSizeX
	ScreenSizeY = mode.ScreenSizeY
)

var SpeedScale float64 = 1.0

var KeySettings = map[int][]input.Key{
	4:               {input.KeyD, input.KeyF, input.KeyJ, input.KeyK},
	5:               {input.KeyD, input.KeyF, input.KeySpace, input.KeyJ, input.KeyK},
	6:               {input.KeyS, input.KeyD, input.KeyF, input.KeyJ, input.KeyK, input.KeyL},
	7:               {input.KeyS, input.KeyD, input.KeyF, input.KeySpace, input.KeyJ, input.KeyK, input.KeyL},
	8 + LeftScratch: {input.KeyA, input.KeyS, input.KeyD, input.KeyF, input.KeySpace, input.KeyJ, input.KeyK, input.KeyL},
	8:               {input.KeyA, input.KeyS, input.KeyD, input.KeyF, input.KeyJ, input.KeyK, input.KeyL, input.KeySemicolon},
	9:               {input.KeyA, input.KeyS, input.KeyD, input.KeyF, input.KeySpace, input.KeyJ, input.KeyK, input.KeyL, input.KeySemicolon},
	10:              {input.KeyA, input.KeyS, input.KeyD, input.KeyF, input.KeyV, input.KeyN, input.KeyJ, input.KeyK, input.KeyL, input.KeySemicolon},
}
var NoteWidthsMap = map[int][3]float64{
	4:  {0.065, 0.065, 0.065},
	5:  {0.065, 0.065, 0.065},
	6:  {0.065, 0.065, 0.065},
	7:  {0.06, 0.06, 0.06},
	8:  {0.06, 0.06, 0.06},
	9:  {0.06, 0.06, 0.06},
	10: {0.06, 0.06, 0.06},
}

// Todo: generalize setting loading function
func init() {
	for k, ws := range NoteWidthsMap {
		ws2 := ws
		for i, w := range ws2 {
			ws2[i] = ScreenSizeX * w
		}
		NoteWidthsMap[k] = ws2
	}
}

// Todo: add note lighting color settings per kind
// Todo: Should NoteHeight be separated into NoteHeight, HeadHeight, TailHeight?
var (
	FieldDarkness float64 = 0.8 // Todo: FieldDarkness -> FieldOpaque
	FieldPosition float64 = ScreenSizeX * 0.5

	HitPosition float64 = ScreenSizeY * 0.90 // The bottom y-value of Hint,  not a middle or top.

	// positionMargin should be larger than MaxSize/2 of all note sprites' width or height.
	positionMargin float64 = 100
	maxPosition    float64 = HitPosition + positionMargin
	minPosition    float64 = HitPosition - ScreenSizeY - positionMargin

	NoteHeigth    float64 = ScreenSizeY * 0.05 // Applies to all notes
	TailExtraTime float64 = 0
	// bodyLoss   float64 = NoteHeigth // Head/2 + Tail/2.

	ComboPosition    float64 = ScreenSizeY * 0.40
	JudgmentPosition float64 = ScreenSizeY * 0.66
)

const (
	BodyStyleStretch = iota
	BodyStyleAttach
)

// Skin-dependent settings.
// Todo: make SkinScaleSettings struct?
var (
	BodyStyle   int  = BodyStyleStretch
	ReverseBody bool = false

	ScoreScale        float64 = 0.65
	ComboScale        float64 = 0.75
	ComboDigitGap     float64 = ScreenSizeX * -0.0008
	JudgmentScale     float64 = 0.33
	HintHeight        float64 = ScreenSizeY * 0.04
	LightingScale     float64 = 1.0
	KeyLightingOpaque float64 = 0.5
	HitLightingOpaque float64 = 1 // Todo: set color per note kind
)

func SwitchDirection() {
	max, min := maxPosition, minPosition
	maxPosition = -min
	minPosition = -max
	ReverseBody = !ReverseBody
}
