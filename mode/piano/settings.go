package piano

import (
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/input"
)

// Logical size of in-game screen
const (
	screenSizeX = gosu.ScreenSizeX
	screenSizeY = gosu.ScreenSizeY
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
var NoteWidthsMap = map[int][4]float64{
	4:  {0.065, 0.065, 0.065, 0.065},
	5:  {0.065, 0.065, 0.065, 0.065},
	6:  {0.065, 0.065, 0.065, 0.065},
	7:  {0.06, 0.06, 0.06, 0.06},
	8:  {0.06, 0.06, 0.06, 0.06},
	9:  {0.06, 0.06, 0.06, 0.06},
	10: {0.06, 0.06, 0.06, 0.06},
}

// Todo: generalize setting loading function
func init() {
	for k, ws := range NoteWidthsMap {
		ws2 := ws
		for i, w := range ws2 {
			ws2[i] = screenSizeX * w
		}
		NoteWidthsMap[k] = ws2
	}
}

var PositionMargin float64 = 100 // It should be larger than MaxSize/2 of all note sprites' width or height.
// Todo: Should NoteHeight be separated into NoteHeight, HeadHeight, TailHeight?
var (
	FieldDarkness float64 = 0.8

	HitPosition float64 = screenSizeY * 0.90 // The bottom y-value of Hint,  not a middle or top.
	maxPosition float64 = HitPosition + PositionMargin
	minPosition float64 = HitPosition - screenSizeY - PositionMargin

	NoteHeigth float64 = screenSizeY * 0.05 // Applies to all notes
	bodyLoss   float64 = NoteHeigth         // Head/2 + Tail/2.

	ComboPosition    float64 = screenSizeY * 0.40
	JudgmentPosition float64 = screenSizeY * 0.66
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

	ScoreScale    float64 = 0.65
	ComboScale    float64 = 0.75
	ComboDigitGap float64 = screenSizeX * -0.0008
	JudgmentScale float64 = 0.33
	HintHeight    float64 = screenSizeY * 0.04
)

func SwitchDirection() {
	max, min := maxPosition, minPosition
	maxPosition = -min
	minPosition = -max
	ReverseBody = !ReverseBody
}
