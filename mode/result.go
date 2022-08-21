package mode

import (
	"time"
)

// Todo: SceneResult
// Todo: implement Replay
type Result struct {
	MD5        [16]byte  // MD5 for raw chart file. md5.Size = 16
	PlayedTime time.Time // Finish time of playing.

	FinalScore float64
	FlowScore  float64
	AccScore   float64
	ExtraScore float64

	JudgmentCounts []int
	MaxNoteWeights float64 // Total NoteWeights. Works as Upper bound.
	Flows          float64 // Sum of Flow
	Accs           float64
	Extras         float64 // Kool rate, for example.
	MaxCombo       int

	FlowMarks []float64 // Length is around 100 ~ 200.
	// KeyLogs []KeyLog // Entire timed-log key strokes.
}
