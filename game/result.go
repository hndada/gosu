package game

import "time"

// Todo: SceneResult. keep playing music when at SceneResult.
// Todo: implement Replay
type Result struct {
	MD5        [16]byte  // MD5 for raw chart file. md5.Size = 16
	PlayedTime time.Time // Finish time of playing.

	ScoreFactors   [3]float64 // Retrieved from the chart.
	Scores         [4]float64
	JudgmentCounts []int
	MaxCombo       int
	// FlowMarks      []float64 // Length is around 100 ~ 200.
	// KeyLogs []KeyLog // Entire timed-log key strokes.
}

func (s Scorer) NewResult(md5 [16]byte) Result {
	return Result{
		MD5:            md5,
		PlayedTime:     time.Now(),
		ScoreFactors:   s.ScoreFactors,
		Scores:         s.Scores,
		JudgmentCounts: s.JudgmentCounts,
		MaxCombo:       s.MaxCombo,
	}
}
